export interface IUserEditForm {
  name: string;
  surname: string;
  email: string;
  phone: string;
  timezone: string;
  country_name: string;
  country_code: string;
  country_state: string;
  city: string;
}

export interface IUserEditFormErrors {
  name?: string;
  surname?: string;
  email?: string;
  phone?: string;
  timezone?: string;
  country_name?: string;
  country_code?: string;
  country_state?: string;
  city?: string;
}
