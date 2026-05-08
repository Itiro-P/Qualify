import { api } from "@/libs/api";
import type {
  Service,
  ServiceResponse,
  ServicesResponse,
  ServiceUpdateRequest,
  ListServicesParams,
} from "@/types/services/service";

export const serviceService = {
  list(params?: ListServicesParams): Promise<ServicesResponse> {
    const search = new URLSearchParams();
    if (params?.status) search.set("status", params.status);
    if (params?.proposal_letter_id) search.set("proposal_letter_id", String(params.proposal_letter_id));
    const qs = search.toString();
    return api.get(`/services${qs ? `?${qs}` : ""}`);
  },

  getById(id: number): Promise<ServiceResponse> {
    return api.get(`/services/${id}`);
  },

  create(data: Service): Promise<ServiceResponse> {
    return api.post("/services", data);
  },

  update(id: number, data: Service): Promise<ServiceResponse> {
    return api.put(`/services/${id}`, data);
  },

  patch(id: number, data: ServiceUpdateRequest): Promise<ServiceResponse> {
    return api.patch(`/services/${id}`, data);
  },

  delete(id: number): Promise<Record<string, string>> {
    return api.delete(`/services/${id}`);
  },
};
