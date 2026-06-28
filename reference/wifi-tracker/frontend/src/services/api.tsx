import axios from "axios";
import type { Device } from "@/types/deviceStore";

type ApiResponse<T> = {
  data: T;
};

export const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL,
});

// Devices

export const getUserDevice = (id: string) =>
  api
    .get<ApiResponse<Device>>(`/user-devices/${id}`)
    .then((res) => res.data.data);

export const getUserDevices = async () => {
  const res = await api.get(`/user-devices`);
  console.log("API Response:", res.data);

  if (res.data.data && Array.isArray(res.data.data)) {
    return res.data.data;
  } else if (Array.isArray(res.data)) {
    return res.data;
  } else {
    return [];
  }
};

export const createUserDevice = (data: Partial<Device>) =>
  api
    .post<ApiResponse<Device>>(`/user-devices`, data)
    .then((res) => res.data.data);

export const updateUserDevice = (id: string, data: Partial<Device>) =>
  api
    .put<ApiResponse<Device>>(`/user-devices/${id}`, data)
    .then((res) => res.data.data);

export const deleteUserDevice = (id: string) =>
  api.delete(`/user-devices/${id}`).then((res) => res.data);
