import { useActiveClients } from "@/hooks/useActiveClients";
import Loading from "../Loading";
import { Users } from "@phosphor-icons/react";

export default function ClientsCountCard() {
  const { data: clients, isLoading, error } = useActiveClients();

  if (isLoading)
    return (
      <>
        <Loading />;
      </>
    );
  if (error)
    return <div className="p-6 text-red-500">Error fetching clients.</div>;

  const activeClients =
    clients?.filter((client) => client.is_active === true) ?? [];

  return (
    <div className="w-full max-w-xs sm:max-w-sm md:max-w-md flex items-center justify-between rounded-xl border-l-4 border-primary p-4 shadow transition hover:shadow-md bg-white">
      <div className="flex flex-col gap-1">
        <span className="text-xs sm:text-sm  font-medium">
          Total Connected Devices
        </span>
        <h1 className="text-3xl sm:text-4xl font-bold ">
          {activeClients.length}
        </h1>
      </div>
      <div className="flex-shrink-0">
        <Users size={48} />
      </div>
    </div>
  );
}
