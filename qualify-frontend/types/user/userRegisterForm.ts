export interface IUserRegisterForm {
  name: string;
  surname: string;
  email: string;
  phone: string;
  timezone: string;
  country_name: string;
  country_code: string;
  country_state: string;
  city: string;
  password: string;
  confirmPassword: string;
}

export interface IUserRegisterFormErrors {
  name?: string;
  surname?: string;
  email?: string;
  phone?: string;
  timezone?: string;
  country_name?: string;
  country_code?: string;
  country_state?: string;
  city?: string;
  password?: string;
  confirmPassword?: string;
}
