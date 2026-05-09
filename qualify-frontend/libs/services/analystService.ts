import { api } from "@/libs/api";
import { Certification } from "@/types/services/certification";
import type {
  AnalystResponse,
  AnalystsResponse,
  AnalystUpdateRequest,
  AnalystProfile,
  AnalystProfileResponse,
  AnalystCertificationResponse,
  AnalystSkillResponse,
  AnalystSkillsResponse,
  ListAnalystsParams,
  Analyst,
} from "@/types/services/analyst";

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

export const analystService = {
  list(params?: ListAnalystsParams): Promise<AnalystsResponse> {
    return api.get(`/analysts${buildQuery(params as Record<string, unknown>)}`);
  },

  getByUserId(userId: number): Promise<AnalystResponse> {
    return api.get(`/users/${userId}/analyst`);
  },

  create(
    userId: number,
    data: { hourly_rate: number },
  ): Promise<AnalystResponse> {
    return api.post(`/users/${userId}/analyst`, data);
  },

  update(userId: number, data: Analyst): Promise<AnalystResponse> {
    return api.put(`/users/${userId}/analyst`, data);
  },

  patch(userId: number, data: AnalystUpdateRequest): Promise<AnalystResponse> {
    return api.patch(`/users/${userId}/analyst`, data);
  },

  delete(userId: number): Promise<Record<string, string>> {
    return api.delete(`/users/${userId}/analyst`);
  },

  // Profile
  getProfile(userId: number): Promise<AnalystProfileResponse> {
    return api.get(`/users/${userId}/analyst/profile`);
  },

  createProfile(
    userId: number,
    data: AnalystProfile,
  ): Promise<AnalystProfileResponse> {
    return api.post(`/users/${userId}/analyst/profile`, data);
  },

  updateProfile(
    userId: number,
    data: AnalystProfile,
  ): Promise<AnalystProfileResponse> {
    return api.put(`/users/${userId}/analyst/profile`, data);
  },

  deleteProfile(userId: number): Promise<Record<string, string>> {
    return api.delete(`/users/${userId}/analyst/profile`);
  },

  // Certifications
  listCertifications(
    userId: number,
  ): Promise<{ certifications: Certification[]; count: number }> {
    return api.get(`/users/${userId}/analyst/certifications`);
  },

  addCertification(
    userId: number,
    data: { certification_id: number },
  ): Promise<AnalystCertificationResponse> {
    return api.post(`/users/${userId}/analyst/certifications`, data);
  },

  removeCertification(
    userId: number,
    certificationId: number,
  ): Promise<Record<string, string>> {
    return api.delete(
      `/users/${userId}/analyst/certifications?certification_id=${certificationId}`,
    );
  },

  // Skills
  listSkills(userId: number): Promise<AnalystSkillsResponse> {
    return api.get(`/users/${userId}/analyst/skills`);
  },

  addSkill(
    userId: number,
    data: { skill_id: number },
  ): Promise<AnalystSkillResponse> {
    return api.post(`/users/${userId}/analyst/skills`, data);
  },

  removeSkill(userId: number): Promise<Record<string, string>> {
    return api.delete(`/users/${userId}/analyst/skills`);
  },
};
