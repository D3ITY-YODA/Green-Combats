// lib/api/reports-api.ts

import { apiFetch } from "@/lib/api/client";

export interface ReportPayload {
  type: string;
  place_id: string;
  location: {
    latitude: number;
    longitude: number;
    accuracy_meters?: number;
  };
  description?: string;
  observed_at: string;
}

export interface ReportResponse {
  success: boolean;
  message: string;
}

interface BackendObservation {
  id: string;
  reporter_id: string;
  category: string;
  description: string;
  status: string;
  created_at: string;
  updated_at: string;
  place_id?: string;
  lat?: number;
  lon?: number;
}

/** Map frontend report type to backend observation category. */
function mapCategory(type: string): string {
  const mapping: Record<string, string> = {
    water_change: "water_change",
    flooding_visible: "flooding",
    unusually_dry: "drought",
    vegetation_stress: "vegetation_stress",
    heat_impact: "heat_impact",
    infrastructure_change: "infrastructure",
    incorrect_information: "data_quality",
    other: "other",
  };
  return mapping[type] ?? type;
}

/**
 * Submit a community observation report to the backend.
 */
export async function submitReport(payload: ReportPayload): Promise<ReportResponse> {
  const body = {
    place_id: payload.place_id,
    lat: payload.location.latitude,
    lon: payload.location.longitude,
    category: mapCategory(payload.type),
    description: payload.description ?? "",
  };

  await apiFetch<BackendObservation>("/api/v1/observations", {
    method: "POST",
    body: JSON.stringify(body),
  });

  return {
    success: true,
    message: "Observation recorded successfully.",
  };
}
