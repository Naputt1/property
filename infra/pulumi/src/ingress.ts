import * as pulumi from "@pulumi/pulumi";
import * as k8s from "@pulumi/kubernetes";
import { Namespace, Service } from "@pulumi/kubernetes/core/v1";

export function createIngress(
  ns: Namespace,
  services: { backend: Service; rustfs: Service },
) {
  const hosts = ["property.napnap.work", "ras-pi.tail0684eb.ts.net"];

  // Middleware to prepend /property to the path for RustFS
  const frontendPrefixMiddleware = new k8s.apiextensions.CustomResource(
    "frontend-prefix",
    {
      apiVersion: "traefik.io/v1alpha1",
      kind: "Middleware",
      metadata: {
        namespace: ns.metadata.name,
        name: "frontend-prefix",
      },
      spec: {
        addPrefix: {
          prefix: "/property",
        },
      },
    },
  );

  const bucketPath = {
    path: "/",
    pathType: "Prefix",
    backend: {
      service: {
        name: services.rustfs.metadata.name,
        port: { number: 9000 },
      },
    },
  };

  // Ingress for API and WebSockets (No Middleware)
  const apiIngress = new k8s.networking.v1.Ingress("property-api-ingress", {
    metadata: {
      namespace: ns.metadata.name,
      annotations: {
        "kubernetes.io/ingress.class": "traefik",
        "traefik.ingress.kubernetes.io/router.entrypoints": "web",
      },
    },
    spec: {
      rules: hosts.map((host) => ({
        host: host,
        http: {
          paths: [
            {
              path: "/api",
              pathType: "Prefix",
              backend: {
                service: {
                  name: services.backend.metadata.name,
                  port: { number: 8080 },
                },
              },
            },
            {
              path: "/ws",
              pathType: "Prefix",
              backend: {
                service: {
                  name: services.backend.metadata.name,
                  port: { number: 8080 },
                },
              },
            },
          ],
        },
      })),
    },
  });

  // Ingress for Frontend (With /property prefix middleware)
  const frontendIngress = new k8s.networking.v1.Ingress(
    "property-frontend-ingress",
    {
      metadata: {
        namespace: ns.metadata.name,
        annotations: {
          "kubernetes.io/ingress.class": "traefik",
          "traefik.ingress.kubernetes.io/router.entrypoints": "web",
          // Format: <namespace>-<name>@kubernetescrd
          "traefik.ingress.kubernetes.io/router.middlewares": pulumi.interpolate`${ns.metadata.name}-frontend-prefix@kubernetescrd`,
        },
      },
      spec: {
        rules: [
          ...hosts.map((host) => ({
            host: host,
            http: {
              paths: [bucketPath],
            },
          })),
          {
            host: "rustfs.property.ras-pi.tail0684eb.ts.net",
            http: {
              paths: [bucketPath],
            },
          },
        ],
      },
    },
  );

  return { apiIngress, frontendIngress, frontendPrefixMiddleware };
}
