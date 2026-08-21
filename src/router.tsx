import { createBrowserRouter } from 'react-router-dom';
import { MainLayout, AuthLayout, DashboardLayout, SettingsLayout } from './components/layout/MainLayout';
import LandingPage from './pages/LandingPage';
import LoginPage from './pages/auth/LoginPage';
import RegisterPage from './pages/auth/RegisterPage';
import ForgotPasswordPage from './pages/auth/ForgotPasswordPage';
import ResetPasswordPage from './pages/auth/ResetPasswordPage';
import TwoFactorSetupPage from './pages/auth/TwoFactorSetupPage';
import DashboardPage from './pages/dashboard/DashboardPage';
import AnalyticsPage from './pages/analytics/AnalyticsPage';
import ReportsPage from './pages/reports/ReportsPage';
import SettingsPage from './components/settings/SettingsPage';
import SavedPlacesPage from './pages/settings/SavedPlacesPage';
import HelpPage from './pages/help/HelpPage';
import NotFoundPage from './pages/NotFoundPage';

export const router = createBrowserRouter([
  {
    path: '/',
    element: <LandingPage />,
  },
  {
    path: '/auth',
    element: <AuthLayout />,
    children: [
      { path: 'login', element: <LoginPage /> },
      { path: 'register', element: <RegisterPage /> },
      { path: 'forgot-password', element: <ForgotPasswordPage /> },
      { path: 'reset-password', element: <ResetPasswordPage /> },
      { path: '2fa-setup', element: <TwoFactorSetupPage /> },
      { path: 'logout', element: <LoginPage /> },
    ],
  },
  {
    path: '/dashboard',
    element: <MainLayout />,
    children: [
      { index: true, element: <DashboardPage /> },
      { path: 'analytics', element: <AnalyticsPage /> },
      { path: 'reports', element: <ReportsPage /> },
    ],
  },
  {
    path: '/settings',
    element: <SettingsLayout />,
    children: [
      { path: '', element: <SettingsPage /> },
      { path: 'saved-places', element: <SavedPlacesPage /> },
    ],
  },
  { path: '/help', element: <HelpPage /> },
  { path: '*', element: <NotFoundPage /> },
]);