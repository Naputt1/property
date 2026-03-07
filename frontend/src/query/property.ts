import type { IPropertyListRes } from "@/types/property";
import { createQueryWrapper } from "@/utils/query";
import { GET } from "@/services/axios";

type PropertyListParams = {
  page: number;
  pageSize: number;
  town_city?: string;
  district?: string;
  county?: string;
  property_type?: string;
  min_price?: number;
  max_price?: number;
};

export const propertyListQuery = createQueryWrapper<
  IPropertyListRes,
  PropertyListParams,
  { param: PropertyListParams }
>({
  queryKey: (params: PropertyListParams) => ["property", "list", params],
  options: {
    queryFn: (params) => {
      const url = new URL("/property", window.location.origin);
      url.searchParams.append("page", params.param.page.toString());
      url.searchParams.append("pageSize", params.param.pageSize.toString());
      if (params.param.town_city)
        url.searchParams.append("town_city", params.param.town_city);
      if (params.param.district)
        url.searchParams.append("district", params.param.district);
      if (params.param.county)
        url.searchParams.append("county", params.param.county);
      if (params.param.property_type)
        url.searchParams.append("property_type", params.param.property_type);
      if (params.param.min_price)
        url.searchParams.append("min_price", params.param.min_price.toString());
      if (params.param.max_price)
        url.searchParams.append("max_price", params.param.max_price.toString());

      return async () => {
        const res = await GET<IPropertyListRes>(url.pathname + url.search);
        if (res.status) {
          // Return the full response so PropertiesPage can access data.data and data.total
          return res as any as IPropertyListRes;
        }
        return Promise.reject(
          new Error(res.err?.data || "Failed to fetch properties"),
        );
      };
    },
    retry: false,
    staleTime: 5000,
  },
});
