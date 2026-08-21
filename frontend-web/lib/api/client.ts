// lib/api/client.ts

import type { ApiResponse, ApiError } from "@/types/api";

const API_BASE_URL =
  process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080";

/**
 * Build the full URL for a backend path.
 * Frontend callers use `/api/v1/...` convention; the Go backend mounts
 * routes at `/v1/...`, so we strip the leading `/api` prefix.
 */
function buildURL(path: string): string {
  const backendPath = path.startsWith("/api") ? path.slice(4) : path;
  return `${API_BASE_URL}${backendPath}`;
}

/**
 * Client-side fetch wrapper for the Go backend.
 *
 * - Automatically prepends the backend base URL.
 * - Sends cookies so the auth middleware can read the session/access token.
 * - Returns the parsed `ApiResponse<T>` envelope.
 * - Throws the parsed `ApiError` on non-OK responses.
 */
export async function apiFetch<T>(
  path: string,
  init?: RequestInit,
): Promise<ApiResponse<T>> {
  const url = buildURL(path);

  const res = await fetch(url, {
    ...init,
    credentials: "include", // send cookies (gc_session, etc.)
    headers: {
      "Content-Type": "application/json",
      Accept: "application/json",
      ...init?.headers,
    },
  });

  const body = await res.json();

  if (!res.ok) {
    const apiError: ApiError = body as ApiError;
    throw apiError;
  }

  return body as ApiResponse<T>;
}
