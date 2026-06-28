import { api } from "./api";
import type { Device } from "@/types/deviceStore";

type ApiResponse<T> = {
  data: T;
};

// Clients
export const getDevice = (id: string) =>
  api
    .get<ApiResponse<Device>>(`/user-device/${id}`)
    .then((res) => res.data.data);

export const getDevices = async () => {
  const res = await api.get(`/user-device`);
  if (res.data.data && Array.isArray(res.data.data)) {
    return res.data.data;
  } else if (Array.isArray(res.data)) {
    return res.data;
  } else {
    return [];
  }
};

export const createDevice = (data: Partial<Device>) =>
  api
    .post<ApiResponse<Device>>(`/user-device`, data)
    .then((res) => res.data.data);

export const updateDevice = (id: string, data: Partial<Device>) =>
  api
    .put<ApiResponse<Device>>(`/user-device/${id}`, data)
    .then((res) => res.data.data);

export const deleteDevice = (id: string) =>
  api.delete(`/user-devices/${id}`).then((res) => res.data);
