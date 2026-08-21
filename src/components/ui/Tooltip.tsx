import { Fragment, useState } from 'react';
import { Popover, Transition } from '@headlessui/react';
import { cn } from '../utils/helpers';

export interface TooltipProps {
  content: React.ReactNode;
  children: React.ReactElement;
  placement?: 'top' | 'bottom' | 'left' | 'right' | 'top-start' | 'top-end' | 'bottom-start' | 'bottom-end' | 'left-start' | 'left-end' | 'right-start' | 'right-end';
  delay?: number;
  className?: string;
  contentClassName?: string;
}

export function Tooltip({
  content,
  children,
  placement = 'top',
  delay = 200,
  className,
  contentClassName,
}: TooltipProps) {
  const [open, setOpen] = useState(false);
  const timeoutRef = useState<ReturnType<typeof setTimeout> | null>(null);

  const showTooltip = () => {
    if (timeoutRef.current) clearTimeout(timeoutRef.current);
    timeoutRef.current = setTimeout(() => setOpen(true), delay);
  };

  const hideTooltip = () => {
    if (timeoutRef.current) clearTimeout(timeoutRef.current);
    setOpen(false);
  };

  const placementClasses = {
    top: 'bottom-full left-1/2 -translate-x-1/2 mb-2',
    bottom: 'top-full left-1/2 -translate-x-1/2 mt-2',
    left: 'right-full top-1/2 -translate-y-1/2 mr-2',
    right: 'left-full top-1/2 -translate-y-1/2 ml-2',
    'top-start': 'bottom-full left-0 mb-2',
    'top-end': 'bottom-full right-0 mb-2',
    'bottom-start': 'top-full left-0 mt-2',
    'bottom-end': 'top-full right-0 mt-2',
    'left-start': 'right-full top-0 mr-2',
    'left-end': 'right-full bottom-0 mr-2',
    'right-start': 'left-full top-0 ml-2',
    'right-end': 'left-full bottom-0 ml-2',
  };

  const arrowPlacementClasses = {
    top: 'top-full left-1/2 -translate-x-1/2 border-t-secondary-900',
    bottom: 'bottom-full left-1/2 -translate-x-1/2 border-b-secondary-900',
    left: 'left-full top-1/2 -translate-y-1/2 border-l-secondary-900',
    right: 'right-full top-1/2 -translate-y-1/2 border-r-secondary-900',
    'top-start': 'top-full left-4 border-t-secondary-900',
    'top-end': 'top-full right-4 border-t-secondary-900',
    'bottom-start': 'bottom-full left-4 border-b-secondary-900',
    'bottom-end': 'bottom-full right-4 border-b-secondary-900',
    'left-start': 'left-full top-4 border-l-secondary-900',
    'left-end': 'left-full bottom-4 border-l-secondary-900',
    'right-start': 'right-full top-4 border-r-secondary-900',
    'right-end': 'right-full bottom-4 border-r-secondary-900',
  };

  const basePlacement = placement.split('-')[0];

  return (
    <Popover open={open} onClose={hideTooltip}>
      <Popover.Button
        as={Fragment}
        onMouseEnter={showTooltip}
        onMouseLeave={hideTooltip}
        onFocus={showTooltip}
        onBlur={hideTooltip}
      >
        {React.cloneElement(children as React.ReactElement, {
          onMouseEnter: showTooltip,
          onMouseLeave: hideTooltip,
          onFocus: showTooltip,
          onBlur: hideTooltip,
        })}
      </Popover.Button>

      <Transition
        as={Fragment}
        appear
        show={open}
        enter="transition ease-out duration-100"
        enterFrom="transform opacity-0 scale-95"
        enterTo="transform opacity-100 scale-100"
        leave="transition ease-in duration-75"
        leaveFrom="transform opacity-100 scale-100"
        leaveTo="transform opacity-0 scale-95"
      >
        <Popover.Panel
          className={cn(
            'absolute z-50 px-3 py-2 text-xs font-medium text-white bg-secondary-900 rounded-lg shadow-lg',
            'dark:bg-secondary-50 dark:text-secondary-900',
            'whitespace-nowrap max-w-xs',
            placementClasses[placement],
            contentClassName
          )}
        >
          <div className="relative">
            {content}
            <div
              className={cn(
                'absolute w-0 h-0 border-4 border-transparent',
                arrowPlacementClasses[basePlacement]
              )}
              aria-hidden="true"
            />
          </div>
        </Popover.Panel>
      </Transition>
    </Popover>
  );
}

export interface TooltipTriggerProps {
  children: React.ReactNode;
  'aria-label'?: string;
  'aria-describedby'?: string;
}

export function TooltipTrigger({ children, ...props }: TooltipTriggerProps) {
  return <span {...props}>{children}</span>;
}

export interface HoverCardProps {
  trigger: React.ReactNode;
  content: React.ReactNode;
  placement?: 'top' | 'bottom' | 'left' | 'right';
  delay?: number;
  className?: string;
  contentClassName?: string;
}

export function HoverCard({
  trigger,
  content,
  placement = 'bottom',
  delay = 300,
  className,
  contentClassName,
}: HoverCardProps) {
  const [open, setOpen] = useState(false);
  const timeoutRef = useState<ReturnType<typeof setTimeout> | null>(null);

  const showCard = () => {
    if (timeoutRef.current) clearTimeout(timeoutRef.current);
    timeoutRef.current = setTimeout(() => setOpen(true), delay);
  };

  const hideCard = () => {
    if (timeoutRef.current) clearTimeout(timeoutRef.current);
    setOpen(false);
  };

  const placementClasses = {
    top: 'bottom-full left-1/2 -translate-x-1/2 mb-2',
    bottom: 'top-full left-1/2 -translate-x-1/2 mt-2',
    left: 'right-full top-1/2 -translate-y-1/2 mr-2',
    right: 'left-full top-1/2 -translate-y-1/2 ml-2',
  };

  return (
    <div className={cn('relative inline-block', className)}>
      <div
        onMouseEnter={showCard}
        onMouseLeave={hideCard}
        onFocus={showCard}
        onBlur={hideCard}
      >
        {trigger}
      </div>

      {open && (
        <div
          className={cn(
            'absolute z-50 w-72 rounded-lg bg-white shadow-lg ring-1 ring-secondary-200 p-4',
            'dark:bg-secondary-800 dark:ring-secondary-700',
            placementClasses[placement]
          )}
          onMouseEnter={showCard}
          onMouseLeave={hideCard}
        >
          {content}
        </div>
      )}
    </div>
  );
}