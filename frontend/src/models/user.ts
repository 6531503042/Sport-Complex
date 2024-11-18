export interface IUser {
  _id: string;
  name: string;
  email: string;
  role_code: number;
  created_at: string;
  updated_at: string;
}

export interface UserResponse {
  data: IUser[];
  status: string;
  message?: string;
} 