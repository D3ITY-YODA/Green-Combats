import { Fragment } from 'react';
import { cn } from '../utils/helpers';
import { ChevronRightIcon, HomeIcon } from 'lucide-react';
import { Link, useLocation, useNavigationType } from 'react-router-dom';

export interface BreadcrumbItem {
  label: string;
  href?: string;
  icon?: React.ReactNode;
  current?: boolean;
}

export interface BreadcrumbProps {
  items: BreadcrumbItem[];
  separator?: React.ReactNode;
  className?: string;
  itemClassName?: string;
  separatorClassName?: string;
  maxItems?: number;
  showHome?: boolean;
  homeHref?: string;
}

export function Breadcrumb({
  items,
  separator = <ChevronRightIcon className="h-4 w-4 text-secondary-400" />,
  className,
  itemClassName,
  separatorClassName,
  maxItems = 5,
  showHome = true,
  homeHref = '/dashboard',
}: BreadcrumbProps) {
  const location = useLocation();
  const navigationType = useNavigationType();

  const displayItems = items.slice(-maxItems);
  const hasMore = items.length > maxItems;

  return (
    <nav
      aria-label="Breadcrumb"
      className={cn('flex items-center gap-1.5 flex-wrap', className)}
    >
      <ol className="flex items-center gap-1.5 flex-wrap" role="list">
        {showHome && (
          <li>
            <Link
              to={homeHref}
              className={cn(
                'flex items-center gap-1.5 px-2 py-1.5 text-sm text-secondary-500 hover:text-secondary-700',
                'dark:text-secondary-400 dark:hover:text-secondary-200',
                'rounded-lg transition-colors',
                itemClassName
              )}
              aria-label="Home"
            >
              <HomeIcon className="h-4 w-4" aria-hidden="true" />
              <span className="hidden sm:inline">Home</span>
            </Link>
            {items.length > 0 && (
              <span className={cn('flex items-center text-secondary-400', separatorClassName)} aria-hidden="true">
                {separator}
              </span>
            )}
          </li>
        )}

        {hasMore && (
          <li className="flex items-center">
            <span className="px-2 py-1.5 text-sm text-secondary-400" aria-hidden="true">...</span>
            <span className={cn('flex items-center text-secondary-400', separatorClassName)} aria-hidden="true">
              {separator}
            </span>
          </li>
        )}

        {displayItems.map((item, index) => (
          <li key={item.href || item.label} className="flex items-center">
            {item.href && !item.current ? (
              <Link
                to={item.href}
                className={cn(
                  'flex items-center gap-1.5 px-2 py-1.5 text-sm',
                  'text-secondary-500 hover:text-secondary-700',
                  'dark:text-secondary-400 dark:hover:text-secondary-200',
                  'rounded-lg transition-colors',
                  itemClassName
                )}
                aria-current={item.current ? 'page' : undefined}
              >
                {item.icon && <span className="flex-shrink-0" aria-hidden="true">{item.icon}</span>}
                <span>{item.label}</span>
              </Link>
            ) : (
              <span
                className={cn(
                  'flex items-center gap-1.5 px-2 py-1.5 text-sm font-medium',
                  'text-secondary-900 dark:text-secondary-100',
                  itemClassName
                )}
                aria-current={item.current ? 'page' : undefined}
              >
                {item.icon && <span className="flex-shrink-0" aria-hidden="true">{item.icon}</span>}
                <span>{item.label}</span>
              </span>
            )}
            {index < displayItems.length - 1 && (
              <span className={cn('flex items-center text-secondary-400', separatorClassName)} aria-hidden="true">
                {separator}
              </span>
            )}
          </li>
        ))}
      </ol>
    </nav>
  );
}

export interface BreadcrumbItemProps {
  item: BreadcrumbItem;
  separator?: React.ReactNode;
  isLast?: boolean;
  className?: string;
}

export function BreadcrumbItemComponent({ item, separator, isLast, className }: BreadcrumbItemProps) {
  return (
    <li className="flex items-center">
      {item.href && !item.current ? (
        <Link
          to={item.href}
          className={cn(
            'flex items-center gap-1.5 px-2 py-1.5 text-sm',
            'text-secondary-500 hover:text-secondary-700',
            'dark:text-secondary-400 dark:hover:text-secondary-200',
            'rounded-lg transition-colors',
            className
          )}
          aria-current={item.current ? 'page' : undefined}
        >
          {item.icon && <span className="flex-shrink-0" aria-hidden="true">{item.icon}</span>}
          <span>{item.label}</span>
        </Link>
      ) : (
        <span
          className={cn(
            'flex items-center gap-1.5 px-2 py-1.5 text-sm font-medium',
            'text-secondary-900 dark:text-secondary-100',
            className
          )}
          aria-current={item.current ? 'page' : undefined}
        >
          {item.icon && <span className="flex-shrink-0" aria-hidden="true">{item.icon}</span>}
          <span>{item.label}</span>
        </span>
      )}
      {!isLast && separator && (
        <span className="flex items-center text-secondary-400" aria-hidden="true">
          {separator}
        </span>
      )}
    </li>
  );
}

export function useBreadcrumbs(pathname: string, routeMap: Record<string, string>): BreadcrumbItem[] {
  const segments = pathname.split('/').filter(Boolean);
  const breadcrumbs: BreadcrumbItem[] = [];

  let currentPath = '';
  for (const segment of segments) {
    currentPath += `/${segment}`;
    const label = routeMap[currentPath] || segment.replace(/-/g, ' ').replace(/\b\w/g, (l) => l.toUpperCase());
    breadcrumbs.push({ label, href: currentPath });
  }

  if (breadcrumbs.length > 0) {
    breadcrumbs[breadcrumbs.length - 1].current = true;
    breadcrumbs[breadcrumbs.length - 1].href = undefined;
  }

  return breadcrumbs;
}