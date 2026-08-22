// types/users.ts

import type { UUID, ISODateTime, Locale } from "./common";

/**
 * User status in the system
 */
export type UserStatus = "active" | "inactive" | "pending_verification" | "suspended" | "deleted";

/**
 * User roles
 */
export type UserRole = "user" | "admin" | "moderator" | "contributor";

/**
 * User preferences
 */
export interface UserPreferences {
  locale: Locale;
  notifications_enabled: boolean;
  email_notifications: boolean;
  push_notifications: boolean;
  sms_notifications: boolean;
  theme: "light" | "dark" | "system";
}

/**
 * User profile from the API
 */
export interface User {
  id: UUID;
  phone: string;
  full_name: string;
  email?: string;
  status: UserStatus;
  role: UserRole;
  preferences: UserPreferences;
  created_at: ISODateTime;
  updated_at: ISODateTime;
  last_login_at?: ISODateTime;
  avatar_url?: string;
}

/**
 * User profile for public display (limited fields)
 */
export interface PublicUser {
  id: UUID;
  full_name: string;
  avatar_url?: string;
}

/**
 * User session data
 */
export interface UserSession {
  user: User;
  session_id: string;
  expires_at: ISODateTime;
  csrf_token?: string;
}