export interface ApiResponse<T = any> {
  data?: T;
  status?: number;
  message?: string;
}

export interface ApiError {
  message?: string;
  status?: number;
  code?: string;
}

export interface RequestConfig {
  url: string;
  params?: Record<string, any>;
  data?: any;
  headers?: Record<string, string>;
}

export interface PaginatedResponse<T> {
  items: T[];
  total: number;
  page: number;
  limit: number;
}