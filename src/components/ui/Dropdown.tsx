import { Fragment, useRef, useState } from 'react';
import { Menu, Transition } from '@headlessui/react';
import { CheckIcon, ChevronDownIcon } from 'lucide-react';
import { cn } from '../utils/helpers';
import { Button } from './Button';
import { IconButton } from './Button';

export interface DropdownItem {
  id: string;
  label: string;
  icon?: React.ReactNode;
  shortcut?: string;
  disabled?: boolean;
  danger?: boolean;
  divider?: boolean;
  subItems?: DropdownItem[];
}

export interface DropdownProps {
  trigger: React.ReactNode;
  items: DropdownItem[];
  align?: 'left' | 'right';
  closeOnSelect?: boolean;
  className?: string;
}

export function Dropdown({ trigger, items, align = 'right', closeOnSelect = true, className }: DropdownProps) {
  const [selectedId, setSelectedId] = useState<string | null>(null);
  const menuRef = useRef<HTMLDivElement>(null);

  const handleSelect = (item: DropdownItem) => {
    if (item.disabled) return;
    setSelectedId(item.id);
    if (closeOnSelect) {
      // Menu closes automatically
    }
  };

  const renderItems = (items: DropdownItem[], depth = 0) => (
    <>
      {items.map((item, index) => {
        if (item.divider) {
          return (
            <div key={`divider-${index}`} className="border-t border-secondary-200 my-1 dark:border-secondary-700" role="separator" />
          );
        }

        if (item.subItems && item.subItems.length > 0) {
          return (
            <Menu.Item key={item.id} disabled={item.disabled}>
              {({ active }) => (
                <div
                  className={cn(
                    'relative flex items-center gap-2 px-3 py-2 text-sm transition-colors',
                    'focus:outline-none focus:bg-secondary-100 dark:focus:bg-secondary-800',
                    item.disabled ? 'opacity-50 cursor-not-allowed' : '',
                    active ? 'bg-secondary-100 dark:bg-secondary-800' : ''
                  )}
                >
                  <span className="flex-1">{item.label}</span>
                  {item.shortcut && (
                    <kbd className="text-xs text-secondary-400 dark:text-secondary-500 px-1.5 py-0.5 rounded bg-secondary-100 dark:bg-secondary-800">
                      {item.shortcut}
                    </kbd>
                  )}
                  <ChevronDownIcon className="h-4 w-4 text-secondary-400" aria-hidden="true" />
                </div>
              )}
            </Menu.Item>
          );
        }

        return (
          <Menu.Item key={item.id} disabled={item.disabled}>
            {({ active }) => (
              <button
                onClick={() => handleSelect(item)}
                className={cn(
                  'w-full flex items-center gap-2 px-3 py-2 text-sm transition-colors text-left',
                  'focus:outline-none focus:bg-secondary-100 dark:focus:bg-secondary-800',
                  item.danger ? 'text-red-600 dark:text-red-400' : 'text-secondary-700 dark:text-secondary-300',
                  item.disabled ? 'opacity-50 cursor-not-allowed' : '',
                  active ? 'bg-secondary-100 dark:bg-secondary-800' : ''
                )}
                disabled={item.disabled}
              >
                {item.icon && <span className="flex-shrink-0 h-5 w-5">{item.icon}</span>}
                <span className="flex-1">{item.label}</span>
                {item.shortcut && (
                  <kbd className="text-xs text-secondary-400 dark:text-secondary-500 px-1.5 py-0.5 rounded bg-secondary-100 dark:bg-secondary-800">
                    {item.shortcut}
                  </kbd>
                )}
                {selectedId === item.id && <CheckIcon className="h-4 w-4 text-primary-600 dark:text-primary-400" />}
              </button>
            )}
          </Menu.Item>
        );
      })}
    </>
  );

  return (
    <Menu as="div" className={cn('relative inline-block', className)}>
      <Menu.Button as={Fragment}>{trigger}</Menu.Button>
      <Transition
        as={Fragment}
        enter="transition ease-out duration-100"
        enterFrom="transform opacity-0 scale-95"
        enterTo="transform opacity-100 scale-100"
        leave="transition ease-in duration-75"
        leaveFrom="transform opacity-100 scale-100"
        leaveTo="transform opacity-0 scale-95"
      >
        <Menu.Items
          ref={menuRef}
          className={cn(
            'absolute z-50 min-w-[160px] origin-top-right rounded-xl bg-white shadow-lg ring-1 ring-secondary-200',
            'dark:bg-secondary-800 dark:ring-secondary-700',
            'focus:outline-none',
            'py-1',
            align === 'right' ? 'right-0' : 'left-0'
          )}
        >
          {renderItems(items)}
        </Menu.Items>
      </Transition>
    </Menu>
  );
}

export interface UserMenuProps {
  user: {
    name: string;
    email: string;
    avatar?: string;
    role: string;
  };
  items: DropdownItem[];
  onSignOut: () => void;
}

export function UserMenu({ user, items, onSignOut }: UserMenuProps) {
  return (
    <Dropdown
      trigger={
        <Button variant="ghost" size="sm" className="gap-2">
          <span className="hidden sm:block text-left">
            <p className="text-sm font-medium text-secondary-900 dark:text-secondary-100">{user.name}</p>
            <p className="text-xs text-secondary-500 dark:text-secondary-400">{user.email}</p>
          </span>
          <ChevronDownIcon className="h-4 w-4 text-secondary-400" />
        </Button>
      }
      items={[
        { id: 'divider-1', label: '', divider: true },
        { id: 'profile', label: 'Profile', icon: <UserIcon className="h-4 w-4" /> },
        { id: 'settings', label: 'Settings', icon: <SettingsIcon className="h-4 w-4" /> },
        { id: 'divider-2', label: '', divider: true },
        { id: 'signout', label: 'Sign out', icon: <LogOutIcon className="h-4 w-4" />, danger: true },
      ]}
      align="right"
    />
  );
}

import { UserIcon, SettingsIcon, LogOutIcon } from 'lucide-react';

export interface SelectDropdownProps {
  options: { value: string; label: string; icon?: React.ReactNode; disabled?: boolean }[];
  value: string;
  onChange: (value: string) => void;
  placeholder?: string;
  className?: string;
  disabled?: boolean;
}

export function SelectDropdown({ options, value, onChange, placeholder, className, disabled }: SelectDropdownProps) {
  const selectedOption = options.find((opt) => opt.value === value);

  return (
    <Dropdown
      trigger={
        <Button
          variant="outline"
          className={cn('w-full justify-between', className)}
          disabled={disabled}
        >
          <span className={cn('truncate', !selectedOption && 'text-secondary-400')}>
            {selectedOption?.label || placeholder}
          </span>
          <ChevronDownIcon className="h-4 w-4 text-secondary-400 flex-shrink-0 ml-2" />
        </Button>
      }
      items={options.map((opt) => ({
        id: opt.value,
        label: opt.label,
        icon: opt.icon,
        disabled: opt.disabled,
      }))}
      onSelect={onChange}
    />
  );
}

import { Menu as HeadlessMenu } from '@headlessui/react';