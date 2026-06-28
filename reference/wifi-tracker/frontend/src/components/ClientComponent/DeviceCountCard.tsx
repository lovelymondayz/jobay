import { useActiveDevices } from "@/hooks/useActiveDevices";
import Loading from "../Loading";
import { Broadcast } from "@phosphor-icons/react";

export default function DeviceCountCard() {
  const { data: devices, isLoading, error } = useActiveDevices();

  if (isLoading) return <Loading />;
  if (error)
    return <div className="p-6 text-red-500">Error fetching devices.</div>;

  const activeDevices = devices?.filter((device) => device.state === 1) ?? [];

  return (
    <div className="w-full max-w-xs sm:max-w-sm md:max-w-md flex items-center justify-between rounded-xl border-l-4 border-primary p-4 shadow transition hover:shadow-md bg-white">
      <div className="flex flex-col gap-1">
        <span className="text-xs sm:text-sm  font-medium">
          Active Access Points
        </span>
        <h1 className="text-3xl sm:text-4xl font-bold ">
          {activeDevices.length}
        </h1>
      </div>
      <div className="flex-shrink-0">
        <Broadcast size={48} className="text-gray-600" />
      </div>
    </div>
  );
}
