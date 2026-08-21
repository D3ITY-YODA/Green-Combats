import { useState, useRef, useEffect } from 'react';
import { Link, useLocation } from 'react-router-dom';
import { SearchIcon, BellIcon, ChevronDownIcon, MoonIcon, SunIcon, HelpCircleIcon, SettingsIcon, UserIcon, LogOutIcon, GlobeIcon, ArrowRightIcon } from 'lucide-react';
import { cn } from '../../utils/helpers';
import { Button, IconButton } from '../ui/Button';
import { Input } from '../ui/Input';
import { Avatar } from '../ui/Avatar';
import { Badge } from '../ui/Badge';
import { Dropdown, DropdownItem } from '../ui/Dropdown';
import { useAuthStore, useUIStore } from '../../store';
import { useClickOutside, useEscapeKey } from '../../hooks';

const notifications = [
  { id: '1', type: 'alert', title: 'High Emissions Alert', message: 'Facility A exceeded CO2 threshold', time: '5 min ago', unread: true },
  { id: '2', type: 'report', title: 'Monthly Report Ready', message: 'Carbon footprint report for January', time: '1 hour ago', unread: true },
  { id: '3', type: 'update', title: 'New Regulation', message: 'EU CSRD reporting requirements updated', time: '3 hours ago', unread: false },
  { id: '4', type: 'mention', title: 'Comment on Project', message: 'Sarah mentioned you in "Solar Farm Phase 2"', time: 'Yesterday', unread: false },
];

const quickActions = [
  { id: 'new-project', label: 'New Project', icon: 'PlusIcon', href: '/projects/new' },
  { id: 'new-report', label: 'Create Report', icon: 'FilePlusIcon', href: '/reports/new' },
  { id: 'new-target', label: 'Set Target', icon: 'TargetIcon', href: '/targets/new' },
  { id: 'upload-data', label: 'Upload Data', icon: 'UploadIcon', href: '/data/upload' },
];

export function Header() {
  const location = useLocation();
  const { user } = useAuthStore();
  const { theme, setTheme, addToast } = useUIStore();
  const [searchQuery, setSearchQuery] = useState('');
  const [notificationsOpen, setNotificationsOpen] = useState(false);
  const [userMenuOpen, setUserMenuOpen] = useState(false);
  const [quickActionsOpen, setQuickActionsOpen] = useState(false);
  const [commandPaletteOpen, setCommandPaletteOpen] = useState(false);
  const notificationsRef = useRef<HTMLDivElement>(null);
  const userMenuRef = useRef<HTMLDivElement>(null);
  const quickActionsRef = useRef<HTMLDivElement>(null);
  const searchRef = useRef<HTMLDivElement>(null);

  useClickOutside(notificationsRef, () => setNotificationsOpen(false));
  useClickOutside(userMenuRef, () => setUserMenuOpen(false));
  useClickOutside(quickActionsRef, () => setQuickActionsOpen(false));
  useClickOutside(searchRef, () => {});

  useEscapeKey(() => {
    setNotificationsOpen(false);
    setUserMenuOpen(false);
    setQuickActionsOpen(false);
    setCommandPaletteOpen(false);
  });

  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if ((e.metaKey || e.ctrlKey) && e.key === 'k') {
        e.preventDefault();
        setCommandPaletteOpen(!commandPaletteOpen);
      }
    };
    document.addEventListener('keydown', handleKeyDown);
    return () => document.removeEventListener('keydown', handleKeyDown);
  }, [commandPaletteOpen]);

  const unreadCount = notifications.filter((n) => n.unread).length;

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    if (searchQuery.trim()) {
      addToast({ type: 'info', title: 'Search', message: `Searching for "${searchQuery}"...` });
    }
  };

  const userMenuItems: DropdownItem[] = [
    { id: 'profile', label: 'Profile', icon: <UserIcon className="h-4 w-4" /> },
    { id: 'settings', label: 'Settings', icon: <SettingsIcon className="h-4 w-4" /> },
    { id: 'divider-1', label: '', divider: true },
    { id: 'help', label: 'Help Center', icon: <HelpCircleIcon className="h-4 w-4" /> },
    { id: 'divider-2', label: '', divider: true },
    { id: 'logout', label: 'Sign Out', icon: <LogOutIcon className="h-4 w-4" />, danger: true },
  ];

  const quickActionItems: DropdownItem[] = quickActions.map((action) => ({
    id: action.id,
    label: action.label,
    icon: <PlusIcon className="h-4 w-4" />,
  }));

  return (
    <header className="sticky top-0 z-30 w-full bg-white/80 dark:bg-secondary-800/80 backdrop-blur-lg border-b border-secondary-200 dark:border-secondary-700">
      <div className="flex h-16 items-center gap-4 px-4 sm:px-6 lg:px-8">
        <div className="relative flex-1 max-w-xl hidden sm:block" ref={searchRef}>
          <form onSubmit={handleSearch} className="relative">
            <SearchIcon className="absolute left-3 top-1/2 h-5 w-5 -translate-y-1/2 text-secondary-400" aria-hidden="true" />
            <Input
              type="search"
              placeholder="Search projects, reports, regulations... (⌘K)"
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="pl-10 pr-10 h-10 text-sm"
              aria-label="Global search"
              autoComplete="off"
            />
            {searchQuery && (
              <button
                type="button"
                className="absolute right-3 top-1/2 -translate-y-1/2 text-secondary-400 hover:text-secondary-600"
                onClick={() => setSearchQuery('')}
                aria-label="Clear search"
              >
                <svg className="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            )}
          </form>
        </div>

        <div className="flex items-center gap-1">
          <IconButton
            variant="ghost"
            size="sm"
            aria-label="Quick actions"
            onClick={() => setQuickActionsOpen(!quickActionsOpen)}
          >
            <PlusIcon className="h-5 w-5" />
          </IconButton>

          <div className="relative" ref={notificationsRef}>
            <IconButton
              variant="ghost"
              size="sm"
              aria-label={`Notifications${unreadCount > 0 ? `, ${unreadCount} unread` : ''}`}
              onClick={() => setNotificationsOpen(!notificationsOpen)}
            >
              <span className="relative">
                <BellIcon className="h-5 w-5" />
                {unreadCount > 0 && (
                  <span className="absolute -top-1 -right-1 flex h-5 w-5 items-center justify-center rounded-full bg-red-500 text-white text-xs font-medium">
                    {unreadCount > 9 ? '9+' : unreadCount}
                  </span>
                )}
              </span>
            </IconButton>

            {notificationsOpen && (
              <div className="absolute right-0 mt-2 w-80 sm:w-96 bg-white dark:bg-secondary-800 rounded-xl shadow-lg border border-secondary-200 dark:border-secondary-700 overflow-hidden">
                <div className="flex items-center justify-between px-4 py-3 border-b border-secondary-200 dark:border-secondary-700">
                  <h3 className="font-semibold text-secondary-900 dark:text-white">Notifications</h3>
                  <Button variant="ghost" size="sm" onClick={() => addToast({ type: 'success', title: 'All marked as read' })}>
                    Mark all read
                  </Button>
                </div>
                <div className="max-h-96 overflow-y-auto">
                  {notifications.length === 0 ? (
                    <div className="px-4 py-8 text-center text-secondary-500 dark:text-secondary-400">
                      No notifications
                    </div>
                  ) : (
                    notifications.map((notification) => (
                      <Link
                        key={notification.id}
                        href="#"
                        className={cn(
                          'flex items-start gap-3 px-4 py-3 border-b border-secondary-100 dark:border-secondary-700',
                          'hover:bg-secondary-50 dark:hover:bg-secondary-700/50',
                          notification.unread && 'bg-primary-50/50 dark:bg-primary-900/10'
                        )}
                      >
                        <div className={cn('flex-shrink-0 h-2 w-2 rounded-full mt-2', notification.unread ? 'bg-primary-500' : 'bg-transparent')}/>
                        <div className="flex-1 min-w-0">
                          <p className={cn('text-sm font-medium', notification.unread ? 'text-secondary-900 dark:text-white' : 'text-secondary-700 dark:text-secondary-300')}>
                            {notification.title}
                          </p>
                          <p className="mt-0.5 text-sm text-secondary-500 dark:text-secondary-400 truncate">
                            {notification.message}
                          </p>
                          <p className="mt-1 text-xs text-secondary-400 dark:text-secondary-500">{notification.time}</p>
                        </div>
                      </Link>
                    ))
                  )}
                </div>
                <div className="px-4 py-3 border-t border-secondary-200 dark:border-secondary-700">
                  <Button variant="ghost" fullWidth asChild>
                    <Link to="/notifications">View all notifications</Link>
                  </Button>
                </div>
              </div>
            )}
          </div>

          <div className="relative" ref={quickActionsRef}>
            <IconButton
              variant="ghost"
              size="sm"
              aria-label="Quick actions"
              onClick={() => setQuickActionsOpen(!quickActionsOpen)}
            >
              <PlusIcon className="h-5 w-5" />
            </IconButton>

            {quickActionsOpen && (
              <div className="absolute right-0 mt-2 w-48 bg-white dark:bg-secondary-800 rounded-xl shadow-lg border border-secondary-200 dark:border-secondary-700 overflow-hidden">
                <div className="p-2">
                  {quickActions.map((action) => (
                    <Button
                      key={action.id}
                      variant="ghost"
                      fullWidth
                      className="justify-start gap-2"
                      asChild
                    >
                      <Link to={action.href} onClick={() => setQuickActionsOpen(false)}>
                        <PlusIcon className="h-4 w-4" />
                        {action.label}
                      </Link>
                    </Button>
                  ))}
                </div>
              </div>
            )}
          </div>

          <IconButton
            variant="ghost"
            size="sm"
            aria-label={theme === 'dark' ? 'Switch to light mode' : 'Switch to dark mode'}
            onClick={() => setTheme(theme === 'dark' ? 'light' : 'dark')}
          >
            {theme === 'dark' ? <SunIcon className="h-5 w-5" /> : <MoonIcon className="h-5 w-5" />}
          </IconButton>

          <div className="relative" ref={userMenuRef}>
            <button
              onClick={() => setUserMenuOpen(!userMenuOpen)}
              className="flex items-center gap-2 px-2 py-1.5 rounded-lg hover:bg-secondary-100 dark:hover:bg-secondary-700 transition-colors"
              aria-label="User menu"
              aria-expanded={userMenuOpen}
              aria-haspopup="true"
            >
              <Avatar
                src={user?.avatar}
                name={user?.name}
                size="sm"
                status="online"
              />
              <span className="hidden md:block text-sm font-medium text-secondary-700 dark:text-secondary-300 truncate max-w-[120px]">
                {user?.name}
              </span>
              <ChevronDownIcon className={cn('h-4 w-4 text-secondary-400 transition-transform', userMenuOpen && 'rotate-180')} />
            </button>

            {userMenuOpen && (
              <div className="absolute right-0 mt-2 w-56 bg-white dark:bg-secondary-800 rounded-xl shadow-lg border border-secondary-200 dark:border-secondary-700 overflow-hidden">
                <div className="px-4 py-3 border-b border-secondary-200 dark:border-secondary-700">
                  <p className="font-medium text-secondary-900 dark:text-white">{user?.name}</p>
                  <p className="text-sm text-secondary-500 dark:text-secondary-400">{user?.email}</p>
                  <Badge variant="primary" size="sm" className="mt-1">
                    {user?.role}
                  </Badge>
                </div>
                <div className="py-1">
                  {userMenuItems.map((item) => (
                    item.divider ? (
                      <div key={item.id} className="border-t border-secondary-200 dark:border-secondary-700 my-1" />
                    ) : (
                      <Link
                        key={item.id}
                        href={item.id === 'logout' ? '/auth/logout' : `/settings/${item.id}`}
                        onClick={() => setUserMenuOpen(false)}
                        className={cn(
                          'flex items-center gap-3 px-4 py-2 text-sm',
                          'text-secondary-700 hover:text-secondary-900 hover:bg-secondary-100',
                          'dark:text-secondary-300 dark:hover:text-secondary-100 dark:hover:bg-secondary-700',
                          item.danger && 'text-red-600 hover:text-red-700'
                        )}
                      >
                        <span className="flex-shrink-0">{item.icon}</span>
                        {item.label}
                      </Link>
                    )
                  ))}
                </div>
              </div>
            )}
          </div>
        </div>

        {commandPaletteOpen && (
          <div className="fixed inset-0 z-50 flex items-start justify-center pt-20 bg-black/50" onClick={() => setCommandPaletteOpen(false)}>
            <div className="w-full max-w-2xl bg-white dark:bg-secondary-800 rounded-xl shadow-xl border border-secondary-200 dark:border-secondary-700 overflow-hidden" onClick={(e) => e.stopPropagation()}>
              <div className="p-4 border-b border-secondary-200 dark:border-secondary-700">
                <div className="relative">
                  <SearchIcon className="absolute left-3 top-1/2 h-5 w-5 -translate-y-1/2 text-secondary-400" />
                  <Input
                    type="search"
                    placeholder="Type a command or search... (⌘K to close)"
                    className="pl-10 pr-10 h-12 text-lg"
                    autoFocus
                    autoComplete="off"
                  />
                  <kbd className="absolute right-3 top-1/2 -translate-y-1/2 px-2 py-1 text-xs text-secondary-400 bg-secondary-100 dark:bg-secondary-700 rounded">
                    ⌘K
                  </kbd>
                </div>
              </div>
              <div className="max-h-96 overflow-y-auto p-4">
                <div className="space-y-4">
                  <div>
                    <h4 className="px-3 py-1 text-xs font-semibold text-secondary-500 uppercase tracking-wider">Suggested</h4>
                    <div className="space-y-1">
                      {[
                        { label: 'New Project', href: '/projects/new', icon: PlusIcon },
                        { label: 'Create Report', href: '/reports/new', icon: FilePlusIcon },
                        { label: 'View Dashboard', href: '/dashboard', icon: LayoutDashboardIcon },
                        { label: 'Emissions Tracking', href: '/emissions', icon: CloudIcon },
                      ].map((item) => (
                        <Link
                          key={item.label}
                          to={item.href}
                          onClick={() => setCommandPaletteOpen(false)}
                          className="flex items-center gap-3 px-3 py-2 rounded-lg text-sm hover:bg-secondary-100 dark:hover:bg-secondary-700"
                        >
                          <item.icon className="h-5 w-5 text-secondary-400" />
                          {item.label}
                        </Link>
                      ))}
</div>
              </div>
            </div>
          )}
        </div>
      </div>
    </header>
  );
}

function PlusIcon({ className }: { className?: string }) {
  return <svg className={className} fill="none" viewBox="0 0 24 24" stroke="currentColor"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" /></svg>;
}
function FilePlusIcon({ className }: { className?: string }) {
  return <svg className={className} fill="none" viewBox="0 0 24 24" stroke="currentColor"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" /></svg>;
}
function TargetIcon({ className }: { className?: string }) {
  return <svg className={className} fill="none" viewBox="0 0 24 24" stroke="currentColor"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" /></svg>;
}
function UploadIcon({ className }: { className?: string }) {
  return <svg className={className} fill="none" viewBox="0 0 24 24" stroke="currentColor"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" /></svg>;
}