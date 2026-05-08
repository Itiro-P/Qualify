export interface SessionUser {
  id: number;
  name: string;
  email: string;
  phone?: string;
  city?: string;
  country_code?: string;
  country_name?: string;
  country_state?: string;
  timezone?: string;
}

const SESSION_KEY = "qualify_user";

export function getSessionUser(): SessionUser | null {
  if (typeof window === "undefined") return null;
  const raw = sessionStorage.getItem(SESSION_KEY);
  if (!raw) return null;
  try {
    return JSON.parse(raw) as SessionUser;
  } catch {
    return null;
  }
}

export function setSessionUser(user: SessionUser): void {
  sessionStorage.setItem(SESSION_KEY, JSON.stringify(user));
}

export function clearSession(): void {
  sessionStorage.removeItem(SESSION_KEY);
}
