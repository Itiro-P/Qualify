interface Location {
  country: string;
  state: string;
  city: string;
}

export interface IFormResponse {
  value: number;
  skills: string[];
  localization: Location;
}
