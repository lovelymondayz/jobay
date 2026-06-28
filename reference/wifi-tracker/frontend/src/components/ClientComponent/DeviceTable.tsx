// import React from "react";
import { useActiveDevices } from "@/hooks/useActiveDevices";

export default function DeviceTable() {
  const { data: devices, isLoading, error } = useActiveDevices();

  if (isLoading) return <div className="p-6">Loading active devices...</div>;
  if (error)
    return <div className="p-6 text-red-500">Error fetching devices.</div>;

  // Dummy grouping floors, assigning first half to Floor 1 and second half to Floor 2
  const floorGroups = [
    {
      name: "Floor 1",
      mesh: devices ? devices[0] : null,
      clients: devices ? devices.slice(1, Math.ceil(devices.length / 2)) : [],
    },
    {
      name: "Floor 2",
      mesh: devices ? devices[Math.ceil(devices.length / 2)] : null,
      clients: devices ? devices.slice(Math.ceil(devices.length / 2) + 1) : [],
    },
  ];

  return (
    <div className="p-6">
      <h1 className="text-2xl font-bold mb-4">Active Mesh Points</h1>

      {floorGroups.map((floor) => (
        <div key={floor.name} className="mb-8">
          <h2 className="text-xl font-semibold mb-2 border-b pb-1">
            {floor.name}
          </h2>

          {floor.mesh ? (
            <div className="mb-4">
              <h3 className="text-lg font-medium">Mesh Point:</h3>
              <table className="min-w-full bg-white border border-gray-300 text-sm mb-2">
                <thead className="bg-gray-100 text-left">
                  <tr>
                    <th className="p-2 border">Hostname</th>
                    <th className="p-2 border">MAC Address</th>
                  </tr>
                </thead>
                <tbody>
                  <tr className="border-t">
                    <td className="p-2 border">{floor.mesh.name}</td>
                    <td className="p-2 border">{floor.mesh.model}</td>
                  </tr>
                </tbody>
              </table>
            </div>
          ) : (
            <p className="text-gray-500 mb-4">No mesh point for this floor.</p>
          )}

          <h3 className="text-lg font-medium">Connected Clients:</h3>
          {floor.clients.length > 0 ? (
            <table className="min-w-full bg-white border border-gray-300 text-sm">
              <thead className="bg-gray-100 text-left">
                <tr>
                  <th className="p-2 border">Hostname</th>
                  <th className="p-2 border">MAC Address</th>
                </tr>
              </thead>
              <tbody>
                {floor.clients.map((client) => (
                  <tr key={client.device_id} className="border-t">
                    <td className="p-2 border">{client.name}</td>
                    <td className="p-2 border">{client.model}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          ) : (
            <p className="text-gray-500">No clients connected on this floor.</p>
          )}
        </div>
      ))}
    </div>
  );
}
