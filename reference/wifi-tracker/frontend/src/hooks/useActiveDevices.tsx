import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/services/api";
import type { Device } from "@/types/deviceStore";
import {
  getDevice,
  // getDevices,
  createDevice,
  updateDevice,
  deleteDevice,
} from "@/services/deviceService";

const fetchActiveDevices = async (): Promise<Device[]> => {
  const response = await api.get("/devices/");
  return response.data.data;
};

export const useActiveDevices = () =>
  useQuery<Device[]>({
    queryKey: ["activeDevices"],
    queryFn: fetchActiveDevices,
  });

export const useUserDevice = (id: string) =>
  useQuery({
    queryKey: ["device", id],
    queryFn: () => getDevice(id),
    enabled: !!id,
  });

export const useCreateDevice = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: createDevice,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["devices"] });
    },
  });
};

export const useUpdateDevice = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<Device> }) =>
      updateDevice(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["devices"] });
    },
  });
};

export const useDeleteDevice = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: deleteDevice,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["devices"] });
    },
  });
};
