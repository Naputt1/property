import * as k8s from "@pulumi/kubernetes";
import { Namespace, Service } from "@pulumi/kubernetes/core/v1";

export function createIngress(
  ns: Namespace,
  services: { backend: Service; rustfs: Service },
) {
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

  const commonPaths = [
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
    bucketPath,
  ];

  const ingress = new k8s.networking.v1.Ingress("property-ingress", {
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
          host: "property.napnap.work",
          http: {
            paths: commonPaths,
          },
        },
        {
          host: "rustfs.property.ras-pi.tail0684eb.ts.net",
          http: {
            paths: [bucketPath],
          },
        },
      ],
    },
  });

  return ingress;
}
