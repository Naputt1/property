import * as pulumi from "@pulumi/pulumi";
import * as k8s from "@pulumi/kubernetes";
import { Namespace } from "@pulumi/kubernetes/core/v1";

export function createPostgres(ns: Namespace, config: pulumi.Config) {
  const dbPassword = config.requireSecret("dbPassword");

  const postgresPvc = new k8s.core.v1.PersistentVolumeClaim("postgres-pvc", {
    metadata: {
      namespace: ns.metadata.name,
      name: "postgres-pvc",
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

  const postgresLabels = { app: "postgres" };
  const postgresDeployment = new k8s.apps.v1.Deployment("postgres", {
    metadata: { namespace: ns.metadata.name },
    spec: {
      selector: { matchLabels: postgresLabels },
      template: {
        metadata: { labels: postgresLabels },
        spec: {
          containers: [
            {
              name: "postgres",
              image: "postgres:15-alpine",
              env: [
                { name: "POSTGRES_USER", value: "postgres" },
                { name: "POSTGRES_PASSWORD", value: dbPassword },
                { name: "POSTGRES_DB", value: "property" },
              ],
              ports: [{ containerPort: 5432 }],
              volumeMounts: [
                {
                  name: "postgres-data",
                  mountPath: "/var/lib/postgresql/data",
                },
              ],
            },
          ],
          volumes: [
            {
              name: "postgres-data",
              persistentVolumeClaim: {
                claimName: postgresPvc.metadata.name,
              },
            },
          ],
        },
      },
    },
  });

  const postgresService = new k8s.core.v1.Service("postgres", {
    metadata: { namespace: ns.metadata.name },
    spec: {
      ports: [{ port: 5432, targetPort: 5432 }],
      selector: postgresLabels,
    },
  });

  return { service: postgresService, deployment: postgresDeployment };
}
