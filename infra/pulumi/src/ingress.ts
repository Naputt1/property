import * as pulumi from "@pulumi/pulumi";
import * as k8s from "@pulumi/kubernetes";
import { Namespace, Service } from "@pulumi/kubernetes/core/v1";

export function createIngress(
  ns: Namespace,
  services: { backend: Service; rustfs: Service },
) {
  const frontendHosts = ["property.napnap.work"];
  const rustfsHost = "rustfs.property.ras-pi.tail0684eb.ts.net";

  const spaFallback = new k8s.apiextensions.CustomResource("spa-fallback", {
    apiVersion: "traefik.io/v1alpha1",
    kind: "Middleware",
    metadata: {
      namespace: ns.metadata.name,
      name: "spa-fallback",
    },
    spec: {
      replacePathRegex: {
        regex: "^/$|^/[^.]+$",
        replacement: "/index.html",
      },
    },
  });

  const bucketPrefix = new k8s.apiextensions.CustomResource("bucket-prefix", {
    apiVersion: "traefik.io/v1alpha1",
    kind: "Middleware",
    metadata: {
      namespace: ns.metadata.name,
      name: "bucket-prefix",
    },
    spec: {
      addPrefix: {
        prefix: "/property",
      },
    },
  });

  const cacheHeaders = new k8s.apiextensions.CustomResource("frontend-cache", {
    apiVersion: "traefik.io/v1alpha1",
    kind: "Middleware",
    metadata: {
      namespace: ns.metadata.name,
      name: "frontend-cache",
    },
    spec: {
      headers: {
        customResponseHeaders: {
          "Cache-Control": "public, max-age=3600",
        },
      },
    },
  });

  const frontendIngress = new k8s.networking.v1.Ingress("property-frontend", {
    metadata: {
      namespace: ns.metadata.name,
      annotations: {
        "kubernetes.io/ingress.class": "traefik",
        "traefik.ingress.kubernetes.io/router.entrypoints": "web",
        "traefik.ingress.kubernetes.io/router.middlewares": pulumi.interpolate`${ns.metadata.name}-spa-fallback@kubernetescrd,${ns.metadata.name}-bucket-prefix@kubernetescrd`,
      },
    },
    spec: {
      rules: frontendHosts.map((host) => ({
        host,
        http: {
          paths: [
            {
              path: "/",
              pathType: "Prefix",
              backend: {
                service: {
                  name: services.rustfs.metadata.name,
                  port: { number: 9000 },
                },
              },
            },
          ],
        },
      })),
    },
  });

  const rustfsIngress = new k8s.networking.v1.Ingress("property-rustfs", {
    metadata: {
      namespace: ns.metadata.name,
      annotations: {
        "kubernetes.io/ingress.class": "traefik",
        "traefik.ingress.kubernetes.io/router.entrypoints": "web",
      },
    },
    spec: {
      rules: [
        {
          host: rustfsHost,
          http: {
            paths: [
              {
                path: "/",
                pathType: "Prefix",
                backend: {
                  service: {
                    name: services.rustfs.metadata.name,
                    port: { number: 9000 },
                  },
                },
              },
            ],
          },
        },
      ],
    },
  });

  const apiIngress = new k8s.networking.v1.Ingress("property-api", {
    metadata: {
      namespace: ns.metadata.name,
      annotations: {
        "kubernetes.io/ingress.class": "traefik",
        "traefik.ingress.kubernetes.io/router.entrypoints": "web",
      },
    },
    spec: {
      rules: frontendHosts.map((host) => ({
        host,
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
              path: "/swagger",
              pathType: "Prefix",
              backend: {
                service: {
                  name: services.backend.metadata.name,
                  port: { number: 8080 },
                },
              },
            },
            {
              path: "/playground",
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

  return {
    apiIngress,
    frontendIngress,
    rustfsIngress,
  };
}
