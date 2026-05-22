export interface Skill {
  id: number;
  name: string;
}

export interface SkillResponse {
  skill: Skill;
}

export interface SkillsResponse {
  skills: Skill[];
  count: number;
}
