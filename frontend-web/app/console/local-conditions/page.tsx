"use client";

import { useState, useEffect, useCallback } from "react";
import { MapPin, TrendingUp, TrendingDown, Minus, Clock } from "lucide-react";
import { getExploreIndicators } from "@/lib/api/updates";
import { getSavedPlacesClient } from "@/lib/api/places";
import { MapView } from "@/components/ui/map-view";
import type { MapPlace } from "@/components/ui/map-view";
import type { ExploreIndicator } from "@/types/updates";
import type { SavedPlace } from "@/types/common";

const CATEGORY_META: Record<string, { label: string; color: string; icon: string }> = {
  weather: { label: "Weather", color: "text-sky", icon: "🌤" },
  water: { label: "Water", color: "text-forest", icon: "💧" },
  agriculture: { label: "Agriculture", color: "text-status-normal", icon: "🌾" },
  air_quality: { label: "Air Quality", color: "text-status-watch", icon: "🌬" },
  other: { label: "Other", color: "text-text-muted", icon: "📊" },
};

function TrendIcon({ trend }: { trend?: string | null }) {
  if (trend === "increasing") return <TrendingUp className="h-4 w-4 text-status-important" />;
  if (trend === "decreasing") return <TrendingDown className="h-4 w-4 text-status-normal" />;
  return <Minus className="h-4 w-4 text-text-muted" />;
}

function UrgencyBar({ score }: { score: number }) {
  let color = "bg-status-normal";
  let label = "Low";
  if (score > 70) { color = "bg-status-emergency"; label = "Critical"; }
  else if (score > 50) { color = "bg-status-important"; label = "High"; }
  else if (score > 30) { color = "bg-status-watch"; label = "Moderate"; }

  return (
    <div>
      <div className="flex items-center justify-between mb-1">
        <span className="text-xs font-medium text-text-muted">Urgency</span>
        <span className="text-xs font-semibold text-text-charcoal">{label} ({score}%)</span>
      </div>
      <div className="h-2 rounded-full bg-background-stone overflow-hidden">
        <div className={`h-full rounded-full ${color} transition-all`} style={{ width: `${score}%` }} />
      </div>
    </div>
  );
}

/** Compute an urgency-like score from indicator trend_confidence values */
function computeUrgencyScore(indicators: ExploreIndicator[]): number {
  if (indicators.length === 0) return 0;
  const avg = indicators.reduce((sum, i) => sum + (i.trend_confidence || 0.5), 0) / indicators.length;
  return Math.round(avg * 100);
}

/** Determine a status label from indicators */
function computeStatus(indicators: ExploreIndicator[]): MapPlace["status"] {
  if (indicators.length === 0) return "normal";
  const avgConf = indicators.reduce((sum, i) => sum + (i.trend_confidence || 0.5), 0) / indicators.length;
  // "increasing" drought_stress or decreasing soil_moisture means worse conditions
  const badIndicators = indicators.filter(
    (i) =>
      (i.trend === "increasing" && (i.code.includes("drought") || i.code.includes("risk"))) ||
      (i.trend === "decreasing" && (i.code.includes("moisture") || i.code.includes("quality")))
  );
  const badRatio = badIndicators.length / indicators.length;
  if (badRatio > 0.5 && avgConf > 0.7) return "emergency";
  if (badRatio > 0.3 || avgConf > 0.8) return "warning";
  if (badRatio > 0.1 || avgConf > 0.6) return "watch";
  return "normal";
}

export default function ConsoleLocalConditions() {
  const [indicators, setIndicators] = useState<ExploreIndicator[]>([]);
  const [savedPlaces, setSavedPlaces] = useState<SavedPlace[]>([]);
  const [selectedPlace, setSelectedPlace] = useState<string>("");
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    setLoading(true);
    Promise.all([
      getExploreIndicators({ limit: 100 }).catch(() => ({ indicators: [], total: 0, page: 1, limit: 100, has_next: false })),
      getSavedPlacesClient().catch(() => []),
    ])
      .then(([exploreResult, places]) => {
        setIndicators(exploreResult.indicators);
        setSavedPlaces(places);
        if (exploreResult.indicators.length > 0 && !selectedPlace) {
          setSelectedPlace(exploreResult.indicators[0].place_id);
        }
      })
      .finally(() => setLoading(false));
  }, []);

  // Build unique place list from indicators (with names)
  const indicatorPlaces = Array.from(
    new Map(indicators.map((i) => [i.place_id, i.place_name])).entries()
  );

  // Merge: prefer saved places (which have lat/lon) with indicator data
  const allPlaceIds = new Set([...indicatorPlaces.map(([id]) => id), ...savedPlaces.map((p) => p.id)]);
  const placesByName = new Map(indicatorPlaces);
  const placesById = new Map(savedPlaces.map((p) => [p.id, p]));

  // Build map places
  const mapPlaces: MapPlace[] = [];
  for (const id of allPlaceIds) {
    const saved = placesById.get(id);
    const name = saved?.name || placesByName.get(id) || "Unknown";
    const placeIndicators = indicators.filter((i) => i.place_id === id);

    // Prefer saved place coordinates, fall back to defaults around Nairobi
    const lat = saved?.lat ?? (-1.2921 + Math.random() * 0.05);
    const lon = saved?.lon ?? (36.8219 + Math.random() * 0.05);

    mapPlaces.push({
      id,
      name,
      lat,
      lon,
      place_type: saved?.place_type,
      is_primary: saved?.is_primary,
      urgencyScore: placeIndicators.length > 0 ? computeUrgencyScore(placeIndicators) : undefined,
      indicatorCount: placeIndicators.length,
      status: placeIndicators.length > 0 ? computeStatus(placeIndicators) : undefined,
    });
  }

  // Filter indicators for selected place
  const placeIndicators = indicators.filter((i) => i.place_id === selectedPlace);

  // Group by category
  const grouped = new Map<string, ExploreIndicator[]>();
  for (const ind of placeIndicators) {
    const cat = ind.category || "other";
    if (!grouped.has(cat)) grouped.set(cat, []);
    grouped.get(cat)!.push(ind);
  }

  const selectedPlaceName = placesByName.get(selectedPlace) || savedPlaces.find((p) => p.id === selectedPlace)?.name || "All places";

  const handleMapPlaceClick = useCallback((place: MapPlace) => {
    setSelectedPlace(place.id);
  }, []);

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-text-charcoal">Local Conditions</h1>
        <p className="text-text-muted mt-1">Monitor environmental indicators and urgency levels for each place.</p>
      </div>

      {/* Map view */}
      {mapPlaces.length > 0 && (
        <div className="rounded-xl border border-background-stone overflow-hidden">
          <MapView
            places={mapPlaces}
            selectedPlaceId={selectedPlace}
            onPlaceClick={handleMapPlaceClick}
            height="380px"
          />
        </div>
      )}

      {/* Place selector */}
      <div className="flex items-center gap-3">
        <MapPin className="h-5 w-5 text-forest" />
        <select
          value={selectedPlace}
          onChange={(e) => setSelectedPlace(e.target.value)}
          className="rounded-lg border border-background-stone bg-background px-3 py-2 text-sm text-text-charcoal focus:border-forest focus:ring-1 focus:ring-forest"
        >
          {indicatorPlaces.map(([id, name]) => (
            <option key={id} value={id}>{name}</option>
          ))}
        </select>
      </div>

      {loading ? (
        <div className="rounded-xl border border-background-stone bg-background p-8 text-center">
          <Clock className="h-6 w-6 text-text-muted mx-auto mb-2 animate-pulse" />
          <p className="text-text-muted text-sm">Loading conditions...</p>
        </div>
      ) : placeIndicators.length === 0 ? (
        <div className="rounded-xl border border-background-stone bg-background p-8 text-center">
          <p className="text-text-muted">No indicator data available for this place.</p>
        </div>
      ) : (
        <>
          {/* Urgency summary */}
          <div className="rounded-xl border border-background-stone bg-background p-5">
            <h2 className="text-sm font-semibold text-text-charcoal mb-3">Condition Summary — {selectedPlaceName}</h2>
            <UrgencyBar score={computeUrgencyScore(placeIndicators)} />
            <p className="mt-3 text-xs text-text-muted">
              Based on {placeIndicators.length} indicators across {grouped.size} categories.
            </p>
          </div>

          {/* Indicator categories */}
          <div className="space-y-4">
            {Array.from(grouped.entries()).map(([category, items]) => {
              const meta = CATEGORY_META[category] || CATEGORY_META.other;
              return (
                <div key={category} className="rounded-xl border border-background-stone bg-background overflow-hidden">
                  <div className="px-5 py-3 bg-background-mist border-b border-background-stone flex items-center gap-2">
                    <span className="text-lg">{meta.icon}</span>
                    <h3 className="text-sm font-semibold text-text-charcoal">{meta.label}</h3>
                    <span className="text-xs text-text-muted ml-auto">{items.length} indicator{items.length !== 1 ? "s" : ""}</span>
                  </div>
                  <div className="divide-y divide-background-stone">
                    {items.map((ind) => (
                      <div key={ind.id} className="px-5 py-3 flex items-center justify-between">
                        <div className="flex-1 min-w-0">
                          <p className="text-sm font-medium text-text-charcoal truncate">{ind.display_name}</p>
                          <p className="text-xs text-text-muted">{ind.code}</p>
                        </div>
                        <div className="flex items-center gap-3 ml-4">
                          <div className="text-right">
                            <span className="text-sm font-semibold text-text-charcoal">{ind.value}</span>
                            <span className="text-xs text-text-muted ml-1">{ind.unit}</span>
                          </div>
                          <TrendIcon trend={ind.trend} />
                          {ind.trend_confidence != null && (
                            <span className="text-xs text-text-muted" title="Confidence">
                              {Math.round(ind.trend_confidence * 100)}%
                            </span>
                          )}
                        </div>
                      </div>
                    ))}
                  </div>
                </div>
              );
            })}
          </div>
        </>
      )}
    </div>
  );
}
