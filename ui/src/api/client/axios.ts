import axios, { type AxiosInstance, type AxiosResponse, AxiosError } from 'axios';
import type { ApiResponse, ApiError, RequestConfig } from './types';

export class ApiClient {
  private instance: AxiosInstance;

  constructor(baseURL: string) {
    this.instance = axios.create({
      baseURL,
      timeout: 10000,
      headers: {
        'Content-Type': 'application/json',
      },
    });

    this.setupInterceptors();
  }

  private setupInterceptors(): void {
    this.instance.interceptors.request.use(
      (config) => {
        const token = localStorage.getItem('token');
        if (token) {
          config.headers.Authorization = `Bearer ${token}`;
        }
        return config;
      },
      (error) => Promise.reject(error)
    );

    this.instance.interceptors.response.use(
      (response) => response,
      (error: AxiosError) => {
        const apiError: ApiError = {
          message: error.message,
          status: error.response?.status,
          code: error.code,
        };
        return Promise.reject(apiError);
      }
    );
  }

  async request<T = any>(config: RequestConfig): Promise<ApiResponse<T>> {
    try {
      const response: AxiosResponse<T> = await this.instance({
        url: config.url,
        method: config.data ? 'POST' : 'GET',
        params: config.params,
        data: config.data,
        headers: config.headers,
      });

      return {
        data: response.data,
        status: response.status,
      };
    } catch (error) {
      throw error;
    }
  }

  async get<T = any>(url: string, params?: Record<string, any>): Promise<ApiResponse<T>> {
    return this.request<T>({ url, params });
  }

  async post<T = any>(url: string, data?: any): Promise<ApiResponse<T>> {
    return this.request<T>({ url, data });
  }

  async put<T = any>(url: string, data?: any): Promise<ApiResponse<T>> {
    return this.request<T>({ url, data, headers: { 'X-HTTP-Method-Override': 'PUT' } });
  }

  async patch<T = any>(url: string, data?: any): Promise<ApiResponse<T>> {
    return this.request<T>({ url, data, headers: { 'X-HTTP-Method-Override': 'PATCH' } });
  }

  async delete<T = any>(url: string): Promise<ApiResponse<T>> {
    return this.request<T>({ url, headers: { 'X-HTTP-Method-Override': 'DELETE' } });
  }
}
