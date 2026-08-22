"use client";

import { useEffect, useRef, useCallback, useState } from "react";
import type { Map as MaplibreMap, Marker, Popup } from "maplibre-gl";

export interface MapPlace {
  id: string;
  name: string;
  lat: number;
  lon: number;
  place_type?: string;
  label?: string;
  is_primary?: boolean;
  urgencyScore?: number;
  indicatorCount?: number;
  status?: "normal" | "watch" | "warning" | "emergency";
}

interface MapViewProps {
  places: MapPlace[];
  selectedPlaceId?: string;
  onPlaceClick?: (place: MapPlace) => void;
  height?: string;
  zoom?: number;
  center?: [number, number];
}

const STATUS_COLORS: Record<string, string> = {
  normal: "#2e7d5b",
  watch: "#b98524",
  warning: "#c76b35",
  emergency: "#a93226",
};

const DEFAULT_STYLE = "https://basemaps.cartocdn.com/gl/positron-gl-style/style.json";

export function MapView({
  places,
  selectedPlaceId,
  onPlaceClick,
  height = "400px",
  zoom = 11,
  center,
}: MapViewProps) {
  const containerRef = useRef<HTMLDivElement>(null);
  const mapRef = useRef<MaplibreMap | null>(null);
  const markersRef = useRef<Marker[]>([]);
  const popupsRef = useRef<Popup[]>([]);
  const [mapReady, setMapReady] = useState(false);
  const [mapError, setMapError] = useState(false);

  // Initialize map
  useEffect(() => {
    if (!containerRef.current || mapRef.current) return;

    let mounted = true;

    const initMap = async () => {
      try {
        const maplibregl = await import("maplibre-gl");
        await import("maplibre-gl/dist/maplibre-gl.css");

        if (!mounted || !containerRef.current) return;

        const mapCenter: [number, number] =
          center || (places.length > 0 ? [places[0].lon, places[0].lat] : [36.8219, -1.2921]);

        const map = new maplibregl.Map({
          container: containerRef.current,
          style: process.env.NEXT_PUBLIC_MAP_STYLE_URL || DEFAULT_STYLE,
          center: mapCenter,
          zoom,
          attributionControl: false,
        });

        map.addControl(new maplibregl.NavigationControl());
        map.addControl(new maplibregl.ScaleControl({ unit: "metric" }));
        map.addControl(new maplibregl.AttributionControl({ compact: true }));

        map.on("load", () => {
          if (mounted) setMapReady(true);
        });

        map.on("error", () => {
          if (mounted) setMapError(true);
        });

        mapRef.current = map;
      } catch {
        if (mounted) setMapError(true);
      }
    };

    initMap();

    return () => {
      mounted = false;
      // Clean up markers and popups
      markersRef.current.forEach((m) => m.remove());
      popupsRef.current.forEach((p) => p.remove());
      markersRef.current = [];
      popupsRef.current = [];
      mapRef.current?.remove();
      mapRef.current = null;
    };
  }, []);

  // Update markers when places change
  useEffect(() => {
    if (!mapRef.current || !mapReady) return;

    let mounted = true;

    const updateMarkers = async () => {
      const maplibregl = await import("maplibre-gl");

      if (!mounted || !mapRef.current) return;

      // Clear existing markers
      markersRef.current.forEach((m) => m.remove());
      popupsRef.current.forEach((p) => p.remove());
      markersRef.current = [];
      popupsRef.current = [];

      const bounds = new maplibregl.LngLatBounds();

      places.forEach((place) => {
        const color = STATUS_COLORS[place.status || "normal"] || STATUS_COLORS.normal;
        const isSelected = place.id === selectedPlaceId;
        const size = isSelected ? 16 : 12;
        const borderWidth = isSelected ? 3 : 2;

        // Create marker element
        const el = document.createElement("div");
        el.style.cssText = `
          width: ${size}px;
          height: ${size}px;
          border-radius: 50%;
          background-color: ${color};
          border: ${borderWidth}px solid white;
          box-shadow: 0 2px 6px rgba(0,0,0,0.3);
          cursor: pointer;
          transition: transform 0.2s ease;
        `;
        if (isSelected) {
          el.style.transform = "scale(1.3)";
          el.style.zIndex = "10";
        }

        // Urgency ring for high-urgency places
        if (place.urgencyScore != null && place.urgencyScore > 50) {
          const ring = document.createElement("div");
          ring.style.cssText = `
            position: absolute;
            top: -6px;
            left: -6px;
            width: ${size + 12}px;
            height: ${size + 12}px;
            border-radius: 50%;
            border: 2px solid ${color};
            opacity: 0.4;
            animation: pulse-ring 2s ease-out infinite;
          `;
          el.style.position = "relative";
          el.appendChild(ring);
        }

        // Popup content
        const popupHtml = `
          <div style="padding: 8px; min-width: 180px; font-family: system-ui, sans-serif;">
            <div style="font-weight: 600; font-size: 14px; color: #17211d; margin-bottom: 4px;">
              ${place.name}
            </div>
            ${place.place_type ? `<div style="font-size: 11px; color: #64716b; text-transform: capitalize; margin-bottom: 6px;">${place.place_type}</div>` : ""}
            <div style="display: flex; gap: 8px; font-size: 12px; color: #64716b;">
              ${place.indicatorCount != null ? `<span>📊 ${place.indicatorCount} indicators</span>` : ""}
              ${place.urgencyScore != null ? `<span style="color: ${color}; font-weight: 500;">⚡ ${place.urgencyScore}%</span>` : ""}
            </div>
            ${place.status ? `
              <div style="margin-top: 6px;">
                <span style="
                  display: inline-flex; align-items: center; gap: 4px;
                  padding: 2px 8px; border-radius: 12px; font-size: 11px; font-weight: 500;
                  background: ${color}15; color: ${color}; border: 1px solid ${color}30;
                ">
                  <span style="width: 6px; height: 6px; border-radius: 50%; background: ${color};"></span>
                  ${place.status.charAt(0).toUpperCase() + place.status.slice(1)}
                </span>
              </div>
            ` : ""}
          </div>
        `;

        const popup = new maplibregl.Popup({
          offset: 15,
          closeButton: false,
          maxWidth: "280px",
        }).setHTML(popupHtml);

        const marker = new maplibregl.Marker({ element: el })
          .setLngLat([place.lon, place.lat])
          .setPopup(popup)
          .addTo(mapRef.current!);

        marker.getElement().addEventListener("click", () => {
          onPlaceClick?.(place);
        });

        markersRef.current.push(marker);
        popupsRef.current.push(popup);
        bounds.extend([place.lon, place.lat]);
      });

      // Fit bounds to show all places
      if (places.length > 1) {
        mapRef.current.fitBounds(bounds, { padding: 50, maxZoom: 14 });
      } else if (places.length === 1) {
        mapRef.current.setCenter([places[0].lon, places[0].lat]);
        mapRef.current.setZoom(zoom);
      }
    };

    updateMarkers();

    return () => {
      mounted = false;
    };
  }, [places, selectedPlaceId, mapReady, onPlaceClick, zoom]);

  // Update selected marker highlight
  useEffect(() => {
    if (!mapRef.current || !mapReady) return;

    // Fly to selected place
    const selected = places.find((p) => p.id === selectedPlaceId);
    if (selected) {
      mapRef.current.flyTo({
        center: [selected.lon, selected.lat],
        zoom: Math.max(mapRef.current.getZoom(), 12),
        essential: true,
      });
    }
  }, [selectedPlaceId, mapReady, places]);

  if (mapError) {
    return (
      <div
        className="rounded-xl border border-background-stone bg-background overflow-hidden"
        style={{ height }}
      >
        <div className="h-full flex flex-col items-center justify-center p-6 text-center">
          <div className="text-4xl mb-3">🗺</div>
          <p className="text-sm font-medium text-text-charcoal">Map unavailable</p>
          <p className="text-xs text-text-muted mt-1 max-w-xs">
            The map could not be loaded. Places are still listed below.
          </p>
        </div>
      </div>
    );
  }

  return (
    <div className="relative rounded-xl border border-background-stone bg-background overflow-hidden" style={{ height }}>
      {/* Map container */}
      <div ref={containerRef} className="h-full w-full" />

      {/* Loading overlay */}
      {!mapReady && !mapError && (
        <div className="absolute inset-0 flex items-center justify-center bg-background-mist">
          <div className="flex items-center gap-2 text-text-muted text-sm">
            <div className="h-4 w-4 border-2 border-forest border-t-transparent rounded-full animate-spin" />
            Loading map...
          </div>
        </div>
      )}

      {/* Legend */}
      {mapReady && !mapError && (
        <div className="absolute bottom-3 left-3 bg-background/90 backdrop-blur-sm rounded-lg border border-background-stone p-2.5 text-xs">
          <p className="font-medium text-text-charcoal mb-1.5">Status</p>
          <div className="space-y-1">
            {Object.entries(STATUS_COLORS).map(([key, color]) => (
              <div key={key} className="flex items-center gap-2">
                <span
                  className="w-2.5 h-2.5 rounded-full flex-shrink-0"
                  style={{ backgroundColor: color }}
                />
                <span className="text-text-muted capitalize">{key}</span>
              </div>
            ))}
          </div>
        </div>
      )}

      {/* Pulse animation */}
      <style jsx>{`
        @keyframes pulse-ring {
          0% { transform: scale(1); opacity: 0.4; }
          100% { transform: scale(1.8); opacity: 0; }
        }
      `}</style>
    </div>
  );
}
