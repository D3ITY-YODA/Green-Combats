// types/auth.ts

/**
 * Payload sent to the backend to authenticate an existing user.
 * Matches the Sign In UI (Screen 6).
 */
export interface SignInInput {
  phone: string;
  password: string;
}

/**
 * Payload sent to the backend to create a new user account.
 * Matches the Sign Up UI (Screen 5).
 */
export interface SignUpInput {
  full_name: string;
  phone: string;
  password: string;
}

/**
 * Payload sent to verify a phone number or email via OTP/Code.
 * Matches the Verify UI flow.
 */
export interface VerifyCodeInput {
  phone: string;
  code: string;
}

/**
 * Payload sent to initiate a password reset flow.
 */
export interface ForgotPasswordInput {
  phone: string;
}

/**
 * Represents the session data returned by the backend upon successful login.
 * Note: The blueprint primarily uses HttpOnly cookies (`gc_session`) for auth,
 * but this interface covers any additional session metadata the Go backend might return.
 */
export interface AuthSession {
  session_id?: string;
  csrf_token?: string;
  expires_at?: string;
}

/**
 * Payload for updating a user's password.
 */
export interface ChangePasswordInput {
  current_password: string;
  new_password: string;
}
