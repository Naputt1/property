import * as pulumi from "@pulumi/pulumi";
import * as k8s from "@pulumi/kubernetes";
import { Namespace, Service } from "@pulumi/kubernetes/core/v1";

export function createIngress(
  ns: Namespace,
  services: { backend: Service; rustfs: Service },
) {
  const frontendHosts = ["property.napnap.work", "ras-pi.tail0684eb.ts.net"];
  const rustfsHost = "rustfs.property.ras-pi.tail0684eb.ts.net";

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

  // Middleware to rewrite / to /index.html for RustFS
  const frontendIndexMiddleware = new k8s.apiextensions.CustomResource(
    "frontend-index",
    {
      apiVersion: "traefik.io/v1alpha1",
      kind: "Middleware",
      metadata: {
        namespace: ns.metadata.name,
        name: "frontend-index",
      },
      spec: {
        replacePath: {
          path: "/index.html",
        },
      },
    },
  );

  // Middleware to set caching headers for frontend
  const frontendCacheMiddleware = new k8s.apiextensions.CustomResource(
    "frontend-cache",
    {
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
    },
  );

  // Middleware to set NO-CACHE headers for index.html
  const frontendNoCacheMiddleware = new k8s.apiextensions.CustomResource(
    "frontend-no-cache",
    {
      apiVersion: "traefik.io/v1alpha1",
      kind: "Middleware",
      metadata: {
        namespace: ns.metadata.name,
        name: "frontend-no-cache",
      },
      spec: {
        headers: {
          customResponseHeaders: {
            "Cache-Control": "no-cache, no-store, must-revalidate",
          },
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
      rules: frontendHosts.map((host) => ({
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

  // Ingress for Frontend INDEX (No caching for index.html)
  const frontendIndexIngress = new k8s.networking.v1.Ingress(
    "property-frontend-index-ingress",
    {
      metadata: {
        namespace: ns.metadata.name,
        annotations: {
          "kubernetes.io/ingress.class": "traefik",
          "traefik.ingress.kubernetes.io/router.entrypoints": "web",
          // Priority: Traefik prioritizes Exact over Prefix, but we can also set priority explicitly if needed
          "traefik.ingress.kubernetes.io/router.middlewares": pulumi.interpolate`${ns.metadata.name}-frontend-prefix@kubernetescrd,${ns.metadata.name}-frontend-index@kubernetescrd,${ns.metadata.name}-frontend-no-cache@kubernetescrd`,
        },
      },
      spec: {
        rules: frontendHosts.flatMap((host) => [
          {
            host: host,
            http: {
              paths: [
                {
                  path: "/",
                  pathType: "Exact",
                  backend: bucketPath.backend,
                },
                {
                  path: "/index.html",
                  pathType: "Exact",
                  backend: bucketPath.backend,
                },
              ],
            },
          },
        ]),
      },
    },
  );

  // Ingress for Frontend Assets (With caching)
  const frontendAssetsIngress = new k8s.networking.v1.Ingress(
    "property-frontend-assets-ingress",
    {
      metadata: {
        namespace: ns.metadata.name,
        annotations: {
          "kubernetes.io/ingress.class": "traefik",
          "traefik.ingress.kubernetes.io/router.entrypoints": "web",
          "traefik.ingress.kubernetes.io/router.middlewares": pulumi.interpolate`${ns.metadata.name}-frontend-prefix@kubernetescrd,${ns.metadata.name}-frontend-cache@kubernetescrd`,
        },
      },
      spec: {
        rules: frontendHosts.map((host) => ({
          host: host,
          http: {
            paths: [
              {
                path: "/",
                pathType: "Prefix",
                backend: bucketPath.backend,
              },
            ],
          },
        })),
      },
    },
  );

  // Ingress for RustFS API (No Middleware)
  // This is used for uploading files (S3 API)
  const rustfsApiIngress = new k8s.networking.v1.Ingress(
    "property-rustfs-api-ingress",
    {
      metadata: {
        namespace: ns.metadata.name,
        annotations: {
          "kubernetes.io/ingress.class": "traefik",
          "traefik.ingress.kubernetes.io/router.entrypoints": "web",
          // Set a higher priority to ensure this matches before other catch-all rules if they exist
          "traefik.ingress.kubernetes.io/router.priority": "100",
        },
      },
      spec: {
        rules: [
          {
            host: "rustfs.property.ras-pi.tail0684eb.ts.net",
            http: {
              paths: [
                {
                  path: "/",
                  pathType: "Prefix",
                  backend: bucketPath.backend,
                },
              ],
            },
          },
        ],
      },
    },
  );

  return {
    apiIngress,
    frontendIndexIngress,
    frontendAssetsIngress,
    rustfsApiIngress,
    frontendPrefixMiddleware,
    frontendIndexMiddleware,
    frontendCacheMiddleware,
    frontendNoCacheMiddleware,
  };
}
