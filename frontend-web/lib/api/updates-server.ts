/**
 * Server-only update/context/explore fetchers.
 * These MUST only be imported by Server Components.
 */

import type { UUID } from "@/types/common";
import type {
  ContextResponse,
  UpdatesListResponse,
  ExploreResponse,
} from "@/types/updates";
import { fetchServer } from "./server-client";

/**
 * Fetch today's content for the current user.
 */
export async function getToday(token?: string): Promise<ContextResponse> {
  return fetchServer<ContextResponse>("/v1/context/today", { token });
}

/**
 * Fetch paginated updates.
 */
export async function getUpdates(params: { page?: number; limit?: number; placeId?: string } = {}, token?: string): Promise<UpdatesListResponse> {
  const searchParams = new URLSearchParams();
  if (params.page) searchParams.set("page", String(params.page));
  if (params.limit) searchParams.set("limit", String(params.limit));
  if (params.placeId) searchParams.set("place_id", params.placeId);
  const qs = searchParams.toString();
  return fetchServer<UpdatesListResponse>(`/v1/updates${qs ? `?${qs}` : ""}`, { token });
}

/**
 * Fetch explore indicators.
 */
export async function getExploreIndicators(params: { category?: string; placeId?: string; page?: number; limit?: number } = {}, token?: string): Promise<ExploreResponse> {
  const searchParams = new URLSearchParams();
  if (params.category) searchParams.set("category", params.category);
  if (params.placeId) searchParams.set("place_id", params.placeId);
  if (params.page) searchParams.set("page", String(params.page));
  if (params.limit) searchParams.set("limit", String(params.limit));
  const qs = searchParams.toString();
  return fetchServer<ExploreResponse>(`/v1/explore${qs ? `?${qs}` : ""}`, { token });
}

/**
 * Fetch topic-specific indicators for a place.
 */
export async function getPlaceTopic(placeId: UUID, topic: string, token?: string): Promise<ExploreResponse> {
  return fetchServer<ExploreResponse>(`/v1/places/${placeId}/${topic}`, { token });
}
