// types/updates.ts

import type { UUID, Locale } from "./common"; // Adjust path to match your common types
import type { TopicKey } from "./context";    // Adjust path to where TopicKey is defined (Section 5.6)

/**
 * The category or severity of the update.
 * Used to determine visual treatment (e.g., color coding) in the UI.
 */
export type UpdateType =
  | "information"
  | "important"
  | "warning"
  | "emergency"
  | "community";

/**
 * The lifecycle status of an update within the institutional console.
 * Public users should only ever see "published" updates.
 */
export type UpdateStatus =
  | "draft"
  | "pending_review"
  | "approved"
  | "published"
  | "rejected"
  | "expired"
  | "cancelled";

/**
 * The complete Update entity returned by the API.
 * This represents a single piece of environmental information or alert.
 */
export interface PublicUpdate {
  id: UUID;
  place_id: UUID;
  topic_key: TopicKey;
  type: UpdateType;
  priority: "normal" | "important" | "urgent";
  title: string;
  message: string;
  summary?: string;
  valid_from: string; // ISO 8601 datetime string
  valid_until?: string; // ISO 8601 datetime string
  updated_at: string; // ISO 8601 datetime string
  source_name?: string;
  source_url?: string;
  status: UpdateStatus;
  display_on_today: boolean;
  display_in_feed: boolean;
  locale: Locale;
  place_name?: string; // Denormalized for easier display
}
