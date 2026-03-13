import { createFileRoute, Link } from "@tanstack/react-router";
import { useState, useMemo } from "react";
import { Map as MapIcon } from "lucide-react";
import { Button } from "@/components/ui/button";
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
  LineChart,
  Line,
  PieChart,
  Pie,
  Cell,
} from "recharts";
import {
  useGetAnalyticsMedianPrice,
  useGetAnalyticsPriceTrend,
  useGetAnalyticsAffordability,
  useGetAnalyticsGrowthHotspots,
  useGetAnalyticsNewBuildPremium,
  useGetAnalyticsPropertyTypeDistribution,
  useGetAnalyticsPriceBracketDistribution,
  useGetAnalyticsTopActiveAreas,
} from "@/gen/hooks";

export const Route = createFileRoute("/analytics/")({
  component: Analytics,
});

const COLORS = [
  "var(--color-chart-1)",
  "var(--color-chart-2)",
  "var(--color-chart-3)",
  "var(--color-chart-4)",
  "var(--color-chart-5)",
];

function Analytics() {
  const [regionType, setRegionType] = useState<any>("county");
  const [trendInterval, setTrendInterval] = useState<any>("month");
  const [activityInterval, setActivityInterval] = useState<any>("month");
  const [premiumRegion, setPremiumRegion] = useState<any>("county");
  const [activeAreaRegion, setActiveAreaRegion] = useState<any>("district");

  const queryOptions = useMemo(
    () => ({
      query: {
        staleTime: 10 * 60 * 1000, // 10 minutes
        gcTime: 15 * 60 * 1000, // 15 minutes
        retry: false,
        refetchOnWindowFocus: false,
      },
    }),
    [],
  );

  // Stabilize all parameter objects
  const medianParams = useMemo(() => ({ by: regionType }), [regionType]);
  const trendParams = useMemo(
    () => ({ interval: trendInterval }),
    [trendInterval],
  );
  const activityTrendParams = useMemo(
    () => ({ interval: activityInterval }),
    [activityInterval],
  );
  const growthParams = useMemo(() => ({ limit: 10 }), []);
  const premiumParams = useMemo(() => ({ by: premiumRegion }), [premiumRegion]);
  const activeParams = useMemo(
    () => ({ by: activeAreaRegion, limit: 10 }),
    [activeAreaRegion],
  );

  const { data: medianPrices, isLoading: loadingMedian } =
    useGetAnalyticsMedianPrice(medianParams, queryOptions);
  const { data: priceTrends, isLoading: loadingTrends } =
    useGetAnalyticsPriceTrend(trendParams, queryOptions);
  const { data: activityTrends, isLoading: loadingActivityTrends } =
    useGetAnalyticsPriceTrend(activityTrendParams, queryOptions);
  const { data: affordability, isLoading: loadingAffordability } =
    useGetAnalyticsAffordability(queryOptions);
  const { data: hotspots, isLoading: loadingHotspots } =
    useGetAnalyticsGrowthHotspots(growthParams, queryOptions);
  const { data: newBuildPremium, isLoading: loadingPremium } =
    useGetAnalyticsNewBuildPremium(premiumParams, queryOptions);
  const { data: typeDistribution, isLoading: loadingTypeDist } =
    useGetAnalyticsPropertyTypeDistribution(queryOptions);
  const { data: bracketDistribution, isLoading: loadingBracketDist } =
    useGetAnalyticsPriceBracketDistribution(queryOptions);
  const { data: activeAreas, isLoading: loadingActiveAreas } =
    useGetAnalyticsTopActiveAreas(activeParams, queryOptions);

  const formatPrice = (value: number | undefined) => {
    if (value === undefined) return "£0";
    return new Intl.NumberFormat("en-GB", {
      style: "currency",
      currency: "GBP",
      maximumFractionDigits: 0,
    }).format(value);
  };

  const getPropertyTypeName = (type: string | undefined) => {
    if (!type) return "Unknown";
    const types: Record<string, string> = {
      D: "Detached",
      S: "Semi-Detached",
      T: "Terraced",
      F: "Flats",
      O: "Other",
    };
    return types[type] || type;
  };

  return (
    <div className="space-y-12 pb-12">
      <header className="flex justify-between items-center">
        <div>
          <h1 className="text-3xl font-bold">UK Housing Market Analytics</h1>
          <p className="text-muted-foreground">
            Comprehensive data insights and market trends.
          </p>
        </div>
        <Link to="/analytics/map">
          <Button className="gap-2">
            <MapIcon className="h-4 w-4" />
            View on Map
          </Button>
        </Link>
      </header>

      {/* Main Trends Grid */}
      <div className="grid grid-cols-1 xl:grid-cols-2 gap-8">
        {/* Price Trend Section */}
        <section className="bg-white p-6 rounded-xl border shadow-sm">
          <div className="flex justify-between items-center mb-6">
            <h2 className="text-xl font-semibold">Market Price Trends</h2>
            <select
              className="border rounded px-3 py-1 text-sm bg-muted/50"
              value={trendInterval}
              onChange={(e) => setTrendInterval(e.target.value)}
            >
              <option value="month">Monthly</option>
              <option value="year">Yearly</option>
            </select>
          </div>
          <div className="h-87.5 w-full">
            {!loadingTrends ? (
              <ResponsiveContainer width="100%" height="100%">
                <LineChart data={priceTrends as any}>
                  <CartesianGrid strokeDasharray="3 3" vertical={false} />
                  <XAxis dataKey="period" fontSize={12} tickMargin={10} />
                  <YAxis
                    tickFormatter={(value) => `£${value / 1000}k`}
                    fontSize={12}
                  />
                  <Tooltip
                    formatter={(value: any) => formatPrice(Number(value))}
                  />
                  <Legend />
                  <Line
                    type="monotone"
                    dataKey="avg_price"
                    stroke="var(--color-chart-1)"
                    strokeWidth={2}
                    dot={false}
                    name="Average Price"
                  />
                  <Line
                    type="monotone"
                    dataKey="median_price"
                    stroke="var(--color-chart-2)"
                    strokeWidth={2}
                    dot={false}
                    name="Median Price"
                  />
                </LineChart>
              </ResponsiveContainer>
            ) : (
              <div className="flex items-center justify-center h-full text-muted-foreground">
                Loading...
              </div>
            )}
          </div>
        </section>

        {/* Activity Trend Section */}
        <section className="bg-white p-6 rounded-xl border shadow-sm">
          <div className="flex justify-between items-center mb-6">
            <h2 className="text-xl font-semibold">Market Activity Trends</h2>
            <select
              className="border rounded px-3 py-1 text-sm bg-muted/50"
              value={activityInterval}
              onChange={(e) => setActivityInterval(e.target.value)}
            >
              <option value="month">Monthly</option>
              <option value="year">Yearly</option>
            </select>
          </div>
          <div className="h-87.5 w-full">
            {!loadingActivityTrends ? (
              <ResponsiveContainer width="100%" height="100%">
                <LineChart data={activityTrends as any}>
                  <CartesianGrid strokeDasharray="3 3" vertical={false} />
                  <XAxis dataKey="period" fontSize={12} tickMargin={10} />
                  <YAxis fontSize={12} />
                  <Tooltip
                    formatter={(value: any) => [
                      value.toLocaleString(),
                      "Transactions",
                    ]}
                  />
                  <Legend />
                  <Line
                    type="monotone"
                    dataKey="transaction_count"
                    stroke="var(--color-chart-3)"
                    strokeWidth={2}
                    dot={false}
                    name="Transaction Count"
                  />
                </LineChart>
              </ResponsiveContainer>
            ) : (
              <div className="flex items-center justify-center h-full text-muted-foreground">
                Loading...
              </div>
            )}
          </div>
        </section>

        {/* Median Price Section */}
        <section className="bg-white p-6 rounded-xl border shadow-sm">
          <div className="flex justify-between items-center mb-6">
            <h2 className="text-xl font-semibold">Median Price by Region</h2>
            <select
              className="border rounded px-3 py-1 text-sm bg-muted/50"
              value={regionType}
              onChange={(e) => setRegionType(e.target.value)}
            >
              <option value="county">County</option>
              <option value="district">District</option>
              <option value="town_city">Town/City</option>
            </select>
          </div>
          <div className="h-87.5 w-full">
            {!loadingMedian ? (
              <ResponsiveContainer width="100%" height="100%">
                <BarChart data={medianPrices as any}>
                  <CartesianGrid strokeDasharray="3 3" vertical={false} />
                  <XAxis dataKey="region" fontSize={12} tickMargin={10} />
                  <YAxis
                    tickFormatter={(value) => `£${value / 1000}k`}
                    fontSize={12}
                  />
                  <Tooltip
                    formatter={(value: any) => formatPrice(Number(value))}
                    contentStyle={{
                      borderRadius: "8px",
                      border: "none",
                      boxShadow: "0 4px 6px -1px rgb(0 0 0 / 0.1)",
                    }}
                  />
                  <Bar
                    dataKey="median_price"
                    fill="var(--color-chart-1)"
                    radius={[4, 4, 0, 0]}
                    name="Median Price"
                  />
                </BarChart>
              </ResponsiveContainer>
            ) : (
              <div className="flex items-center justify-center h-full text-muted-foreground">
                Loading...
              </div>
            )}
          </div>
        </section>
      </div>

      {/* Complex Analytics Grid */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
        {/* New Build Premium */}
        <section className="bg-white p-6 rounded-xl border shadow-sm">
          <div className="flex justify-between items-center mb-6">
            <h2 className="text-xl font-semibold">New Build Premium</h2>
            <select
              className="border rounded px-3 py-1 text-sm bg-muted/50"
              value={premiumRegion}
              onChange={(e) => setPremiumRegion(e.target.value)}
            >
              <option value="county">County</option>
              <option value="district">District</option>
            </select>
          </div>
          <div className="h-87.5 w-full">
            {!loadingPremium ? (
              <ResponsiveContainer width="100%" height="100%">
                <BarChart data={newBuildPremium as any}>
                  <CartesianGrid strokeDasharray="3 3" vertical={false} />
                  <XAxis dataKey="region" fontSize={10} />
                  <YAxis
                    tickFormatter={(value) => `£${value / 1000}k`}
                    fontSize={12}
                  />
                  <Tooltip
                    formatter={(value: any) => formatPrice(Number(value))}
                    contentStyle={{
                      borderRadius: "8px",
                      border: "none",
                      boxShadow: "0 4px 6px -1px rgb(0 0 0 / 0.1)",
                    }}
                  />
                  <Bar
                    dataKey="premium_percentage"
                    fill="var(--color-chart-4)"
                    radius={[4, 4, 0, 0]}
                    name="Premium %"
                  />
                </BarChart>
              </ResponsiveContainer>
            ) : (
              <div className="flex items-center justify-center h-full text-muted-foreground">
                Loading...
              </div>
            )}
          </div>
        </section>

        {/* Top Active Areas */}
        <section className="bg-white p-6 rounded-xl border shadow-sm">
          <div className="flex justify-between items-center mb-6">
            <h2 className="text-xl font-semibold">Most Active Areas</h2>
            <select
              className="border rounded px-3 py-1 text-sm bg-muted/50"
              value={activeAreaRegion}
              onChange={(e) => setActiveAreaRegion(e.target.value)}
            >
              <option value="district">District</option>
              <option value="town_city">Town/City</option>
              <option value="county">County</option>
            </select>
          </div>
          <div className="overflow-x-auto">
            <table className="w-full text-sm text-left">
              <thead className="bg-muted/50 text-muted-foreground uppercase text-xs font-medium">
                <tr>
                  <th className="px-4 py-3">Region</th>
                  <th className="px-4 py-3 text-right">Market Activity</th>
                  <th className="px-4 py-3 text-right">Total Value</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-border">
                {!loadingActiveAreas ? (
                  (activeAreas as any)?.map((area: any, i: number) => (
                    <tr key={i} className="hover:bg-muted/50 transition-colors">
                      <td className="px-4 py-3 font-medium">{area.region}</td>
                      <td className="px-4 py-3 text-right font-mono">
                        {area.transaction_count?.toLocaleString()}
                      </td>
                      <td className="px-4 py-3 text-right font-mono">
                        {formatPrice(area.total_value)}
                      </td>
                    </tr>
                  ))
                ) : (
                  <tr>
                    <td
                      colSpan={3}
                      className="px-4 py-8 text-center text-gray-400"
                    >
                      Loading...
                    </td>
                  </tr>
                )}
              </tbody>
            </table>
          </div>
        </section>
      </div>

      {/* Distribution Grid */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
        {/* Property Type Distribution */}
        <section className="bg-white p-6 rounded-xl border shadow-sm">
          <h2 className="text-lg font-semibold mb-6">Property Type Stock</h2>
          <div className="h-75 w-full">
            {!loadingTypeDist ? (
              <ResponsiveContainer width="100%" height="100%">
                <PieChart>
                  <Pie
                    data={typeDistribution as any}
                    cx="50%"
                    cy="50%"
                    innerRadius={60}
                    outerRadius={80}
                    paddingAngle={5}
                    dataKey="count"
                    nameKey="property_type"
                    label={({ name }) => getPropertyTypeName(name)}
                  >
                    {(typeDistribution as any)?.map((_: any, index: number) => (
                      <Cell
                        key={`cell-${index}`}
                        fill={COLORS[index % COLORS.length]}
                      />
                    ))}
                  </Pie>
                  <Tooltip
                    formatter={(value: any, name: any) => [
                      value.toLocaleString(),
                      getPropertyTypeName(name as string),
                    ]}
                  />
                </PieChart>
              </ResponsiveContainer>
            ) : (
              <div className="flex items-center justify-center h-full text-muted-foreground">
                Loading...
              </div>
            )}
          </div>
        </section>

        {/* Price Bracket Distribution */}
        <section className="bg-white p-6 rounded-xl border shadow-sm md:col-span-2">
          <h2 className="text-lg font-semibold mb-6">Market Price Segments</h2>
          <div className="h-75 w-full">
            {!loadingBracketDist ? (
              <ResponsiveContainer width="100%" height="100%">
                <BarChart data={bracketDistribution as any}>
                  <CartesianGrid strokeDasharray="3 3" vertical={false} />
                  <XAxis dataKey="bracket" fontSize={12} />
                  <YAxis fontSize={12} />
                  <Tooltip
                    formatter={(value: any) => [
                      `${Number(value).toFixed(1)}%`,
                      "Market Share",
                    ]}
                  />
                  <Bar
                    dataKey="percentage"
                    fill="var(--color-primary)"
                    radius={[4, 4, 0, 0]}
                    name="Market Share"
                  />
                </BarChart>
              </ResponsiveContainer>
            ) : (
              <div className="flex items-center justify-center h-full text-muted-foreground">
                Loading...
              </div>
            )}
          </div>
        </section>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
        {/* Affordability Index */}
        <section className="bg-white p-6 rounded-xl border shadow-sm">
          <h2 className="text-xl font-semibold mb-6">Affordability by Type</h2>
          <div className="h-87.5 w-full">
            {!loadingAffordability ? (
              <ResponsiveContainer width="100%" height="100%">
                <BarChart data={affordability as any} layout="vertical">
                  <CartesianGrid strokeDasharray="3 3" horizontal={false} />
                  <XAxis type="number" hide />
                  <YAxis
                    dataKey="property_type"
                    type="category"
                    width={100}
                    fontSize={12}
                    tickFormatter={(value) => getPropertyTypeName(value)}
                  />
                  <Tooltip
                    formatter={(value: any, name: any) => {
                      if (name === "avg_price")
                        return formatPrice(Number(value));
                      return Number(value).toFixed(2);
                    }}
                  />
                  <Bar
                    dataKey="avg_price"
                    fill="var(--color-chart-1)"
                    name="Average Price"
                    radius={[0, 4, 4, 0]}
                  />
                </BarChart>
              </ResponsiveContainer>
            ) : (
              <div className="flex items-center justify-center h-full text-muted-foreground">
                Loading...
              </div>
            )}
          </div>
        </section>

        {/* Growth Hotspots */}
        <section className="bg-white p-6 rounded-xl border shadow-sm">
          <h2 className="text-xl font-semibold mb-6">Top Growth Hotspots</h2>
          <div className="overflow-x-auto">
            <table className="w-full text-sm text-left">
              <thead className="bg-muted/50 text-muted-foreground uppercase text-xs font-medium">
                <tr>
                  <th className="px-4 py-3">Region</th>
                  <th className="px-4 py-3">Growth</th>
                  <th className="px-4 py-3 text-right">Current Median</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-border">
                {!loadingHotspots ? (
                  (hotspots as any)?.map((h: any, i: number) => (
                    <tr key={i} className="hover:bg-muted/50 transition-colors">
                      <td className="px-4 py-3 font-medium">{h.region}</td>
                      <td className="px-4 py-3">
                        <span className="text-chart-2 font-semibold">
                          +{h.growth_rate?.toFixed(1)}%
                        </span>
                      </td>
                      <td className="px-4 py-3 text-right font-mono">
                        {formatPrice(h.current_median)}
                      </td>
                    </tr>
                  ))
                ) : (
                  <tr>
                    <td
                      colSpan={3}
                      className="px-4 py-8 text-center text-muted-foreground"
                    >
                      Loading hotspots...
                    </td>
                  </tr>
                )}
              </tbody>
            </table>
          </div>
        </section>
      </div>
    </div>
  );
}
