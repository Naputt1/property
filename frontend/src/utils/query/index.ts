import { queryOptions } from "@tanstack/react-query";
import type {
  QueryFunction,
  UseMutationOptions,
  UseQueryOptions,
} from "@tanstack/react-query";

export interface QueryOptionsBase<T, TError = Error> extends Partial<
  Omit<UseQueryOptions<T, TError, T, Array<unknown>>, "queryFn">
> {
  queryKey?: Array<unknown>;
  queryFn?: QueryFunction<T>;
}

interface QueryOptionsFin<T, TError = Error> extends QueryOptionsBase<
  T,
  TError
> {
  queryKey: Array<unknown>;
}

export interface IQueryOptionsQueryFn<T, TArgs> extends Omit<
  QueryOptionsBase<T>,
  "queryFn"
> {
  queryFn: (args: TArgs) => QueryFunction<T>;
}

type GetQueryKey<TArgs extends {}> = (args: TArgs) => Array<unknown>;

type QueryConstructorProps<TQK, TOption> = {
  queryKey: TQK;
  options: TOption;
};

type OptionParam<TKeyArgs> = keyof TKeyArgs extends never
  ? {
      param?: undefined;
    }
  : {
      param: TKeyArgs;
    };

export type QueryWrapper<
  T = any,
  TKeyArgs extends object = object,
  TOptionsArgs extends object = object,
> = {
  getKey: (
    ...args: keyof TKeyArgs extends never ? [] : [args: TKeyArgs]
  ) => Array<unknown>;
  getOptions: <TData = T, TError = Error>(
    options: TOptionsArgs & OptionParam<TKeyArgs>,
    overrideOptions?: QueryOptionsBase<TData>,
  ) => UseQueryOptions<TData, TError, TData, Array<unknown>>;
};

export function createQueryWrapper<
  T = any,
  TKeyArgs extends object = object,
  TOptionsArgs extends object = object,
>({
  queryKey,
  options: { queryFn, ..._options },
}: QueryConstructorProps<
  GetQueryKey<TKeyArgs>,
  IQueryOptionsQueryFn<T, TOptionsArgs>
>) {
  return {
    getKey: (
      ...args: keyof TKeyArgs extends never ? [] : [args: TKeyArgs]
    ): Array<unknown> => {
      return queryKey(args[0] ?? ({} as TKeyArgs));
    },
    getOptions: (
      options: TOptionsArgs & OptionParam<TKeyArgs>,
      overrideOptions?: QueryOptionsBase<T>,
    ) => {
      const key = queryKey(options.param ?? ({} as TKeyArgs));
      const option = {
        ..._options,
        ...options,
        queryKey: key,
        queryFn: queryFn(options) as any as QueryFunction<T>,
        ...overrideOptions,
      };

      return queryOptions({
        ...option,
        ...overrideOptions,
      } as QueryOptionsFin<T, Error>);
    },
  };
}

export function createMutationWrapper<
  TData = unknown,
  TError = Error,
  TVariables = void,
  TContext = unknown,
>(options: UseMutationOptions<TData, TError, TVariables, TContext>) {
  return {
    getOptions: (
      overrideOptions?: UseMutationOptions<TData, TError, TVariables, TContext>,
    ) => ({
      ...options,
      ...overrideOptions,
    }),
  };
}
