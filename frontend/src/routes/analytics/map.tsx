import { createFileRoute } from "@tanstack/react-router";
import { useState, useMemo, useCallback } from "react";
import UKMap from "@/components/UKMap";
import { useGetAnalyticsMedianPrice } from "@/gen/hooks/useGetAnalyticsMedianPrice";
import { useGetAnalyticsGrowthHotspots } from "@/gen/hooks/useGetAnalyticsGrowthHotspots";
import { useGetAnalyticsTopActiveAreas } from "@/gen/hooks/useGetAnalyticsTopActiveAreas";
import { Button } from "@/components/ui/button";
import { Map as MapIcon, BarChart3, TrendingUp, Activity } from "lucide-react";
import { Link } from "@tanstack/react-router";

export const Route = createFileRoute("/analytics/map")({
  component: MapAnalytics,
});

type MetricType = "median_price" | "growth_rate" | "transaction_count";

function MapAnalytics() {
  const [metric, setMetric] = useState<MetricType>("median_price");
  const [regionType, setRegionType] = useState<"county" | "district">("county");

  const queryOptions = useMemo(() => ({
    query: {
      retry: false,
      refetchOnWindowFocus: false,
      staleTime: 10 * 60 * 1000,
      gcTime: 10 * 60 * 1000,
    }
  }), []);

  // Stabilize parameters
  const medianParams = useMemo(() => ({ by: regionType }), [regionType]);
  const activeParams = useMemo(() => ({ by: regionType, limit: 100 }), [regionType]);
  const growthParams = useMemo(() => ({ limit: 100 }), []);

  const { data: medianPrices, isLoading: loadingMedian } = useGetAnalyticsMedianPrice(medianParams, queryOptions);
  const { data: hotspots, isLoading: loadingHotspots } = useGetAnalyticsGrowthHotspots(growthParams, queryOptions);
  const { data: activeAreas, isLoading: loadingActiveAreas } = useGetAnalyticsTopActiveAreas(activeParams, queryOptions);

  // Debug logging
  useMemo(() => {
    if (medianPrices || hotspots || activeAreas) {
      console.log("Analytics data received:", {
        medianPrices: medianPrices?.length || 0,
        hotspots: hotspots?.length || 0,
        activeAreas: activeAreas?.length || 0,
      });
      if (medianPrices && medianPrices.length > 0) {
        console.log("Sample region from API:", medianPrices[0].region);
      }
    }
  }, [medianPrices, hotspots, activeAreas]);

  const mapData = useMemo(() => {
    let rawData: Array<{ region?: string; value?: number | bigint }> = [];
    if (metric === "median_price") {
      rawData = (medianPrices || []).map((d) => ({
        region: d.region,
        value: d.median_price,
      }));
    } else if (metric === "growth_rate") {
      rawData = (hotspots || []).map((d) => ({
        region: d.region,
        value: d.growth_rate,
      }));
    } else {
      rawData = (activeAreas || []).map((d) => ({
        region: d.region,
        value: Number(d.transaction_count),
      }));
    }
    
    return rawData
      .filter((d): d is { region: string; value: number } => 
        d && d.region !== undefined && d.value !== undefined
      )
      .map(d => ({ region: d.region, value: Number(d.value) }));
  }, [metric, medianPrices, hotspots, activeAreas]);

  // useCallback for prop stability
  const formatValue = useCallback((value: number) => {
    if (metric === "median_price") {
      return new Intl.NumberFormat("en-GB", {
        style: "currency",
        currency: "GBP",
        maximumFractionDigits: 0,
      }).format(value);
    } else if (metric === "growth_rate") {
      return `${value.toFixed(1)}%`;
    } else {
      return value.toLocaleString();
    }
  }, [metric]);

  const metricLabel = useMemo(() => {
    switch (metric) {
      case "median_price":
        return "Median Price";
      case "growth_rate":
        return "Growth Rate";
      case "transaction_count":
        return "Transactions";
      default:
        return "";
    }
  }, [metric]);

  const isLoading = loadingMedian || loadingHotspots || loadingActiveAreas;

  return (
    <div className="space-y-6 pb-12">
      <header className="flex flex-col md:flex-row md:items-center justify-between gap-4">
        <div>
          <h1 className="text-3xl font-bold flex items-center gap-2">
            <MapIcon className="h-8 w-8 text-primary" />
            UK Market Map
          </h1>
          <p className="text-muted-foreground">
            Geographic distribution of housing market metrics.
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Link to="/analytics">
            <Button variant="outline" size="sm" className="gap-2">
              <BarChart3 className="h-4 w-4" />
              Chart View
            </Button>
          </Link>
        </div>
      </header>

      <div className="grid grid-cols-1 lg:grid-cols-4 gap-6">
        <aside className="lg:col-span-1 space-y-6">
          <section className="bg-white p-4 rounded-xl border shadow-sm space-y-4">
            <h2 className="font-semibold text-sm uppercase tracking-wider text-muted-foreground">
              Map Controls
            </h2>
            
            <div className="space-y-2">
              <label className="text-sm font-medium">Select Metric</label>
              <div className="grid grid-cols-1 gap-2">
                <Button 
                  variant={metric === "median_price" ? "default" : "outline"}
                  className="justify-start gap-2 h-10"
                  onClick={() => setMetric("median_price")}
                >
                  <TrendingUp className="h-4 w-4" />
                  Median Price
                </Button>
                <Button 
                  variant={metric === "growth_rate" ? "default" : "outline"}
                  className="justify-start gap-2 h-10"
                  onClick={() => setMetric("growth_rate")}
                >
                  <TrendingUp className="h-4 w-4" />
                  Growth Hotspots
                </Button>
                <Button 
                  variant={metric === "transaction_count" ? "default" : "outline"}
                  className="justify-start gap-2 h-10"
                  onClick={() => setMetric("transaction_count")}
                >
                  <Activity className="h-4 w-4" />
                  Market Activity
                </Button>
              </div>
            </div>

            <div className="space-y-2">
              <label className="text-sm font-medium">Region Level</label>
              <select
                className="w-full border rounded-md px-3 py-2 text-sm bg-muted/50 focus:outline-none focus:ring-2 focus:ring-primary"
                value={regionType}
                onChange={(e) => setRegionType(e.target.value as any)}
              >
                <option value="county">County</option>
              </select>
              <p className="text-[10px] text-muted-foreground">
                Note: Mapping works best with County level data currently.
              </p>
            </div>
          </section>

          <section className="bg-blue-50 p-4 rounded-xl border border-blue-100">
            <h3 className="text-blue-900 font-semibold text-sm mb-2">How to read</h3>
            <p className="text-blue-800 text-xs leading-relaxed">
              Hover over regions to see exact values. Darker areas indicate higher {metricLabel.toLowerCase()}. 
              Use the scroll wheel to zoom into specific areas.
            </p>
          </section>
        </aside>

        <main className="lg:col-span-3">
          {isLoading ? (
            <div className="h-[600px] w-full bg-muted/10 animate-pulse rounded-xl border flex items-center justify-center">
              <p className="text-muted-foreground">Loading analytics data...</p>
            </div>
          ) : (
            <UKMap 
              data={mapData} 
              regionType={regionType} 
              metricLabel={metricLabel}
              formatValue={formatValue}
            />
          )}
        </main>
      </div>
    </div>
  );
}
