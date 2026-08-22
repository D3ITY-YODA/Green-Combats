/**
 * Server-only reporting fetchers.
 * These MUST only be imported by Server Components.
 */

import { fetchServer } from "./server-client";
import type { DashboardData } from "./reporting";

export async function getDashboard(token?: string): Promise<DashboardData> {
  return fetchServer<DashboardData>("/v1/reporting/dashboard", { token });
}
