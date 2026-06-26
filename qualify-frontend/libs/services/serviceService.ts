import { api } from "@/libs/api";
import type {
  Service,
  ServiceResponse,
  ServicesResponse,
  ServiceUpdateRequest,
  ListServicesParams,
} from "@/types/services/service";

export const serviceService = {
  list(params?: ListServicesParams): Promise<Service[] | null> {
    const search = new URLSearchParams();
    if (params?.status) search.set("status", params.status);
    if (params?.proposal_id)
      search.set("proposal_id", String(params.proposal_id));
    const qs = search.toString();
    return api.get<ServicesResponse>(`/services${qs ? `?${qs}` : ""}`).then(
      (resp) => {
        return resp.services;
      },
      () => {
        return null;
      },
    );
  },

  getById(id: number): Promise<Service | null> {
    return api.get<ServiceResponse>(`/services/${id}`).then(
      (resp) => {
        return resp.service;
      },
      () => {
        return null;
      },
    );
  },

  create(data: Service): Promise<Service | null> {
    return api.post<ServiceResponse>("/services", data).then(
      (resp) => {
        return resp.service;
      },
      () => {
        return null;
      },
    );
  },

  update(id: number, data: Service): Promise<Service | null> {
    return api.put<ServiceResponse>(`/services/${id}`, data).then(
      (resp) => {
        return resp.service;
      },
      () => {
        return null;
      },
    );
  },

  patch(id: number, data: ServiceUpdateRequest): Promise<Service | null> {
    return api.patch<ServiceResponse>(`/services/${id}`, data).then(
      (resp) => {
        return resp.service;
      },
      () => {
        return null;
      },
    );
  },

  delete(id: number): Promise<Record<string, string>> {
    return api.delete(`/services/${id}`);
  },

  listServicesByClient(userId: number): Promise<Service[] | null> {
    return api.get<ServicesResponse>(`/users/${userId}/client/services`).then(
      (resp) => {
        return resp.services;
      },
      () => {
        return null;
      },
    );
  },
}