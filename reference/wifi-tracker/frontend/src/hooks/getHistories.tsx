import { useQuery } from "@tanstack/react-query";
import { api } from "@/services/api";
import type { ApiResponse } from "@/types/historiesStore";

const fetchGetHistories = async (
  page: number,
  size: number,
  search?: string,
  startDate?: string,
  endDate?: string
): Promise<ApiResponse["data"]> => {
  const response = await api.get("/histories", {
    params: {
      search,
      page,
      size,
      ...(startDate && { start_date: startDate }),
      ...(endDate && { end_date: endDate }),
    },
  });
  return response.data.data;
};

export const getHistories = (
  pageIndex: number,
  pageSize: number,
  search?: string,
  startDate?: string,
  endDate?: string
) =>
  useQuery({
    queryKey: ["getHistories", pageIndex, pageSize, startDate, endDate, search],
    queryFn: () =>
      fetchGetHistories(pageIndex + 1, pageSize, search, startDate, endDate),
    // keepPreviousData: true,
  });
