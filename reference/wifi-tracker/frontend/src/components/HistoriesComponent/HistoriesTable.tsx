import React from "react";
import {
  useReactTable,
  getCoreRowModel,
  flexRender,
} from "@tanstack/react-table";
import type { ColumnDef } from "@tanstack/react-table";
import type { History } from "@/types/historiesStore";
import { useHistoriesStore } from "@/types/historiesStore";

export default function HistoriesTable({ data }: { data: any }) {
  const { sorting, pagination, setSorting, setPagination } =
    useHistoriesStore();

  const columns = React.useMemo<ColumnDef<History, any>[]>(
    () => [
      { accessorKey: "user_name", header: "User Name" },
      { accessorKey: "device_name", header: "Device Name" },
      { accessorKey: "to_aps", header: "To AP's" },
      { accessorKey: "from_aps", header: "From AP's" },
      {
        accessorKey: "created_at",
        header: "Created At",
        cell: ({ getValue }) => {
          const date = new Date(getValue<string>());
          return date.toLocaleString();
        },
      },
    ],
    []
  );

  const table = useReactTable({
    data: data?.items ?? [],
    columns,
    pageCount: data?.max_pages ?? 0,
    state: {
      sorting,
      pagination,
    },
    onSortingChange: setSorting,
    onPaginationChange: setPagination,
    getCoreRowModel: getCoreRowModel(),
    manualPagination: true,
    manualSorting: true,
  });

  return (
    <>
      <div className="mb-2 text-sm text-gray-600">
        Showing {data?.visible || 0} of {data?.total_data || 0} records
      </div>

      <div className=" overflow-x-auto w-5/6 lg:w-full lg:min-w-max">
        <table className=" w-auto lg:w-full border text-sm">
          <thead className="text-left sticky top-0 z-10 bg-based">
            {table.getHeaderGroups().map((headerGroup) => (
              <tr key={headerGroup.id}>
                {headerGroup.headers.map((header) => (
                  <th
                    key={header.id}
                    onClick={header.column.getToggleSortingHandler()}
                    className="p-2  cursor-pointer select-none whitespace-nowrap border-2"
                  >
                    {flexRender(
                      header.column.columnDef.header,
                      header.getContext()
                    )}
                    {{
                      asc: " 🔼",
                      desc: " 🔽",
                    }[header.column.getIsSorted() as string] ?? ""}
                  </th>
                ))}
              </tr>
            ))}
          </thead>
          <tbody>
            {table.getRowModel().rows.map((row) => (
              <tr key={row.id} className="border-t ">
                {row.getVisibleCells().map((cell) => (
                  <td key={cell.id} className="p-2 border whitespace-nowrap">
                    {flexRender(cell.column.columnDef.cell, cell.getContext())}
                  </td>
                ))}
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </>
  );
}
