/**
 * Auth API — login, register, logout, getCurrentUser.
 */

import type { LoginInput, RegisterInput, LoginResponse, RegisterResponse, SessionUser } from "@/types/auth";
import { fetchAPI, setAccessToken, setRefreshToken, clearTokens } from "./client";

// ── Login ─────────────────────────────────────────────────────────

export async function login(input: LoginInput): Promise<SessionUser> {
  const { tokens } = await fetchAPI<LoginResponse>("/v1/auth/login", {
    method: "POST",
    body: input,
    noAuth: true,
  });

  setAccessToken(tokens.access_token, tokens.access_expires_at);
  setRefreshToken(tokens.refresh_token, tokens.refresh_expires_at);

  // Fetch user profile after login
  return getMe();
}

// ── Register ──────────────────────────────────────────────────────

export async function register(input: RegisterInput): Promise<SessionUser> {
  const { user, tokens } = await fetchAPI<RegisterResponse>("/v1/auth/register", {
    method: "POST",
    body: input,
    noAuth: true,
  });

  setAccessToken(tokens.access_token, tokens.access_expires_at);
  setRefreshToken(tokens.refresh_token, tokens.refresh_expires_at);

  return user;
}

// ── Logout ────────────────────────────────────────────────────────

export async function logout(): Promise<void> {
  const refreshToken = getRefreshToken();
  try {
    await fetchAPI("/v1/auth/logout", {
      method: "POST",
      body: { refresh_token: refreshToken },
    });
  } finally {
    clearTokens();
  }
}

// ── Get current user ─────────────────────────────────────────────

export async function getMe(): Promise<SessionUser> {
  return fetchAPI<SessionUser>("/v1/auth/me");
}

// ── Helpers ───────────────────────────────────────────────────────

function getRefreshToken(): string | null {
  if (typeof document === "undefined") return null;
  const match = document.cookie.match(/(?:^|; )gc_refresh=([^;]*)/);
  return match ? decodeURIComponent(match[1]) : null;
}
