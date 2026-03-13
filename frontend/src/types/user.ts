export interface User {
  id: number;
  username: string;
  name?: string;
  is_admin: boolean;
  created_at?: string;
  updated_at?: string;
}
