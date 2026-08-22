// lib/api/updates-api.ts

import { fetchAPI } from "./client";

export type StatusLevel =
  | "normal"
  | "watch"
  | "important"
  | "emergency"
  | "delayed"
  | "unavailable";

export interface UpdateItem {
  id: string;
  status: StatusLevel;
  title: string;
  description: string;
  location: string;
  updated: string;
}

// Used only for local development when no backend is configured.
const devFallback: UpdateItem[] = [
  {
    id: "1",
    status: "emergency",
    title: "Flood warning",
    description:
      "Heavy rainfall may cause flooding within 24 hours in low-lying areas.",
    location: "Lower Valley",
    updated: "10:00",
  },
  {
    id: "2",
    status: "watch",
    title: "Water outlook",
    description:
      "Water availability may decline slightly over the next two weeks.",
    location: "Lower Valley",
    updated: "Today",
  },
  {
    id: "3",
    status: "normal",
    title: "Seasonal update",
    description:
      "Planting season is approaching. Soil moisture levels are currently adequate.",
    location: "Lower Valley",
    updated: "Yesterday",
  },
];

function mapStatus(contentType?: string): StatusLevel {
  switch (contentType) {
    case "alert":
      return "important";
    case "forecast":
      return "watch";
    default:
      return "normal";
  }
}

export async function getUpdates(): Promise<UpdateItem[]> {
  if (!process.env.NEXT_PUBLIC_API_URL) {
    return devFallback;
  }

  try {
    const data = await fetchAPI<{
      updates?: Array<{
        id: string;
        content_type?: string;
        headline?: string;
        body_text?: string;
        place_name?: string;
        generated_at?: string;
      }>;
    }>("/v1/updates?limit=20");

    const items = (data.updates ?? []).map((u) => ({
      id: u.id,
      status: mapStatus(u.content_type),
      title: u.headline ?? "Update",
      description: u.body_text ?? "",
      location: u.place_name ?? "",
      updated: u.generated_at
        ? new Date(u.generated_at).toLocaleDateString()
        : "",
    }));

    return items;
  } catch {
    return [];
  }
}
