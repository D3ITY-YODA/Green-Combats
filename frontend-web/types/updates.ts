import { UUID, Locale, Place, DataStatus, TopicSection } from "./common";
export type UpdateType = "information" | "important" | "warning" | "emergency" | "community";
export type UpdateStatus = "draft" | "pending_review" | "approved" | "published" | "rejected" | "expired" | "cancelled";
export interface PublicUpdate {
  id: UUID;
  place_id: UUID;
  topic_key: string;
  type: UpdateType;
  priority: "normal" | "important" | "urgent";
  title: string;
  message: string;
  summary?: string;
  valid_from: string;
  valid_until?: string;
  updated_at: string;
  source_name?: string;
  source_url?: string;
  status: UpdateStatus;
  display_on_today: boolean;
  display_in_feed: boolean;
  locale: Locale;
  place_name?: string;
}
export interface TodayResponse {
  place: Place;
  status: { title: string; message: string; updated_at?: string; data_status: DataStatus; };
  updates: PublicUpdate[];
  sections: TopicSection[];
  updated_at: string;
}
export interface ExploreResponse {
  place: Place;
  sections: TopicSection[];
  updated_at: string;
}
