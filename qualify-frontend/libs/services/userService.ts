import { api } from "@/libs/api";
import type {
  User,
  UserResponse,
  UserRegister,
  UserLogin,
  UserUpdateRequest,
} from "@/types/services/user";

export const userService = {
  register(data: UserRegister): Promise<User | null> {
    return api.post<UserResponse>("/register", data).then(
      (resp) => {
        return resp.user;
      },
      () => {
        return null;
      },
    );
  },

  login(data: UserLogin): Promise<User | null> {
    return api.post<UserResponse>("/auth/login", data).then(
      (resp) => {
        return resp.user;
      },
      () => {
        return null;
      },
    );
  },

  getById(id: number): Promise<User | null> {
    return api.get<UserResponse>(`/users/${id}`).then(
      (resp) => {
        return resp.user;
      },
      () => {
        return null;
      },
    );
  },

  update(id: number, data: User): Promise<User | null> {
    return api.put<UserResponse>(`/users/${id}`, data).then(
      (resp) => {
        return resp.user;
      },
      () => {
        return null;
      },
    );
  },

  patch(id: number, data: UserUpdateRequest): Promise<User | null> {
    return api.patch<UserResponse>(`/users/${id}`, data).then(
      (resp) => {
        return resp.user;
      },
      () => {
        return null;
      },
    );
  },

  delete(id: number): Promise<Record<string, string>> {
    return api.delete(`/users/${id}`);
  },
};
