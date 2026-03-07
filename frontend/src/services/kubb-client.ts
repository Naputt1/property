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

  const { status, message, error, ...payload } = response.data;

  if (status === false) {
    // If the backend says status is false, treat it as an error
    const errorMessage = error || message || "Unknown API error";
    const axiosError = new Error(errorMessage) as any;
    axiosError.response = response;
    axiosError.isAxiosError = true;
    axiosError.status = response.status;
    axiosError.code = "API_ERROR";
    throw axiosError;
  }

  // Smart unwrapping
  // If 'data' is the only field left in payload, return just the data
  // Otherwise return the whole payload (e.g. { data, total })
  const keys = Object.keys(payload);
  let resultData = payload;
  
  if (keys.length === 1 && keys[0] === 'data') {
    resultData = payload.data as any;
  } else if (keys.length === 0 && 'data' in response.data === false) {
    // Fallback if the body didn't even have data/status (shouldn't happen with our API)
    resultData = response.data as any;
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
