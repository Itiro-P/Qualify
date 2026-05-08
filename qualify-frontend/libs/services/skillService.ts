import { api } from "@/libs/api";
import type {
  SkillResponse,
  SkillsResponse,
} from "@/types/services/skill";

export const skillService = {
  list(name?: string): Promise<SkillsResponse> {
    const query = name ? `?name=${encodeURIComponent(name)}` : "";
    return api.get(`/skills${query}`);
  },

  getById(id: number): Promise<SkillResponse> {
    return api.get(`/skills/${id}`);
  },

  create(data: { name: string }): Promise<SkillResponse> {
    return api.post("/skills", data);
  },

  update(id: number, data: { name: string }): Promise<SkillResponse> {
    return api.put(`/skills/${id}`, data);
  },

  delete(id: number): Promise<Record<string, string>> {
    return api.delete(`/skills/${id}`);
  },
};
