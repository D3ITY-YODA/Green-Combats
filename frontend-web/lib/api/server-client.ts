/**
 * Server-side API client for Server Components and Route Handlers.
 *
 * Reads the access token from cookies (forwarded by Next.js).
 * Does NOT auto-refresh — the browser client handles that.
 */

import { cookies } from "next/headers";
import type { ApiResponse } from "@/types/api";

const API_BASE = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

interface ServerFetchOptions extends Omit<RequestInit, "method" | "body"> {
  method?: string;
  body?: unknown;
  token?: string; // Override token (e.g. from middleware context)
}

/**
 * Fetch from the Go backend on the server side.
 *
 * @param path  API path, e.g. "/v1/context/today"
 * @param options  fetch options
 * @returns Unwrapped `data` from the response envelope
 */
export async function fetchServer<T>(path: string, options: ServerFetchOptions = {}): Promise<T> {
  const { body, token, headers: extraHeaders, ...rest } = options;

  // Get token from cookies or use override
  let authToken: string | undefined = token || undefined;
  if (!authToken) {
    const cookieStore = await cookies();
    authToken = cookieStore.get("gc_session")?.value || undefined;
  }

  const headers: Record<string, string> = {
    "Content-Type": "application/json",
    ...((extraHeaders as Record<string, string>) || {}),
  };

  if (authToken) {
    headers["Authorization"] = `Bearer ${authToken}`;
  }

  const res = await fetch(`${API_BASE}${path}`, {
    ...rest,
    headers,
    body: body != null ? JSON.stringify(body) : undefined,
    cache: "no-store",
  });

  if (res.status === 204) {
    return undefined as T;
  }

  if (!res.ok) {
    const errorBody = await res.json().catch(() => null);
    const message = errorBody?.error?.message || `API error ${res.status}`;
    throw new Error(message);
  }

  const json: ApiResponse<T> = await res.json();
  return json.data;
}
