export interface Service {
  id?: number;
  proposal_letter_id?: number;
  title?: string;
  content?: string;
  status?: string;
  hourly_rate?: number;
  time_created?: string;
}

export interface ServiceResponse {
  service: Service;
}

export interface ServicesResponse {
  services: Service[];
  count: number;
}

export interface ServiceUpdateRequest {
  title?: string;
  content?: string;
  status?: string;
  hourly_rate?: number;
}

export interface ListServicesParams {
  status?: string;
  proposal_letter_id?: number;
}
