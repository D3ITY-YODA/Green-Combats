// lib/api/updates-api.ts

import { apiFetch } from "@/lib/api/client";

export type StatusLevel = "normal" | "watch" | "important" | "emergency" | "delayed" | "unavailable";

export interface UpdateItem {
  id: string;
  status: StatusLevel;
  title: string;
  description: string;
  location: string;
  updated: string;
}

interface BackendUpdate {
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
}

/** Map backend content_type to a UI status level. */
function mapStatus(contentType: string): StatusLevel {
  switch (contentType) {
    case "alert":
      return "emergency";
    case "forecast":
      return "watch";
    case "today":
      return "normal";
    default:
      return "normal";
  }
}

/**
 * Fetch updates from the backend for the authenticated user's places.
 * @param placeId - Optional place UUID to filter updates.
 */
export async function getUpdates(placeId?: string): Promise<UpdateItem[]> {
  const params = new URLSearchParams();
  if (placeId) params.set("place_id", placeId);
  const query = params.toString();
  const path = `/api/v1/updates${query ? `?${query}` : ""}`;

  const response = await apiFetch<{ updates: BackendUpdate[] }>(path);
  const updates = response.data.updates ?? [];

  return updates.map((u) => ({
    id: u.id,
    status: mapStatus(u.content_type),
    title: u.headline,
    description: u.body_text,
    location: u.place_name,
    updated: new Date(u.generated_at).toLocaleTimeString(),
  }));
}
