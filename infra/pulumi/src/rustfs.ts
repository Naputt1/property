import * as pulumi from "@pulumi/pulumi";
import * as k8s from "@pulumi/kubernetes";
import { Namespace } from "@pulumi/kubernetes/core/v1";

export function createRustFS(ns: Namespace, config: pulumi.Config) {
  const rustfsPassword = config.requireSecret("rustfsPassword");

  const rustfsPvc = new k8s.core.v1.PersistentVolumeClaim("rustfs-pvc", {
    metadata: {
      namespace: ns.metadata.name,
      name: "rustfs-pvc",
    },
    spec: {
      accessModes: ["ReadWriteOnce"],
      storageClassName: "longhorn",
      resources: {
        requests: {
          storage: "10Gi",
        },
      },
    },
  });

  const rustfsLabels = { app: "rustfs" };
  const rustfsDeployment = new k8s.apps.v1.Deployment("rustfs", {
    metadata: { namespace: ns.metadata.name },
    spec: {
      selector: { matchLabels: rustfsLabels },
      template: {
        metadata: { labels: rustfsLabels },
        spec: {
          containers: [
            {
              name: "rustfs",
              image: "rustfs/rustfs:latest",
              args: ["/data"],
              env: [
                { name: "RUSTFS_ACCESS_KEY", value: "rustfsadmin" },
                { name: "RUSTFS_SECRET_KEY", value: rustfsPassword },
              ],
              ports: [{ containerPort: 9000 }, { containerPort: 9001 }],
            },
          ],
          volumes: [
            {
              name: "rustfs-data",
              persistentVolumeClaim: {
                claimName: rustfsPvc.metadata.name,
              },
            },
          ],
        },
      },
    },
  });

  const rustfsService = new k8s.core.v1.Service("rustfs", {
    metadata: { namespace: ns.metadata.name },
    spec: {
      type: "LoadBalancer",
      ports: [
        { name: "api", port: 9000, targetPort: 9000 },
        { name: "console", port: 9001, targetPort: 9001 },
      ],
      selector: rustfsLabels,
    },
  });

  return { service: rustfsService, deployment: rustfsDeployment };
}
