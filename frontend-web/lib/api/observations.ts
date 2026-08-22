/**
 * Observations API — community reports.
 */

import type { Observation, SubmitObservationInput } from "@/types/reports";
import { fetchAPI } from "./client";

/**
 * Submit a community observation/report.
 */
export async function submitObservation(input: SubmitObservationInput): Promise<Observation> {
  return fetchAPI<Observation>("/v1/observations", {
    method: "POST",
    body: input,
  });
}

/**
 * List the current user's observations.
 */
export async function getMyObservations(limit?: number): Promise<Observation[]> {
  const params = limit ? `?limit=${limit}` : "";
  const result = await fetchAPI<{ observations: Observation[] }>(`/v1/observations${params}`);
  return result.observations;
}

/**
 * Get a single observation by ID.
 */
export async function getObservation(id: string): Promise<Observation> {
  return fetchAPI<Observation>(`/v1/observations/${id}`);
}

/**
 * List pending observations (institutional view).
 */
export async function getPendingObservations(limit?: number): Promise<Observation[]> {
  const params = limit ? `?limit=${limit}` : "";
  const result = await fetchAPI<{ observations: Observation[] }>(`/v1/observations/institutional/pending${params}`);
  return result.observations;
}

/**
 * Verify/reject an observation.
 */
export async function verifyObservation(id: string, status: string, notes?: string): Promise<Observation> {
  return fetchAPI<Observation>(`/v1/observations/institutional/${id}/verify`, {
    method: "PATCH",
    body: { status, notes },
  });
}
