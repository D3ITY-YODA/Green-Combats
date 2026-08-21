// lib/api/server-client.ts

import type { ApiResponse, ApiError } from "@/types/api";

const API_BASE_URL =
  process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080";

/**
 * Build the full URL for a backend path.
 * Strips the `/api` prefix that the frontend convention uses.
 */
function buildURL(path: string): string {
  const backendPath = path.startsWith("/api") ? path.slice(4) : path;
  return `${API_BASE_URL}${backendPath}`;
}

/**
 * Server-side fetch wrapper for the Go backend.
 *
 * Use this in Next.js Server Components and Route Handlers.
 * It forwards the `cookie` header from the incoming request so the
 * backend auth middleware can verify the session.
 *
 * @param path - Backend path (e.g. `/api/v1/me/places`)
 * @param init - Optional fetch options (headers, method, body, etc.)
 * @returns Parsed `ApiResponse<T>` envelope
 * @throws `ApiError` on non-OK responses
 */
export async function serverApiFetch<T>(
  path: string,
  init?: RequestInit & { cookies?: string },
): Promise<ApiResponse<T>> {
  const url = buildURL(path);

  // Forward cookies from the incoming request when provided
  const headers: Record<string, string> = {
    "Content-Type": "application/json",
    Accept: "application/json",
  };

  if (init?.cookies) {
    headers.Cookie = init.cookies;
  }

  const res = await fetch(url, {
    ...init,
    headers: {
      ...headers,
      ...init?.headers,
    },
    cache: "no-store",
  });

  const body = await res.json();

  if (!res.ok) {
    const apiError: ApiError = body as ApiError;
    throw apiError;
  }

  return body as ApiResponse<T>;
}
