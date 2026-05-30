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
  AnalystCertification,
  AnalystSkill,
} from "@/types/services/analyst";
import { Service } from "@/types/services/service";

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

  getByUserId(userId: number): Promise<Analyst | null> {
    return api.get<AnalystResponse>(`/users/${userId}/analyst`).then(
      (resp) => {
        return resp.analyst;
      },
      () => {
        return null;
      },
    );
  },

  create(
    userId: number,
    data: { hourly_rate: number },
  ): Promise<Analyst | null> {
    return api.post<AnalystResponse>(`/users/${userId}/analyst`, data).then(
      (resp) => {
        return resp.analyst;
      },
      () => {
        return null;
      },
    );
  },

  update(userId: number, data: Analyst): Promise<Analyst | null> {
    return api.put<AnalystResponse>(`/users/${userId}/analyst`, data).then(
      (resp) => {
        return resp.analyst;
      },
      () => {
        return null;
      },
    );
  },

  patch(userId: number, data: AnalystUpdateRequest): Promise<Analyst | null> {
    return api.patch<AnalystResponse>(`/users/${userId}/analyst`, data).then(
      (resp) => {
        return resp.analyst;
      },
      () => {
        return null;
      },
    );
  },

  delete(userId: number): Promise<Record<string, string>> {
    return api.delete(`/users/${userId}/analyst`);
  },

  // Profile
  getProfile(userId: number): Promise<AnalystProfile | null> {
    return api
      .get<AnalystProfileResponse>(`/users/${userId}/analyst/profile`)
      .then(
        (resp) => {
          return resp.analyst_profile;
        },
        () => {
          return null;
        },
      );
  },

  createProfile(
    userId: number,
    data: AnalystProfile,
  ): Promise<AnalystProfile | null> {
    return api
      .post<AnalystProfileResponse>(`/users/${userId}/analyst/profile`, data)
      .then(
        (resp) => {
          return resp.analyst_profile;
        },
        () => {
          return null;
        },
      );
  },

  updateProfile(
    userId: number,
    data: AnalystProfile,
  ): Promise<AnalystProfile | null> {
    return api
      .put<AnalystProfileResponse>(`/users/${userId}/analyst/profile`, data)
      .then(
        (resp) => {
          return resp.analyst_profile;
        },
        () => {
          return null;
        },
      );
  },

  deleteProfile(userId: number): Promise<Record<string, string>> {
    return api.delete(`/users/${userId}/analyst/profile`);
  },

  // Certifications
  listCertifications(userId: number): Promise<Certification[] | null> {
    return api
      .get<{
        certifications: Certification[];
        count: number;
      }>(`/users/${userId}/analyst/certifications`)
      .then(
        (resp) => {
          if (resp.count === 0) {
            return [];
          }

          return resp.certifications;
        },
        () => {
          return null;
        },
      );
  },

  addCertification(
    userId: number,
    data: { certification_id: number },
  ): Promise<AnalystCertification | null> {
    return api
      .post<AnalystCertificationResponse>(
        `/users/${userId}/analyst/certifications`,
        data,
      )
      .then(
        (resp) => {
          return resp.analyst_certification;
        },
        () => {
          return null;
        },
      );
  },

  removeCertification(
    userId: number,
    certificationId: number,
  ): Promise<boolean> {
    return api
      .delete(
        `/users/${userId}/analyst/certifications?certification_id=${certificationId}`,
      )
      .then(() => true)
      .catch(() => false);
  },

  // Skills
  listSkills(userId: number): Promise<AnalystSkill[] | null> {
    return api
      .get<AnalystSkillsResponse>(`/users/${userId}/analyst/skills`)
      .then(
        (resp) => {
          return resp.analyst_skills;
        },
        () => {
          return null;
        },
      );
  },

  addSkill(
    userId: number,
    data: { skill_id: number },
  ): Promise<AnalystSkill | null> {
    return api
      .post<AnalystSkillResponse>(`/users/${userId}/analyst/skills`, data)
      .then(
        (resp) => {
          return resp.analyst_skill;
        },
        () => {
          return null;
        },
      );
  },

  removeSkill(userId: number, skillId: number): Promise<boolean> {
    return api
      .delete(`/users/${userId}/analyst/skills?skill_id=${skillId}`)
      .then(() => true)
      .catch(() => false);
  },

  //Services
  listServices(userId: number): Promise<Service[] | null> {
    return api
      .get<{
        services: Service[];
        count: number;
      }>(`/analyst/${userId}/services`)
      .then(
        (resp) => {
          return resp.services;
        },
        () => {
          return null;
        },
      );
  },
};
