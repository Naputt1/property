import * as pulumi from "@pulumi/pulumi";
import * as k8s from "@pulumi/kubernetes";
import { Namespace, Service } from "@pulumi/kubernetes/core/v1";

export function createObservability(
  ns: Namespace,
  config: pulumi.Config,
  services: { backend: Service; redis: Service },
  dependOn?: any[],
) {
  const redisHost = pulumi.interpolate`${services.redis.metadata.name}.${ns.metadata.name}.svc.cluster.local`;
  const backendHost = pulumi.interpolate`${services.backend.metadata.name}.${ns.metadata.name}.svc.cluster.local`;

  const asynqmonLabels = { app: "asynqmon" };
  const asynqmonDeployment = new k8s.apps.v1.Deployment(
    "asynqmon",
    {
      metadata: { namespace: ns.metadata.name },
      spec: {
        selector: { matchLabels: asynqmonLabels },
        template: {
          metadata: { labels: asynqmonLabels },
          spec: {
            containers: [
              {
                name: "asynqmon",
                image: "naputt/asynqmon:latest",
                args: [pulumi.interpolate`--redis-addr=${redisHost}:6379`],
                ports: [{ containerPort: 8080 }],
              },
            ],
          },
        },
      },
    },
    { dependsOn: dependOn },
  );

  const asynqmonService = new k8s.core.v1.Service("asynqmon", {
    metadata: { namespace: ns.metadata.name },
    spec: {
      type: "LoadBalancer",
      ports: [{ port: 8082, targetPort: 8080 }],
      selector: asynqmonLabels,
    },
  });

  // Prometheus persistence
  const prometheusPvc = new k8s.core.v1.PersistentVolumeClaim(
    "prometheus-pvc",
    {
      metadata: { namespace: ns.metadata.name, name: "prometheus-pvc" },
      spec: {
        accessModes: ["ReadWriteOnce"],
        storageClassName: "longhorn",
        resources: { requests: { storage: "10Gi" } },
      },
    },
  );

  const prometheusConfig = new k8s.core.v1.ConfigMap("prometheus-config", {
    metadata: { namespace: ns.metadata.name },
    data: {
      "prometheus.yml": pulumi.interpolate`
global:
  scrape_interval: 5s

scrape_configs:
  - job_name: "backend"
    metrics_path: /metrics
    static_configs:
      - targets:
          - ${backendHost}:${services.backend.spec.ports[0].port}
`,
    },
  });

  const prometheusLabels = { app: "prometheus" };
  const prometheusDeployment = new k8s.apps.v1.Deployment(
    "prometheus",
    {
      metadata: { namespace: ns.metadata.name },
      spec: {
        selector: { matchLabels: prometheusLabels },
        template: {
          metadata: { labels: prometheusLabels },
          spec: {
            securityContext: {
              fsGroup: 65534,
              fsGroupChangePolicy: "OnRootMismatch",
            },
            containers: [
              {
                name: "prometheus",
                image: "prom/prometheus:latest",
                args: [
                  "--config.file=/etc/prometheus/prometheus.yml",
                  "--storage.tsdb.path=/prometheus",
                ],
                ports: [{ containerPort: 9090 }],
                volumeMounts: [
                  { name: "config", mountPath: "/etc/prometheus" },
                  { name: "prometheus-data", mountPath: "/prometheus" },
                ],
              },
            ],
            volumes: [
              {
                name: "config",
                configMap: { name: prometheusConfig.metadata.name },
              },
              {
                name: "prometheus-data",
                persistentVolumeClaim: {
                  claimName: prometheusPvc.metadata.name,
                },
              },
            ],
          },
        },
      },
    },
    { dependsOn: dependOn },
  );

  const prometheusService = new k8s.core.v1.Service("prometheus", {
    metadata: { namespace: ns.metadata.name },
    spec: {
      type: "LoadBalancer",
      ports: [{ port: 9090, targetPort: 9090 }],
      selector: prometheusLabels,
    },
  });

  // Loki
  const lokiPvc = new k8s.core.v1.PersistentVolumeClaim("loki-pvc", {
    metadata: { namespace: ns.metadata.name, name: "loki-pvc" },
    spec: {
      accessModes: ["ReadWriteOnce"],
      storageClassName: "longhorn",
      resources: { requests: { storage: "10Gi" } },
    },
  });

  const lokiConfig = new k8s.core.v1.ConfigMap("loki-config", {
    metadata: { namespace: ns.metadata.name },
    data: {
      "local-config.yaml": `
auth_enabled: false
server:
  http_listen_port: 3100
common:
  instance_addr: 127.0.0.1
  path_prefix: /loki
  storage:
    filesystem:
      chunks_directory: /loki/chunks
      rules_directory: /loki/rules
  replication_factor: 1
  ring:
    kvstore:
      store: inmemory
schema_config:
  configs:
    - from: 2020-10-24
      store: tsdb
      object_store: filesystem
      schema: v13
      index:
        prefix: index_
        period: 24h
`,
    },
  });

  const lokiLabels = { app: "loki" };
  const lokiDeployment = new k8s.apps.v1.Deployment(
    "loki",
    {
      metadata: { namespace: ns.metadata.name },
      spec: {
        selector: { matchLabels: lokiLabels },
        template: {
          metadata: { labels: lokiLabels },
          spec: {
            securityContext: {
              fsGroup: 10001,
              fsGroupChangePolicy: "OnRootMismatch",
            },
            containers: [
              {
                name: "loki",
                image: "grafana/loki:latest",
                args: ["-config.file=/etc/loki/local-config.yaml"],
                ports: [{ containerPort: 3100 }],
                volumeMounts: [
                  { name: "config", mountPath: "/etc/loki" },
                  { name: "loki-data", mountPath: "/loki" },
                ],
              },
            ],
            volumes: [
              { name: "config", configMap: { name: lokiConfig.metadata.name } },
              {
                name: "loki-data",
                persistentVolumeClaim: { claimName: lokiPvc.metadata.name },
              },
            ],
          },
        },
      },
    },
    { dependsOn: dependOn },
  );

  const lokiService = new k8s.core.v1.Service("loki", {
    metadata: { namespace: ns.metadata.name },
    spec: {
      type: "LoadBalancer",
      ports: [{ port: 3100, targetPort: 3100 }],
      selector: lokiLabels,
    },
  });

  // Grafana provisioning
  const grafanaDatasourcesConfig = new k8s.core.v1.ConfigMap(
    "grafana-datasources",
    {
      metadata: { namespace: ns.metadata.name },
      data: {
        "datasources.yml": pulumi.interpolate`
apiVersion: 1
datasources:
  - name: Prometheus
    type: prometheus
    access: proxy
    url: http://prometheus.${ns.metadata.name}.svc.cluster.local:9090
    isDefault: true
    editable: true
  - name: Loki
    type: loki
    access: proxy
    url: http://loki.${ns.metadata.name}.svc.cluster.local:3100
    editable: true
`,
      },
    },
  );

  const grafanaDashboardsProviderConfig = new k8s.core.v1.ConfigMap(
    "grafana-dashboards-provider",
    {
      metadata: { namespace: ns.metadata.name },
      data: {
        "dashboard-provider.yml": `
apiVersion: 1
providers:
  - name: 'Default'
    orgId: 1
    folder: ''
    type: file
    disableDeletion: false
    editable: true
    options:
      path: /var/lib/grafana/dashboards
`,
      },
    },
  );

  const grafanaDashboardsConfig = new k8s.core.v1.ConfigMap(
    "grafana-dashboards",
    {
      metadata: { namespace: ns.metadata.name },
      data: {
        "system-overview.json": JSON.stringify({
          annotations: {
            list: [
              {
                builtIn: 1,
                datasource: { type: "datasource", uid: "grafana" },
                enable: true,
                hide: true,
                iconColor: "rgba(0, 211, 255, 1)",
                name: "Annotations & Alerts",
                type: "dashboard",
              },
            ],
          },
          editable: true,
          fiscalYearStartMonth: 0,
          graphTooltip: 0,
          id: 1,
          links: [],
          liveNow: false,
          panels: [
            {
              datasource: { type: "prometheus", uid: "Prometheus" },
              fieldConfig: {
                defaults: {
                  color: { mode: "palette-classic" },
                  custom: {
                    axisCenteredZero: false,
                    axisColorMode: "text",
                    axisLabel: "",
                    axisPlacement: "auto",
                    barAlignment: 0,
                    drawStyle: "line",
                    fillOpacity: 0,
                    gradientMode: "none",
                    hideFrom: { legend: false, tooltip: false, viz: false },
                    insertNulls: false,
                    lineInterpolation: "linear",
                    lineWidth: 1,
                    pointSize: 5,
                    scaleDistribution: { type: "linear" },
                    showPoints: "auto",
                    spanNulls: false,
                    stacking: { group: "A", mode: "none" },
                    thresholdsStyle: { mode: "off" },
                  },
                  mappings: [],
                  thresholds: {
                    mode: "absolute",
                    steps: [
                      { color: "green", value: null },
                      { color: "red", value: 80 },
                    ],
                  },
                },
                overrides: [],
              },
              gridPos: { h: 8, w: 12, x: 0, y: 0 },
              id: 1,
              options: {
                legend: {
                  calcs: [],
                  displayMode: "list",
                  placement: "bottom",
                  showLegend: true,
                },
                tooltip: { mode: "single", sort: "none" },
              },
              targets: [
                {
                  datasource: { type: "prometheus", uid: "Prometheus" },
                  editorMode: "code",
                  expr: "rate(gin_http_requests_total[5m])",
                  format: "time_series",
                  range: true,
                  refId: "A",
                },
              ],
              title: "Requests per Second",
              type: "timeseries",
            },
            {
              datasource: { type: "prometheus", uid: "Prometheus" },
              fieldConfig: {
                defaults: {
                  color: { mode: "palette-classic" },
                  custom: {
                    axisCenteredZero: false,
                    axisColorMode: "text",
                    axisLabel: "",
                    axisPlacement: "auto",
                    barAlignment: 0,
                    drawStyle: "line",
                    fillOpacity: 0,
                    gradientMode: "none",
                    hideFrom: { legend: false, tooltip: false, viz: false },
                    insertNulls: false,
                    lineInterpolation: "linear",
                    lineWidth: 1,
                    pointSize: 5,
                    scaleDistribution: { type: "linear" },
                    showPoints: "auto",
                    spanNulls: false,
                    stacking: { group: "A", mode: "none" },
                    thresholdsStyle: { mode: "off" },
                  },
                  mappings: [],
                  thresholds: {
                    mode: "absolute",
                    steps: [
                      { color: "green", value: null },
                      { color: "red", value: 80 },
                    ],
                  },
                },
                overrides: [],
              },
              gridPos: { h: 8, w: 12, x: 12, y: 0 },
              id: 2,
              options: {
                legend: {
                  calcs: [],
                  displayMode: "list",
                  placement: "bottom",
                  showLegend: true,
                },
                tooltip: { mode: "single", sort: "none" },
              },
              targets: [
                {
                  datasource: { type: "prometheus", uid: "Prometheus" },
                  editorMode: "code",
                  expr: "sum(rate(gin_http_request_duration_seconds_sum[5m])) / sum(rate(gin_http_request_duration_seconds_count[5m]))",
                  format: "time_series",
                  range: true,
                  refId: "A",
                },
              ],
              title: "Average Response Latency",
              type: "timeseries",
            },
            {
              datasource: { type: "loki", uid: "Loki" },
              gridPos: { h: 12, w: 24, x: 0, y: 8 },
              id: 3,
              options: {
                dedupStrategy: "none",
                enableLogDetails: true,
                prettifyLogMessage: false,
                showCommonLabels: false,
                showLabels: false,
                showTime: true,
                sortOrder: "Descending",
                wrapLogMessage: false,
              },
              targets: [
                {
                  datasource: { type: "loki", uid: "Loki" },
                  expr: '{job="backend"} | json',
                  refId: "A",
                },
              ],
              title: "Backend Logs",
              type: "logs",
            },
          ],
          schemaVersion: 38,
          style: "dark",
          tags: [],
          templating: { list: [] },
          time: { from: "now-6h", to: "now" },
          timepicker: {},
          timezone: "",
          title: "System Overview",
          uid: "system-overview",
          version: 1,
          weekStart: "",
        }),
      },
    },
  );

  // Grafana persistence
  const grafanaPvc = new k8s.core.v1.PersistentVolumeClaim("grafana-pvc", {
    metadata: { namespace: ns.metadata.name, name: "grafana-pvc" },
    spec: {
      accessModes: ["ReadWriteOnce"],
      storageClassName: "longhorn",
      resources: { requests: { storage: "2Gi" } },
    },
  });

  const grafanaLabels = { app: "grafana" };
  const grafanaDeployment = new k8s.apps.v1.Deployment(
    "grafana",
    {
      metadata: { namespace: ns.metadata.name },
      spec: {
        selector: { matchLabels: grafanaLabels },
        template: {
          metadata: { labels: grafanaLabels },
          spec: {
            securityContext: {
              fsGroup: 472, // Grafana default user group
              fsGroupChangePolicy: "OnRootMismatch",
            },
            containers: [
              {
                name: "grafana",
                image: "grafana/grafana:latest",
                env: [
                  { name: "GF_SECURITY_ADMIN_USER", value: "admin" },
                  { name: "GF_SECURITY_ADMIN_PASSWORD", value: "admin" },
                ],
                ports: [{ containerPort: 3000 }],
                volumeMounts: [
                  { name: "grafana-data", mountPath: "/var/lib/grafana" },
                  {
                    name: "datasources",
                    mountPath: "/etc/grafana/provisioning/datasources",
                  },
                  {
                    name: "dashboards-provider",
                    mountPath: "/etc/grafana/provisioning/dashboards",
                  },
                  {
                    name: "dashboards",
                    mountPath: "/var/lib/grafana/dashboards",
                  },
                ],
              },
            ],
            volumes: [
              {
                name: "grafana-data",
                persistentVolumeClaim: { claimName: grafanaPvc.metadata.name },
              },
              {
                name: "datasources",
                configMap: { name: grafanaDatasourcesConfig.metadata.name },
              },
              {
                name: "dashboards-provider",
                configMap: {
                  name: grafanaDashboardsProviderConfig.metadata.name,
                },
              },
              {
                name: "dashboards",
                configMap: { name: grafanaDashboardsConfig.metadata.name },
              },
            ],
          },
        },
      },
    },
    { dependsOn: dependOn },
  );

  const grafanaService = new k8s.core.v1.Service("grafana", {
    metadata: { namespace: ns.metadata.name },
    spec: {
      type: "LoadBalancer",
      ports: [{ port: 3000, targetPort: 3000 }],
      selector: grafanaLabels,
    },
  });

  return {
    asynqmonService,
    prometheusService,
    grafanaService,
    lokiService,
  };
}
