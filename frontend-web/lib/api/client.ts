/**
 * Browser-side API client.
 *
 * Features:
 *  - Prepends NEXT_PUBLIC_API_URL to all paths
 *  - Attaches Bearer token from cookie
 *  - Unwraps the { data, meta } envelope
 *  - Auto-refreshes on 401, then retries once
 *  - Redirects to /sign-in on refresh failure
 */

import type { ApiResponse } from "@/types/api";

const API_BASE = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

// ── Cookie helpers ────────────────────────────────────────────────

function getCookie(name: string): string | null {
  if (typeof document === "undefined") return null;
  const match = document.cookie.match(new RegExp(`(?:^|; )${name}=([^;]*)`));
  return match ? decodeURIComponent(match[1]) : null;
}

export function setAccessToken(token: string, expiresAt: string) {
  const expires = new Date(expiresAt).toUTCString();
  document.cookie = `gc_session=${encodeURIComponent(token)}; Path=/; SameSite=Lax; Expires=${expires}`;
}

export function setRefreshToken(token: string, expiresAt: string) {
  const expires = new Date(expiresAt).toUTCString();
  document.cookie = `gc_refresh=${encodeURIComponent(token)}; Path=/; SameSite=Lax; Expires=${expires}`;
}

export function getAccessToken(): string | null {
  return getCookie("gc_session");
}

export function getRefreshToken(): string | null {
  return getCookie("gc_refresh");
}

export function clearTokens() {
  document.cookie = "gc_session=; Path=/; Max-Age=0";
  document.cookie = "gc_refresh=; Path=/; Max-Age=0";
}

// ── Core fetch ────────────────────────────────────────────────────

interface FetchOptions extends Omit<RequestInit, "method" | "body"> {
  method?: string;
  body?: unknown;
  /** Skip token attachment (for public endpoints) */
  noAuth?: boolean;
}

/**
 * Make an authenticated API request and unwrap the envelope.
 *
 * @param path  API path, e.g. "/v1/auth/me"
 * @param options  fetch options (body is auto-JSON-serialized)
 * @returns Unwrapped `data` from the response envelope
 */
export async function fetchAPI<T>(path: string, options: FetchOptions = {}): Promise<T> {
  const { noAuth, body, headers: extraHeaders, ...rest } = options;

  const headers: Record<string, string> = {
    "Content-Type": "application/json",
    ...((extraHeaders as Record<string, string>) || {}),
  };

  if (!noAuth) {
    const token = getAccessToken();
    if (token) {
      headers["Authorization"] = `Bearer ${token}`;
    }
  }

  const res = await fetch(`${API_BASE}${path}`, {
    ...rest,
    headers,
    body: body != null ? JSON.stringify(body) : undefined,
  });

  // 204 No Content (e.g. logout, save/unsave)
  if (res.status === 204) {
    return undefined as T;
  }

  // Handle errors
  if (!res.ok) {
    const errorBody = await res.json().catch(() => null);

    // 401 — try refresh, then retry
    if (res.status === 401 && !noAuth) {
      const refreshed = await tryRefresh();
      if (refreshed) {
        return fetchAPI<T>(path, options);
      }
      clearTokens();
      if (typeof window !== "undefined") {
        window.location.href = "/sign-in";
      }
      throw new Error("Session expired");
    }

    const message = errorBody?.error?.message || `API error ${res.status}`;
    throw new Error(message);
  }

  const json: ApiResponse<T> = await res.json();
  return json.data;
}

// ── Token refresh ─────────────────────────────────────────────────

async function tryRefresh(): Promise<boolean> {
  const refreshToken = getRefreshToken();
  if (!refreshToken) return false;

  try {
    const res = await fetch(`${API_BASE}/v1/auth/refresh`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ refresh_token: refreshToken }),
    });

    if (!res.ok) return false;

    const json: ApiResponse<{ tokens: { access_token: string; refresh_token: string; access_expires_at: string; refresh_expires_at: string } }> =
      await res.json();

    const { tokens } = json.data;
    setAccessToken(tokens.access_token, tokens.access_expires_at);
    setRefreshToken(tokens.refresh_token, tokens.refresh_expires_at);
    return true;
  } catch {
    return false;
  }
}
