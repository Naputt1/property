import { createQueryWrapper } from "@/utils/query";
import { defaultQueryFn } from "@/utils/query/defaultQuery";

export interface MedianPriceResult {
  region: string;
  median_price: number;
}

export interface PriceTrendResult {
  period: string;
  avg_price: number;
  median_price: number;
  transaction_count: number;
}

export interface AffordabilityResult {
  property_type: string;
  avg_price: number;
  relative_affordability: number;
}

export interface GrowthHotspotResult {
  region: string;
  growth_rate: number;
  prev_median: number;
  current_median: number;
}

export const medianPriceQuery = createQueryWrapper<
  MedianPriceResult[],
  { by: string },
  { param: { by: string } }
>({
  queryKey: (params) => ["analytics", "median-price", params.by],
  options: {
    queryFn: defaultQueryFn({
      url: "/analytics/median-price?by=$by",
    })<MedianPriceResult[]>(),
  },
});

export const priceTrendQuery = createQueryWrapper<
  PriceTrendResult[],
  { interval: string },
  { param: { interval: string } }
>({
  queryKey: (params) => ["analytics", "price-trend", params.interval],
  options: {
    queryFn: defaultQueryFn({
      url: "/analytics/price-trend?interval=$interval",
    })<PriceTrendResult[]>(),
  },
});

export const affordabilityQuery = createQueryWrapper<AffordabilityResult[]>({
  queryKey: () => ["analytics", "affordability"],
  options: {
    queryFn: defaultQueryFn({
      url: "/analytics/affordability",
    })<AffordabilityResult[]>(),
  },
});

export const growthHotspotsQuery = createQueryWrapper<
  GrowthHotspotResult[],
  { limit: number },
  { param: { limit: number } }
>({
  queryKey: (params) => ["analytics", "growth-hotspots", params.limit],
  options: {
    queryFn: defaultQueryFn({
      url: "/analytics/growth-hotspots?limit=$limit",
    })<GrowthHotspotResult[]>(),
  },
});
