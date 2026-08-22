/**
 * Server-only place fetchers.
 * These MUST only be imported by Server Components.
 */

import type { UUID } from "@/types/common";
import type { Place, SavedPlace } from "@/types/places";
import { fetchServer } from "./server-client";

export async function getSavedPlaces(token?: string): Promise<SavedPlace[]> {
  const result = await fetchServer<{ places: SavedPlace[] }>("/v1/me/places", { token });
  return result.places;
}

export async function getPlace(id: UUID, token?: string): Promise<Place> {
  return fetchServer<Place>(`/v1/places/${id}`, { token });
}
