// types/users.ts

/**
 * The authenticated user object returned by GET /v1/auth/me.
 * Matches the sessionUserJSON shape from the Go backend.
 */
export interface User {
  id: string;
  display_name: string;
  email: string | null;
  phone_number: string | null;
  language: string;
  is_platform_admin: boolean;
}
