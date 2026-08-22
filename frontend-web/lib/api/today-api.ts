// lib/api/today-api.ts

import { fetchAPI } from "./client";

export interface TodayData {
  location: string;
  lastUpdated: string;
  hasImportantUpdates: boolean;
  localOutlook: string;
  waterOutlook: string;
}

// Used only for local development when no backend is configured.
const devFallback: TodayData = {
  location: "Lower Valley",
  lastUpdated: "10:00",
  hasImportantUpdates: false,
  localOutlook:
    "Conditions for the coming days are stable with mild temperatures.",
  waterOutlook:
    "Water availability may decline slightly over the next two weeks.",
};

export async function getTodayData(): Promise<TodayData> {
  // Local dev without a backend: serve mock data so the UI is usable.
  if (!process.env.NEXT_PUBLIC_API_URL) {
    return devFallback;
  }

  try {
    const data = await fetchAPI<{
      place?: { name?: string };
      update?: { content_type?: string; body_text?: string };
    }>("/v1/context");

    return {
      location: data.place?.name ?? "",
      lastUpdated: new Date().toLocaleTimeString([], {
        hour: "2-digit",
        minute: "2-digit",
      }),
      hasImportantUpdates: data.update?.content_type === "alert",
      localOutlook: data.update?.body_text ?? "",
      waterOutlook: "",
    };
  } catch {
    // No fabricated data in production — surface an empty (honest) state.
    return {
      location: "",
      lastUpdated: "",
      hasImportantUpdates: false,
      localOutlook: "",
      waterOutlook: "",
    };
  }
}
