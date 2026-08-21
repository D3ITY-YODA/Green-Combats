import { create } from 'zustand';
import { persist, devtools } from 'zustand/middleware';
import type { User, AuthState, LoginCredentials, RegisterData, UserPreferences, UserSecurity, NavSection } from '../types';

interface AuthStore extends AuthState {
  login: (credentials: LoginCredentials) => Promise<void>;
  register: (data: RegisterData) => Promise<void>;
  logout: () => void;
  refreshAccessToken: () => Promise<void>;
  setUser: (user: User) => void;
  updateUser: (updates: Partial<User>) => void;
}

const mockUser: User = {
  id: 'usr_001',
  email: 'alex.morgan@greencombats.io',
  name: 'Alex Morgan',
  avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=alex',
  role: 'admin',
  organization: {
    id: 'org_001',
    name: 'Green Combats Global',
    slug: 'green-combats-global',
    plan: 'enterprise',
    settings: {
      timezone: 'UTC',
      dateFormat: 'YYYY-MM-DD',
      language: 'en',
      dataRetentionDays: 2555,
      allowPublicSharing: true,
      require2FA: true,
    },
  },
  preferences: {
    theme: 'system',
    language: 'en',
    timezone: 'UTC',
    dateFormat: 'YYYY-MM-DD',
    notifications: {
      email: true,
      push: true,
      inApp: true,
      frequency: 'immediate',
      types: {
        alerts: true,
        reports: true,
        updates: true,
        mentions: true,
        assignments: true,
      },
    },
    dashboard: {
      defaultView: 'overview',
      autoRefresh: true,
      refreshInterval: 300000,
      compactMode: false,
      showTooltips: true,
    },
    maps: {
      defaultBasemap: 'satellite',
      defaultZoom: 3,
      defaultCenter: [20, 0],
      showCoordinates: true,
      measureUnit: 'metric',
    },
  },
  security: {
    twoFactorEnabled: true,
    twoFactorMethod: 'authenticator',
    backupCodes: ['BC-1234-5678', 'BC-8765-4321', 'BC-1111-2222'],
    sessions: [
      {
        id: 'ses_001',
        device: 'MacBook Pro',
        browser: 'Chrome 120.0',
        os: 'macOS 14.2',
        ip: '192.168.1.100',
        location: 'San Francisco, CA, US',
        current: true,
        createdAt: '2024-01-15T08:30:00Z',
        lastActiveAt: '2024-01-20T14:22:00Z',
      },
      {
        id: 'ses_002',
        device: 'iPhone 15 Pro',
        browser: 'Safari 17.2',
        os: 'iOS 17.2',
        ip: '10.0.0.50',
        location: 'New York, NY, US',
        current: false,
        createdAt: '2024-01-18T12:15:00Z',
        lastActiveAt: '2024-01-19T22:10:00Z',
      },
    ],
    passwordLastChanged: '2023-11-15T10:00:00Z',
    loginHistory: [
      {
        id: 'log_001',
        timestamp: '2024-01-20T14:22:00Z',
        ip: '192.168.1.100',
        location: 'San Francisco, CA, US',
        device: 'MacBook Pro',
        browser: 'Chrome 120.0',
        success: true,
      },
      {
        id: 'log_002',
        timestamp: '2024-01-20T08:15:00Z',
        ip: '192.168.1.100',
        location: 'San Francisco, CA, US',
        device: 'MacBook Pro',
        browser: 'Chrome 120.0',
        success: true,
      },
      {
        id: 'log_003',
        timestamp: '2024-01-19T22:10:00Z',
        ip: '10.0.0.50',
        location: 'New York, NY, US',
        device: 'iPhone 15 Pro',
        browser: 'Safari 17.2',
        success: true,
      },
    ],
  },
  createdAt: '2023-06-15T10:00:00Z',
  updatedAt: '2024-01-15T08:30:00Z',
  lastLoginAt: '2024-01-20T14:22:00Z',
};

export const useAuthStore = create<AuthStore>()(
  devtools(
    persist(
      (set, get) => ({
        user: null,
        isAuthenticated: false,
        isLoading: false,
        accessToken: null,
        refreshToken: null,

        login: async (credentials: LoginCredentials) => {
          set({ isLoading: true });
          await new Promise((resolve) => setTimeout(resolve, 1000));
          
          if (credentials.email === 'demo@greencombats.io' && credentials.password === 'demo123') {
            const tokens = {
              accessToken: 'mock_access_token_' + Date.now(),
              refreshToken: 'mock_refresh_token_' + Date.now(),
            };
            set({
              user: mockUser,
              isAuthenticated: true,
              accessToken: tokens.accessToken,
              refreshToken: tokens.refreshToken,
              isLoading: false,
            });
          } else {
            set({ isLoading: false });
            throw new Error('Invalid credentials');
          }
        },

        register: async (data: RegisterData) => {
          set({ isLoading: true });
          await new Promise((resolve) => setTimeout(resolve, 1500));
          
          const newUser: User = {
            ...mockUser,
            id: 'usr_' + Date.now(),
            email: data.email,
            name: data.name,
            organization: data.organizationName ? {
              id: 'org_' + Date.now(),
              name: data.organizationName,
              slug: data.organizationName.toLowerCase().replace(/\s+/g, '-'),
              plan: 'free',
              settings: {
                timezone: 'UTC',
                dateFormat: 'YYYY-MM-DD',
                language: 'en',
                dataRetentionDays: 365,
                allowPublicSharing: false,
                require2FA: false,
              },
            } : undefined,
            createdAt: new Date().toISOString(),
            updatedAt: new Date().toISOString(),
          };
          
          const tokens = {
            accessToken: 'mock_access_token_' + Date.now(),
            refreshToken: 'mock_refresh_token_' + Date.now(),
          };
          
          set({
            user: newUser,
            isAuthenticated: true,
            accessToken: tokens.accessToken,
            refreshToken: tokens.refreshToken,
            isLoading: false,
          });
        },

        logout: () => {
          set({
            user: null,
            isAuthenticated: false,
            accessToken: null,
            refreshToken: null,
          });
        },

        refreshAccessToken: async () => {
          const { refreshToken } = get();
          if (!refreshToken) throw new Error('No refresh token');
          
          await new Promise((resolve) => setTimeout(resolve, 500));
          set({ accessToken: 'mock_access_token_' + Date.now() });
        },

        setUser: (user: User) => set({ user, isAuthenticated: true }),

        updateUser: (updates: Partial<User>) => {
          const { user } = get();
          if (user) {
            set({ user: { ...user, ...updates, updatedAt: new Date().toISOString() } });
          }
        },
      }),
      {
        name: 'auth-storage',
        partialize: (state) => ({
          user: state.user,
          isAuthenticated: state.isAuthenticated,
          accessToken: state.accessToken,
          refreshToken: state.refreshToken,
        }),
      }
    ),
    { name: 'AuthStore' }
  )
);

interface UIStore {
  sidebarOpen: boolean;
  sidebarCollapsed: boolean;
  mobileMenuOpen: boolean;
  theme: 'light' | 'dark' | 'system';
  toasts: Toast[];
  modals: Record<string, boolean>;
  toggleSidebar: () => void;
  setSidebarOpen: (open: boolean) => void;
  setSidebarCollapsed: (collapsed: boolean) => void;
  toggleMobileMenu: () => void;
  setTheme: (theme: 'light' | 'dark' | 'system') => void;
  addToast: (toast: Omit<Toast, 'id'>) => string;
  removeToast: (id: string) => void;
  openModal: (id: string) => void;
  closeModal: (id: string) => void;
}

export interface Toast {
  id: string;
  type: 'success' | 'error' | 'warning' | 'info';
  title: string;
  message?: string;
  duration?: number;
  action?: { label: string; onClick: () => void };
}

export const useUIStore = create<UIStore>()(
  devtools(
    persist(
      (set, get) => ({
        sidebarOpen: true,
        sidebarCollapsed: false,
        mobileMenuOpen: false,
        theme: 'system',
        toasts: [],
        modals: {},

        toggleSidebar: () => set((state) => ({ sidebarOpen: !state.sidebarOpen })),
        setSidebarOpen: (open: boolean) => set({ sidebarOpen: open }),
        setSidebarCollapsed: (collapsed: boolean) => set({ sidebarCollapsed: collapsed }),
        toggleMobileMenu: () => set((state) => ({ mobileMenuOpen: !state.mobileMenuOpen })),
        setTheme: (theme: 'light' | 'dark' | 'system') => set({ theme }),

        addToast: (toast) => {
          const id = 'toast_' + Date.now() + '_' + Math.random().toString(36).substr(2, 9);
          const newToast = { ...toast, id, duration: toast.duration ?? 5000 };
          set((state) => ({ toasts: [...state.toasts, newToast] }));
          
          if (newToast.duration && newToast.duration > 0) {
            setTimeout(() => get().removeToast(id), newToast.duration);
          }
          return id;
        },

        removeToast: (id: string) => set((state) => ({ toasts: state.toasts.filter((t) => t.id !== id) })),

        openModal: (id: string) => set((state) => ({ modals: { ...state.modals, [id]: true } })),
        closeModal: (id: string) => set((state) => ({ modals: { ...state.modals, [id]: false } })),
      }),
      {
        name: 'ui-storage',
        partialize: (state) => ({
          sidebarCollapsed: state.sidebarCollapsed,
          theme: state.theme,
        }),
      }
    ),
    { name: 'UIStore' }
  )
);

interface NavigationStore {
  currentPath: string;
  breadcrumbs: BreadcrumbItem[];
  navSections: NavSection[];
  setCurrentPath: (path: string) => void;
  setBreadcrumbs: (breadcrumbs: BreadcrumbItem[]) => void;
  getNavItemsForRole: (role: string) => NavItem[];
}

export interface BreadcrumbItem {
  label: string;
  href?: string;
  icon?: React.ComponentType<{ className?: string }>;
}

export const useNavigationStore = create<NavigationStore>()(
  devtools(
    (set, get) => ({
      currentPath: '/',
      breadcrumbs: [],
      navSections: [],

      setCurrentPath: (path: string) => set({ currentPath: path }),
      setBreadcrumbs: (breadcrumbs: BreadcrumbItem[]) => set({ breadcrumbs }),

      getNavItemsForRole: (role: string) => {
        const { navSections } = get();
        return navSections
          .flatMap((section) => section.items)
          .filter((item) => !item.roles || item.roles.includes(role as any));
      },
    }),
    { name: 'NavigationStore' }
  )
);

interface DataStore {
  emissionsData: Record<string, unknown> | null;
  carbonFootprint: Record<string, unknown> | null;
  energyData: Record<string, unknown> | null;
  waterData: Record<string, unknown> | null;
  biodiversityData: Record<string, unknown> | null;
  airQualityData: Record<string, unknown> | null;
  projects: Record<string, unknown>[];
  initiatives: Record<string, unknown>[];
  scenarios: Record<string, unknown>[];
  regulations: Record<string, unknown>[];
  certifications: Record<string, unknown>[];
  loading: Record<string, boolean>;
  error: Record<string, string | null>;
  setData: <T>(key: string, data: T) => void;
  setLoading: (key: string, loading: boolean) => void;
  setError: (key: string, error: string | null) => void;
  clearData: (key: string) => void;
}

export const useDataStore = create<DataStore>()(
  devtools(
    (set) => ({
      emissionsData: null,
      carbonFootprint: null,
      energyData: null,
      waterData: null,
      biodiversityData: null,
      airQualityData: null,
      projects: [],
      initiatives: [],
      scenarios: [],
      regulations: [],
      certifications: [],
      loading: {},
      error: {},

      setData: <T>(key: string, data: T) => set((state) => ({ [key]: data, loading: { ...state.loading, [key]: false }, error: { ...state.error, [key]: null } })),
      setLoading: (key: string, loading: boolean) => set((state) => ({ loading: { ...state.loading, [key]: loading } })),
      setError: (key: string, error: string | null) => set((state) => ({ error: { ...state.error, [key]: error }, loading: { ...state.loading, [key]: false } })),
      clearData: (key: string) => set((state) => ({ [key]: null, loading: { ...state.loading, [key]: false }, error: { ...state.error, [key]: null } })),
    }),
    { name: 'DataStore' }
  )
);