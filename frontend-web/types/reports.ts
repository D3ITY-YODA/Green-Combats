import { UUID } from "./common";

/**
 * Matches the backend Observation model from /v1/observations.
 */
export interface Observation {
  id: UUID;
  reporter_id: UUID;
  place_id?: UUID | null;
  lat?: number | null;
  lon?: number | null;
  category: ObservationCategory;
  description: string;
  photo_object_key?: string | null;
  status: ObservationStatus;
  verified_by?: UUID | null;
  verified_at?: string | null;
  created_at: string;
  updated_at: string;
}

export type ObservationCategory =
  | "flood"
  | "drought"
  | "water_quality"
  | "crop_damage"
  | "air_quality"
  | "other";

export type ObservationStatus = "pending" | "verified" | "rejected" | "flagged";

/**
 * Input for POST /v1/observations
 */
export interface SubmitObservationInput {
  place_id?: UUID;
  lat?: number;
  lon?: number;
  category: ObservationCategory;
  description: string;
  photo_object_key?: string;
}

/**
 * Legacy report types (kept for report form compatibility)
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

export interface ReportLocation {
  latitude: number;
  longitude: number;
  accuracy_meters?: number;
}

export interface SubmitReportInput {
  type: ReportType;
  place_id: string;
  location: ReportLocation;
  observed_at: string;
  description?: string;
}
