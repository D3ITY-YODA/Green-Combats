import { Fragment } from 'react';
import { TabGroup, TabList, Tab, TabPanels, TabPanel } from '@headlessui/react';
import { cn } from '../utils/helpers';

export interface TabItem {
  id: string;
  label: string;
  icon?: React.ReactNode;
  count?: number;
  disabled?: boolean;
  badge?: string | number;
  badgeVariant?: 'primary' | 'success' | 'warning' | 'danger' | 'info';
}

export interface TabsProps {
  tabs: TabItem[];
  activeTab: string;
  onChange: (tabId: string) => void;
  variant?: 'line' | 'enclosed' | 'soft' | 'pills';
  orientation?: 'horizontal' | 'vertical';
  fullWidth?: boolean;
  className?: string;
  tabClassName?: string;
  panelClassName?: string;
}

export function Tabs({
  tabs,
  activeTab,
  onChange,
  variant = 'line',
  orientation = 'horizontal',
  fullWidth = false,
  className,
  tabClassName,
  panelClassName,
}: TabsProps) {
  const variantClasses = {
    line: {
      list: 'border-b border-secondary-200 dark:border-secondary-700',
      tab: 'border-b-2 -mb-px',
      active: 'border-primary-600 text-primary-600',
      inactive: 'border-transparent text-secondary-500 hover:text-secondary-700 hover:border-secondary-300 dark:hover:text-secondary-300 dark:hover:border-secondary-600',
    },
    enclosed: {
      list: 'bg-secondary-100 dark:bg-secondary-800 p-1 rounded-lg',
      tab: 'rounded-md',
      active: 'bg-white dark:bg-secondary-700 text-primary-600 dark:text-primary-400 shadow-sm',
      inactive: 'text-secondary-600 hover:text-secondary-900 dark:text-secondary-400 dark:hover:text-secondary-100',
    },
    soft: {
      list: '',
      tab: 'rounded-lg',
      active: 'bg-primary-50 dark:bg-primary-900/30 text-primary-600 dark:text-primary-400',
      inactive: 'text-secondary-600 hover:text-secondary-900 hover:bg-secondary-100 dark:text-secondary-400 dark:hover:text-secondary-100 dark:hover:bg-secondary-800',
    },
    pills: {
      list: 'gap-1',
      tab: 'rounded-lg',
      active: 'bg-primary-600 text-white shadow-sm',
      inactive: 'text-secondary-600 hover:text-secondary-900 hover:bg-secondary-100 dark:text-secondary-400 dark:hover:text-secondary-100 dark:hover:bg-secondary-800',
    },
  };

  const v = variantClasses[variant];

  return (
    <TabGroup className={cn(className)}>
      <TabList
        className={cn(
          'flex',
          v.list,
          orientation === 'vertical' ? 'flex-col' : 'flex-row',
          fullWidth && 'w-full',
          variant === 'pills' && 'bg-secondary-100 dark:bg-secondary-800 p-1 rounded-lg'
        )}
        role="tablist"
        aria-orientation={orientation}
      >
        {tabs.map((tab) => (
          <Tab
            key={tab.id}
            disabled={tab.disabled}
            className={({ selected }) =>
              cn(
                'relative inline-flex items-center justify-center gap-2 px-4 py-2.5 text-sm font-medium transition-all duration-200',
                'focus:outline-none focus-visible:ring-2 focus-visible:ring-primary-500 focus-visible:ring-offset-2',
                'disabled:opacity-50 disabled:cursor-not-allowed',
                selected ? v.active : v.inactive,
                tabClassName
              )
            }
            onClick={() => onChange(tab.id)}
          >
            {tab.icon && <span className="flex-shrink-0">{tab.icon}</span>}
            {tab.label}
            {(tab.count !== undefined || tab.badge !== undefined) && (
              <span
                className={cn(
                  'flex-shrink-0 px-2 py-0.5 text-xs font-medium rounded-full',
                  selected
                    ? variant === 'pills'
                      ? 'bg-white/20 text-white'
                      : 'bg-primary-100 text-primary-700 dark:bg-primary-900/30 dark:text-primary-300'
                    : 'bg-secondary-100 text-secondary-600 dark:bg-secondary-700 dark:text-secondary-400'
                )}
              >
                {tab.badge ?? tab.count}
              </span>
            )}
          </Tab>
        ))}
      </TabList>
      <TabPanels className={cn('mt-4', panelClassName)}>
        {tabs.map((tab) => (
          <TabPanel key={tab.id} className="focus:outline-none">
            <div id={`${tab.id}-panel`} role="tabpanel">
              {tab.children}
            </div>
          </TabPanel>
        ))}
      </TabPanels>
    </TabGroup>
  );
}

export interface SimpleTabsProps {
  tabs: TabItem[];
  activeTab: string;
  onChange: (tabId: string) => void;
  variant?: 'line' | 'enclosed' | 'soft' | 'pills';
  className?: string;
}

export function SimpleTabs({ tabs, activeTab, onChange, variant = 'line', className }: SimpleTabsProps) {
  const variantClasses = {
    line: {
      container: 'border-b border-secondary-200 dark:border-secondary-700',
      tab: 'border-b-2 -mb-px',
      active: 'border-primary-600 text-primary-600',
      inactive: 'border-transparent text-secondary-500 hover:text-secondary-700 hover:border-secondary-300 dark:hover:text-secondary-300 dark:hover:border-secondary-600',
    },
    enclosed: {
      container: 'bg-secondary-100 dark:bg-secondary-800 rounded-lg p-1',
      tab: 'rounded-md',
      active: 'bg-white dark:bg-secondary-700 text-primary-600 dark:text-primary-400 shadow-sm',
      inactive: 'text-secondary-600 hover:text-secondary-900 dark:text-secondary-400 dark:hover:text-secondary-100',
    },
    soft: {
      container: '',
      tab: 'rounded-lg',
      active: 'bg-primary-50 dark:bg-primary-900/30 text-primary-600 dark:text-primary-400',
      inactive: 'text-secondary-600 hover:text-secondary-900 hover:bg-secondary-100 dark:text-secondary-400 dark:hover:text-secondary-100 dark:hover:bg-secondary-800',
    },
    pills: {
      container: 'bg-secondary-100 dark:bg-secondary-800 rounded-lg p-1',
      tab: 'rounded-lg',
      active: 'bg-primary-600 text-white shadow-sm',
      inactive: 'text-secondary-600 hover:text-secondary-900 hover:bg-secondary-100 dark:text-secondary-400 dark:hover:text-secondary-100 dark:hover:bg-secondary-800',
    },
  };

  const v = variantClasses[variant];

  return (
    <div className={cn(v.container, className)}>
      <div className="flex gap-1" role="tablist">
        {tabs.map((tab) => (
          <button
            key={tab.id}
            role="tab"
            aria-selected={activeTab === tab.id}
            aria-controls={`${tab.id}-panel`}
            id={`${tab.id}-tab`}
            disabled={tab.disabled}
            onClick={() => !tab.disabled && onChange(tab.id)}
            className={cn(
              'flex items-center justify-center gap-2 px-4 py-2.5 text-sm font-medium transition-all duration-200',
              'focus:outline-none focus-visible:ring-2 focus-visible:ring-primary-500 focus-visible:ring-offset-2',
              'disabled:opacity-50 disabled:cursor-not-allowed',
              fullWidth && 'flex-1',
              activeTab === tab.id ? v.active : v.inactive
            )}
          >
            {tab.icon && <span className="flex-shrink-0">{tab.icon}</span>}
            {tab.label}
            {(tab.count !== undefined || tab.badge !== undefined) && (
              <span
                className={cn(
                  'flex-shrink-0 px-2 py-0.5 text-xs font-medium rounded-full',
                  activeTab === tab.id
                    ? variant === 'pills'
                      ? 'bg-white/20 text-white'
                      : 'bg-primary-100 text-primary-700 dark:bg-primary-900/30 dark:text-primary-300'
                    : 'bg-secondary-100 text-secondary-600 dark:bg-secondary-700 dark:text-secondary-400'
                )}
              >
                {tab.badge ?? tab.count}
              </span>
            )}
          </button>
        ))}
      </div>
    </div>
  );
}

export interface VerticalTabsProps {
  tabs: TabItem[];
  activeTab: string;
  onChange: (tabId: string) => void;
  className?: string;
}

export function VerticalTabs({ tabs, activeTab, onChange, className }: VerticalTabsProps) {
  return (
    <div className={cn('flex gap-6', className)}>
      <nav className="flex-1 min-w-0" aria-label="Vertical tabs">
        <ul className="flex flex-col gap-1" role="tablist" aria-orientation="vertical">
          {tabs.map((tab) => (
            <li key={tab.id} role="presentation">
              <button
                role="tab"
                aria-selected={activeTab === tab.id}
                aria-controls={`${tab.id}-panel`}
                id={`${tab.id}-tab`}
                disabled={tab.disabled}
                onClick={() => !tab.disabled && onChange(tab.id)}
                className={cn(
                  'flex items-center gap-3 px-3 py-2.5 text-sm font-medium rounded-lg transition-all duration-200',
                  'focus:outline-none focus-visible:ring-2 focus-visible:ring-primary-500 focus-visible:ring-offset-2',
                  'disabled:opacity-50 disabled:cursor-not-allowed',
                  activeTab === tab.id
                    ? 'bg-primary-50 dark:bg-primary-900/30 text-primary-600 dark:text-primary-400'
                    : 'text-secondary-600 hover:text-secondary-900 hover:bg-secondary-100 dark:text-secondary-400 dark:hover:text-secondary-100 dark:hover:bg-secondary-800'
                )}
              >
                {tab.icon && <span className="flex-shrink-0">{tab.icon}</span>}
                <span className="flex-1 text-left">{tab.label}</span>
                {(tab.count !== undefined || tab.badge !== undefined) && (
                  <span
                    className={cn(
                      'flex-shrink-0 px-2 py-0.5 text-xs font-medium rounded-full',
                      activeTab === tab.id
                        ? 'bg-primary-100 text-primary-700 dark:bg-primary-900/30 dark:text-primary-300'
                        : 'bg-secondary-100 text-secondary-600 dark:bg-secondary-700 dark:text-secondary-400'
                    )}
                  >
                    {tab.badge ?? tab.count}
                  </span>
                )}
              </button>
            </li>
          ))}
        </ul>
      </nav>
      <div className="flex-1 min-w-0">
        {tabs.map((tab) => (
          <div
            key={tab.id}
            id={`${tab.id}-panel`}
            role="tabpanel"
            aria-labelledby={`${tab.id}-tab`}
            hidden={activeTab !== tab.id}
            className="focus:outline-none"
          >
            {tab.children}
          </div>
        ))}
      </div>
    </div>
  );
}