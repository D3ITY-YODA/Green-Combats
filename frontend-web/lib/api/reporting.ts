/**
 * Client-safe reporting module — types and client-side fetchers only.
 * Server-side fetchers are in ./reporting-server.ts
 */

import { fetchAPI } from "./client";

export async function getDashboardClient(): Promise<DashboardData> {
  return fetchAPI<DashboardData>("/v1/reporting/dashboard");
}

export async function getDeliveryStats(): Promise<DeliveryStats> {
  return fetchAPI<DeliveryStats>("/v1/reporting/delivery");
}

export async function getReportStats(): Promise<ReportStats> {
  return fetchAPI<ReportStats>("/v1/reporting/reports");
}

export interface DashboardData {
  important_updates?: number;
  community_reports?: number;
  information_delayed?: number;
  pending_review?: number;
}

export interface DeliveryStats {
  total_sent?: number;
  total_delivered?: number;
  delivery_rate?: number;
}

export interface ReportStats {
  total_reports?: number;
  verified?: number;
  pending?: number;
  rejected?: number;
}
