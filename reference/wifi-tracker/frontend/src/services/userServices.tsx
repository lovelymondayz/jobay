import { api } from "./api";

import type { User } from "@/types/userStore";

type ApiResponse<T> = {
  data: T;
};

// User
export const getUser = (id: string) =>
  api.get<ApiResponse<User>>(`/users/${id}`).then((res) => res.data.data);

export const getUsers = async () => {
  const res = await api.get(`/users/`);
  if (res.data.data && Array.isArray(res.data.data)) {
    return res.data.data;
  } else if (Array.isArray(res.data)) {
    return res.data;
  } else {
    return [];
  }
};

export const createUser = (data: Partial<User>) =>
  api.post<ApiResponse<User>>(`/users`, data).then((res) => res.data.data);

export const updateUser = (id: string, data: Partial<User>) =>
  api.put<ApiResponse<User>>(`/users/${id}`, data).then((res) => res.data.data);

export const deleteUser = (id: string) =>
  api.delete(`/users/${id}`).then((res) => res.data);
