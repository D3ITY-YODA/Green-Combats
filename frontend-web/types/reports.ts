// types/reports.ts

import type { UUID } from "./common"; // Adjust path to match your common types

/**
 * The specific category of environmental observation being reported.
 * Maps directly to the options in the Report Form UI (Screen 16).
 */
export type ReportType =
  | "water_change"
  | "flooding_visible"
  | "unusually_dry"
  | "vegetation_stress"
  | "heat_impact"
  | "infrastructure_change"
  | "incorrect_information"
  | "other";

/**
 * The lifecycle status of a submitted report.
 */
export type ReportStatus =
  | "pending_sync"
  | "submitted"
  | "under_review"
  | "verified"
  | "rejected"
  | "resolved";

/**
 * Geographic coordinates and accuracy metadata for a report.
 */
export interface ReportLocation {
  latitude: number;
  longitude: number;
  accuracy_meters?: number;
}

/**
 * Represents a media attachment (photo/video) associated with a report.
 */
export interface ReportMedia {
  id: UUID;
  report_id: UUID;
  content_type: string; // e.g., "image/jpeg", "video/mp4"
  url?: string;
  status: "uploading" | "uploaded" | "failed";
}

/**
 * The complete Community Report entity returned by the API.
 */
export interface CommunityReport {
  id: UUID;
  reporter_id?: UUID;
  place_id?: UUID;
  type: ReportType;
  location: ReportLocation;
  description?: string;
  observed_at: string; // ISO 8601 datetime string
  submitted_at: string; // ISO 8601 datetime string
  status: ReportStatus;
  visibility: "private" | "institutional" | "public";
  media?: ReportMedia[];
}
