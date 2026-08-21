// lib/api/explore-api.ts

import { apiFetch } from "@/lib/api/client";

export interface ExploreSection {
  id: string;
  title: string;
  description: string;
  icon: "local" | "water" | "seasonal" | "community";
}

interface BackendExploreIndicator {
  id: string;
  place_id: string;
  place_name: string;
  code: string;
  display_name: string;
  category: string;
  value: number;
  unit: string;
  trend?: string;
  trend_confidence?: number;
  computed_at: string;
}

/** Map backend category to a frontend icon. */
function mapIcon(category: string): ExploreSection["icon"] {
  switch (category) {
    case "weather":
    case "local_outlook":
      return "local";
    case "water":
    case "water_outlook":
      return "water";
    case "agriculture":
    case "seasonal_information":
    case "land_and_ecosystems":
    case "food_and_agriculture":
      return "seasonal";
    case "community_updates":
      return "community";
    default:
      return "local";
  }
}

/** Build a human-readable description from an indicator. */
function describeIndicator(ind: BackendExploreIndicator): string {
  const trendText = ind.trend ? ` Trending ${ind.trend}.` : "";
  return `${ind.display_name}: ${ind.value}${ind.unit}.${trendText}`;
}

/**
 * Fetch explore indicators from the backend for the authenticated user.
 * Groups indicators by category into ExploreSections.
 */
export async function getExploreSections(): Promise<ExploreSection[]> {
  const response = await apiFetch<{ indicators: BackendExploreIndicator[] }>(
    "/api/v1/explore?limit=50",
  );
  const indicators = response.data.indicators ?? [];

  // Group by category
  const grouped = new Map<string, BackendExploreIndicator[]>();
  for (const ind of indicators) {
    const existing = grouped.get(ind.category) ?? [];
    existing.push(ind);
    grouped.set(ind.category, existing);
  }

  // Convert groups to sections
  const sections: ExploreSection[] = [];
  for (const [category, items] of grouped) {
    const descriptions = items.map(describeIndicator);
    sections.push({
      id: category,
      title: items[0].display_name,
      description: descriptions.join(" "),
      icon: mapIcon(category),
    });
  }

  return sections;
}
