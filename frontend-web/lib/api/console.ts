/**
 * Console API — institutional dashboard endpoints.
 * Client-side fetchers for projects, audit, users, delivery, and report stats.
 */

import { fetchAPI } from "./client";

// ── Types ─────────────────────────────────────────────────────────

export interface Project {
  id: string;
  org_id: string;
  name: string;
  description: string;
  status: string; // "active" | "archived"
  start_date: string | null;
  end_date: string | null;
  created_at: string;
  updated_at: string;
}

export interface ProjectsListResponse {
  projects: Project[];
  total: number;
  page: number;
  limit: number;
  has_next: boolean;
}

export interface AuditEntry {
  id: string;
  actor_id: string | null;
  actor_org_id: string | null;
  action: string;
  resource_type: string;
  resource_id: string | null;
  details: Record<string, unknown> | null;
  ip_address: string | null;
  created_at: string;
}

export interface AuditListResponse {
  entries: AuditEntry[];
  total: number;
  page: number;
  limit: number;
  has_next: boolean;
}

export interface UserSummary {
  id: string;
  display_name: string;
  email: string | null;
  phone_number: string | null;
  language: string;
  is_platform_admin: boolean;
}

export interface UsersListResponse {
  users: UserSummary[];
  page: number;
  limit: number;
}

export interface DeliveryStats {
  place_id: string;
  total_sent: number;
  total_failed: number;
  by_channel: Record<string, number>;
  period_start: string;
  period_end: string;
}

export interface ReportStats {
  org_id: string;
  total_submitted: number;
  total_verified: number;
  total_rejected: number;
  pending_review: number;
  period_start: string;
  period_end: string;
}

export interface DashboardResponse {
  delivery?: DeliveryStats | null;
  interaction?: {
    total_views: number;
    total_clicks: number;
    unique_users: number;
  } | null;
  reports?: ReportStats | null;
}

// ── Client-side fetchers ──────────────────────────────────────────

export async function getProjects(orgId: string, page = 1, limit = 20): Promise<ProjectsListResponse> {
  const params = new URLSearchParams({ org_id: orgId, page: String(page), limit: String(limit) });
  return fetchAPI<ProjectsListResponse>(`/v1/projects?${params}`);
}

export async function getAuditEntries(params: { page?: number; limit?: number; action?: string; resourceType?: string } = {}): Promise<AuditListResponse> {
  const searchParams = new URLSearchParams();
  if (params.page) searchParams.set("page", String(params.page));
  if (params.limit) searchParams.set("limit", String(params.limit));
  if (params.action) searchParams.set("action", params.action);
  if (params.resourceType) searchParams.set("resource_type", params.resourceType);
  const qs = searchParams.toString();
  return fetchAPI<AuditListResponse>(`/v1/audit${qs ? `?${qs}` : ""}`);
}

export async function getUsers(page = 1, limit = 20): Promise<UsersListResponse> {
  const params = new URLSearchParams({ page: String(page), limit: String(limit) });
  return fetchAPI<UsersListResponse>(`/v1/users?${params}`);
}

export async function getDeliveryStats(placeId: string): Promise<DeliveryStats> {
  const result = await fetchAPI<{ delivery_stats: DeliveryStats }>(`/v1/reporting/delivery?place_id=${placeId}`);
  return result.delivery_stats;
}

export async function getReportStats(orgId: string): Promise<ReportStats> {
  const result = await fetchAPI<{ report_stats: ReportStats }>(`/v1/reporting/reports?org_id=${orgId}`);
  return result.report_stats;
}

export async function getDashboard(params: { orgId?: string; placeId?: string } = {}): Promise<DashboardResponse> {
  const searchParams = new URLSearchParams();
  if (params.orgId) searchParams.set("org_id", params.orgId);
  if (params.placeId) searchParams.set("place_id", params.placeId);
  const qs = searchParams.toString();
  const result = await fetchAPI<{ dashboard: DashboardResponse }>(`/v1/reporting/dashboard${qs ? `?${qs}` : ""}`);
  return result.dashboard;
}
