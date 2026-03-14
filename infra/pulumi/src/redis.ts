import * as pulumi from "@pulumi/pulumi";
import * as k8s from "@pulumi/kubernetes";
import { Namespace } from "@pulumi/kubernetes/core/v1";

export function createRedis(ns: Namespace, config: pulumi.Config) {
  const redisPvc = new k8s.core.v1.PersistentVolumeClaim("redis-pvc", {
    metadata: {
      namespace: ns.metadata.name,
      name: "redis-pvc",
    },
    spec: {
      accessModes: ["ReadWriteOnce"],
      storageClassName: "longhorn",
      resources: {
        requests: {
          storage: "2Gi",
        },
      },
    },
  });

  const redisLabels = { app: "redis" };
  const redisDeployment = new k8s.apps.v1.Deployment("redis", {
    metadata: { namespace: ns.metadata.name },
    spec: {
      selector: { matchLabels: redisLabels },
      template: {
        metadata: { labels: redisLabels },
        spec: {
          containers: [
            {
              name: "redis",
              image: "redis:7-alpine",
              ports: [{ containerPort: 6379 }],
              volumeMounts: [
                {
                  name: "redis-data",
                  mountPath: "/data",
                },
              ],
            },
          ],
          volumes: [
            {
              name: "redis-data",
              persistentVolumeClaim: {
                claimName: redisPvc.metadata.name,
              },
            },
          ],
        },
      },
    },
  });

  const redisService = new k8s.core.v1.Service("redis", {
    metadata: { namespace: ns.metadata.name },
    spec: {
      ports: [{ port: 6379, targetPort: 6379 }],
      selector: redisLabels,
    },
  });

  return { service: redisService, deployment: redisDeployment };
}
