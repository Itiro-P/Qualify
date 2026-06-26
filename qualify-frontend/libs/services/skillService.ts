import { api } from "@/libs/api";
import type {
  Skill,
  SkillResponse,
  SkillsResponse,
} from "@/types/services/skill";

export const skillService = {
  list(name?: string): Promise<Skill[] | null> {
    const query = name ? `?name=${encodeURIComponent(name)}` : "";
    return api.get<SkillsResponse>(`/skills${query}`).then(
      (resp) => {
        return resp.skills;
      },
      () => {
        return null;
      },
    );
  },

  getById(id: number): Promise<Skill | null> {
    return api.get<SkillResponse>(`/skills/${id}`).then(
      (resp) => {
        return resp.skill;
      },
      () => {
        return null;
      },
    );
  },

  create(data: { name: string }): Promise<Skill | null> {
    return api.post<SkillResponse>("/skills", data).then(
      (resp) => {
        return resp.skill;
      },
      () => {
        return null;
      },
    );
  },

  update(id: number, data: { name: string }): Promise<Skill | null> {
    return api.put<SkillResponse>(`/skills/${id}`, data).then(
      (resp) => {
        return resp.skill;
      },
      () => {
        return null;
      },
    );
  },

  delete(id: number): Promise<Record<string, string>> {
    return api.delete(`/skills/${id}`);
  },
};
