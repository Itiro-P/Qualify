export interface Analyst {
  id: number;
  name: string;
  email: string;
  phone: string;
  city: string;
  country_code: string;
  country_name: string;
  country_state: string;
  timezone: string;
  hourly_rate: number;
  mean_rating: number;
  total_reviews: number;
  time_created: string;
}

export interface AnalystResponse {
  analyst: Analyst;
}

export interface AnalystsResponse {
  analysts: Analyst[];
  count: number;
}

export interface AnalystUpdateRequest {
  name?: string;
  email?: string;
  phone?: string;
  city?: string;
  country_code?: string;
  country_name?: string;
  country_state?: string;
  timezone?: string;
  hourly_rate?: number;
  mean_rating?: number;
  total_reviews?: number;
}

export interface AnalystProfile {
  user_id?: number;
  biography: string;
}

export interface AnalystProfileResponse {
  analyst_profile: AnalystProfile;
}

export interface AnalystCertification {
  analyst_id: number;
  certification_id: number;
}

export interface AnalystCertificationResponse {
  analyst_certification: AnalystCertification;
}

export interface AnalystSkill {
  analyst_id: number;
  skill_id: number;
}

export interface AnalystSkillResponse {
  analyst_skill: AnalystSkill;
}

export interface AnalystSkillsResponse {
  analyst_skills: AnalystSkill[];
  count: number;
}

export interface ListAnalystsParams {
  name?: string;
  email?: string;
  country?: string;
  country_code?: string;
  country_state?: string;
  city?: string;
  time_zone?: string;
  min_hourly_rate?: number;
  max_hourly_rate?: number;
  min_total_reviews?: number;
  max_total_reviews?: number;
  min_rating?: number;
  max_rating?: number;
  skills?: string;
  sort_by?: string;
  order?: "ASC" | "DESC";
}
