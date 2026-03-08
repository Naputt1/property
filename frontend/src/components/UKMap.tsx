import React, { useEffect, useState, useMemo } from "react";
import { MapContainer, TileLayer, GeoJSON, useMap } from "react-leaflet";
import L from "leaflet";
import "leaflet/dist/leaflet.css";

// Fix Leaflet marker icon issue with Vite
// @ts-ignore
delete L.Icon.Default.prototype._getIconUrl;
L.Icon.Default.mergeOptions({
  iconRetinaUrl:
    "https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.7.1/images/marker-icon-2x.png",
  iconUrl:
    "https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.7.1/images/marker-icon.png",
  shadowUrl:
    "https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.7.1/images/marker-shadow.png",
});

interface UKMapProps {
  data: Array<{ region: string; value: number }>;
  regionType: "county" | "district";
  metricLabel: string;
  formatValue: (value: number) => string;
}

const COUNTY_GEOJSON_URL =
  "https://raw.githubusercontent.com/evansd/uk-ceremonial-counties/master/uk-ceremonial-counties.geojson";
// const DISTRICT_GEOJSON_URL = "https://raw.githubusercontent.com/martinjc/UK-GeoJSON/master/json/administrative/gb/lad.json";

const UKMap: React.FC<UKMapProps> = ({
  data,
  regionType,
  metricLabel,
  formatValue,
}) => {
  const [geoJson, setGeoJson] = useState<any>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const controller = new AbortController();
    const signal = controller.signal;

    setLoading(true);
    const url =
      regionType === "county" ? COUNTY_GEOJSON_URL : COUNTY_GEOJSON_URL; // Default to county for now

    fetch(url, { signal })
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

  const normalizeName = (name: string) => {
    if (!name) return "";
    return name
      .toUpperCase()
      .replace(/^THE\s+/, "")
      .replace(/\s+COUNTY$/, "")
      .replace(/COUNTY\s+OF\s+/, "")
      .replace(/\s+CITY\s+OF$/, "")
      .replace(/CITY\s+OF\s+/, "")
      .replace(/\s+BOROUGH\s+OF$/, "")
      .replace(/BOROUGH\s+OF\s+/, "")
      .replace(/&/g, "AND")
      .replace(/\bST\./g, "SAINT")
      .replace(/\bST\s/g, "SAINT ")
      .replace(/-/g, " ")
      .replace(/\'/g, "") // remove apostrophes (e.g. King's Lynn)
      .replace(/\s+/g, " ") // normalize multiple spaces
      .trim();
  };

  const dataMap = useMemo(() => {
    const map = new Map<string, number>();
    data.forEach((item) => {
      map.set(normalizeName(item.region), item.value);
    });
    console.log(map, data);
    return map;
  }, [data]);

  const maxValue = useMemo(() => {
    if (data.length === 0) return 0;
    return Math.max(...data.map((d) => d.value));
  }, [data]);

  const getColor = (value: number | undefined) => {
    if (value === undefined) return "#f1f5f9"; // Slate-100 for no data

    // Simple color scale: light to dark blue/indigo
    const ratio = maxValue > 0 ? value / maxValue : 0;

    if (ratio > 0.8) return "#312e81"; // Indigo-900
    if (ratio > 0.6) return "#4338ca"; // Indigo-700
    if (ratio > 0.4) return "#6366f1"; // Indigo-500
    if (ratio > 0.2) return "#818cf8"; // Indigo-400
    return "#c7d2fe"; // Indigo-200
  };

  const style = (feature: any) => {
    const name =
      feature.properties.name ||
      feature.properties.NAME_1 ||
      feature.properties.LAD13NM ||
      "";
    const value = dataMap.get(normalizeName(name));

    return {
      fillColor: getColor(value),
      weight: 1,
      opacity: 1,
      color: "white",
      fillOpacity: 0.7,
    };
  };

  const onEachFeature = (feature: any, layer: L.Layer) => {
    const rawName =
      feature.properties.name ||
      feature.properties.NAME_1 ||
      feature.properties.LAD13NM ||
      "Unknown";
    const normalizedName = normalizeName(rawName);
    const value = dataMap.get(normalizedName);
    const formattedValue = value !== undefined ? formatValue(value) : "No data";

    // Debugging: Log names that don't match
    if (value === undefined && rawName !== "Unknown") {
      console.log(
        `No match for GeoJSON region: "${rawName}" (normalized: "${normalizedName}")`,
      );
    }

    layer.bindTooltip(
      `
      <div class="p-1">
        <div class="font-bold">${rawName}</div>
        <div class="text-sm">${metricLabel}: ${formattedValue}</div>
      </div>
    `,
      { sticky: true },
    );

    layer.on({
      mouseover: (e) => {
        const layer = e.target;
        layer.setStyle({
          weight: 2,
          color: "#666",
          fillOpacity: 0.9,
        });
        layer.bringToFront();
      },
      mouseout: (e) => {
        const layer = e.target;
        layer.setStyle({
          weight: 1,
          color: "white",
          fillOpacity: 0.7,
        });
      },
    });
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center h-[500px] bg-muted/20 rounded-lg border border-dashed">
        <div className="text-center">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary mx-auto mb-2"></div>
          <p className="text-muted-foreground">Loading UK Map data...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="h-[600px] w-full rounded-xl overflow-hidden border shadow-sm relative z-0">
      <MapContainer
        center={[54.5, -2.5] as any} // Center of UK
        zoom={6}
        scrollWheelZoom={true}
        className="h-full w-full"
      >
        <TileLayer
          attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
          url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
        />
        {geoJson && data.length > 0 && (
          <GeoJSON
            key={`${regionType}-${metricLabel}`}
            data={geoJson}
            style={style}
            onEachFeature={onEachFeature}
          />
        )}
      </MapContainer>

      {geoJson && data.length === 0 && (
        <div className="absolute top-4 left-1/2 -translate-x-1/2 bg-yellow-50 text-yellow-800 px-4 py-2 rounded-full border border-yellow-200 shadow-md z-[1000] text-sm font-medium">
          No data available for the current selection.
        </div>
      )}

      {/* Legend */}
      <div className="absolute bottom-4 right-4 bg-white/90 p-3 rounded-lg border shadow-md z-[1000] text-xs">
        <h4 className="font-semibold mb-2">{metricLabel}</h4>
        <div className="space-y-1">
          <div className="flex items-center gap-2">
            <div
              className="w-4 h-4 rounded"
              style={{ backgroundColor: "#312e81" }}
            ></div>
            <span>High</span>
          </div>
          <div className="flex items-center gap-2">
            <div
              className="w-4 h-4 rounded"
              style={{ backgroundColor: "#6366f1" }}
            ></div>
            <span>Medium</span>
          </div>
          <div className="flex items-center gap-2">
            <div
              className="w-4 h-4 rounded"
              style={{ backgroundColor: "#c7d2fe" }}
            ></div>
            <span>Low</span>
          </div>
          <div className="flex items-center gap-2">
            <div
              className="w-4 h-4 rounded"
              style={{ backgroundColor: "#f1f5f9" }}
            ></div>
            <span>No data</span>
          </div>
        </div>
      </div>
    </div>
  );
};

export default UKMap;
