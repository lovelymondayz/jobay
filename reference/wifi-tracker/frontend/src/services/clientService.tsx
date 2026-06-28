import { api } from "./api";
import type { Clients } from "@/types/clientsStore";

type ApiResponse<T> = {
  data: T;
};

// Clients
export const getClient = (id: string) =>
  api.get<ApiResponse<Clients>>(`/clients/${id}`).then((res) => res.data.data);

export const getClients = async () => {
  const res = await api.get(`/clients`);
  if (res.data.data && Array.isArray(res.data.data)) {
    return res.data.data;
  } else if (Array.isArray(res.data)) {
    return res.data;
  } else {
    return [];
  }
};

export const createClient = (data: Partial<Clients>) =>
  api.post<ApiResponse<Clients>>(`/clients`, data).then((res) => res.data.data);

export const updateClient = (id: string, data: Partial<Clients>) =>
  api
    .put<ApiResponse<Clients>>(`/clients/${id}`, data)
    .then((res) => res.data.data);

export const deleteClient = (id: string) =>
  api.delete(`/clients/${id}`).then((res) => res.data);
