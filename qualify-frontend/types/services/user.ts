export interface User {
  id: number;
  name: string;
  email: string;
  phone: string;
  city: string;
  country_code: string;
  country_name: string;
  country_state: string;
  timezone: string;
  time_created: string;
}

export interface UserProfile {
  user_id?: number;
  picture?: string;
  biography: string;
}

export interface UserProfileResponse {
  user_profile: UserProfile;
}

export interface HTTPUserResponse {
  user: User;
}

export interface HTTPLoginResponse {
  access_token: string;
  refresh_token: string;
  token_type: string;
  expires_in: number;
  expires_at: string;
}

export interface HTTPPasswordChangeResponse {
  success: boolean;
  message: string;
  changed_at: string;
}

export interface HTTPPasswordResetResponse {
  success: boolean;
  message: string;
}

export interface HTTPSuccessResponse {
  success: boolean;
  message: string;
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

export interface PasswordChangeRequest {
  current_password: string;
  new_password: string;
}

export interface PasswordResetConfirmRequest {
  token: string;
  new_password: string;
}
