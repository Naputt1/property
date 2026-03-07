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

export interface NewBuildPremiumResult {
  region: string;
  new_avg: number;
  old_avg: number;
  premium_percent: number;
}

export interface PropertyTypeDistributionResult {
  property_type: string;
  count: number;
  percentage: number;
}

export interface PriceBracketResult {
  bracket: string;
  count: number;
  percentage: number;
}

export interface TopActiveAreaResult {
  region: string;
  transaction_count: number;
  total_value: number;
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

export const newBuildPremiumQuery = createQueryWrapper<
  NewBuildPremiumResult[],
  { by: string },
  { param: { by: string } }
>({
  queryKey: (params) => ["analytics", "new-build-premium", params.by],
  options: {
    queryFn: defaultQueryFn({
      url: "/analytics/new-build-premium?by=$by",
    })<NewBuildPremiumResult[]>(),
  },
});

export const propertyTypeDistributionQuery = createQueryWrapper<
  PropertyTypeDistributionResult[]
>({
  queryKey: () => ["analytics", "property-type-distribution"],
  options: {
    queryFn: defaultQueryFn({
      url: "/analytics/property-type-distribution",
    })<PropertyTypeDistributionResult[]>(),
  },
});

export const priceBracketDistributionQuery = createQueryWrapper<
  PriceBracketResult[]
>({
  queryKey: () => ["analytics", "price-bracket-distribution"],
  options: {
    queryFn: defaultQueryFn({
      url: "/analytics/price-bracket-distribution",
    })<PriceBracketResult[]>(),
  },
});

export const topActiveAreasQuery = createQueryWrapper<
  TopActiveAreaResult[],
  { by: string; limit: number },
  { param: { by: string; limit: number } }
>({
  queryKey: (params) => ["analytics", "top-active-areas", params.by, params.limit],
  options: {
    queryFn: defaultQueryFn({
      url: "/analytics/top-active-areas?by=$by&limit=$limit",
    })<TopActiveAreaResult[]>(),
  },
});
