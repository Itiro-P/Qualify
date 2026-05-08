export interface User {
  id?: number;
  name?: string;
  email?: string;
  phone?: string;
  city?: string;
  country_code?: string;
  country_name?: string;
  country_state?: string;
  timezone?: string;
  time_created?: string;
}

export interface UserResponse {
  user: User;
}

export interface UserRegister {
  name: string;
  email: string;
  password: string;
  phone: string;
  city: string;
  country_code: string;
  country_name: string;
  country_state: string;
  timezone: string;
}

export interface UserLogin {
  email: string;
  password: string;
}

export interface UserUpdateRequest {
  name?: string;
  email?: string;
  phone?: string;
  city?: string;
  country_code?: string;
  country_name?: string;
  country_state?: string;
  timezone?: string;
}
