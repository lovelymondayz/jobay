import { create } from "zustand";
import type { SortingState } from "@tanstack/react-table";

interface DateRange {
  start?: string;
  end?: string;
}
export interface History {
  connection_history_id: string;
  user_name: string;
  device_name: string;
  to_aps: string;
  from_aps: string;
  created_at: string;
}

export interface ApiResponse {
  status: boolean;
  message: string;
  data: {
    items: History[];
    page: number;
    size: number;
    total_data: number;
    max_pages: number;
    visible: number;
  };
}

interface HistoriesState {
  globalFilterInput: string;
  globalFilter: string;
  dateRange: DateRange;
  sorting: SortingState;
  pagination: {
    pageIndex: number;
    pageSize: number;
  };
  setGlobalFilterInput: (v: string) => void;
  setGlobalFilter: (v: string) => void;
  setDateRange: (v: DateRange) => void;
  setSorting: (
    updaterOrValue: SortingState | ((old: SortingState) => SortingState)
  ) => void;
  setPagination: (
    value:
      | { pageIndex: number; pageSize: number }
      | ((old: { pageIndex: number; pageSize: number }) => {
          pageIndex: number;
          pageSize: number;
        })
  ) => void;
}

export const useHistoriesStore = create<HistoriesState>((set) => ({
  globalFilterInput: "",
  globalFilter: "",
  dateRange: {},
  sorting: [],
  pagination: {
    pageIndex: 0,
    pageSize: 10,
  },
  setGlobalFilterInput: (v) => set({ globalFilterInput: v }),
  setGlobalFilter: (v) => set({ globalFilter: v }),
  setDateRange: (v) => set({ dateRange: v }),
  setSorting: (value) =>
    set((state) => ({
      sorting: typeof value === "function" ? value(state.sorting) : value,
    })),
  setPagination: (
    value:
      | { pageIndex: number; pageSize: number }
      | ((old: { pageIndex: number; pageSize: number }) => {
          pageIndex: number;
          pageSize: number;
        })
  ) =>
    set((state) => ({
      pagination: typeof value === "function" ? value(state.pagination) : value,
    })),
}));
