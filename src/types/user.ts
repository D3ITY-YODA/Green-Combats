export interface User {
  id: string;
  email: string;
  name: string;
  avatar?: string;
  role: UserRole;
  organization?: Organization;
  preferences: UserPreferences;
  security: UserSecurity;
  createdAt: string;
  updatedAt: string;
  lastLoginAt?: string;
}

export type UserRole = 'admin' | 'analyst' | 'viewer' | 'contributor';

export interface Organization {
  id: string;
  name: string;
  slug: string;
  logo?: string;
  plan: 'free' | 'pro' | 'enterprise';
  settings: OrganizationSettings;
}

export interface OrganizationSettings {
  timezone: string;
  dateFormat: string;
  language: string;
  dataRetentionDays: number;
  allowPublicSharing: boolean;
  require2FA: boolean;
}

export interface UserPreferences {
  theme: 'light' | 'dark' | 'system';
  language: string;
  timezone: string;
  dateFormat: string;
  notifications: NotificationPreferences;
  dashboard: DashboardPreferences;
  maps: MapPreferences;
}

export interface NotificationPreferences {
  email: boolean;
  push: boolean;
  inApp: boolean;
  frequency: 'immediate' | 'hourly' | 'daily' | 'weekly';
  types: {
    alerts: boolean;
    reports: boolean;
    updates: boolean;
    mentions: boolean;
    assignments: boolean;
  };
}

export interface DashboardPreferences {
  defaultView: string;
  autoRefresh: boolean;
  refreshInterval: number;
  compactMode: boolean;
  showTooltips: boolean;
}

export interface MapPreferences {
  defaultBasemap: string;
  defaultZoom: number;
  defaultCenter: [number, number];
  showCoordinates: boolean;
  measureUnit: 'metric' | 'imperial';
}

export interface UserSecurity {
  twoFactorEnabled: boolean;
  twoFactorMethod?: 'authenticator' | 'sms' | 'email';
  backupCodes?: string[];
  sessions: UserSession[];
  passwordLastChanged: string;
  loginHistory: LoginEvent[];
}

export interface UserSession {
  id: string;
  device: string;
  browser: string;
  os: string;
  ip: string;
  location?: string;
  current: boolean;
  createdAt: string;
  lastActiveAt: string;
}

export interface LoginEvent {
  id: string;
  timestamp: string;
  ip: string;
  location?: string;
  device: string;
  browser: string;
  success: boolean;
  failureReason?: string;
}

export interface AuthState {
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  accessToken: string | null;
  refreshToken: string | null;
}

export interface LoginCredentials {
  email: string;
  password: string;
  rememberMe?: boolean;
  twoFactorCode?: string;
}

export interface RegisterData {
  name: string;
  email: string;
  password: string;
  confirmPassword: string;
  organizationName?: string;
  acceptTerms: boolean;
}

export interface PasswordResetRequest {
  email: string;
}

export interface PasswordResetConfirm {
  token: string;
  password: string;
  confirmPassword: string;
}

export interface TwoFactorSetup {
  secret: string;
  qrCode: string;
  backupCodes: string[];
}

export interface TwoFactorVerify {
  code: string;
  method: 'authenticator' | 'sms' | 'email';
}