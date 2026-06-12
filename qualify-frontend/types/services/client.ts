export interface Client {
  id?: number;
  name?: string;
  email?: string;
  phone?: string;
  city?: string;
  country_code?: string;
  country_name?: string;
  country_state?: string;
  timezone?: string;
  proposed_budget?: number;
  time_created?: string;
}

export interface ClientResponse {
  client: Client;
}

export interface ClientsResponse {
  clients: Client[];
  count: number;
}

export interface ClientUpdateRequest {
  name?: string;
  email?: string;
  phone?: string;
  city?: string;
  country_code?: string;
  country_name?: string;
  country_state?: string;
  timezone?: string;
  proposed_budget?: number;
}

export interface ClientProfile {
  user_id?: number;
  picture?: string;
  biography?: string;
}

export interface ClientProfileResponse {
  client_profile: ClientProfile;
}

export interface ListClientsParams {
  name?: string;
  country?: string;
  city?: string;
  min_proposed_budget?: number;
  max_proposed_budget?: number;
}
