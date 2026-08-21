import '@testing-library/jest-dom'
import { vi } from 'vitest'

vi.mock('react-router-dom', () => ({
  useNavigate: () => vi.fn(),
  useLocation: () => ({ pathname: '/dashboard' }),
  useParams: () => ({}),
  Link: ({ children, to, ...props }: any) => <a href={to} {...props}>{children}</a>,
  NavLink: ({ children, to, ...props }: any) => <a href={to} {...props}>{children}</a>,
  Outlet: () => null,
}))

vi.mock('../../store', () => ({
  useAuthStore: () => ({
    user: {
      id: 'usr_001',
      name: 'Test User',
      email: 'test@example.com',
      role: 'admin',
      avatar: null,
      preferences: {
        theme: 'system',
        language: 'en',
        timezone: 'UTC',
        dateFormat: 'YYYY-MM-DD',
        notifications: { email: true, push: true, inApp: true, frequency: 'immediate', types: { alerts: true, reports: true, updates: true, mentions: true, assignments: true } },
        dashboard: { defaultView: 'overview', autoRefresh: true, refreshInterval: 300000, compactMode: false, showTooltips: true },
        maps: { defaultBasemap: 'satellite', defaultZoom: 3, defaultCenter: [20, 0], showCoordinates: true, measureUnit: 'metric' },
      },
      security: { twoFactorEnabled: false, sessions: [], passwordLastChanged: '2024-01-01', loginHistory: [] },
      createdAt: '2024-01-01',
      updatedAt: '2024-01-01',
    },
    isAuthenticated: true,
    login: vi.fn(),
    logout: vi.fn(),
    refreshAccessToken: vi.fn(),
    setUser: vi.fn(),
    updateUser: vi.fn(),
  }),
  useUIStore: () => ({
    sidebarOpen: true,
    sidebarCollapsed: false,
    mobileMenuOpen: false,
    theme: 'system',
    setSidebarOpen: vi.fn(),
    setSidebarCollapsed: vi.fn(),
    setTheme: vi.fn(),
    addToast: vi.fn(),
  }),
  useNavigationStore: () => ({
    currentPath: '/dashboard',
    breadcrumbs: [],
    setCurrentPath: vi.fn(),
    setBreadcrumbs: vi.fn(),
  }),
}))

Object.defineProperty(window, 'matchMedia', {
  writable: true,
  value: vi.fn().mockImplementation(query => ({
    matches: false,
    media: query,
    onchange: null,
    addListener: vi.fn(),
    removeListener: vi.fn(),
    addEventListener: vi.fn(),
    removeEventListener: vi.fn(),
    dispatchEvent: vi.fn(),
  })),
})

Object.defineProperty(window, 'localStorage', {
  writable: true,
  value: {
    getItem: vi.fn(),
    setItem: vi.fn(),
    removeItem: vi.fn(),
    clear: vi.fn(),
  },
})

Object.defineProperty(window, 'sessionStorage', {
  writable: true,
  value: {
    getItem: vi.fn(),
    setItem: vi.fn(),
    removeItem: vi.fn(),
    clear: vi.fn(),
  },
})

HTMLCanvasElement.prototype.getContext = vi.fn()