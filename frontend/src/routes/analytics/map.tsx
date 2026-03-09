import { createFileRoute } from "@tanstack/react-router";
import { useState, useMemo, useCallback, useEffect } from "react";
import UKMap from "@/components/UKMap";
import {
  useGetAnalyticsMedianPrice,
  useGetAnalyticsGrowthHotspots,
  useGetAnalyticsTopActiveAreas,
  useGetAnalyticsTimeRange,
} from "@/gen/hooks";
import { Button } from "@/components/ui/button";
import { Slider } from "@/components/ui/slider";
import {
  Map as MapIcon,
  BarChart3,
  TrendingUp,
  Activity,
  Calendar,
  Clock,
} from "lucide-react";
import { Link } from "@tanstack/react-router";

export const Route = createFileRoute("/analytics/map")({
  component: MapAnalytics,
});

type MetricType = "median_price" | "growth_rate" | "transaction_count";

function MapAnalytics() {
  const [metric, setMetric] = useState<MetricType>("median_price");
  const [regionType, setRegionType] = useState<"county" | "district">("county");
  const [selectedYear, setSelectedYear] = useState<number | null>(null);

  const queryOptions = useMemo(
    () => ({
      query: {
        retry: false,
        refetchOnWindowFocus: false,
        staleTime: 10 * 60 * 1000, // 10 minutes
        gcTime: 15 * 60 * 1000, // 15 minutes
      },
    }),
    [],
  );

  const { data: timeRange } = useGetAnalyticsTimeRange(queryOptions);

  // Initialize selected year from time range
  useEffect(() => {
    if (timeRange && selectedYear === null && timeRange.max_year) {
      setSelectedYear(timeRange.max_year);
    }
  }, [timeRange, selectedYear]);

  // Stabilize parameters
  const medianParams = useMemo(
    () => ({
      by: regionType,
      year: selectedYear || undefined,
    }),
    [regionType, selectedYear],
  );
  const activeParams = useMemo(
    () => ({
      by: regionType,
      limit: 0,
      year: selectedYear || undefined,
    }),
    [regionType, selectedYear],
  );
  const growthParams = useMemo(
    () => ({
      by: regionType,
      limit: 0,
      year: selectedYear || undefined,
    }),
    [regionType, selectedYear],
  );

  const { data: medianPrices, isLoading: loadingMedian } =
    useGetAnalyticsMedianPrice(medianParams, queryOptions);
  const { data: hotspots, isLoading: loadingHotspots } =
    useGetAnalyticsGrowthHotspots(growthParams, queryOptions);
  const { data: activeAreas, isLoading: loadingActiveAreas } =
    useGetAnalyticsTopActiveAreas(activeParams, queryOptions);

  const mapData = useMemo(() => {
    let rawData: Array<{ region?: string; value?: number | bigint }> = [];
    if (metric === "median_price") {
      rawData = (medianPrices || []).map((d: any) => ({
        region: d.region,
        value: d.median_price,
      }));
    } else if (metric === "growth_rate") {
      rawData = (hotspots || []).map((d: any) => ({
        region: d.region,
        value: d.growth_rate,
      }));
    } else {
      rawData = (activeAreas || []).map((d: any) => ({
        region: d.region,
        value: Number(d.transaction_count),
      }));
    }

    return rawData
      .filter(
        (d): d is { region: string; value: number } =>
          d && d.region !== undefined && d.value !== undefined,
      )
      .map((d) => ({ region: d.region, value: Number(d.value) }));
  }, [metric, medianPrices, hotspots, activeAreas]);

  // useCallback for prop stability
  const formatValue = useCallback(
    (value: number) => {
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
    },
    [metric],
  );

  const metricLabel = useMemo(() => {
    switch (metric) {
      case "median_price":
        return "Median Price";
      case "growth_rate":
        return "Growth Rate";
      case "transaction_count":
        return "Market Activity";
      default:
        return "";
    }
  }, [metric]);

  const isLoading = loadingMedian || loadingHotspots || loadingActiveAreas;

  return (
    <div className="space-y-8 pb-16">
      <header className="flex flex-col md:flex-row md:items-end justify-between gap-6 border-b pb-8">
        <div>
          <div className="flex items-center gap-3 mb-2">
            <div className="p-2 bg-indigo-600 rounded-lg text-white">
              <MapIcon className="h-6 w-6" />
            </div>
            <h1 className="text-4xl font-black tracking-tight text-slate-900">
              UK Market Intelligence
            </h1>
          </div>
          <p className="text-slate-500 font-medium max-w-2xl leading-relaxed">
            Geospatial visualization of housing market dynamics across the
            United Kingdom. Identify hotspots, analyze affordability, and track
            regional growth trends.
          </p>
        </div>
        <div className="flex items-center gap-3">
          <Link to="/analytics">
            <Button
              variant="outline"
              size="lg"
              className="gap-2 border-slate-200 hover:bg-slate-50 shadow-sm"
            >
              <BarChart3 className="h-5 w-5 text-indigo-600" />
              Chart Insights
            </Button>
          </Link>
        </div>
      </header>

      <div className="grid grid-cols-1 lg:grid-cols-12 gap-8 items-start">
        <aside className="lg:col-span-3 space-y-8">
          <section className="bg-white p-6 rounded-2xl border border-slate-200 shadow-sm space-y-6">
            <div>
              <h2 className="font-bold text-xs uppercase tracking-[0.2em] text-slate-400 mb-4">
                Configuration
              </h2>

              <div className="space-y-4">
                <div className="space-y-2">
                  <label className="text-sm font-bold text-slate-700">
                    Market Metric
                  </label>
                  <div className="grid grid-cols-1 gap-2">
                    <Button
                      variant={
                        metric === "median_price" ? "default" : "outline"
                      }
                      className={`justify-start gap-3 h-12 px-4 transition-all ${metric === "median_price" ? "bg-indigo-600 shadow-lg shadow-indigo-100" : "hover:bg-slate-50 border-slate-200"}`}
                      onClick={() => setMetric("median_price")}
                    >
                      <TrendingUp className="h-4 w-4" />
                      Median Price
                    </Button>
                    <Button
                      variant={metric === "growth_rate" ? "default" : "outline"}
                      className={`justify-start gap-3 h-12 px-4 transition-all ${metric === "growth_rate" ? "bg-indigo-600 shadow-lg shadow-indigo-100" : "hover:bg-slate-50 border-slate-200"}`}
                      onClick={() => setMetric("growth_rate")}
                    >
                      <Activity className="h-4 w-4" />
                      Growth Hotspots
                    </Button>
                    <Button
                      variant={
                        metric === "transaction_count" ? "default" : "outline"
                      }
                      className={`justify-start gap-3 h-12 px-4 transition-all ${metric === "transaction_count" ? "bg-indigo-600 shadow-lg shadow-indigo-100" : "hover:bg-slate-50 border-slate-200"}`}
                      onClick={() => setMetric("transaction_count")}
                    >
                      <BarChart3 className="h-4 w-4" />
                      Market Activity
                    </Button>
                  </div>
                </div>

                <div className="space-y-2 pt-2">
                  <label className="text-sm font-bold text-slate-700">
                    Granularity
                  </label>
                  <div className="relative group">
                    <select
                      className="w-full border border-slate-200 rounded-xl px-4 py-3 text-sm bg-slate-50/50 appearance-none focus:outline-none focus:ring-2 focus:ring-indigo-500/20 focus:border-indigo-500 transition-all cursor-pointer font-medium"
                      value={regionType}
                      onChange={(e) => setRegionType(e.target.value as any)}
                    >
                      <option value="county">County Level</option>
                      <option value="district">District Level</option>
                    </select>
                    <div className="absolute right-4 top-1/2 -translate-y-1/2 pointer-events-none text-slate-400 group-hover:text-indigo-500 transition-colors">
                      <svg
                        xmlns="http://www.w3.org/2000/svg"
                        width="16"
                        height="16"
                        viewBox="0 0 24 24"
                        fill="none"
                        stroke="currentColor"
                        strokeWidth="2"
                        strokeLinecap="round"
                        strokeLinejoin="round"
                      >
                        <path d="m6 9 6 6 6-6" />
                      </svg>
                    </div>
                  </div>
                </div>

                {timeRange && (
                  <div className="space-y-6 pt-4 border-t border-slate-100">
                    {metric === "growth_rate" && (
                      <p className="text-[10px] text-slate-500 font-medium bg-slate-50 p-2 rounded-lg border border-slate-100 italic">
                        Note: Growth rate compares the selected period to the
                        same period in the previous year.
                      </p>
                    )}
                    <div className="space-y-4">
                      <div className="flex items-center justify-between">
                        <label className="text-sm font-bold text-slate-700 flex items-center gap-2">
                          <Calendar className="h-4 w-4 text-indigo-500" />
                          Year: {selectedYear}
                        </label>
                        <button
                          onClick={() => setSelectedYear(null)}
                          className={`text-[10px] font-bold uppercase tracking-wider ${selectedYear === null ? "text-indigo-600" : "text-slate-400 hover:text-indigo-600"}`}
                        >
                          All Time
                        </button>
                      </div>
                      <Slider
                        value={[
                          selectedYear ||
                            (timeRange?.max_year ?? new Date().getFullYear()),
                        ]}
                        min={timeRange?.min_year ?? 1995}
                        max={timeRange?.max_year ?? new Date().getFullYear()}
                        step={1}
                        onValueChange={([val]) => setSelectedYear(val)}
                        className={selectedYear === null ? "opacity-40" : ""}
                      />
                      <div className="flex justify-between text-[10px] font-bold text-slate-400">
                        <span>{timeRange?.min_year ?? 1995}</span>
                        <span>
                          {timeRange?.max_year ?? new Date().getFullYear()}
                        </span>
                      </div>
                    </div>
                  </div>
                )}
              </div>
            </div>
          </section>

          <section className="bg-indigo-50/50 p-6 rounded-2xl border border-indigo-100 space-y-4">
            <h3 className="text-indigo-900 font-bold text-sm flex items-center gap-2">
              <Activity className="h-4 w-4" />
              Quick Intelligence
            </h3>
            <p className="text-indigo-800/80 text-xs leading-relaxed font-medium">
              Interpreting the {metricLabel}: Darker indigo shades represent
              high-value clusters. The current dataset covers{" "}
              <span className="font-bold">{mapData.length}</span> active
              regions.
            </p>
            <div className="pt-2">
              <div className="p-3 bg-white/60 backdrop-blur rounded-xl border border-indigo-100/50 shadow-sm">
                <div className="text-[10px] uppercase tracking-wider text-indigo-400 font-bold mb-1">
                  Interactive Help
                </div>
                <ul className="text-[10px] text-indigo-700 space-y-1.5 font-medium">
                  <li className="flex gap-2">
                    <span className="text-indigo-400">●</span>
                    Hover to inspect regional data
                  </li>
                  <li className="flex gap-2">
                    <span className="text-indigo-400">●</span>
                    Scroll to zoom into specific clusters
                  </li>
                  <li className="flex gap-2">
                    <span className="text-indigo-400">●</span>
                    Click labels to toggle perspectives
                  </li>
                </ul>
              </div>
            </div>
          </section>
        </aside>

        <main className="lg:col-span-9 space-y-6">
          {isLoading ? (
            <div className="h-162.5 w-full bg-slate-50 animate-pulse rounded-2xl border border-slate-200 flex flex-col items-center justify-center gap-4">
              <div className="w-12 h-12 border-4 border-indigo-100 border-t-indigo-600 rounded-full animate-spin"></div>
              <p className="text-slate-400 font-bold text-sm uppercase tracking-widest">
                Aggregating Metrics...
              </p>
            </div>
          ) : (
            <div className="relative group">
              <UKMap
                data={mapData}
                regionType={regionType}
                metricLabel={metricLabel}
                formatValue={formatValue}
              />
              <div className="absolute top-6 right-6 z-1000 flex gap-2">
                <div className="bg-white/80 backdrop-blur-md px-3 py-1.5 rounded-lg border border-slate-200 shadow-sm text-[10px] font-bold text-slate-500 uppercase tracking-tight flex items-center gap-2">
                  <Clock className="h-3 w-3 text-indigo-500" />
                  Period: {selectedYear ? selectedYear : "All Time"}
                </div>
              </div>
            </div>
          )}

          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            <div className="bg-white p-6 rounded-2xl border border-slate-200 shadow-sm">
              <div className="text-[10px] font-bold text-slate-400 uppercase tracking-widest mb-1">
                National Pulse
              </div>
              <div className="text-2xl font-black text-slate-900">
                {formatValue(
                  mapData.reduce((acc, d) => acc + d.value, 0) /
                    (mapData.length || 1),
                )}
              </div>
              <div className="text-xs text-slate-500 font-medium mt-1">
                Average {metricLabel} across UK
              </div>
            </div>
            <div className="bg-white p-6 rounded-2xl border border-slate-200 shadow-sm">
              <div className="text-[10px] font-bold text-slate-400 uppercase tracking-widest mb-1">
                Top Performer
              </div>
              <div className="text-2xl font-black text-indigo-600">
                {mapData.length > 0
                  ? [...mapData].sort((a, b) => b.value - a.value)[0].region
                  : "N/A"}
              </div>
              <div className="text-xs text-slate-500 font-medium mt-1">
                Highest regional {metricLabel.toLowerCase()}
              </div>
            </div>
            <div className="bg-white p-6 rounded-2xl border border-slate-200 shadow-sm">
              <div className="text-[10px] font-bold text-slate-400 uppercase tracking-widest mb-1">
                Sample Coverage
              </div>
              <div className="text-2xl font-black text-slate-900">
                {mapData.length}
              </div>
              <div className="text-xs text-slate-500 font-medium mt-1">
                Data points mapped for current view
              </div>
            </div>
          </div>
        </main>
      </div>
    </div>
  );
}
