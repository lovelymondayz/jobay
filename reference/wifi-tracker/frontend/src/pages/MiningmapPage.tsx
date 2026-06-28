import { Outlet } from "@tanstack/react-router";
import ClientTable from "@/components/ClientComponent/ClientTable";
import DeviceCountCard from "@/components/ClientComponent/DeviceCountCard";
import ClientsCountCard from "@/components/ClientComponent/ClientsCountCard";
import { TimeComponent } from "../components/ClientComponent/TimeComponent";

export function MiningmapPage() {
  return (
    <div className="starter">
      <div className="flex-col gap-2 lg:flex-row flex h-screen w-screen lg:w-full">
        <TimeComponent />
        <DeviceCountCard />
        <ClientsCountCard />
      </div>
      <Outlet />
      <div>
        <ClientTable />
      </div>
    </div>
  );
}
