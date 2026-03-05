import type { IPropertyListRes } from "@/types/property";
import { createQueryWrapper } from "@/utils/query";
import { defaultQueryFn } from "@/utils/query/defaultQuery";

type PropertyListParams = { page: number; pageSize: number; search?: string };

export const propertyListQuery = createQueryWrapper<
  IPropertyListRes,
  PropertyListParams,
  { param: PropertyListParams }
>({
  queryKey: (params: PropertyListParams) => [
    "property",
    "list",
    JSON.stringify(params),
  ],
  options: {
    queryFn: defaultQueryFn({
      url: "/property?page=$page&pageSize=$pageSize&search=$search",
    })<IPropertyListRes>(),
    retry: false,
    staleTime: 5000,
  },
});
