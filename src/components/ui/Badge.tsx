import { forwardRef, HTMLAttributes } from 'react';
import { cn } from '../utils/helpers';

export interface BadgeProps extends HTMLAttributes<HTMLSpanElement> {
  variant?: 'primary' | 'success' | 'warning' | 'danger' | 'info' | 'secondary' | 'outline';
  size?: 'sm' | 'md' | 'lg';
  dot?: boolean;
  dotColor?: string;
}

export const Badge = forwardRef<HTMLSpanElement, BadgeProps>(
  ({ children, variant = 'primary', size = 'md', dot = false, dotColor, className, ...props }, ref) => {
    const variantClasses = {
      primary: 'bg-primary-100 text-primary-800 dark:bg-primary-900/30 dark:text-primary-300',
      success: 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-300',
      warning: 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900/30 dark:text-yellow-300',
      danger: 'bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-300',
      info: 'bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-300',
      secondary: 'bg-secondary-100 text-secondary-800 dark:bg-secondary-700 dark:text-secondary-300',
      outline: 'border border-secondary-300 bg-transparent text-secondary-700 dark:border-secondary-600 dark:text-secondary-300',
    };

    const sizeClasses = {
      sm: 'px-2 py-0.5 text-xs',
      md: 'px-2.5 py-0.5 text-xs',
      lg: 'px-3 py-1 text-sm',
    };

    return (
      <span
        ref={ref}
        className={cn(
          'inline-flex items-center gap-1.5 rounded-full font-medium',
          variantClasses[variant],
          sizeClasses[size],
          className
        )}
        {...props}
      >
        {dot && <span className={cn('relative flex h-1.5 w-1.5 rounded-full', dotColor && `bg-[${dotColor}]`)} />}
        {children}
      </span>
    );
  }
);

Badge.displayName = 'Badge';

export interface StatusBadgeProps {
  status: 'active' | 'inactive' | 'pending' | 'completed' | 'failed' | 'warning' | 'draft' | 'published' | 'archived';
  size?: 'sm' | 'md' | 'lg';
  showDot?: boolean;
}

export const StatusBadge = ({ status, size = 'md', showDot = true }: StatusBadgeProps) => {
  const statusConfig: Record<string, { label: string; variant: BadgeProps['variant']; dotColor?: string }> = {
    active: { label: 'Active', variant: 'success', dotColor: '#22c55e' },
    inactive: { label: 'Inactive', variant: 'secondary', dotColor: '#9ca3af' },
    pending: { label: 'Pending', variant: 'warning', dotColor: '#f59e0b' },
    completed: { label: 'Completed', variant: 'primary', dotColor: '#3b82f6' },
    failed: { label: 'Failed', variant: 'danger', dotColor: '#ef4444' },
    warning: { label: 'Warning', variant: 'warning', dotColor: '#f59e0b' },
    draft: { label: 'Draft', variant: 'outline', dotColor: '#6b7280' },
    published: { label: 'Published', variant: 'success', dotColor: '#22c55e' },
    archived: { label: 'Archived', variant: 'secondary', dotColor: '#9ca3af' },
  };

  const config = statusConfig[status] || { label: status, variant: 'secondary' };

  return (
    <Badge variant={config.variant} size={size} dot={showDot} dotColor={config.dotColor}>
      {config.label}
    </Badge>
  );
};

export interface MetricBadgeProps {
  value: number | string;
  trend?: 'up' | 'down' | 'neutral';
  trendValue?: number;
  prefix?: string;
  suffix?: string;
  variant?: 'primary' | 'success' | 'warning' | 'danger' | 'info';
}

export const MetricBadge = ({
  value,
  trend,
  trendValue,
  prefix = '',
  suffix = '',
  variant = 'primary',
}: MetricBadgeProps) => {
  const trendIcons = {
    up: '↑',
    down: '↓',
    neutral: '→',
  };

  const trendColors = {
    up: 'text-green-600 dark:text-green-400',
    down: 'text-red-600 dark:text-red-400',
    neutral: 'text-secondary-500 dark:text-secondary-400',
  };

  return (
    <div className="flex items-end gap-1">
      <span className={cn('font-mono font-semibold', variant === 'primary' ? 'text-primary-700 dark:text-primary-300' : '')}>
        {prefix}{value}{suffix}
      </span>
      {trend && trend !== 'neutral' && (
        <span className={cn('text-xs font-medium', trendColors[trend])}>
          {trendIcons[trend]}
          {trendValue !== undefined && (
            <span className="ml-0.5">{trendValue >= 0 ? '+' : ''}{trendValue}%</span>
          )}
        </span>
      )}
    </div>
  );
};