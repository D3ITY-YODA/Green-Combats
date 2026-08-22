// lib/api/explore-api.ts

import { fetchAPI } from "./client";

export interface ExploreSection {
  id: string;
  title: string;
  description: string;
  icon: "local" | "water" | "seasonal" | "community";
}

// Used only for local development when no backend is configured.
const devFallback: ExploreSection[] = [
  {
    id: "1",
    title: "Local outlook",
    description:
      "Conditions for the coming days are stable with mild temperatures.",
    icon: "local",
  },
  {
    id: "2",
    title: "Water outlook",
    description:
      "Water availability may decline slightly over the next two weeks.",
    icon: "water",
  },
  {
    id: "3",
    title: "Seasonal information",
    description:
      "Planting season is approaching. Soil moisture levels are adequate.",
    icon: "seasonal",
  },
];

export async function getExploreSections(): Promise<ExploreSection[]> {
  if (!process.env.NEXT_PUBLIC_API_URL) {
    return devFallback;
  }

  try {
    const data = await fetchAPI<{
      indicators?: Array<{ display_name?: string; code?: string }>;
    }>("/v1/explore");

    const indicators = data.indicators ?? [];
    if (indicators.length === 0) return [];

    return indicators.slice(0, 6).map((ind, i) => ({
      id: String(i + 1),
      title: ind.display_name ?? "Explore",
      description: "Relevant environmental information for your area.",
      icon: (["local", "water", "seasonal", "community"] as const)[i % 4],
    }));
  } catch {
    return [];
  }
}
