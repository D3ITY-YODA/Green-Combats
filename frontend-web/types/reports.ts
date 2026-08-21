import { UUID } from "./common";

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
