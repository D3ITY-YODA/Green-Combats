/**
 * Metadata included in paginated or list API responses.
 */
export interface ApiMeta {
  request_id: string;
  generated_at: string;
  pagination?: PaginationMeta;
}

export interface PaginationMeta {
  page: number;
  page_size: number;
  total: number;
  has_next: boolean;
  next_cursor?: string;
}

/**
 * Standard envelope for all successful API responses.
 * The `data` property is generic and will be populated with the specific resource type.
 */
export interface ApiResponse<T> {
  data: T;
  meta?: ApiMeta;
}

/**
 * Standard envelope for all API errors returned by the backend.
 */
export interface ApiError {
  error: {
    code: string;
    message: string;
    request_id: string;
    details?: Record<string, unknown>;
  };
}
