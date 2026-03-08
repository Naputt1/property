export interface IProperty {
  id: string;
  created_at: string;
  updated_at: string;
  price: number;
  date_of_transfer: string;
  postcode_outward: string;
  postcode_inward: string;
  property_type: string;
  old_new: string;
  duration: string;
  address: string;
  town_city: string;
  district: string;
  county: string;
  ppd_category_type: string;
  record_status: string;
}

export interface IPropertyListRes {
  status: boolean;
  data: IProperty[];
  total: number;
}
