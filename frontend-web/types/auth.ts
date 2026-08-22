// types/auth.ts — Matches backend auth handler responses

import { UUID } from "./common";

/**
 * Payload for POST /v1/auth/register
 */
export interface RegisterInput {
  phone_number?: string;
  email?: string;
  display_name: string;
  password: string;
  language?: string;
}

/**
 * Payload for POST /v1/auth/login
 */
export interface LoginInput {
  identifier: string; // email or phone
  password: string;
}

/**
 * Token pair returned by login/register/refresh.
 */
export interface TokenPair {
  access_token: string;
  refresh_token: string;
  access_expires_at: string;
  refresh_expires_at: string;
}

/**
 * Response from POST /v1/auth/register
 */
export interface RegisterResponse {
  user: SessionUser;
  tokens: TokenPair;
}

/**
 * Response from POST /v1/auth/login and POST /v1/auth/refresh
 */
export interface LoginResponse {
  tokens: TokenPair;
}

/**
 * Current user profile from GET /v1/auth/me
 */
export interface SessionUser {
  id: UUID;
  display_name: string;
  email?: string | null;
  phone_number?: string | null;
  language: string;
  is_platform_admin: boolean;
}

/**
 * VerifyCodeInput (kept for future OTP flow)
 */
export interface VerifyCodeInput {
  phone: string;
  code: string;
}

/**
 * ForgotPasswordInput (kept for future password reset)
 */
export interface ForgotPasswordInput {
  phone: string;
}

/**
 * ChangePasswordInput (kept for future password change)
 */
export interface ChangePasswordInput {
  current_password: string;
  new_password: string;
}
