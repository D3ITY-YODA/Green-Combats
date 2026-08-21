// types/common.ts

/**
 * Shared primitive types used across the frontend.
 */
export type UUID = string;

/**
 * ISO 8601 datetime string.
 */
export type Locale = string;

/**
 * Represents the freshness status of data for a place.
 */
export type DataStatus = "current" | "delayed" | "stale" | "unavailable";

/**
 * Public update as served by the backend (spec §9).
 * Mirrors the Go `updates.Update` JSON shape.
 */
export type UpdateType =
  | "information"
  | "important"
  | "warning"
  | "emergency"
  | "community";

export interface PublicUpdate {
  id: UUID;
  place_id?: UUID;
  topic_key?: string;
  type: UpdateType;
  priority: string;
  title: string;
  message: string;
  summary?: string;
  place_name?: string;
  valid_from: string;
  valid_until?: string;
  updated_at: string;
  source_name?: string;
  source_url?: string;
  status: string;
  display_on_today: boolean;
  display_in_feed: boolean;
  locale: Locale;
}
