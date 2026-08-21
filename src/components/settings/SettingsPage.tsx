import { useState } from 'react';
import { NavLink, Outlet } from 'react-router-dom';
import { UserIcon, BellIcon, ShieldIcon, SettingsIcon, ChevronRightIcon } from 'lucide-react';
import { cn } from '../../utils/helpers';
import { Card, CardBody } from '../ui/Card';
import { Badge } from '../ui/Badge';
import { ProfileSettings } from './ProfileSettings';
import { PreferencesSettings } from './PreferencesSettings';
import { NotificationSettings } from './NotificationSettings';
import { SecuritySettings } from './SecuritySettings';

const settingsSections = [
  { id: 'profile', label: 'Profile', icon: UserIcon, description: 'Personal information and avatar' },
  { id: 'preferences', label: 'Preferences', icon: SettingsIcon, description: 'Appearance, dashboard, and map settings' },
  { id: 'notifications', label: 'Notifications', icon: BellIcon, description: 'Email, push, and in-app notifications' },
  { id: 'security', label: 'Security', icon: ShieldIcon, description: 'Password, 2FA, sessions, and API keys' },
];

export function SettingsPage() {
  const [activeSection, setActiveSection] = useState('profile');

  const sectionComponents: Record<string, React.ComponentType> = {
    profile: ProfileSettings,
    preferences: PreferencesSettings,
    notifications: NotificationSettings,
    security: SecuritySettings,
  };

  const ActiveComponent = sectionComponents[activeSection] || ProfileSettings;

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-secondary-900 dark:text-white">Settings</h1>
        <p className="mt-1 text-secondary-600 dark:text-secondary-400">
          Manage your account settings and preferences
        </p>
      </div>

      <div className="grid gap-6 lg:grid-cols-4">
        <aside className="lg:col-span-1">
          <nav className="space-y-1" aria-label="Settings navigation">
            {settingsSections.map((section) => (
              <NavLink
                key={section.id}
                to={`/settings/${section.id}`}
                onClick={() => setActiveSection(section.id)}
                className={({ isActive }) => cn(
                  'flex items-center gap-3 px-4 py-3 rounded-xl transition-colors',
                  isActive
                    ? 'bg-primary-50 dark:bg-primary-900/20 text-primary-700 dark:text-primary-300 border border-primary-200 dark:border-primary-800'
                    : 'text-secondary-600 hover:text-secondary-900 hover:bg-secondary-100 dark:text-secondary-400 dark:hover:text-secondary-100 dark:hover:bg-secondary-800'
                )}
              >
                <section.icon className="h-5 w-5 flex-shrink-0" aria-hidden="true" />
                <div className="flex-1 min-w-0">
                  <p className="font-medium truncate">{section.label}</p>
                  <p className="text-xs text-secondary-500 dark:text-secondary-400 truncate">{section.description}</p>
                </div>
                <ChevronRightIcon className="h-4 w-4 text-secondary-400" />
              </NavLink>
            ))}
          </nav>

          <div className="mt-6 p-4 bg-secondary-50 dark:bg-secondary-800 rounded-xl">
            <p className="text-sm font-medium text-secondary-900 dark:text-white mb-2">Need help?</p>
            <p className="text-sm text-secondary-600 dark:text-secondary-400 mb-4">
              Visit our help center or contact support for assistance.
            </p>
            <div className="flex gap-2">
              <a href="/help" className="text-sm text-primary-600 hover:text-primary-700 dark:text-primary-400 font-medium">
                Help Center
              </a>
              <span className="text-secondary-400">·</span>
              <a href="/support" className="text-sm text-primary-600 hover:text-primary-700 dark:text-primary-400 font-medium">
                Contact Support
              </a>
            </div>
          </div>
        </aside>

        <div className="lg:col-span-3">
          <Card>
            <CardBody className="p-0">
              <ActiveComponent />
            </CardBody>
          </Card>
        </div>
      </div>
    </div>
  );
}

export function SettingsLayout({ children }: { children: React.ReactNode }) {
  return (
    <div className="min-h-screen bg-secondary-50 dark:bg-secondary-900">
      {children}
    </div>
  );
}