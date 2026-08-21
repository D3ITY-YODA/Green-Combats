import { Outlet } from 'react-router-dom';
import { Sidebar } from './Sidebar';
import { Header } from './Header';
import { useUIStore } from '../../store';
import { cn } from '../../utils/helpers';

export function MainLayout() {
  const { sidebarCollapsed, mobileMenuOpen } = useUIStore();

  return (
    <div className="min-h-screen bg-secondary-50 dark:bg-secondary-900">
      <Sidebar />
      <div
        className={cn(
          'transition-all duration-300',
          'lg:pl-64',
          sidebarCollapsed && 'lg:pl-16',
          mobileMenuOpen && 'lg:pl-64'
        )}
      >
        <Header />
        <main className="p-4 sm:p-6 lg:p-8" id="main-content" tabIndex={-1}>
          <Outlet />
        </main>
        <footer className="border-t border-secondary-200 dark:border-secondary-700 bg-white dark:bg-secondary-800">
          <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
            <div className="flex flex-col md:flex-row items-center justify-between gap-4">
              <p className="text-sm text-secondary-500 dark:text-secondary-400">
                © 2024 Green Combats. All rights reserved.
              </p>
              <div className="flex items-center gap-6 text-sm text-secondary-500 dark:text-secondary-400">
                <a href="/privacy" className="hover:text-secondary-700 dark:hover:text-secondary-300 transition-colors">
                  Privacy Policy
                </a>
                <a href="/terms" className="hover:text-secondary-700 dark:hover:text-secondary-300 transition-colors">
                  Terms of Service
                </a>
                <a href="/cookies" className="hover:text-secondary-700 dark:hover:text-secondary-300 transition-colors">
                  Cookie Policy
                </a>
                <a href="/security" className="hover:text-secondary-700 dark:hover:text-secondary-300 transition-colors">
                  Security
                </a>
              </div>
            </div>
          </div>
        </footer>
      </div>
    </div>
  );
}

export function AuthLayout({ children }: { children: React.ReactNode }) {
  return (
    <div className="min-h-screen bg-secondary-50 dark:bg-secondary-900">
      <div className="flex min-h-screen items-center justify-center px-4 py-12">
        {children}
      </div>
    </div>
  );
}

export function BlankLayout({ children }: { children: React.ReactNode }) {
  return (
    <div className="min-h-screen bg-secondary-50 dark:bg-secondary-900">
      {children}
    </div>
  );
}

export function DashboardLayout({ children }: { children: React.ReactNode }) {
  const { sidebarCollapsed, mobileMenuOpen } = useUIStore();

  return (
    <div className="min-h-screen bg-secondary-50 dark:bg-secondary-900">
      <Sidebar />
      <div
        className={cn(
          'transition-all duration-300',
          'lg:pl-64',
          sidebarCollapsed && 'lg:pl-16',
          mobileMenuOpen && 'lg:pl-64'
        )}
      >
        <Header />
        <div className="p-4 sm:p-6 lg:p-8">
          {children}
        </div>
      </div>
    </div>
  );
}