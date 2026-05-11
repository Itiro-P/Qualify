import { api } from "@/libs/api";
import type {
  Client,
  ClientResponse,
  ClientsResponse,
  ClientUpdateRequest,
  ClientProfile,
  ClientProfileResponse,
  ListClientsParams,
} from "@/types/services/client";

function buildQuery(params?: Record<string, unknown>): string {
  if (!params) return "";
  const search = new URLSearchParams();
  for (const [key, value] of Object.entries(params)) {
    if (value !== undefined && value !== null && value !== "") {
      search.set(key, String(value));
    }
  }
  const qs = search.toString();
  return qs ? `?${qs}` : "";
}

export const clientService = {
  list(params?: ListClientsParams): Promise<ClientsResponse> {
    return api.get(`/clients${buildQuery(params as Record<string, unknown>)}`);
  },

  getByUserId(userId: number): Promise<ClientResponse> {
    return api.get(`/users/${userId}/client`);
  },

  create(
    userId: number,
    data: { proposed_budget: number },
  ): Promise<ClientResponse> {
    return api.post(`/users/${userId}/client`, data);
  },

  update(userId: number, data: Client): Promise<ClientResponse> {
    return api.put(`/users/${userId}/client`, data);
  },

  patch(userId: number, data: ClientUpdateRequest): Promise<ClientResponse> {
    return api.patch(`/users/${userId}/client`, data);
  },

  delete(userId: number): Promise<Record<string, string>> {
    return api.delete(`/users/${userId}/client`);
  },

  // Profile
  getProfile(userId: number): Promise<ClientProfileResponse> {
    return api.get(`/users/${userId}/client/profile`);
  },

  createProfile(
    userId: number,
    data: ClientProfile,
  ): Promise<ClientProfileResponse> {
    return api.post(`/users/${userId}/client/profile`, data);
  },

  updateProfile(
    userId: number,
    data: ClientProfile,
  ): Promise<ClientProfileResponse> {
    return api.put(`/users/${userId}/client/profile`, data);
  },

  deleteProfile(userId: number): Promise<Record<string, string>> {
    return api.delete(`/users/${userId}/client/profile`);
  },
};
