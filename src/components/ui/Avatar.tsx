import { forwardRef, HTMLAttributes, useMemo } from 'react';
import { cn, getInitials } from '../utils/helpers';

export interface AvatarProps extends HTMLAttributes<HTMLDivElement> {
  src?: string | null;
  alt?: string;
  name?: string;
  size?: 'xs' | 'sm' | 'md' | 'lg' | 'xl' | '2xl';
  shape?: 'circle' | 'square';
  status?: 'online' | 'offline' | 'busy' | 'away';
  statusPosition?: 'bottom-right' | 'bottom-left' | 'top-right' | 'top-left';
  border?: boolean;
}

export const Avatar = forwardRef<HTMLDivElement, AvatarProps>(
  (
    {
      src,
      alt,
      name,
      size = 'md',
      shape = 'circle',
      status,
      statusPosition = 'bottom-right',
      border = false,
      className,
      ...props
    },
    ref
  ) => {
    const sizeClasses = {
      xs: 'h-6 w-6 text-xs',
      sm: 'h-8 w-8 text-sm',
      md: 'h-10 w-10 text-base',
      lg: 'h-12 w-12 text-lg',
      xl: 'h-16 w-16 text-xl',
      '2xl': 'h-20 w-20 text-2xl',
    };

    const shapeClasses = {
      circle: 'rounded-full',
      square: 'rounded-xl',
    };

    const statusSizeClasses = {
      xs: 'h-1.5 w-1.5',
      sm: 'h-2 w-2',
      md: 'h-2.5 w-2.5',
      lg: 'h-3 w-3',
      xl: 'h-3.5 w-3.5',
      '2xl': 'h-4 w-4',
    };

    const statusPositionClasses = {
      'bottom-right': 'bottom-0 right-0',
      'bottom-left': 'bottom-0 left-0',
      'top-right': 'top-0 right-0',
      'top-left': 'top-0 left-0',
    };

    const statusColors = {
      online: 'bg-green-500',
      offline: 'bg-secondary-400',
      busy: 'bg-red-500',
      away: 'bg-yellow-500',
    };

    const initials = useMemo(() => (name ? getInitials(name) : '?'), [name]);
    const backgroundColor = useMemo(() => {
      if (!name) return 'bg-secondary-300 dark:bg-secondary-600';
      let hash = 0;
      for (let i = 0; i < name.length; i++) {
        hash = name.charCodeAt(i) + ((hash << 5) - hash);
      }
      const hue = hash % 360;
      return `hsl(${hue}, 60%, 45%)`;
    }, [name]);

    return (
      <div
        ref={ref}
        className={cn('relative inline-flex items-center justify-center font-medium text-white bg-secondary-100 overflow-hidden', sizeClasses[size], shapeClasses[shape], border && 'ring-2 ring-white dark:ring-secondary-800', className)}
        {...props}
      >
        {src ? (
          <img
            src={src}
            alt={alt || name || 'Avatar'}
            className="h-full w-full object-cover"
            loading="lazy"
          />
        ) : (
          <div
            className={cn('h-full w-full flex items-center justify-center', backgroundColor)}
            style={{ backgroundColor: backgroundColor.replace('bg-', '') }}
          >
            {initials}
          </div>
        )}
        {status && (
          <span
            className={cn(
              'absolute rounded-full border-2 border-white dark:border-secondary-800',
              statusColors[status],
              statusSizeClasses[size],
              statusPositionClasses[statusPosition]
            )}
            aria-label={`${status} status`}
          />
        )}
      </div>
    );
  }
);

Avatar.displayName = 'Avatar';

export interface AvatarGroupProps extends HTMLAttributes<HTMLDivElement> {
  max?: number;
  size?: 'xs' | 'sm' | 'md' | 'lg' | 'xl' | '2xl';
  spacing?: number;
}

export const AvatarGroup = forwardRef<HTMLDivElement, AvatarGroupProps>(
  ({ children, max = 5, size = 'md', spacing = -8, className, ...props }, ref) => {
    const childArray = Array.isArray(children) ? children : [children];
    const validChildren = childArray.filter((child) => child != null);
    const visibleChildren = validChildren.slice(0, max);
    const remainingCount = validChildren.length - max;

    return (
      <div ref={ref} className={cn('flex items-center', className)} {...props}>
        <div className="flex" role="group" aria-label={`${validChildren.length} users`}>
          {visibleChildren.map((child, index) =>
            React.cloneElement(child as React.ReactElement, {
              key: (child as React.ReactElement).key || index,
              size,
              className: cn(
                index > 0 && `ml-${Math.abs(spacing)}`,
                'border-2 border-white dark:border-secondary-800'
              ),
            })
          )}
          {remainingCount > 0 && (
            <div
              className={cn(
                'ml-2 flex items-center justify-center font-medium text-secondary-600 bg-secondary-100 border-2 border-white dark:border-secondary-800 dark:bg-secondary-700 dark:text-secondary-300',
                sizeClasses[size],
                shapeClasses.circle
              )}
              aria-label={`${remainingCount} more users`}
            >
              +{remainingCount}
            </div>
          )}
        </div>
      </div>
    );
  }
);

AvatarGroup.displayName = 'AvatarGroup';

import React from 'react';

const sizeClasses = {
  xs: 'h-6 w-6 text-xs',
  sm: 'h-8 w-8 text-sm',
  md: 'h-10 w-10 text-base',
  lg: 'h-12 w-12 text-lg',
  xl: 'h-16 w-16 text-xl',
  '2xl': 'h-20 w-20 text-2xl',
};

const shapeClasses = {
  circle: 'rounded-full',
  square: 'rounded-xl',
};