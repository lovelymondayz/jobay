import HistoriesFilter from "../components/HistoriesComponent/HistoriesFilter";
import HistoriesTable from "../components/HistoriesComponent/HistoriesTable";
import Loading from "@/components/Loading";
import { getHistories } from "@/hooks/getHistories";
import { useHistoriesStore } from "@/types/historiesStore";
import { convertToWIBTimeString } from "@/helpers/dateTimeHelper";

export default function HistoriesPage() {
  const { pagination, globalFilter, dateRange } = useHistoriesStore();

  const { data, isLoading, error } = getHistories(
    pagination.pageIndex,
    pagination.pageSize,
    globalFilter,
    convertToWIBTimeString(dateRange.start),
    convertToWIBTimeString(dateRange.end)
  );

  return (
    <div className="h-screen w-screen overflow-hidden flex flex-col p-4 lg:w-full">
      <h1 className="text-2xl font-bold mb-2">Histories</h1>
      <HistoriesFilter maxPages={data?.max_pages || 0} />
      <div className="flex-1 overflow-hidden flex flex-col">
        {isLoading ? (
          <div className="flex-1 flex items-center justify-center">
            <Loading />
          </div>
        ) : error ? (
          <div className="flex-1 flex items-center justify-center text-red-500">
            Error loading data
          </div>
        ) : (
          <HistoriesTable data={data} />
        )}
      </div>
    </div>
  );
}
