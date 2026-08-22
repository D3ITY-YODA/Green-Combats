"use client";

import { useState, useEffect } from "react";
import { Database, BarChart3, Clock, Layers, Info } from "lucide-react";
import { getExploreIndicators } from "@/lib/api/updates";
import { getDashboard } from "@/lib/api/console";
import type { ExploreIndicator } from "@/types/updates";

const CATEGORY_INFO: Record<string, { label: string; description: string; icon: string }> = {
  weather: { label: "Weather", description: "Temperature, precipitation, wind speed and other meteorological data.", icon: "🌤" },
  water: { label: "Water", description: "Water quality, levels, turbidity and availability metrics.", icon: "💧" },
  agriculture: { label: "Agriculture", description: "Soil moisture, drought stress, and crop health indicators.", icon: "🌾" },
  air_quality: { label: "Air Quality", description: "PM2.5, ozone, and composite air quality index measurements.", icon: "🌬" },
  other: { label: "Other", description: "Additional environmental monitoring data.", icon: "📊" },
};

export default function ConsoleInformation() {
  const [indicators, setIndicators] = useState<ExploreIndicator[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    setLoading(true);
    getExploreIndicators({ limit: 100 })
      .then((result) => setIndicators(result.indicators))
      .catch(() => {})
      .finally(() => setLoading(false));
  }, []);

  // Group indicators by category
  const grouped = new Map<string, ExploreIndicator[]>();
  for (const ind of indicators) {
    const cat = ind.category || "other";
    if (!grouped.has(cat)) grouped.set(cat, []);
    grouped.get(cat)!.push(ind);
  }

  // Count unique places being monitored
  const uniquePlaces = new Set(indicators.map((i) => i.place_id));

  // Count unique indicator types
  const uniqueIndicatorCodes = new Set(indicators.map((i) => i.code));

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-text-charcoal">Information</h1>
        <p className="text-text-muted mt-1">Indicator definitions, data sources, and monitoring infrastructure.</p>
      </div>

      {/* Summary stats */}
      <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
        <div className="bg-background rounded-xl border border-background-stone p-5 flex items-center gap-4">
          <div className="p-3 rounded-lg bg-forest/10 text-forest">
            <Layers className="h-5 w-5" />
          </div>
          <div>
            <p className="text-2xl font-bold text-text-charcoal">{uniqueIndicatorCodes.size}</p>
            <p className="text-xs text-text-muted">Indicator Types</p>
          </div>
        </div>
        <div className="bg-background rounded-xl border border-background-stone p-5 flex items-center gap-4">
          <div className="p-3 rounded-lg bg-sky/10 text-sky">
            <Database className="h-5 w-5" />
          </div>
          <div>
            <p className="text-2xl font-bold text-text-charcoal">{uniquePlaces.size}</p>
            <p className="text-xs text-text-muted">Places Monitored</p>
          </div>
        </div>
        <div className="bg-background rounded-xl border border-background-stone p-5 flex items-center gap-4">
          <div className="p-3 rounded-lg bg-status-normal/10 text-status-normal">
            <BarChart3 className="h-5 w-5" />
          </div>
          <div>
            <p className="text-2xl font-bold text-text-charcoal">{grouped.size}</p>
            <p className="text-xs text-text-muted">Categories</p>
          </div>
        </div>
      </div>

      {/* Indicator categories */}
      {loading ? (
        <div className="rounded-xl border border-background-stone bg-background p-8 text-center">
          <Clock className="h-6 w-6 text-text-muted mx-auto mb-2 animate-pulse" />
          <p className="text-text-muted text-sm">Loading indicator definitions...</p>
        </div>
      ) : (
        <div className="space-y-4">
          {Array.from(grouped.entries()).map(([category, items]) => {
            const info = CATEGORY_INFO[category] || CATEGORY_INFO.other;
            // Deduplicate by code
            const uniqueItems = Array.from(new Map(items.map((i) => [i.code, i])).values());

            return (
              <div key={category} className="rounded-xl border border-background-stone bg-background overflow-hidden">
                <div className="px-5 py-4 bg-background-mist border-b border-background-stone">
                  <div className="flex items-center gap-2">
                    <span className="text-lg">{info.icon}</span>
                    <h3 className="text-sm font-semibold text-text-charcoal">{info.label}</h3>
                    <span className="text-xs text-text-muted ml-auto">{uniqueItems.length} indicator{uniqueItems.length !== 1 ? "s" : ""}</span>
                  </div>
                  <p className="text-xs text-text-muted mt-1">{info.description}</p>
                </div>
                <div className="divide-y divide-background-stone">
                  {uniqueItems.map((ind) => (
                    <div key={ind.code} className="px-5 py-3">
                      <div className="flex items-center justify-between">
                        <div>
                          <p className="text-sm font-medium text-text-charcoal">{ind.display_name}</p>
                          <p className="text-xs text-text-muted font-mono">{ind.code}</p>
                        </div>
                        <div className="text-right">
                          <p className="text-sm font-semibold text-text-charcoal">{ind.value} <span className="text-xs font-normal text-text-muted">{ind.unit}</span></p>
                          {ind.trend && (
                            <p className="text-xs text-text-muted capitalize">Trend: {ind.trend}</p>
                          )}
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            );
          })}

          {indicators.length === 0 && (
            <div className="rounded-xl border border-background-stone bg-background p-8 text-center">
              <Info className="h-8 w-8 text-text-muted mx-auto mb-2" />
              <p className="text-text-muted">No indicator data available. Run the seed script to populate monitoring data.</p>
            </div>
          )}
        </div>
      )}

      {/* Data source info */}
      <div className="rounded-xl border border-background-stone bg-background p-5">
        <h3 className="text-sm font-semibold text-text-charcoal mb-2">Data Sources</h3>
        <div className="space-y-2 text-xs text-text-muted">
          <p>Indicator values are computed from ingestion pipelines that aggregate raw sensor data, satellite imagery, and community observations.</p>
          <p>Confidence scores reflect data freshness and source reliability. Higher confidence indicates more recent and reliable data.</p>
          <p>All indicators are recalculated on each ingestion cycle (typically every 6 hours).</p>
        </div>
      </div>
    </div>
  );
}
