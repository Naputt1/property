import React, { useEffect, useState, useMemo } from "react";
import { MapContainer, TileLayer, GeoJSON } from "react-leaflet";
import L from "leaflet";
import "leaflet/dist/leaflet.css";

// UK GeoJSON Sources
const COUNTY_GEOJSON_URL =
  "https://raw.githubusercontent.com/evansd/uk-ceremonial-counties/master/uk-ceremonial-counties.geojson";
const DISTRICT_GEOJSON_URL =
  "https://raw.githubusercontent.com/martinjc/UK-GeoJSON/master/json/administrative/gb/lad.json";

interface UKMapProps {
  data: Array<{ region: string; value: number }>;
  regionType: "county" | "district";
  metricLabel: string;
  formatValue: (value: number) => string;
}

const UKMap: React.FC<UKMapProps> = ({
  data,
  regionType,
  metricLabel,
  formatValue,
}) => {
  const [geoJson, setGeoJson] = useState<any>(null);
  const [loading, setLoading] = useState(true);

  // Fetch GeoJSON based on region type
  useEffect(() => {
    const controller = new AbortController();
    setLoading(true);

    const url =
      regionType === "county" ? COUNTY_GEOJSON_URL : DISTRICT_GEOJSON_URL;

    fetch(url, { signal: controller.signal })
      .then((res) => res.json())
      .then((data) => {
        setGeoJson(data);
        setLoading(false);
      })
      .catch((err) => {
        if (err.name === "AbortError") return;
        console.error("Failed to load GeoJSON", err);
        setLoading(false);
      });

    return () => controller.abort();
  }, [regionType]);

  // Robust normalization for UK regional names to fix mismatches
  // Handles: "Bristol, City of" vs "City of Bristol", "St Albans" vs "St. Albans", etc.
  const normalizeRegionName = (name: string): string => {
    if (!name) return "";
    return name
      .toLowerCase()
      .replace(/[.,]/g, "") // Remove commas and dots
      .replace(
        /\b(city of|city|of|and|the|upon|on|borough|royal|ceremonial)\b/g,
        "",
      ) // Remove common filler words
      .replace(/\s+/g, " ")
      .trim()
      .split(" ")
      .sort() // Sort words alphabetically to handle reordering
      .join(" ");
  };

  // Map data for O(1) lookup with normalized keys
  const dataLookup = useMemo(() => {
    const map = new Map<string, number>();
    data.forEach((item) => {
      map.set(normalizeRegionName(item.region), item.value);
    });
    return map;
  }, [data]);

  // Calculate stats for scaling
  const stats = useMemo(() => {
    if (data.length === 0) return { min: 0, max: 0, avg: 0 };
    const values = data.map((d) => d.value);
    return {
      min: Math.min(...values),
      max: Math.max(...values),
      avg: values.reduce((a, b) => a + b, 0) / values.length,
    };
  }, [data]);

  // Dynamic Color Scaling - Linear Interpolation
  const getColor = (value: number | undefined) => {
    if (value === undefined || stats.max === stats.min) return "#f8fafc"; // slate-50

    // Normalize value between 0 and 1
    const normalize = (value - stats.min) / (stats.max - stats.min);

    // Three-stop linear interpolation (Indigo Scale)
    // 0.0 -> #eef2ff (50)
    // 0.5 -> #6366f1 (500)
    // 1.0 -> #312e81 (900)

    const hexToRgb = (hex: string) => {
      const r = parseInt(hex.slice(1, 3), 16);
      const g = parseInt(hex.slice(3, 5), 16);
      const b = parseInt(hex.slice(5, 7), 16);
      return [r, g, b];
    };

    const rgbToHex = (r: number, g: number, b: number) =>
      "#" +
      [r, g, b]
        .map((x) => Math.round(x).toString(16).padStart(2, "0"))
        .join("");

    const low = hexToRgb("#eef2ff");
    const mid = hexToRgb("#6366f1");
    const high = hexToRgb("#312e81");

    let r, g, b;
    if (normalize < 0.5) {
      const t = normalize * 2;
      r = low[0] + (mid[0] - low[0]) * t;
      g = low[1] + (mid[1] - low[1]) * t;
      b = low[2] + (mid[2] - low[2]) * t;
    } else {
      const t = (normalize - 0.5) * 2;
      r = mid[0] + (high[0] - mid[0]) * t;
      g = mid[1] + (high[1] - mid[1]) * t;
      b = mid[2] + (high[2] - mid[2]) * t;
    }

    return rgbToHex(r, g, b);
  };

  const style = (feature: any) => {
    const rawName =
      feature.properties.name ||
      feature.properties.county ||
      feature.properties.LAD13NM ||
      "";
    const value = dataLookup.get(normalizeRegionName(rawName));

    return {
      fillColor: getColor(value),
      weight: 0.5,
      opacity: 1,
      color: "white",
      fillOpacity: 0.8,
    };
  };

  const onEachFeature = (feature: any, layer: L.Layer) => {
    const rawName =
      feature.properties.name ||
      feature.properties.county ||
      feature.properties.LAD13NM ||
      "";
    const value = dataLookup.get(normalizeRegionName(rawName));
    const formattedValue =
      value !== undefined ? formatValue(value) : "No data available";

    layer.bindTooltip(
      `
      <div class="p-3 bg-white shadow-xl rounded-lg border border-slate-100 min-w-40">
        <div class="text-[10px] font-bold text-slate-400 uppercase tracking-widest mb-1">UK Region</div>
        <div class="font-bold text-slate-900 text-sm mb-2 pb-2 border-b border-slate-50">${rawName}</div>
        <div class="flex items-end justify-between gap-2">
          <div>
            <div class="text-[10px] text-slate-500 uppercase font-semibold">${metricLabel}</div>
            <div class="text-lg font-bold text-indigo-600 leading-none">${formattedValue}</div>
          </div>
          <div class="w-8 h-8 rounded bg-indigo-50 flex items-center justify-center text-indigo-400">
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z"></path><circle cx="12" cy="10" r="3"></circle></svg>
          </div>
        </div>
      </div>
    `,
      {
        sticky: true,
        className: "custom-leaflet-tooltip",
        direction: "top",
        offset: [0, -10],
      },
    );

    layer.on({
      mouseover: (e) => {
        const layer = e.target;
        layer.setStyle({
          weight: 2.5,
          color: "#312e81",
          fillOpacity: 0.95,
        });
      },
      mouseout: (e) => {
        const layer = e.target;
        layer.setStyle({
          weight: 0.5,
          color: "white",
          fillOpacity: 0.8,
        });
        if (layer.closeTooltip) {
          layer.closeTooltip();
        }
      },
    });
  };

  if (loading) {
    return (
      <div className="flex flex-col items-center justify-center h-162.5 bg-slate-50/50 rounded-2xl border border-dashed border-slate-200">
        <div className="relative">
          <div className="w-16 h-16 border-4 border-indigo-100 rounded-full"></div>
          <div className="absolute top-0 left-0 w-16 h-16 border-4 border-indigo-600 rounded-full border-t-transparent animate-spin"></div>
        </div>
        <p className="mt-6 text-slate-500 font-semibold text-sm tracking-wide">
          SYNCHRONIZING GEOSPATIAL DATA
        </p>
      </div>
    );
  }

  return (
    <div className="h-162.5 w-full rounded-2xl overflow-hidden border border-slate-200 bg-white shadow-sm relative z-0 group">
      <MapContainer
        center={[54.5, -2.5] as any}
        zoom={6}
        scrollWheelZoom={true}
        className="h-full w-full"
        zoomControl={false}
      >
        <TileLayer
          attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> &copy; <a href="https://carto.com/attributions">CARTO</a>'
          url="https://{s}.basemaps.cartocdn.com/light_all/{z}/{x}/{y}{r}.png"
        />
        {geoJson && (
          <GeoJSON
            key={`${regionType}-${metricLabel}-${data.length}`}
            data={geoJson}
            style={style}
            onEachFeature={onEachFeature}
          />
        )}
      </MapContainer>

      {/* Modern Legend */}
      <div className="absolute bottom-8 right-8 bg-white/95 backdrop-blur-md p-5 rounded-2xl border border-slate-200 shadow-2xl z-1000 min-w-50 transition-all duration-300 group-hover:-translate-y-1">
        <div className="flex items-center justify-between mb-4">
          <h4 className="font-bold text-slate-800 text-xs uppercase tracking-widest flex items-center gap-2">
            <span className="w-2 h-2 rounded-full bg-indigo-500 animate-pulse"></span>
            {metricLabel}
          </h4>
        </div>

        <div className="space-y-4">
          {/* Continuous Gradient Bar */}
          <div className="relative pt-1">
            <div className="flex items-center justify-between text-[10px] font-bold text-slate-400 mb-1 px-0.5">
              <span>{formatValue(stats.min)}</span>
              <span>{formatValue(stats.max)}</span>
            </div>
            <div
              className="h-2.5 w-full rounded-full shadow-inner"
              style={{
                background:
                  "linear-gradient(to right, #eef2ff, #6366f1, #312e81)",
              }}
            ></div>
          </div>

          <div className="grid grid-cols-1 gap-2 border-t border-slate-50 pt-3">
            <div className="flex items-center justify-between text-[10px]">
              <span className="text-slate-500 font-medium">
                Regional Average
              </span>
              <span className="text-slate-900 font-bold">
                {formatValue(stats.avg)}
              </span>
            </div>
            <div className="flex items-center justify-between text-[10px]">
              <span className="text-slate-500 font-medium">Mapped Regions</span>
              <span className="text-slate-900 font-bold">{data.length}</span>
            </div>
          </div>
        </div>
      </div>

      {/* Floating Indicator */}
      <div className="absolute top-6 left-6 z-1000 pointer-events-none">
        <div className="bg-slate-900 text-white px-4 py-2 rounded-full shadow-lg flex items-center gap-3 border border-slate-800">
          <div className="w-2 h-2 rounded-full bg-emerald-400"></div>
          <span className="text-[10px] font-bold uppercase tracking-widest opacity-90">
            {regionType} Level View
          </span>
        </div>
      </div>

      {data.length === 0 && !loading && (
        <div className="absolute inset-0 flex items-center justify-center bg-slate-900/10 backdrop-blur-[2px] z-2000">
          <div className="bg-white p-6 rounded-2xl border shadow-2xl text-center max-w-xs animate-in zoom-in-95 duration-300">
            <div className="w-12 h-12 bg-amber-100 text-amber-600 rounded-full flex items-center justify-center mx-auto mb-4">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                width="24"
                height="24"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                strokeWidth="2"
                strokeLinecap="round"
                strokeLinejoin="round"
              >
                <path d="m21.73 18-8-14a2 2 0 0 0-3.48 0l-8 14A2 2 0 0 0 4 21h16a2 2 0 0 0 1.73-3Z" />
                <path d="M12 9v4" />
                <path d="M12 17h.01" />
              </svg>
            </div>
            <h4 className="font-bold text-slate-900 mb-1">No Market Data</h4>
            <p className="text-sm text-slate-500 mb-4">
              We couldn't find any {metricLabel.toLowerCase()} records for the
              selected region level.
            </p>
            <button
              onClick={() => window.location.reload()}
              className="px-4 py-2 bg-slate-900 text-white text-xs font-bold rounded-lg hover:bg-slate-800 transition-colors"
            >
              Retry Sync
            </button>
          </div>
        </div>
      )}
    </div>
  );
};

export default UKMap;
