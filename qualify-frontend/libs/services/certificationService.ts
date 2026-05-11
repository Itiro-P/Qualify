import { api } from "@/libs/api";
import type {
  Certification,
  CertificationResponse,
  CertificationsResponse,
  CertificationUpdateRequest,
} from "@/types/services/certification";

export const certificationService = {
  list(name?: string): Promise<CertificationsResponse> {
    const query = name ? `?name=${encodeURIComponent(name)}` : "";
    return api.get(`/certifications${query}`);
  },

  getById(id: number): Promise<CertificationResponse> {
    return api.get(`/certifications/${id}`);
  },

  create(data: Omit<Certification, "id">): Promise<CertificationResponse> {
    return api.post("/certifications", data);
  },

  update(
    id: number,
    data: Omit<Certification, "id">,
  ): Promise<CertificationResponse> {
    return api.put(`/certifications/${id}`, data);
  },

  patch(
    id: number,
    data: CertificationUpdateRequest,
  ): Promise<CertificationResponse> {
    return api.patch(`/certifications/${id}`, data);
  },

  delete(id: number): Promise<Record<string, string>> {
    return api.delete(`/certifications/${id}`);
  },
};
