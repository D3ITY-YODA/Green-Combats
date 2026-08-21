// types/places.ts

// Note: Adjust these import paths based on where you store your shared types.
// If you haven't created them yet, see the "Dependencies" section below.
import type { UUID } from "./common"; 
import type { GeoJSONGeometry } from "./geojson"; 

/**
 * Represents the hierarchical or thematic type of a place.
 */
export type PlaceType =
  | "country"
  | "region"
  | "district"
  | "ward"
  | "village"
  | "community"
  | "watershed"
  | "waterbody"
  | "custom_area";

/**
 * A geographic location or area within the Green Compass system.
 */
export interface Place {
  id: UUID;
  name: string;
  type: PlaceType;
  country_code: string;
  geometry?: GeoJSONGeometry;
  parent_place_id?: UUID;
}

/**
 * Represents a user's saved connection to a specific place.
 * This links a User to a Place and tracks preferences like primary status.
 */
export interface UserPlace {
  id: UUID;
  user_id: UUID;
  place_id: UUID;
  label?: string; // Optional user-defined nickname for the place
  is_primary: boolean;
  place: Place; // The fully populated Place object
}
