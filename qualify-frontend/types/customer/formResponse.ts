interface Location {
  country: string;
  state: string;
  city: string;
}

export interface IFormResponse {
  min_hourly_rate: number;
  max_hourly_rate: number;
  rating: number;
  skills: string[];
  localization: Location;
}
