import { UUID, Locale, Place, DataStatus, TopicSection } from "./common";

// --- Backend-aligned types ---

/**
 * Matches the backend Update model from internal/updates/model.go.
 */
export interface Update {
  id: UUID;
  place_id: UUID;
  place_name: string;
  content_type: "today" | "forecast" | "alert";
  headline: string;
  body_text: string;
  call_to_action?: string | null;
  indicators?: Indicator[];
  generated_at: string;
  period_start: string;
  period_end: string;
}

/**
 * Matches the backend Indicator model.
 */
export interface Indicator {
  id: UUID;
  code: string;
  display_name: string;
  value: number;
  unit: string;
  trend?: "increasing" | "stable" | "decreasing" | null;
  relevance_score: number;
}

/**
 * Matches the backend ExploreIndicator model.
 */
export interface ExploreIndicator {
  id: UUID;
  place_id: UUID;
  place_name: string;
  code: string;
  display_name: string;
  category: string;
  value: number;
  unit: string;
  trend?: string | null;
  trend_confidence?: number | null;
  computed_at: string;
}

/**
 * Matches the backend ContextResponse.
 */
export interface ContextResponse {
  place: Place;
  update?: Update | null;
}

/**
 * Matches the backend ListResponse for updates.
 */
export interface UpdatesListResponse {
  updates: Update[];
  total: number;
  page: number;
  limit: number;
  has_next: boolean;
}

/**
 * Matches the backend ExploreResponse.
 */
export interface ExploreResponse {
  indicators: ExploreIndicator[];
  total: number;
  page: number;
  limit: number;
  has_next: boolean;
}

// --- Frontend adapter types (for backward compat with existing components) ---

/**
 * Adapted TodayResponse for the Today page component.
 * Maps from backend ContextResponse to what the UI expects.
 */
export interface TodayResponse {
  place: Place;
  status: { title: string; message: string; updated_at?: string; data_status: DataStatus };
  updates: PublicUpdate[];
  sections: TopicSection[];
  updated_at: string;
}

/**
 * PublicUpdate — adapted from backend Update for card components.
 */
export interface PublicUpdate {
  id: UUID;
  place_id: UUID;
  topic_key: string;
  type: "information" | "important" | "warning" | "emergency" | "community";
  priority: "normal" | "important" | "urgent";
  title: string;
  message: string;
  summary?: string;
  valid_from: string;
  valid_until?: string;
  updated_at: string;
  source_name?: string;
  source_url?: string;
  status: "draft" | "pending_review" | "approved" | "published" | "rejected" | "expired" | "cancelled";
  display_on_today: boolean;
  display_in_feed: boolean;
  locale: Locale;
  place_name?: string;
}

export type UpdateType = "information" | "important" | "warning" | "emergency" | "community";
export type UpdateStatus = "draft" | "pending_review" | "approved" | "published" | "rejected" | "expired" | "cancelled";
