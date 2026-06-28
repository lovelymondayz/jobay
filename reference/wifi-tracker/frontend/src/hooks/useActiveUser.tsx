import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";

import type { User } from "@/types/userStore";

import {
  getUsers,
  getUser,
  createUser,
  updateUser,
  deleteUser,
} from "@/services/userServices";

// Get all Users
export const useUsers = () =>
  useQuery<User[]>({
    queryKey: ["users"],
    queryFn: getUsers,
  });

// Get single User
export const useUser = (id: string) =>
  useQuery<User>({
    queryKey: ["users", id],
    queryFn: () => getUser(id),
    enabled: !!id,
  });

// Create
export const useCreateUser = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: createUser,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["users"] });
    },
  });
};

// Update
export const useUpdateUser = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<User> }) =>
      updateUser(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["users"] });
    },
  });
};

// Delete
export const useDeleteUser = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: deleteUser,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["users"] });
    },
  });
};
