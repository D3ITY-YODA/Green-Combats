/**
 * Client-safe places module — client-side fetchers only.
 * Server-side fetchers are in ./places-server.ts
 */

import type { UUID } from "@/types/common";
import type { Place, SavedPlace, NearbyPlace, CreatePlaceInput } from "@/types/places";
import { fetchAPI } from "./client";

export async function getSavedPlacesClient(): Promise<SavedPlace[]> {
  const result = await fetchAPI<{ places: SavedPlace[] }>('/v1/me/places');
  return result.places;
}

export async function searchNearbyPlaces(lat: number, lon: number, radiusM?: number, limit?: number): Promise<NearbyPlace[]> {
  const params = new URLSearchParams({ lat: String(lat), lon: String(lon) });
  if (radiusM) params.set("radius_m", String(radiusM));
  if (limit) params.set("limit", String(limit));
  const result = await fetchAPI<{ places: NearbyPlace[] }>(`/v1/places/nearby?${params}`);
  return result.places;
}

export async function savePlace(placeId: UUID, label?: string): Promise<void> {
  await fetchAPI(`/v1/me/places/${placeId}`, { method: "PUT", body: label ? { label } : undefined });
}

export async function unsavePlace(placeId: UUID): Promise<void> {
  await fetchAPI(`/v1/me/places/${placeId}`, { method: "DELETE" });
}

export async function setPrimaryPlace(placeId: UUID): Promise<void> {
  await fetchAPI(`/v1/me/places/${placeId}/primary`, { method: "PUT" });
}

export async function createPlace(input: CreatePlaceInput): Promise<Place> {
  return fetchAPI<Place>("/v1/places", { method: "POST", body: input });
}
