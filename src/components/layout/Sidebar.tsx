import { useState } from 'react';
import { Link, useLocation, NavLink } from 'react-router-dom';
import { ChevronLeftIcon, ChevronRightIcon, ChevronDownIcon, ChevronUpIcon, LayoutDashboardIcon, BarChart2Icon, MapIcon, FileTextIcon, CloudIcon, LeafIcon, ZapIcon, DropletIcon, BugIcon, WindIcon, TargetIcon, RocketIcon, FolderIcon, GitBranchIcon, ArrowUpRightIcon, ShieldCheckIcon, SearchIcon, TrophyIcon, UsersIcon, Building2Icon, SettingsIcon, PlugIcon, CodeIcon, MenuIcon, XIcon, BellIcon, HelpCircleIcon, LogOutIcon, UserIcon, MoonIcon, SunIcon } from 'lucide-react';
import { cn } from '../../utils/helpers';
import { Button, IconButton } from '../ui/Button';
import { Avatar } from '../ui/Avatar';
import { Badge } from '../ui/Badge';
import { useAuthStore, useUIStore } from '../../store';
import { NAV_SECTIONS } from '../../types/navigation';
import { usePermissions } from '../../hooks';

interface NavItem {
  id: string;
  label: string;
  href: string;
  icon: React.ComponentType<{ className?: string }>;
  badge?: string | number;
  badgeVariant?: 'primary' | 'success' | 'warning' | 'danger' | 'info';
  disabled?: boolean;
  roles?: string[];
  matchPaths?: string[];
  children?: NavItem[];
}

interface NavSection {
  id: string;
  label?: string;
  items: NavItem[];
  collapsible?: boolean;
  defaultOpen?: boolean;
}

export function Sidebar() {
  const { user } = useAuthStore();
  const { sidebarCollapsed, setSidebarCollapsed, mobileMenuOpen, setMobileMenuOpen, theme, setTheme } = useUIStore();
  const location = useLocation();
  const { hasRole } = usePermissions();
  const [openSections, setOpenSections] = useState<string[]>(['main', 'data', 'climate', 'compliance']);

  const filteredSections = NAV_SECTIONS.filter((section) => {
    const visibleItems = section.items.filter((item) => !item.roles || (user && item.roles.includes(user.role)));
    return visibleItems.length > 0;
  }).map((section) => ({
    ...section,
    items: section.items.filter((item) => !item.roles || (user && item.roles.includes(user.role))),
  }));

  const isActive = (href: string, matchPaths?: string[]) => {
    if (matchPaths) {
      return matchPaths.some((path) => location.pathname.startsWith(path));
    }
    return location.pathname === href || (href !== '/' && location.pathname.startsWith(href));
  };

  const toggleSection = (sectionId: string) => {
    setOpenSections((prev) => (prev.includes(sectionId) ? prev.filter((id) => id !== sectionId) : [...prev, sectionId]));
  };

  const handleMobileLinkClick = () => {
    setMobileMenuOpen(false);
  };

  if (sidebarCollapsed) {
    return (
      <>
        <button
          className="fixed top-4 left-4 z-50 lg:hidden p-2 rounded-lg bg-white dark:bg-secondary-800 shadow-lg"
          onClick={() => setMobileMenuOpen(true)}
          aria-label="Open menu"
        >
          <MenuIcon className="h-6 w-6" />
        </button>
        <aside
          className={cn(
            'fixed inset-y-0 left-0 z-40 w-16 bg-white dark:bg-secondary-800 border-r border-secondary-200 dark:border-secondary-700 transition-all duration-300',
            'transform lg:translate-x-0',
            mobileMenuOpen ? 'translate-x-0' : '-translate-x-full lg:translate-x-0'
          )}
          aria-label="Sidebar navigation"
        >
          <div className="flex h-full flex-col">
            <div className="flex h-16 items-center justify-center border-b border-secondary-200 dark:border-secondary-700">
              <Link to="/dashboard" className="text-xl font-bold text-primary-600" aria-label="Green Combats Home">
                GC
              </Link>
            </div>

            <nav className="flex-1 space-y-1 p-2 overflow-y-auto" aria-label="Main navigation">
              {filteredSections.map((section) => (
                <div key={section.id}>
                  {section.label && (
                    <div className="px-3 py-2">
                      <span className="sr-only">{section.label}</span>
                    </div>
                  )}
                  {section.items.map((item) => {
                    const active = isActive(item.href, item.matchPaths);
                    return (
                      <Link
                        key={item.id}
                        to={item.href}
                        onClick={handleMobileLinkClick}
                        className={cn(
                          'relative flex items-center justify-center gap-2 px-3 py-2.5 rounded-lg',
                          'text-secondary-500 hover:text-secondary-700 hover:bg-secondary-100',
                          'dark:text-secondary-400 dark:hover:text-secondary-200 dark:hover:bg-secondary-700',
                          'transition-colors',
                          active && 'bg-primary-50 text-primary-600 dark:bg-primary-900/30 dark:text-primary-400',
                          item.disabled && 'opacity-50 cursor-not-allowed'
                        )}
                        aria-current={active ? 'page' : undefined}
                        aria-label={item.label}
                        title={sidebarCollapsed ? item.label : undefined}
                      >
                        <item.icon className="h-5 w-5 flex-shrink-0" aria-hidden="true" />
                        {item.badge && (
                          <Badge variant={item.badgeVariant || 'primary'} size="sm" className="absolute -top-1 -right-1 min-w-[18px] h-5 px-1.5">
                            {item.badge}
                          </Badge>
                        )}
                      </Link>
                    );
                  })}
                </div>
              ))}
            </nav>

            <div className="border-t border-secondary-200 dark:border-secondary-700 p-2">
              <div className="flex items-center justify-center gap-2 px-3 py-2">
                <IconButton
                  variant="ghost"
                  size="sm"
                  aria-label={theme === 'dark' ? 'Switch to light mode' : 'Switch to dark mode'}
                  onClick={() => setTheme(theme === 'dark' ? 'light' : 'dark')}
                >
                  {theme === 'dark' ? <SunIcon className="h-5 w-5" /> : <MoonIcon className="h-5 w-5" />}
                </IconButton>
                <IconButton
                  variant="ghost"
                  size="sm"
                  aria-label="Expand sidebar"
                  onClick={() => setSidebarCollapsed(false)}
                >
                  <ChevronRightIcon className="h-5 w-5" />
                </IconButton>
              </div>
            </div>
          </div>
        </aside>
        {mobileMenuOpen && (
          <div
            className="fixed inset-0 z-30 bg-black/50 lg:hidden"
            onClick={() => setMobileMenuOpen(false)}
            aria-hidden="true"
          />
        )}
      </>
    );
  }

  return (
    <>
      <button
        className="fixed top-4 left-4 z-50 lg:hidden p-2 rounded-lg bg-white dark:bg-secondary-800 shadow-lg"
        onClick={() => setMobileMenuOpen(true)}
        aria-label="Open menu"
      >
        <MenuIcon className="h-6 w-6" />
      </button>
      <aside
        className={cn(
          'fixed inset-y-0 left-0 z-40 w-64 bg-white dark:bg-secondary-800 border-r border-secondary-200 dark:border-secondary-700 transition-all duration-300',
          'transform lg:translate-x-0',
          mobileMenuOpen ? 'translate-x-0' : '-translate-x-full lg:translate-x-0'
        )}
        aria-label="Sidebar navigation"
      >
        <div className="flex h-full flex-col">
          <div className="flex h-16 items-center justify-between border-b border-secondary-200 dark:border-secondary-700 px-4">
            <Link to="/dashboard" className="flex items-center gap-2" aria-label="Green Combats Home">
              <div className="h-8 w-8 rounded-xl bg-primary-600 flex items-center justify-center">
                <span className="text-white font-bold text-lg">GC</span>
              </div>
              <span className="text-xl font-bold text-secondary-900 dark:text-white">Green Combats</span>
            </Link>
            <IconButton
              variant="ghost"
              size="sm"
              aria-label="Collapse sidebar"
              onClick={() => setSidebarCollapsed(true)}
            >
              <ChevronLeftIcon className="h-5 w-5" />
            </IconButton>
          </div>

          <nav className="flex-1 space-y-1 p-2 overflow-y-auto" aria-label="Main navigation">
            {filteredSections.map((section) => (
              <div key={section.id}>
                {section.label && (
                  <div className="px-3 py-2">
                    <h3 className="text-xs font-semibold text-secondary-500 uppercase tracking-wider dark:text-secondary-400">
                      {section.label}
                    </h3>
                  </div>
                )}
                {section.collapsible ? (
                  <>
                    <button
                      type="button"
                      className={cn(
                        'w-full flex items-center justify-between px-3 py-2.5 rounded-lg',
                        'text-secondary-500 hover:text-secondary-700 hover:bg-secondary-100',
                        'dark:text-secondary-400 dark:hover:text-secondary-200 dark:hover:bg-secondary-700',
                        'transition-colors',
                        openSections.includes(section.id) ? 'bg-secondary-100 dark:bg-secondary-700' : ''
                      )}
                      onClick={() => toggleSection(section.id)}
                      aria-expanded={openSections.includes(section.id)}
                    >
                      <span className="flex items-center gap-2">
                        <span className="text-sm font-medium">{section.label}</span>
                      </span>
                      {openSections.includes(section.id) ? (
                        <ChevronUpIcon className="h-4 w-4 text-secondary-400" />
                      ) : (
                        <ChevronDownIcon className="h-4 w-4 text-secondary-400" />
                      )}
                    </button>
                    <div
                      className={cn(
                        'overflow-hidden transition-all duration-200',
                        openSections.includes(section.id) ? 'max-h-96 opacity-100' : 'max-h-0 opacity-0'
                      )}
                    >
                      <div className="space-y-1 mt-1 ml-6">
                        {section.items.map((item) => {
                          const active = isActive(item.href, item.matchPaths);
                          return (
                            <Link
                              key={item.id}
                              to={item.href}
                              onClick={handleMobileLinkClick}
                              className={cn(
                                'flex items-center gap-2 px-3 py-2 rounded-lg text-sm',
                                'text-secondary-500 hover:text-secondary-700 hover:bg-secondary-100',
                                'dark:text-secondary-400 dark:hover:text-secondary-200 dark:hover:bg-secondary-700',
                                'transition-colors',
                                active && 'bg-primary-50 text-primary-600 dark:bg-primary-900/30 dark:text-primary-400 font-medium',
                                item.disabled && 'opacity-50 cursor-not-allowed'
                              )}
                              aria-current={active ? 'page' : undefined}
                            >
                              <item.icon className="h-4 w-4 flex-shrink-0" aria-hidden="true" />
                              <span className="truncate">{item.label}</span>
                              {item.badge && (
                                <Badge variant={item.badgeVariant || 'primary'} size="sm" className="ml-auto">
                                  {item.badge}
                                </Badge>
                              )}
                            </Link>
                          );
                        })}
                      </div>
                    </div>
                  </>
                ) : (
                  section.items.map((item) => {
                    const active = isActive(item.href, item.matchPaths);
                    return (
                      <Link
                        key={item.id}
                        to={item.href}
                        onClick={handleMobileLinkClick}
                        className={cn(
                          'flex items-center gap-3 px-3 py-2.5 rounded-lg',
                          'text-secondary-500 hover:text-secondary-700 hover:bg-secondary-100',
                          'dark:text-secondary-400 dark:hover:text-secondary-200 dark:hover:bg-secondary-700',
                          'transition-colors',
                          active && 'bg-primary-50 text-primary-600 dark:bg-primary-900/30 dark:text-primary-400',
                          item.disabled && 'opacity-50 cursor-not-allowed'
                        )}
                        aria-current={active ? 'page' : undefined}
                      >
                        <item.icon className="h-5 w-5 flex-shrink-0" aria-hidden="true" />
                        <span className="truncate font-medium">{item.label}</span>
                        {item.badge && (
                          <Badge variant={item.badgeVariant || 'primary'} size="sm" className="ml-auto">
                            {item.badge}
                          </Badge>
                        )}
                      </Link>
                    );
                  )}
              </div>
            ))}
          </nav>

          <div className="border-t border-secondary-200 dark:border-secondary-700 p-4 space-y-3">
            <div className="flex items-center gap-3">
              <Avatar
                src={user?.avatar}
                name={user?.name}
                size="md"
                status="online"
              />
              <div className="flex-1 min-w-0">
                <p className="text-sm font-medium text-secondary-900 dark:text-white truncate">{user?.name}</p>
                <p className="text-xs text-secondary-500 dark:text-secondary-400 truncate">{user?.email}</p>
              </div>
            </div>

            <div className="flex items-center gap-2">
              <Button variant="ghost" size="sm" className="flex-1 justify-start" asChild>
                <Link to="/settings/profile">
                  <UserIcon className="h-4 w-4 mr-2" />
                  Profile
                </Link>
              </Button>
              <Button variant="ghost" size="sm" className="flex-1 justify-start" asChild>
                <Link to="/settings">
                  <SettingsIcon className="h-4 w-4 mr-2" />
                  Settings
                </Link>
              </Button>
            </div>

            <div className="flex items-center gap-2 pt-2">
              <IconButton
                variant="ghost"
                size="sm"
                aria-label={theme === 'dark' ? 'Switch to light mode' : 'Switch to dark mode'}
                onClick={() => setTheme(theme === 'dark' ? 'light' : 'dark')}
                className="flex-1 justify-start"
              >
                {theme === 'dark' ? <SunIcon className="h-5 w-5 mr-2" /> : <MoonIcon className="h-5 w-5 mr-2" />}
                <span>{theme === 'dark' ? 'Light' : 'Dark'}</span>
              </IconButton>
              <Button variant="ghost" size="sm" className="flex-1 justify-start text-red-600 hover:text-red-700" asChild>
                <Link to="/auth/logout">
                  <LogOutIcon className="h-4 w-4 mr-2" />
                  Sign Out
                </Link>
              </Button>
            </div>
          </div>
        </div>
      </aside>
      {mobileMenuOpen && (
        <div
          className="fixed inset-0 z-30 bg-black/50 lg:hidden"
          onClick={() => setMobileMenuOpen(false)}
          aria-hidden="true"
        />
      )}
    </>
  );
}