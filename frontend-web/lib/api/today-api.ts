// lib/api/today-api.ts

import { apiFetch } from "@/lib/api/client";

export interface TodayData {
  location: string;
  lastUpdated: string;
  hasImportantUpdates: boolean;
  localOutlook: string;
  waterOutlook: string;
}

interface BackendContextResponse {
  place: {
    id: string;
    name: string;
    place_type: string;
    lat: number;
    lon: number;
    label?: string;
    is_primary: boolean;
  };
  update?: {
    id: string;
    place_id: string;
    place_name: string;
    content_type: string;
    headline: string;
    body_text: string;
    call_to_action?: string;
    generated_at: string;
    period_start: string;
    period_end: string;
  };
}

/**
 * Fetch today's context from the backend.
 * @param placeId - Optional place UUID to fetch context for a specific place.
 */
export async function getTodayData(placeId?: string): Promise<TodayData> {
  const path = placeId
    ? `/api/v1/context/places/${placeId}/today`
    : "/api/v1/context/today";

  const response = await apiFetch<BackendContextResponse>(path);
  const ctx = response.data;

  return {
    location: ctx.place.label ?? ctx.place.name,
    lastUpdated: ctx.update?.generated_at
      ? new Date(ctx.update.generated_at).toLocaleTimeString()
      : "N/A",
    hasImportantUpdates: ctx.update?.content_type === "alert",
    localOutlook: ctx.update?.headline ?? "No data available",
    waterOutlook: ctx.update?.body_text ?? "No data available",
  };
}
