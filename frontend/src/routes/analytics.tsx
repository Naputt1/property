import { createFileRoute } from "@tanstack/react-router";
import { useState } from "react";
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
} from "recharts";
import {
  medianPriceQuery,
  priceTrendQuery,
  affordabilityQuery,
  growthHotspotsQuery,
} from "@/query/analytics";
import { useQuery } from "@tanstack/react-query";

export const Route = createFileRoute("/analytics")({
  component: Analytics,
});

function Analytics() {
  const [regionType, setRegionType] = useState("county");
  const [trendInterval, setTrendInterval] = useState("month");

  const { data: medianPrices, isLoading: loadingMedian } = useQuery(
    medianPriceQuery.getOptions({ param: { by: regionType } }),
  );
  const { data: priceTrends, isLoading: loadingTrends } = useQuery(
    priceTrendQuery.getOptions({ param: { interval: trendInterval } }),
  );
  const { data: affordability, isLoading: loadingAffordability } = useQuery(
    affordabilityQuery.getOptions({}),
  );
  const { data: hotspots, isLoading: loadingHotspots } = useQuery(
    growthHotspotsQuery.getOptions({ param: { limit: 10 } }),
  );

  // const medianPrices = (medianPricesRes as any)?.data || [];
  // const priceTrends = (priceTrendsRes as any)?.data || [];
  // const affordability = (affordabilityRes as any)?.data || [];
  // const hotspots = (hotspotsRes as any)?.data || [];

  const formatPrice = (value: number | undefined) => {
    if (value === undefined) return "£0";
    return new Intl.NumberFormat("en-GB", {
      style: "currency",
      currency: "GBP",
      maximumFractionDigits: 0,
    }).format(value);
  };

  return (
    <div className="space-y-12 pb-12">
      <header>
        <h1 className="text-3xl font-bold">UK Housing Market Analytics</h1>
        <p className="text-gray-500">
          Comprehensive data insights and market trends.
        </p>
      </header>

      {/* Median Price Section */}
      <section className="bg-white p-6 rounded-xl border shadow-sm">
        <div className="flex justify-between items-center mb-6">
          <h2 className="text-xl font-semibold">Median Price by Region</h2>
          <select
            className="border rounded px-3 py-1 text-sm bg-gray-50"
            value={regionType}
            onChange={(e) => setRegionType(e.target.value)}
          >
            <option value="county">County</option>
            <option value="district">District</option>
            <option value="town_city">Town/City</option>
          </select>
        </div>
        <div className="h-[400px] w-full">
          {!loadingMedian ? (
            <ResponsiveContainer width="100%" height="100%">
              <BarChart data={medianPrices?.slice(0, 15)}>
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
                  fill="#2563eb"
                  radius={[4, 4, 0, 0]}
                  name="Median Price"
                />
              </BarChart>
            </ResponsiveContainer>
          ) : (
            <div className="flex items-center justify-center h-full text-gray-400">
              Loading...
            </div>
          )}
        </div>
      </section>

      {/* Price Trend Section */}
      <section className="bg-white p-6 rounded-xl border shadow-sm">
        <div className="flex justify-between items-center mb-6">
          <h2 className="text-xl font-semibold">Market Price Trends</h2>
          <select
            className="border rounded px-3 py-1 text-sm bg-gray-50"
            value={trendInterval}
            onChange={(e) => setTrendInterval(e.target.value)}
          >
            <option value="month">Monthly</option>
            <option value="year">Yearly</option>
          </select>
        </div>
        <div className="h-[400px] w-full">
          {!loadingTrends ? (
            <ResponsiveContainer width="100%" height="100%">
              <LineChart data={priceTrends}>
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
                  stroke="#2563eb"
                  strokeWidth={2}
                  dot={false}
                  name="Average Price"
                />
                <Line
                  type="monotone"
                  dataKey="median_price"
                  stroke="#10b981"
                  strokeWidth={2}
                  dot={false}
                  name="Median Price"
                />
              </LineChart>
            </ResponsiveContainer>
          ) : (
            <div className="flex items-center justify-center h-full text-gray-400">
              Loading...
            </div>
          )}
        </div>
      </section>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
        {/* Affordability Index */}
        <section className="bg-white p-6 rounded-xl border shadow-sm">
          <h2 className="text-xl font-semibold mb-6">Affordability by Type</h2>
          <div className="h-[350px] w-full">
            {!loadingAffordability ? (
              <ResponsiveContainer width="100%" height="100%">
                <BarChart data={affordability} layout="vertical">
                  <CartesianGrid strokeDasharray="3 3" horizontal={false} />
                  <XAxis type="number" hide />
                  <YAxis
                    dataKey="property_type"
                    type="category"
                    width={100}
                    fontSize={12}
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
                    fill="#3b82f6"
                    name="Average Price"
                    radius={[0, 4, 4, 0]}
                  />
                </BarChart>
              </ResponsiveContainer>
            ) : (
              <div className="flex items-center justify-center h-full text-gray-400">
                Loading...
              </div>
            )}
          </div>
          <p className="text-xs text-gray-400 mt-4">
            * Lower relative affordability index indicates more accessible
            pricing compared to market average.
          </p>
        </section>

        {/* Growth Hotspots */}
        <section className="bg-white p-6 rounded-xl border shadow-sm">
          <h2 className="text-xl font-semibold mb-6">Top Growth Hotspots</h2>
          <div className="overflow-x-auto">
            <table className="w-full text-sm text-left">
              <thead className="bg-gray-50 text-gray-600 uppercase text-xs font-medium">
                <tr>
                  <th className="px-4 py-3">District</th>
                  <th className="px-4 py-3">Growth</th>
                  <th className="px-4 py-3 text-right">Current Median</th>
                </tr>
              </thead>
              <tbody className="divide-y">
                {!loadingHotspots ? (
                  hotspots?.map((h: any, i: number) => (
                    <tr key={i} className="hover:bg-gray-50">
                      <td className="px-4 py-3 font-medium">{h.region}</td>
                      <td className="px-4 py-3">
                        <span className="text-green-600 font-semibold">
                          +{h.growth_rate.toFixed(1)}%
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
                      className="px-4 py-8 text-center text-gray-400"
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
