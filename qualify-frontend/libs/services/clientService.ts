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
  list(params?: ListClientsParams): Promise<Client[] | null> {
    return api
      .get<ClientsResponse>(
        `/clients${buildQuery(params as Record<string, unknown>)}`,
      )
      .then(
        (resp) => {
          return resp.clients;
        },
        () => {
          return null;
        },
      );
  },

  getByUserId(userId: number): Promise<Client | null> {
    return api.get<ClientResponse>(`/users/${userId}/client`).then(
      (resp) => {
        return resp.client;
      },
      () => {
        return null;
      },
    );
  },

  create(
    userId: number,
    data: { proposed_budget: number },
  ): Promise<Client | null> {
    return api.post<ClientResponse>(`/users/${userId}/client`, data).then(
      (resp) => {
        return resp.client;
      },
      () => {
        return null;
      },
    );
  },

  update(userId: number, data: Client): Promise<Client | null> {
    return api.put<ClientResponse>(`/users/${userId}/client`, data).then(
      (resp) => {
        return resp.client;
      },
      () => {
        return null;
      },
    );
  },

  patch(userId: number, data: ClientUpdateRequest): Promise<Client | null> {
    return api.patch<ClientResponse>(`/users/${userId}/client`, data).then(
      (resp) => {
        return resp.client;
      },
      () => {
        return null;
      },
    );
  },

  delete(userId: number): Promise<Record<string, string>> {
    return api.delete(`/users/${userId}/client`);
  },

  // Profile
  getProfile(userId: number): Promise<ClientProfile | null> {
    return api
      .get<ClientProfileResponse>(`/users/${userId}/client/profile`)
      .then(
        (resp) => {
          return resp.client_profile;
        },
        () => {
          return null;
        },
      );
  },

  createProfile(
    userId: number,
    data: ClientProfile,
  ): Promise<ClientProfile | null> {
    return api
      .post<ClientProfileResponse>(`/users/${userId}/client/profile`, data)
      .then(
        (resp) => {
          return resp.client_profile;
        },
        () => {
          return null;
        },
      );
  },

  updateProfile(
    userId: number,
    data: ClientProfile,
  ): Promise<ClientProfile | null> {
    return api
      .put<ClientProfileResponse>(`/users/${userId}/client/profile`, data)
      .then(
        (resp) => {
          return resp.client_profile;
        },
        () => {
          return null;
        },
      );
  },

  deleteProfile(userId: number): Promise<Record<string, string>> {
    return api.delete(`/users/${userId}/client/profile`);
  },
};
