import * as pulumi from "@pulumi/pulumi";
import * as k8s from "@pulumi/kubernetes";
import { Namespace, Service } from "@pulumi/kubernetes/core/v1";

export function createBackend(
  ns: Namespace,
  config: pulumi.Config,
  services: { postgres: Service; rustfs: Service; redis: Service },
  dependOn?: any[],
) {
  const secretKey = config.requireSecret("secretKey");
  const dbPassword = config.requireSecret("dbPassword"); // postgres
  const rustfsPassword = config.requireSecret("rustfsPassword");
  const turnstileSiteKey = config.requireSecret("turnstileSiteKey");
  const turnstileSecretKey = config.requireSecret("turnstileSecretKey");

  const dbHost = pulumi.interpolate`${services.postgres.metadata.name}.${ns.metadata.name}.svc.cluster.local`;
  const rustfsHost = pulumi.interpolate`${services.rustfs.metadata.name}.${ns.metadata.name}.svc.cluster.local`;
  const redisHost = pulumi.interpolate`${services.redis.metadata.name}.${ns.metadata.name}.svc.cluster.local`;

  const backendLabels = { app: "backend" };
  const backendDeployment = new k8s.apps.v1.Deployment(
    "backend",
    {
      metadata: { namespace: ns.metadata.name },
      spec: {
        selector: { matchLabels: backendLabels },
        template: {
          metadata: { labels: backendLabels },
          spec: {
            containers: [
              {
                name: "backend",
                image: "naputt/git:property-manage-backend-latest",
                imagePullPolicy: "Always",
                env: [
                  { name: "PORT", value: "8080" },
                  { name: "IS_PROD", value: "true" },
                  { name: "SECRET_KEY", value: secretKey },

                  { name: "DATABASE_HOST", value: dbHost },
                  { name: "DATABASE_PORT", value: "5432" },
                  { name: "DATABASE_USERNAME", value: "postgres" },
                  { name: "DATABASE_PASSWORD", value: dbPassword },
                  { name: "DATABASE_NAME", value: "property" },

                  {
                    name: "REDIS_URL",
                    value: pulumi.interpolate`${redisHost}:6379`,
                  },
                  {
                    name: "ASYNQ_URL",
                    value: pulumi.interpolate`${redisHost}:6379`,
                  },

                  {
                    name: "RUSTFS_URL",
                    value: pulumi.interpolate`http://${rustfsHost}:${services.rustfs.spec.ports[0].port}`,
                  },
                  { name: "RUSTFS_ACCESS_KEY", value: "rustfsadmin" },
                  { name: "RUSTFS_SECRET_KEY", value: rustfsPassword },

                  { name: "TURNSTILE_SITE_KEY", value: turnstileSiteKey },
                  { name: "TURNSTILE_SECRET_KEY", value: turnstileSecretKey },
                ],
                ports: [{ containerPort: 8080 }],
              },
            ],
          },
        },
      },
    },
    { dependsOn: dependOn },
  );

  const backendService = new k8s.core.v1.Service("backend", {
    metadata: { namespace: ns.metadata.name },
    spec: {
      type: "LoadBalancer",
      ports: [{ port: 8080, targetPort: 8080 }],
      selector: backendLabels,
    },
  });

  return { service: backendService, deployment: backendDeployment };
}
