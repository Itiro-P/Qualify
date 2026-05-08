export interface Certification {
  id: number;
  name: string;
  description: string;
  institution: string;
  year: number;
}

export interface CertificationResponse {
  certification: Certification;
}

export interface CertificationsResponse {
  certifications: Certification[];
  count: number;
}

export interface CertificationUpdateRequest {
  name?: string;
  description?: string;
  institution?: string;
  year?: number;
}
