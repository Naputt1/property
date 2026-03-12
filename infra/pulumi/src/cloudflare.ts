import * as pulumi from "@pulumi/pulumi";
import * as cloudflare from "@pulumi/cloudflare";
import * as k8s from "@pulumi/kubernetes";
import { Namespace } from "@pulumi/kubernetes/core/v1";

export async function createCloudflareTunnel(
  ns: Namespace,
  config: pulumi.Config,
  backendService: k8s.core.v1.Service,
) {
  // Use .require() instead of .requireSecret() for IDs to allow using them in get() calls.
  // These values must be provided as plain strings in the Pulumi config.
  const accountId = config.require("cloudflareAccountId");
  const zoneId = config.require("cloudflareZoneId");
  const tunnelName = config.require("cloudflareTunnelname");
  const tunnelId = config.require("cloudflareTunnelId");

  const cfConfig = new pulumi.Config("cloudflare");
  const apiKey = cfConfig.requireSecret("apiKey");
  const email = cfConfig.require("email");

  const provider = new cloudflare.Provider("cloudflare-provider", {
    apiKey: apiKey,
    email: email,
  });

  // 1. Create the Zero Trust Tunnel (Reference)
  const tunnel = cloudflare.ZeroTrustTunnelCloudflared.get(
    "backend-tunnel",
    pulumi.interpolate`${accountId}/${tunnelId}`,
    {},
    { provider: provider },
  );

  const tunnelToken = cloudflare.getZeroTrustTunnelCloudflaredTokenOutput(
    {
      accountId: accountId,
      tunnelId: tunnelId,
    },
    { provider: provider },
  );

  // 2. Create the DNS record pointing to the tunnel
  // Check if it exists to avoid "Record already exists" error
  const recordHostname = "property.napnap.work";
  const existingRecords = await cloudflare
    .getDnsRecords(
      {
        zoneId: zoneId,
        name: { exact: recordHostname },
      },
      { provider: provider },
    )
    .catch(() => undefined);

  const existingRecord = existingRecords?.results?.[0];

  const dnsRecord = new cloudflare.DnsRecord(
    "backend-dns-v2",
    {
      zoneId: zoneId,
      name: "property",
      type: "CNAME",
      content: pulumi.interpolate`${tunnelId}.cfargotunnel.com`,
      ttl: 1,
      proxied: true,
    },
    {
      provider: provider,
      import: existingRecord?.id, // Adopt existing record if found
    },
  );

  // 3. Configure the Zero Trust Tunnel Ingress rules
  // Fetch existing config to preserve other routes
  let ingresses: any[] = [];
  try {
    const currentConfig = await cloudflare.getZeroTrustTunnelCloudflaredConfig(
      {
        accountId,
        tunnelId,
      },
      { provider: provider },
    );
    if (currentConfig.config?.ingresses) {
      ingresses = [...currentConfig.config.ingresses];
    }
  } catch (e) {
    // Ignore if no config exists or fetch fails
    console.warn("Could not fetch existing tunnel config, starting fresh.");
  }

  // Filter out existing rules for our service and the catch-all to avoid duplicates/ordering issues
  ingresses = ingresses.filter(
    (r) => r.hostname !== recordHostname && r.service !== "http_status:404",
  );

  // Add our rule
  // Using the internal ClusterIP service name
  ingresses.push({
    hostname: recordHostname,
    service: pulumi.interpolate`http://${backendService.spec.clusterIP}:${backendService.spec.ports[0].port}`,
  });

  // Append catch-all rule at the end
  ingresses.push({
    service: "http_status:404",
  });

  const tunnelConfig = new cloudflare.ZeroTrustTunnelCloudflaredConfig(
    "backend-tunnel-config",
    {
      accountId: accountId,
      tunnelId: tunnelId,
      config: {
        ingresses: ingresses,
      },
    },
    { provider: provider },
  );

  return { tunnel, dnsRecord, tunnelId, tunnelToken: tunnelToken.token };
}
