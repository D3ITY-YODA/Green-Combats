/**
 * Client-safe updates module — adapters and client-side fetchers.
 * Server-side fetchers are in ./updates-server.ts
 */

import type { TopicSection, DataStatus } from "@/types/common";
import type {
  Update,
  PublicUpdate,
  TodayResponse,
  ContextResponse,
  UpdatesListResponse,
  ExploreIndicator,
} from "@/types/updates";
import { fetchAPI } from "./client";

// ── Client-side fetchers ──────────────────────────────────────────

export async function getTodayClient(): Promise<TodayResponse> {
  const ctx = await fetchAPI<ContextResponse>("/v1/context/today");
  return adaptContextToTodayResponse(ctx);
}

export async function getUpdatesClient(params: { page?: number; limit?: number; placeId?: string } = {}): Promise<UpdatesListResponse> {
  const searchParams = new URLSearchParams();
  if (params.page) searchParams.set("page", String(params.page));
  if (params.limit) searchParams.set("limit", String(params.limit));
  if (params.placeId) searchParams.set("place_id", params.placeId);
  const qs = searchParams.toString();
  return fetchAPI<UpdatesListResponse>(`/v1/updates${qs ? `?${qs}` : ""}`);
}

export async function getExploreIndicators(params: { category?: string; placeId?: string; page?: number; limit?: number } = {}): Promise<{ indicators: ExploreIndicator[]; total: number; page: number; limit: number; has_next: boolean }> {
  const searchParams = new URLSearchParams();
  if (params.category) searchParams.set("category", params.category);
  if (params.placeId) searchParams.set("place_id", params.placeId);
  if (params.page) searchParams.set("page", String(params.page));
  if (params.limit) searchParams.set("limit", String(params.limit));
  const qs = searchParams.toString();
  return fetchAPI(`/v1/explore${qs ? `?${qs}` : ""}`);
}

// ── Adapters ──────────────────────────────────────────────────────

/**
 * Map backend ContextResponse → frontend TodayResponse.
 */
export function adaptContextToTodayResponse(ctx: ContextResponse): TodayResponse {
  const update = ctx.update;
  const dataStatus: DataStatus = update ? "current" : "unavailable";

  const status = update
    ? { title: update.headline, message: update.body_text.slice(0, 200), updated_at: update.generated_at, data_status: dataStatus }
    : { title: "No updates", message: "There are no updates for this location yet.", data_status: dataStatus };

  const updates: PublicUpdate[] = update ? [adaptUpdateToPublicUpdate(update)] : [];

  const sections: TopicSection[] = [
    { key: "local_outlook", title: "Local outlook", description: "Conditions for the coming days.", available: true, href: "/explore/local-outlook" },
    { key: "water_outlook", title: "Water outlook", description: "Water levels and availability.", available: true, href: "/explore/water" },
  ];

  return { place: ctx.place, status, updates, sections, updated_at: update?.generated_at || new Date().toISOString() };
}

/**
 * Map backend Update → frontend PublicUpdate.
 */
export function adaptUpdateToPublicUpdate(update: Update): PublicUpdate {
  const typeMap: Record<string, PublicUpdate["type"]> = { alert: "warning", forecast: "information", today: "information" };
  const priorityMap: Record<string, PublicUpdate["priority"]> = { alert: "important", forecast: "normal", today: "normal" };

  return {
    id: update.id,
    place_id: update.place_id,
    topic_key: update.content_type,
    type: typeMap[update.content_type] || "information",
    priority: priorityMap[update.content_type] || "normal",
    title: update.headline,
    message: update.body_text,
    valid_from: update.period_start,
    valid_until: update.period_end,
    updated_at: update.generated_at,
    status: "published",
    display_on_today: true,
    display_in_feed: true,
    locale: "en",
    place_name: update.place_name,
  };
}

/**
 * Map backend ExploreIndicator[] → frontend TopicSection[].
 */
export function adaptExploreToTopicSections(indicators: ExploreIndicator[]): TopicSection[] {
  const grouped = new Map<string, ExploreIndicator[]>();
  for (const ind of indicators) {
    const cat = ind.category || "other";
    if (!grouped.has(cat)) grouped.set(cat, []);
    grouped.get(cat)!.push(ind);
  }

  const topicMap: Record<string, { title: string; description: string; href: string }> = {
    weather: { title: "Weather", description: "Current weather conditions and forecasts.", href: "/explore/local-outlook" },
    water: { title: "Water outlook", description: "Water levels, availability and quality.", href: "/explore/water" },
    agriculture: { title: "Agriculture", description: "Crop and farming conditions.", href: "/explore/seasonal" },
    air_quality: { title: "Air quality", description: "Air quality measurements.", href: "/explore/local-outlook" },
    other: { title: "Other indicators", description: "Additional environmental data.", href: "/explore/local-outlook" },
  };

  return Array.from(grouped.entries()).map(([key, items]) => {
    const meta = topicMap[key] || topicMap.other;
    return { key, title: meta.title, description: `${items.length} indicator${items.length !== 1 ? "s" : ""}. ${meta.description}`, available: true, href: meta.href };
  });
}
