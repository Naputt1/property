import axios from "axios";
import { toast } from "sonner";
import { useAuthStore } from "@/store/auth";
import type {
  AxiosError,
  AxiosHeaderValue,
  AxiosRequestConfig,
  AxiosResponse,
} from "axios";

export const serverRoute = "/api";
export const axiosConfig = {
  baseURL: serverRoute,
};

export interface IApiResponseDataBase {
  status: boolean;
}

export interface IErrorResponse {
  status: false;
  error: string;
}

export type ApiError<T = any> = {
  status?: number;
  data?: T;
};

type ApiResponseErr<TError = any, TMeta = any> = {
  status?: number;
  data?: TError;
  meta?: TMeta;
};

export type ApiHeader = {
  [key: string]: AxiosHeaderValue;
};

export type ApiResponseAxios<T = any, TError = any> =
  | {
      status: true;
      data: T;
      header?: ApiHeader;
    }
  | {
      status: false;
      data: TError;
      header?: ApiHeader;
    };

export type ApiResponseSuccess<T = any> = {
  status: true;
  data: T;
  header?: ApiHeader;
};

export type ApiResponseError<TError = any, TMeta = any> = {
  status: false;
  err: ApiResponseErr<TError, TMeta>;
  header?: ApiHeader;
};

export type ApiResponse<T = any, TError = any, TMeta = any> =
  | ApiResponseSuccess<T>
  | ApiResponseError<TError, TMeta>;

export function initAxios() {
  const { clearAuth } = useAuthStore.getState();

  const axiosInstance = axios.create({
    baseURL: axiosConfig.baseURL,
    timeout: 120000,
    headers: {
      "Content-Type": "application/json",
    },
    withCredentials: true,
  });

  axiosInstance.interceptors.response.use(
    (response) => {
      return response;
    },
    (error: AxiosError) => {
      if (error.code === "ERR_NETWORK") {
        toast.error("Network Error: Please try again later.");
        return Promise.reject(error);
      }

      if (error.response) {
        if (error.response.status === 401) {
          clearAuth();
          // Handle unauthorized, maybe redirect to login
          if (window.location.pathname !== "/login") {
            window.location.href = "/login";
          }
        }

        if (error.response.status >= 500 && error.response.status < 600) {
          const err: string | undefined = (error.response.data as any)?.error;
          if (typeof err === "string") {
            toast.error(err);
          } else {
            toast.error("Server Error: Please try again later.");
          }
        }
      }

      return Promise.reject(error);
    },
  );

  return axiosInstance;
}

function geHeader(header?: AxiosResponse["headers"]): ApiHeader {
  if (header == null) {
    return {};
  }

  return header as ApiHeader;
}

function responseCatch<TError, TMeta>(
  err: AxiosError<any, any>,
): ApiResponse<any, TError, TMeta> {
  return {
    status: false,
    err: {
      status: err.response?.status,
      data: (err.response?.data?.error || err.message) as TError,
      meta: err.response?.data,
    },
    header: geHeader(err.response?.headers),
  };
}

export async function GET<T = any, TMeta = any>(
  url: string,
  config?: AxiosRequestConfig,
): Promise<ApiResponse<T>> {
  const api = initAxios();
  try {
    const res = await api.get<ApiResponse<T>>(url, config);
    if (res.data.status === false) {
      return {
        status: false,
        err: {
          status: res.status,
          data: res.data,
        },
        header: geHeader(res.headers),
      };
    }

    return {
      // eslint-disable-next-line @typescript-eslint/ban-ts-comment
      // @ts-ignore
      status: true,
      ...res.data,
      header: geHeader(res.headers),
    };
  } catch (err) {
    return responseCatch<T, TMeta>(err as AxiosError);
  }
}

export async function POST<T = any, TError = any, TData = any, TMeta = any>(
  url: string,
  data?: TData,
  config?: AxiosRequestConfig,
): Promise<ApiResponse<T, TError>> {
  try {
    const api = initAxios();
    const res = await api.post<ApiResponseAxios<T, TError>>(url, data, config);
    if (res.data.status === false) {
      return {
        status: false,
        err: {
          status: res.status,
          data: res.data.data,
        },
        header: geHeader(res.headers),
      };
    }

    return {
      // eslint-disable-next-line @typescript-eslint/ban-ts-comment
      // @ts-ignore
      status: true,
      ...res.data,
      header: geHeader(res.headers),
    };
  } catch (err) {
    console.log("axios error", err);
    if (axios.isAxiosError(err)) {
      return responseCatch<TError, TMeta>(err);
    }

    return responseCatch<TError, TMeta>(err as AxiosError);
  }
}

export async function DELETE<T = any, TMeta = any>(
  url: string,
  config?: AxiosRequestConfig,
): Promise<ApiResponse<T>> {
  const api = initAxios();
  try {
    const res = await api.delete<ApiResponse<T>>(url, config);
    if (res.data.status === false) {
      return {
        status: false,
        err: {
          status: res.status,
          data: res.data,
        },
        header: geHeader(res.headers),
      };
    }

    return {
      // eslint-disable-next-line @typescript-eslint/ban-ts-comment
      // @ts-ignore
      status: true,
      ...res.data,
      header: geHeader(res.headers),
    };
  } catch (err) {
    return responseCatch<T, TMeta>(err as AxiosError);
  }
}

export default initAxios;
