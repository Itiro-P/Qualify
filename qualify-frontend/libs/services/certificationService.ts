import { api } from "@/libs/api";
import type {
  Certification,
  CertificationResponse,
  CertificationsResponse,
  CertificationUpdateRequest,
} from "@/types/services/certification";

export const certificationService = {
  list(name?: string): Promise<Certification[] | null> {
    const query = name ? `?name=${encodeURIComponent(name)}` : "";
    return api.get<CertificationsResponse>(`/certifications${query}`).then(
      (resp) => {
        return resp.certifications;
      },
      () => {
        return null;
      },
    );
  },

  getById(id: number): Promise<Certification | null> {
    return api.get<CertificationResponse>(`/certifications/${id}`).then(
      (resp) => {
        return resp.certification;
      },
      () => {
        return null;
      },
    );
  },

  create(data: Omit<Certification, "id">): Promise<Certification | null> {
    return api.post<CertificationResponse>("/certifications", data).then(
      (resp) => {
        return resp.certification;
      },
      () => {
        return null;
      },
    );
  },

  update(
    id: number,
    data: Omit<Certification, "id">,
  ): Promise<Certification | null> {
    return api.put<CertificationResponse>(`/certifications/${id}`, data).then(
      (resp) => {
        return resp.certification;
      },
      () => {
        return null;
      },
    );
  },

  patch(
    id: number,
    data: CertificationUpdateRequest,
  ): Promise<Certification | null> {
    return api.patch<CertificationResponse>(`/certifications/${id}`, data).then(
      (resp) => {
        return resp.certification;
      },
      () => {
        return null;
      },
    );
  },

  delete(id: number): Promise<Record<string, string>> {
    return api.delete(`/certifications/${id}`);
  },
};
