// components/HistoriesFilter.tsx
import React from "react";
import { useHistoriesStore } from "@/types/historiesStore";
import {
  CaretDoubleLeft,
  CaretDoubleRight,
  CaretLeft,
  CaretRight,
} from "@phosphor-icons/react";
export default function HistoriesFilter({ maxPages }: { maxPages: number }) {
  const {
    globalFilterInput,
    setGlobalFilterInput,
    setGlobalFilter,
    dateRange,
    setDateRange,
    pagination,
    setPagination,
  } = useHistoriesStore();

  // debounce global filter input
  React.useEffect(() => {
    const timeout = setTimeout(() => {
      setGlobalFilter(globalFilterInput);
    }, 500);
    return () => clearTimeout(timeout);
  }, [globalFilterInput]);

  return (
    <div className="flex flex-wrap items-start justify-between gap-4 z-10 mb-4  w-5/6 lg:w-full flex-col ">
      {/* Search */}
      <input
        value={globalFilterInput}
        onChange={(e) => setGlobalFilterInput(e.target.value)}
        placeholder="Search..."
        className="p-2 border rounded w-full lg:w-1/3"
      />

      {/* Date Filters */}
      <div className="flex gap-2 flex-wrap w-full lg:w-1/3 flex-col ">
        <label htmlFor="">Start Date</label>
        <input
          type="datetime-local"
          value={dateRange.start || ""}
          onChange={(e) =>
            setDateRange({
              ...dateRange,
              start: e.target.value,
            })
          }
          className="p-2 border rounded w-full"
        />
        <label htmlFor="">End Date</label>
        <input
          type="datetime-local"
          value={dateRange.end || ""}
          onChange={(e) =>
            setDateRange({
              ...dateRange,
              end: e.target.value,
            })
          }
          className="p-2 border rounded w-full "
        />
      </div>

      {/* Pagination Controls */}
      <div className="flex flex-col sm:flex-row items-center justify-between gap-4 w-full">
        <div className="flex items-center gap-1">
          <button
            onClick={() => setPagination({ ...pagination, pageIndex: 0 })}
            disabled={pagination.pageIndex === 0}
            className="w-9 h-9 flex items-center justify-center rounded border hover:bg-gray-100 disabled:opacity-50 transition"
          >
            <CaretDoubleLeft />
          </button>
          <button
            onClick={() =>
              setPagination({
                ...pagination,
                pageIndex: Math.max(pagination.pageIndex - 1, 0),
              })
            }
            disabled={pagination.pageIndex === 0}
            className="w-9 h-9 flex items-center justify-center rounded border hover:bg-gray-100 disabled:opacity-50 transition"
          >
            <CaretLeft />
          </button>

          <span className="px-3 text-sm whitespace-nowrap">
            Page <span className="font-medium">{pagination.pageIndex + 1}</span>{" "}
            of <span className="font-medium">{maxPages}</span>
          </span>

          <button
            onClick={() =>
              setPagination({
                ...pagination,
                pageIndex: pagination.pageIndex + 1,
              })
            }
            disabled={pagination.pageIndex >= maxPages - 1}
            className="w-9 h-9 flex items-center justify-center rounded border hover:bg-gray-100 disabled:opacity-50 transition"
          >
            <CaretRight />
          </button>
          <button
            onClick={() =>
              setPagination({
                ...pagination,
                pageIndex: maxPages - 1,
              })
            }
            disabled={pagination.pageIndex >= maxPages - 1}
            className="w-9 h-9 flex items-center justify-center rounded border hover:bg-gray-100 disabled:opacity-50 transition"
          >
            <CaretDoubleRight />
          </button>
        </div>

        <div className="flex items-center gap-2">
          <label htmlFor="page-size" className="text-sm whitespace-nowrap">
            Rows per page:
          </label>
          <select
            id="page-size"
            value={pagination.pageSize}
            onChange={(e) =>
              setPagination({
                pageIndex: 0,
                pageSize: Number(e.target.value),
              })
            }
            className="border rounded px-2 py-1 text-sm bg-white focus:ring-2 focus:ring-primary"
          >
            {[10, 50, 100].map((size) => (
              <option key={size} value={size}>
                {size}
              </option>
            ))}
          </select>
        </div>
      </div>
    </div>
  );
}
