import { initAxios } from "./axios";
import axios from "axios";
import type { AxiosRequestConfig, AxiosResponse, AxiosError } from "axios";

const axiosInstance = initAxios();

/**
 * Standard backend response format
 */
interface BackendResponse<T = any> {
  status: boolean;
  data: T;
  message?: string;
  error?: string;
  total?: number;
  [key: string]: any;
}

/**
 * Custom client for Kubb that handles the "status" field and extract payload
 */
export async function client<TData = any, TError = any, TVariables = any>(
  config: AxiosRequestConfig<TVariables>,
): Promise<AxiosResponse<TData>> {
  const response = await axiosInstance.request<BackendResponse<TData>>(config);

  // If response.data is missing or not an object, return raw response
  if (!response.data || typeof response.data !== 'object') {
    return response as any;
  }

  const { status, data, message, error, ...rest } = response.data;

  if (status === false) {
    const errorMessage = error || message || "Unknown API error";
    const axiosError = new Error(errorMessage) as any;
    axiosError.response = response;
    axiosError.isAxiosError = true;
    axiosError.status = response.status;
    axiosError.code = "API_ERROR";
    throw axiosError;
  }

  // Robust unwrapping:
  // 1. If 'data' field exists, use it (even if it's null or empty array)
  // 2. If 'data' is missing but we have other fields in 'rest', use 'rest'
  // 3. Fallback to the whole response body (excluding status)
  let resultData: any;
  
  if (response.data && 'data' in response.data) {
    resultData = data;
  } else if (Object.keys(rest).length > 0) {
    resultData = rest;
  } else {
    resultData = response.data;
  }

  return {
    ...response,
    data: resultData as TData,
  };
}

export type Client = typeof client;
export type RequestConfig<TVariables = any> = AxiosRequestConfig<TVariables>;
export type ResponseConfig<TData = any> = AxiosResponse<TData>;
export type ResponseErrorConfig<TError = any> = AxiosError<TError>;

export default client;
