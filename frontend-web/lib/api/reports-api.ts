// lib/api/reports-api.ts

import { fetchAPI } from "./client";

export interface ReportPayload {
  observationType: string;
  details?: string;
  location: string;
}

export interface ReportResponse {
  success: boolean;
  message: string;
}

export async function submitReport(
  payload: ReportPayload,
): Promise<ReportResponse> {
  // Local dev without a backend: simulate a successful submission.
  if (!process.env.NEXT_PUBLIC_API_URL) {
    return {
      success: true,
      message: "Observation recorded successfully.",
    };
  }

  try {
    await fetchAPI("/v1/reports", {
      method: "POST",
      body: {
        observation_type: payload.observationType,
        details: payload.details,
        location: payload.location,
      },
    });

    return {
      success: true,
      message: "Observation recorded successfully.",
    };
  } catch (err) {
    return {
      success: false,
      message: err instanceof Error ? err.message : "Unable to submit observation.",
    };
  }
}
