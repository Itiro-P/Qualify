import { api } from "@/libs/api";
import type {
  User,
  UserResponse,
  UserRegister,
  UserLogin,
  UserUpdateRequest,
} from "@/types/services/user";

export const userService = {
  register(data: UserRegister): Promise<UserResponse> {
    return api.post("/register", data);
  },

  login(data: UserLogin): Promise<UserResponse> {
    return api.post("/login", data);
  },

  getById(id: number): Promise<UserResponse> {
    return api.get(`/users/${id}`);
  },

  update(id: number, data: User): Promise<UserResponse> {
    return api.put(`/users/${id}`, data);
  },

  patch(id: number, data: UserUpdateRequest): Promise<UserResponse> {
    return api.patch(`/users/${id}`, data);
  },

  delete(id: number): Promise<Record<string, string>> {
    return api.delete(`/users/${id}`);
  },
};
