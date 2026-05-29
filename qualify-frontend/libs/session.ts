import type { User, HTTPUserResponse } from "@/types/services/user";
import { userService } from "@/libs/services";

export type { User };

interface StoredTokens {
  access_token: string;
  refresh_token: string;
}

const TOKEN_KEY = "qualify_tokens";

export function getStoredTokens(): StoredTokens | null {
  if (typeof window === "undefined") return null;
  const raw = localStorage.getItem(TOKEN_KEY);
  if (!raw) return null;
  try {
    return JSON.parse(raw) as StoredTokens;
  } catch {
    return null;
  }
}

export function setStoredTokens(tokens: StoredTokens): void {
  localStorage.setItem(TOKEN_KEY, JSON.stringify(tokens));
}

export function getAccessToken(): string | null {
  const tokens = getStoredTokens();
  return tokens?.access_token ?? null;
}

export async function getSessionUser(): Promise<User | null> {
  if (typeof window === "undefined") return null;
  const tokens = getStoredTokens();
  if (!tokens) return null;

  try {
    const user = await userService.me();
    if (!user) {
      clearSession();
      return null;
    }

    return user;
  } catch {
    return null;
  }
}

export function clearSession(): void {
  localStorage.removeItem(TOKEN_KEY);
}
