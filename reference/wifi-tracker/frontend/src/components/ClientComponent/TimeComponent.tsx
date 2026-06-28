import { useQuery } from "@tanstack/react-query";
import Loading from "../Loading";

// Custom hook to get current time
const useCurrentTime = () => {
  return useQuery({
    queryKey: ["currentTime"],
    queryFn: () => {
      const now = new Date();
      return {
        time: now.toLocaleTimeString("en-US", {
          hour: "2-digit",
          minute: "2-digit",
          second: "2-digit",
          hour12: true,
        }),
        date: now.toLocaleDateString("en-US", {
          weekday: "long",
          year: "numeric",
          month: "long",
          day: "numeric",
        }),
        timestamp: now.getTime(),
      };
    },
    refetchInterval: 1000, // Update every second
    staleTime: 0, // Always consider stale to ensure updates
  });
};

// Time component
export const TimeComponent = () => {
  const { data: timeData, isLoading, error } = useCurrentTime();

  if (isLoading) {
    return <Loading />;
  }
  if (error) {
    return <div className="text-red-500 p-4">Error loading time</div>;
  }
  ``;
  return (
    <div className="w-full max-w-xs sm:max-w-sm md:max-w-md flex items-center justify-between rounded-xl border-l-4 border-primary p-4 shadow transition hover:shadow-md bg-white">
      <div className="flex flex-col">
        <span className="text-xl font-medium">Local Time</span>
        <h1 className="text-3xl font-mono font-bold">{timeData?.time}</h1>
        <span className="text-xl">{timeData?.date}</span>
      </div>
    </div>
  );
};
