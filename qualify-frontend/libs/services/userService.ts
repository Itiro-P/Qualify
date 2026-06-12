import { api } from "@/libs/api";
import type {
  User,
  UserProfile,
  UserProfileResponse,
  HTTPUserResponse,
  HTTPLoginResponse,
  HTTPPasswordChangeResponse,
  HTTPPasswordResetResponse,
  HTTPSuccessResponse,
  UserRegister,
  UserLogin,
  UserUpdateRequest,
  PasswordChangeRequest,
  PasswordResetConfirmRequest,
} from "@/types/services/user";

export const userService = {
  register(data: UserRegister): Promise<User | null> {
    return api.post<HTTPUserResponse>("/register", data).then(
      (res) => res.user,
      () => null,
    );
  },

  login(data: UserLogin): Promise<HTTPLoginResponse | null> {
    return api.post<HTTPLoginResponse>("/auth/login", data).then(
      (res) => res,
      () => null,
    );
  },

  me(): Promise<User | null> {
    // Não precisa colocar o token, porque o método final (request) adiciona automaticamente o token de acesso se ele existir
    return api.get<HTTPUserResponse>("/auth/me").then(
      (res) => res.user,
      () => null,
    );
  },

  logout(refreshToken?: string): Promise<boolean> {
    const query = refreshToken
      ? `?refresh_token=${encodeURIComponent(refreshToken)}`
      : "";
    return api.post<HTTPSuccessResponse>(`/auth/logout${query}`, {}).then(
      () => true,
      () => false,
    );
  },

  refreshToken(token: string): Promise<HTTPLoginResponse | null> {
    return api.post<HTTPLoginResponse>(
      `/auth/refresh`,
      { refresh_token: token },
    ).then(
      (res) => res,
      () => null,
    );
  },

  resetPassword(email: string): Promise<HTTPPasswordResetResponse | null> {
    return api.post<HTTPPasswordResetResponse>(
      `/auth/reset-password?email=${encodeURIComponent(email)}`,
      {},
    ).then(
      (res) => res,
      () => null,
    );
  },

  resetPasswordConfirm(data: PasswordResetConfirmRequest): Promise<HTTPPasswordResetResponse | null> {
    return api.post<HTTPPasswordResetResponse>(
      `/auth/reset-password/confirm`,
      data,
    ).then(
      (res) => res,
      () => null,
    );
  },

  changePassword(data: PasswordChangeRequest): Promise<HTTPPasswordChangeResponse | null> {
    return api.post<HTTPPasswordChangeResponse>("/auth/change-password", data).then(
      (res) => res,
      () => null,
    );
  },

  getById(id: number): Promise<User | null> {
    return api.get<HTTPUserResponse>(`/users/${id}`).then(
      (res) => res.user,
      () => null,
    );
  },

  update(id: number, data: UserUpdateRequest): Promise<User | null> {
    return api.put<HTTPUserResponse>(`/users/${id}`, data).then(
      (res) => res.user,
      () => null,
    );
  },

  patch(id: number, data: UserUpdateRequest): Promise<User | null> {
    return api.patch<HTTPUserResponse>(`/users/${id}`, data).then(
      (res) => res.user,
      () => null,
    );
  },

  delete(id: number): Promise<boolean> {
    return api.delete(`/users/${id}`).then(
      () => true,
      () => false,
    );
  },

  // Profile
  getProfile(userId: number): Promise<UserProfile | null> {
    return api.get<UserProfileResponse>(`/users/${userId}/profile`).then(
      (res) => res.user_profile,
      () => null,
    );
  },

  createProfile(userId: number, data: UserProfile): Promise<UserProfile | null> {
    return api.post<UserProfileResponse>(`/users/${userId}/profile`, data).then(
      (res) => res.user_profile,
      () => null,
    );
  },

  updateProfile(userId: number, data: UserProfile): Promise<UserProfile | null> {
    return api.put<UserProfileResponse>(`/users/${userId}/profile`, data).then(
      (res) => res.user_profile,
      () => null,
    );
  },

  deleteProfile(userId: number): Promise<boolean> {
    return api.delete(`/users/${userId}/profile`).then(
      () => true,
      () => false,
    );
  },

  uploadProfilePicture(userId: number, file: File): Promise<UserProfile | null> {
    const form = new FormData();
    form.append("picture", file);
    return api.post<UserProfileResponse>(`/users/${userId}/profile/picture`, form).then(
      (res) => res.user_profile,
      () => null,
    );
  },
};
