import { forwardRef, HTMLAttributes, useMemo, useState } from 'react';
import { cn } from '../utils/helpers';
import { ChevronUpIcon, ChevronDownIcon, ChevronUpDownIcon } from 'lucide-react';
import { Button } from './Button';
import { Input } from './Input';
import { Select } from './Input';

export interface Column<T> {
  key: string;
  header: string;
  accessor: keyof T | ((row: T) => React.ReactNode);
  width?: string;
  minWidth?: string;
  maxWidth?: string;
  align?: 'left' | 'center' | 'right';
  sortable?: boolean;
  filterable?: boolean;
  filterOptions?: { value: string; label: string }[];
  render?: (value: unknown, row: T, index: number) => React.ReactNode;
  headerRender?: () => React.ReactNode;
  className?: string;
  headerClassName?: string;
  cellClassName?: string;
}

export interface TableProps<T> {
  columns: Column<T>[];
  data: T[];
  keyExtractor: (row: T) => string;
  loading?: boolean;
  emptyMessage?: string;
  striped?: boolean;
  hoverable?: boolean;
  bordered?: boolean;
  compact?: boolean;
  sortable?: boolean;
  filterable?: boolean;
  selectable?: boolean;
  selectedKeys?: Set<string>;
  onSelectionChange?: (keys: Set<string>) => void;
  onRowClick?: (row: T, index: number) => void;
  pagination?: {
    page: number;
    pageSize: number;
    total: number;
    onPageChange: (page: number) => void;
    onPageSizeChange: (pageSize: number) => void;
    pageSizeOptions?: number[];
  };
  className?: string;
}

export function Table<T>({
  columns,
  data,
  keyExtractor,
  loading = false,
  emptyMessage = 'No data available',
  striped = true,
  hoverable = true,
  bordered = false,
  compact = false,
  sortable = false,
  filterable = false,
  selectable = false,
  selectedKeys = new Set(),
  onSelectionChange,
  onRowClick,
  pagination,
  className,
}: TableProps<T>) {
  const [sortConfig, setSortConfig] = useState<{ key: string; direction: 'asc' | 'desc' } | null>(null);
  const [filters, setFilters] = useState<Record<string, string>>({});

  const sortedAndFilteredData = useMemo(() => {
    let result = [...data];

    // Apply filters
    if (filterable) {
      Object.entries(filters).forEach(([key, value]) => {
        if (value) {
          const column = columns.find((c) => c.key === key);
          if (column) {
            result = result.filter((row) => {
              const cellValue = typeof column.accessor === 'function' 
                ? column.accessor(row) 
                : String(row[column.accessor as keyof T]);
              return String(cellValue).toLowerCase().includes(value.toLowerCase());
            });
          }
        }
      });
    }

    // Apply sorting
    if (sortable && sortConfig) {
      const column = columns.find((c) => c.key === sortConfig.key);
      if (column) {
        result.sort((a, b) => {
          const aVal = typeof column.accessor === 'function' ? column.accessor(a) : a[column.accessor as keyof T];
          const bVal = typeof column.accessor === 'function' ? column.accessor(b) : b[column.accessor as keyof T];
          
          if (aVal < bVal) return sortConfig.direction === 'asc' ? -1 : 1;
          if (aVal > bVal) return sortConfig.direction === 'asc' ? 1 : -1;
          return 0;
        });
      }
    }

    return result;
  }, [data, columns, filters, sortConfig, filterable, sortable]);

  const handleSort = (key: string) => {
    const column = columns.find((c) => c.key === key);
    if (!column?.sortable) return;

    setSortConfig((prev) => {
      if (prev?.key === key) {
        return { key, direction: prev.direction === 'asc' ? 'desc' : 'asc' };
      }
      return { key, direction: 'asc' };
    });
  };

  const handleFilterChange = (key: string, value: string) => {
    setFilters((prev) => ({ ...prev, [key]: value }));
  };

  const handleSelectAll = () => {
    if (selectedKeys.size === sortedAndFilteredData.length) {
      onSelectionChange?.(new Set());
    } else {
      onSelectionChange?.(new Set(sortedAndFilteredData.map(keyExtractor)));
    }
  };

  const handleSelectRow = (key: string) => {
    const newSelection = new Set(selectedKeys);
    if (newSelection.has(key)) {
      newSelection.delete(key);
    } else {
      newSelection.add(key);
    }
    onSelectionChange?.(newSelection);
  };

  const allSelected = sortedAndFilteredData.length > 0 && 
    sortedAndFilteredData.every((row) => selectedKeys.has(keyExtractor(row)));
  const someSelected = sortedAndFilteredData.some((row) => selectedKeys.has(keyExtractor(row)));

  const displayedData = pagination 
    ? sortedAndFilteredData.slice((pagination.page - 1) * pagination.pageSize, pagination.page * pagination.pageSize)
    : sortedAndFilteredData;

  return (
    <div className={cn('overflow-x-auto rounded-xl border border-secondary-200 dark:border-secondary-700', className)}>
      <table className="w-full text-sm" role="grid">
        <thead className="bg-secondary-50 dark:bg-secondary-800/50">
          <tr>
            {selectable && (
              <th
                className={cn(
                  'px-4 py-3 text-left font-medium text-secondary-500',
                  'w-12',
                  compact && 'py-2'
                )}
                scope="col"
              >
                <input
                  type="checkbox"
                  checked={allSelected}
                  indeterminate={someSelected && !allSelected}
                  onChange={handleSelectAll}
                  className="h-4 w-4 rounded border-secondary-300 text-primary-600 focus:ring-2 focus:ring-primary-500/20"
                  aria-label="Select all rows"
                />
              </th>
            )}
            {columns.map((column) => (
              <th
                key={column.key}
                scope="col"
                style={{
                  width: column.width,
                  minWidth: column.minWidth,
                  maxWidth: column.maxWidth,
                  textAlign: column.align,
                }}
                className={cn(
                  'px-4 py-3 font-medium text-secondary-500',
                  'text-left',
                  compact && 'py-2',
                  column.sortable && 'cursor-pointer select-none hover:text-secondary-700 dark:hover:text-secondary-300',
                  column.headerClassName
                )}
                onClick={() => column.sortable && handleSort(column.key)}
              >
                <div className="flex items-center gap-1">
                  {column.headerRender ? column.headerRender() : column.header}
                  {column.sortable && (
                    <span className="flex items-center">
                      {sortConfig?.key === column.key ? (
                        sortConfig.direction === 'asc' ? (
                          <ChevronUpIcon className="h-4 w-4 text-primary-600" />
                        ) : (
                          <ChevronDownIcon className="h-4 w-4 text-primary-600" />
                        )
                      ) : (
                        <ChevronUpDownIcon className="h-4 w-4 text-secondary-400" />
                      )}
                    </span>
                  )}
                </div>
                {filterable && column.filterable && (
                  <div className="mt-2">
                    {column.filterOptions ? (
                      <Select
                        value={filters[column.key] || ''}
                        onChange={(e) => handleFilterChange(column.key, e.target.value)}
                        options={[{ value: '', label: 'All' }, ...column.filterOptions]}
                        placeholder="Filter"
                        className="w-full"
                      />
                    ) : (
                      <Input
                        type="text"
                        value={filters[column.key] || ''}
                        onChange={(e) => handleFilterChange(column.key, e.target.value)}
                        placeholder="Filter..."
                        className="w-full"
                      />
                    )}
                  </div>
                )}
              </th>
            ))}
          </tr>
        </thead>
        <tbody className="divide-y divide-secondary-200 dark:divide-secondary-700">
          {loading ? (
            <tr>
              <td colSpan={columns.length + (selectable ? 1 : 0)} className="px-4 py-12 text-center text-secondary-500">
                <div className="flex items-center justify-center gap-2">
                  <svg className="animate-spin h-5 w-5 text-primary-600" viewBox="0 0 24 24">
                    <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" fill="none" />
                    <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
                  </svg>
                  Loading...
                </div>
              </td>
            </tr>
          ) : displayedData.length === 0 ? (
            <tr>
              <td colSpan={columns.length + (selectable ? 1 : 0)} className="px-4 py-12 text-center text-secondary-500">
                {emptyMessage}
              </td>
            </tr>
          ) : (
            displayedData.map((row, rowIndex) => {
              const rowKey = keyExtractor(row);
              const isSelected = selectedKeys.has(rowKey);
              
              return (
                <tr
                  key={rowKey}
                  className={cn(
                    'transition-colors',
                    striped && rowIndex % 2 === 1 && 'bg-secondary-50/50 dark:bg-secondary-800/50',
                    hoverable && 'hover:bg-secondary-100/50 dark:hover:bg-secondary-700/50',
                    isSelected && 'bg-primary-50 dark:bg-primary-900/20',
                    onRowClick && 'cursor-pointer'
                  )}
                  onClick={() => onRowClick?.(row, rowIndex)}
                >
                  {selectable && (
                    <td className={cn('px-4 py-3', compact && 'py-2')}>
                      <input
                        type="checkbox"
                        checked={isSelected}
                        onChange={() => handleSelectRow(rowKey)}
                        onClick={(e) => e.stopPropagation()}
                        className="h-4 w-4 rounded border-secondary-300 text-primary-600 focus:ring-2 focus:ring-primary-500/20"
                        aria-label="Select row"
                      />
                    </td>
                  )}
                  {columns.map((column) => {
                    const value = typeof column.accessor === 'function' 
                      ? column.accessor(row) 
                      : row[column.accessor as keyof T];
                    
                    return (
                      <td
                        key={column.key}
                        style={{ textAlign: column.align }}
                        className={cn(
                          'px-4 py-3 text-secondary-900 dark:text-secondary-100',
                          compact && 'py-2',
                          column.cellClassName
                        )}
                      >
                        {column.render 
                          ? column.render(value, row, rowIndex)
                          : typeof value === 'object' && value !== null
                            ? JSON.stringify(value)
                            : String(value ?? '')
                        }
                      </td>
                    );
                  })}
                </tr>
              );
            })
          )}
        </tbody>
      </table>
      
      {pagination && (
        <div className="border-t border-secondary-200 p-4 dark:border-secondary-700 flex items-center justify-between gap-4">
          <div className="flex items-center gap-2">
            <span className="text-sm text-secondary-600 dark:text-secondary-400">
              Showing {(pagination.page - 1) * pagination.pageSize + 1} to {Math.min(pagination.page * pagination.pageSize, pagination.total)} of {pagination.total}
            </span>
            <Select
              value={String(pagination.pageSize)}
              onChange={(e) => pagination.onPageSizeChange(Number(e.target.value))}
              options={pagination.pageSizeOptions?.map((size) => ({ value: String(size), label: String(size) })) || [
                { value: '10', label: '10' },
                { value: '25', label: '25' },
                { value: '50', label: '50' },
                { value: '100', label: '100' },
              ]}
              className="w-auto"
            />
          </div>
          <div className="flex items-center gap-2">
            <Button
              variant="outline"
              size="sm"
              onClick={() => pagination.onPageChange(pagination.page - 1)}
              disabled={pagination.page === 1}
              aria-label="Previous page"
            >
              <ChevronUpIcon className="h-4 w-4 rotate-90" />
            </Button>
            <span className="text-sm text-secondary-600 dark:text-secondary-400">
              Page {pagination.page} of {Math.ceil(pagination.total / pagination.pageSize)}
            </span>
            <Button
              variant="outline"
              size="sm"
              onClick={() => pagination.onPageChange(pagination.page + 1)}
              disabled={pagination.page >= Math.ceil(pagination.total / pagination.pageSize)}
              aria-label="Next page"
            >
              <ChevronDownIcon className="h-4 w-4 rotate-90" />
            </Button>
          </div>
        </div>
      )}
    </div>
  );
}

export interface DataGridProps<T> extends TableProps<T> {
  rowActions?: (row: T) => React.ReactNode;
  rowExpansion?: {
    render: (row: T) => React.ReactNode;
    expandedKeys?: Set<string>;
    onExpandChange?: (keys: Set<string>) => void;
  };
}

export function DataGrid<T>({
  rowActions,
  rowExpansion,
  ...tableProps
}: DataGridProps<T>) {
  const [expandedKeys, setExpandedKeys] = useState<Set<string>>(rowExpansion?.expandedKeys || new Set());

  const columnsWithActions = useMemo(() => {
    const cols = [...tableProps.columns];
    if (rowActions) {
      cols.push({
        key: 'actions',
        header: 'Actions',
        accessor: 'actions' as keyof T,
        width: '120px',
        render: (_, row) => rowActions(row),
        sortable: false,
        filterable: false,
      });
    }
    if (rowExpansion) {
      cols.unshift({
        key: 'expand',
        header: '',
        accessor: 'expand' as keyof T,
        width: '40px',
        render: (_, row, index) => {
          const key = tableProps.keyExtractor(row);
          const expanded = expandedKeys.has(key);
          return (
            <Button
              variant="ghost"
              size="sm"
              onClick={(e) => {
                e.stopPropagation();
                const newKeys = new Set(expandedKeys);
                if (expanded) newKeys.delete(key);
                else newKeys.add(key);
                setExpandedKeys(newKeys);
                rowExpansion.onExpandChange?.(newKeys);
              }}
              aria-label={expanded ? 'Collapse row' : 'Expand row'}
              aria-expanded={expanded}
            >
              {expanded ? <ChevronUpIcon className="h-4 w-4" /> : <ChevronDownIcon className="h-4 w-4" />}
            </Button>
          );
        },
        sortable: false,
        filterable: false,
      });
    }
    return cols;
  }, [tableProps.columns, rowActions, rowExpansion, expandedKeys, tableProps.keyExtractor]);

  return (
    <>
      <Table
        {...tableProps}
        columns={columnsWithActions}
        onRowClick={rowExpansion ? undefined : tableProps.onRowClick}
      />
      {rowExpansion && (
        <style jsx>{`
          .expanded-row td {
            padding: 0 !important;
            border: none !important;
          }
        `}</style>
      )}
    </>
  );
}