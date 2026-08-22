export type UUID = string;
export type Locale = "en" | "sw" | "fr" | string;

export type DataStatus = "current" | "delayed" | "stale" | "unavailable" | "not_configured";
export type RiskState = "normal" | "watch" | "warning" | "emergency" | "stale" | "unavailable";

/**
 * Matches the backend PlaceInfo model.
 */
export interface Place {
  id: UUID;
  name: string;
  place_type: string;
  lat: number;
  lon: number;
  label?: string | null;
  is_primary: boolean;
}

/**
 * Matches the backend SavedPlace model (from /v1/me/places).
 */
export interface SavedPlace {
  id: UUID;
  name: string;
  place_type: string;
  lat: number;
  lon: number;
  label?: string | null;
  is_primary: boolean;
  saved_at: string;
}

/**
 * Legacy TopicSection used by explore cards — kept for component compatibility.
 */
export interface TopicSection {
  key: string;
  title: string;
  description: string;
  available: boolean;
  data_status?: DataStatus;
  href?: string;
}
