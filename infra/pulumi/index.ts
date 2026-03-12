import * as pulumi from "@pulumi/pulumi";
import * as k8s from "@pulumi/kubernetes";
import { createPostgres } from "./src/pg";
import { createRustFS } from "./src/rustfs";
import { createRedis } from "./src/redis";
import { createBackend } from "./src/backend";
import { createObservability } from "./src/observability";
import { createCloudflareTunnel } from "./src/cloudflare";
import { createIngress } from "./src/ingress";

// Config
const config = new pulumi.Config();

const namespaceName = "property";

// Create Namespace
const ns = new k8s.core.v1.Namespace(namespaceName, {
  metadata: { name: namespaceName },
});

// Redis Deployment
const { service: redisService, deployment: redisDeployment } = createRedis(
  ns,
  config,
);

// Postgres Deployment
const { service: postgresService, deployment: postgresDeployment } =
  createPostgres(ns, config);

// RustFS Deployment
const { service: rustfsService, deployment: rustfsDeployment } = createRustFS(
  ns,
  config,
);

// Backend Deployment
const { service: backendService, deployment: backendDeployment } =
  createBackend(
    ns,
    config,
    {
      postgres: postgresService,
      rustfs: rustfsService,
      redis: redisService,
    },
    [
      redisService,
      redisDeployment,
      postgresService,
      postgresDeployment,
      rustfsService,
      rustfsDeployment,
    ],
  );

// Ingress
const { apiIngress, frontendIngress, rustfsApiIngress } = createIngress(ns, {
  backend: backendService,
  rustfs: rustfsService,
});

// Observability Deployment
const { asynqmonService, prometheusService, grafanaService } =
  createObservability(
    ns,
    config,
    {
      backend: backendService,
      redis: redisService,
    },
    [backendService, backendDeployment],
  );

// Cloudflare Tunnel
const tunnelResources = createCloudflareTunnel(ns, config);

// Outputs
export const backendUrl = "https://property.napnap.work/api";
export const frontendUrl = "https://property.napnap.work";
export const rustfsUploadUrl = "http://ras-pi.tail0684eb.ts.net";

export const rustfsPort = rustfsService.spec.ports.apply((ports) => {
  const apiPort = ports.find((p) => p.name === "api");

  return apiPort?.port;
});

export const rustfsIP = rustfsService.spec.clusterIP;

export const asynqmonUrl = asynqmonService.status.loadBalancer.ingress[0].ip;
export const prometheusUrl = prometheusService.status.loadBalancer.ingress[0].ip;
export const grafanaUrl = grafanaService.status.loadBalancer.ingress[0].ip;

export const tunnelId = tunnelResources.then((r) => r.tunnelId);
export const rustfsAccessKey = "rustfsadmin";
export const rustfsPassword = config.requireSecret("rustfsPassword");
