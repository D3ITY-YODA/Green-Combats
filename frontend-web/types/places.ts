import { UUID } from "./common";

/**
 * Matches the backend Place model.
 */
export interface Place {
  id: UUID;
  name: string;
  place_type: PlaceType;
  lat: number;
  lon: number;
  external_code?: string | null;
  created_by?: UUID | null;
  created_at?: string;
  updated_at?: string;
}

export type PlaceType = "country" | "region" | "district" | "ward" | "village" | "community" | "watershed" | "waterbody" | "custom_area";

/**
 * Matches the backend SavedPlace from /v1/me/places.
 */
export interface SavedPlace {
  id: UUID;
  name: string;
  place_type: PlaceType;
  lat: number;
  lon: number;
  label?: string | null;
  is_primary: boolean;
  saved_at: string;
}

/**
 * Nearby place with distance — from /v1/places/nearby.
 */
export interface NearbyPlace extends Place {
  distance_m: number;
}

/**
 * Input for creating a place.
 */
export interface CreatePlaceInput {
  name: string;
  place_type: string;
  lat: number;
  lon: number;
}
