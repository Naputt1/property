import { GET } from "@/services/axios";
import {
  redirect,
  type LinkProps,
  type UseNavigateResult,
} from "@tanstack/react-router";

function genUrl(url: string, param?: { [key: string]: any }) {
  if (param == null) {
    return url;
  }

  const { query, ...restParam } = param;
  let finalUrl = url;

  for (const key in restParam) {
    const value = restParam[key];
    const placeholder = new RegExp(`(=?)\\$${key}\\b`, "g");

    if (Array.isArray(value)) {
      const arrayString = value.map((v) => `${key}[]=${v}`).join("&");
      finalUrl = finalUrl.replace(placeholder, () => arrayString);
    } else {
      finalUrl = finalUrl.replace(placeholder, (_match, equals) => {
        if (value == null) {
          return equals ? "" : "";
        }
        return equals ? `=${value}` : `${value}`;
      });
    }
  }

  // Clean up remaining placeholders and empty query params
  finalUrl = finalUrl.replace(/[&?][^&?]+=\$[^&?]+/g, "");
  finalUrl = finalUrl.replace(/\/\$[^/?]+/g, "");
  finalUrl = finalUrl.replace(/\?&/g, "?");
  finalUrl = finalUrl.replace(/[&?]$/g, "");

  if (query) {
    const queryParts: Array<string> = [];
    for (const key in query) {
      const value = query[key];
      if (value == null || value === "") continue;
      if (Array.isArray(value)) {
        value.forEach((v) =>
          queryParts.push(`${key}[]=${encodeURIComponent(v)}`),
        );
      } else {
        queryParts.push(`${key}=${encodeURIComponent(value)}`);
      }
    }
    if (queryParts.length > 0) {
      const separator = finalUrl.includes("?") ? "&" : "?";
      finalUrl += separator + queryParts.join("&");
    }
  }

  return finalUrl;
}

type ExtractParams<TString extends string> =
  TString extends `${string}$${infer Param}&${infer Rest}`
    ? { [K in Param]?: any } & ExtractParams<Rest>
    : TString extends `${string}$${infer Param}/${infer Rest}`
      ? { [K in Param]?: any } & ExtractParams<Rest>
      : TString extends `${string}$${infer Param}`
        ? { [K in Param]?: any }
        : {};

export interface IDefaultQueryBase {
  navigate?: UseNavigateResult<string>;
  redirect?: Pick<LinkProps, "to" | "params" | "search">;
}

export type IDefailtQuery<TString extends string> = IDefaultQueryBase &
  (keyof ExtractParams<TString> extends never
    ? {
        param?: { query?: any };
      }
    : {
        param: ExtractParams<TString> & { query?: any };
      });

type DefaultQueryFacOptionsBase<TString extends string = string> = {
  url: TString;
  redirect?: Pick<LinkProps, "to" | "params" | "search">;
};

export const defaultQueryFn =
  <TString extends string>({
    url,
    redirect: redirectOp,
  }: DefaultQueryFacOptionsBase<TString>) =>
  <T, TError = Error>() =>
  ({ param, navigate, redirect: redirectOptions }: IDefailtQuery<TString>) =>
  async () => {
    const res = await GET<T>(genUrl(url, param));

    if (res.status) {
      return res.data;
    }

    if (res.err.status === 404) {
      if ((redirectOp ?? redirectOptions) != null) {
        const _redirectOp = {
          ...redirectOp,
          param,
          ...redirectOptions,
        };

        if (navigate == null) {
          redirect(_redirectOp);
        } else {
          navigate(_redirectOp);
        }
      }
    }

    return Promise.reject(new Error(res.err.data) as TError);
  };
