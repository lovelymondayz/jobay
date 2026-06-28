import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/services/api";
import type { Clients } from "@/types/clientsStore";

import {
  getClients,
  getClient,
  createClient,
  updateClient,
  deleteClient,
} from "@/services/clientService";
const fetchActiveClients = async (): Promise<Clients[]> => {
  const response = await api.get("/clients/active");
  return response.data.data;
};

export const useActiveClients = () =>
  useQuery<Clients[]>({
    queryKey: ["activeClients"],
    queryFn: fetchActiveClients,
    // refetchInterval: 5000,
  });

// Get all clients
export const useClients = () =>
  useQuery<Clients[]>({
    queryKey: ["clients"],
    queryFn: getClients,
  });

// Get single client
export const useClient = (id: string) =>
  useQuery<Clients>({
    queryKey: ["client", id],
    queryFn: () => getClient(id),
    enabled: !!id,
  });

// Create
export const useCreateClient = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: createClient,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["clients"] });
    },
  });
};

// Update
export const useUpdateClient = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<Clients> }) =>
      updateClient(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["clients"] });
    },
  });
};

// Delete
export const useDeleteClient = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: deleteClient,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["clients"] });
    },
  });
};
