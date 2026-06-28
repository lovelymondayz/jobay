import { useState } from "react";
import { useActiveClients } from "@/hooks/useActiveClients";
import { useActiveDevices } from "@/hooks/useActiveDevices";

import {
  CaretDown,
  CaretRight,
  WifiHigh,
  User,
  Clock,
  MapPin,
  Pulse,
} from "@phosphor-icons/react";

import type { ClientRowProps, Clients, NearbyAP } from "@/types/clientsStore";

import type {
  Device,
  DeviceRowProps,
  FloorSectionProps,
} from "@/types/deviceStore";

const getSignalColor = (strength: number) => {
  if (strength > -50) return "bg-emerald-500";
  if (strength > -60) return "bg-green-500";
  if (strength > -70) return "bg-yellow-500";
  if (strength > -80) return "bg-orange-500";
  return "bg-red-500";
};

const getSignalBars = (strength: number) => {
  if (strength > -50) return 4;
  if (strength > -60) return 3;
  if (strength > -70) return 2;
  if (strength > -80) return 1;
  return 0;
};

const SignalIndicator = ({ strength }: { strength: number }) => {
  const bars = getSignalBars(strength);
  const color = getSignalColor(strength);

  return (
    <div className="flex items-center gap-1">
      <div className="flex items-end gap-0.5">
        {[1, 2, 3, 4].map((bar) => (
          <div
            key={bar}
            className={`w-1 ${bar === 1 ? "h-1" : bar === 2 ? "h-2" : bar === 3 ? "h-3" : "h-4"} rounded-sm ${
              bar <= bars ? color : "bg-gray-300"
            }`}
          />
        ))}
      </div>
      <span className="text-xs text-gray-600 ml-1">{strength} dBm</span>
    </div>
  );
};

const ClientRow = ({ client, isExpanded, onToggle }: ClientRowProps) => {
  return (
    <div className="border border-gray-200 rounded-lg mb-3 bg-white shadow-sm hover:shadow-md transition-shadow min-w-full">
      <div className="flex items-center justify-between p-4 hover:bg-gray-50/50">
        <div className="flex items-center gap-3 min-w-0 flex-1">
          <div className="flex-shrink-0">
            <div className="w-10 h-10 bg-blue-100 rounded-full flex items-center justify-center">
              <User className="w-5 h-5 text-blue-600" />
            </div>
          </div>
          <div className="min-w-0 flex-1">
            <div className="font-medium text-gray-900 truncate">
              {client.hostname}
            </div>
            <div className="text-sm text-gray-500 font-mono">
              {client.mac_address_client}
            </div>
          </div>
        </div>

        <div className="flex items-center gap-4">
          <SignalIndicator strength={client.signal_strength} />
          <button
            onClick={onToggle}
            className="flex items-center gap-2 px-3 py-1.5 text-sm text-blue-600 hover:text-blue-700 hover:bg-blue-50 rounded-md transition-colors"
          >
            {isExpanded ? (
              <>
                <CaretDown className="w-4 h-4" />
              </>
            ) : (
              <>
                <CaretRight className="w-4 h-4" />
              </>
            )}
          </button>
        </div>
      </div>

      {isExpanded && (
        <div className="border-t bg-gray-50/30">
          <div className="p-4 space-y-4">
            {/* Last Seen Section */}
            <div className="flex items-center gap-2 text-sm">
              <Clock className="w-4 h-4 text-gray-500" />
              <span className="text-gray-600">Last seen:</span>
              <span className="inline-flex items-center px-2 py-1 rounded-md text-xs font-medium bg-gray-100 text-gray-800 font-mono">
                {client.last_seen}
              </span>
            </div>

            {/* Nearby Access Points */}
            <div>
              <div className="flex items-center gap-2 mb-3">
                <MapPin className="w-4 h-4 text-gray-600" />
                <h4 className="font-medium text-gray-800">
                  Nearby Access Points
                </h4>
              </div>

              {!client.nearby_aps || client.nearby_aps.length === 0 ? (
                <div className="text-gray-500 italic text-sm bg-white p-3 rounded-md border">
                  No nearby access points detected
                </div>
              ) : (
                <div className="space-y-2">
                  {client.nearby_aps.map((ap: NearbyAP, idx: number) => (
                    <div
                      key={idx}
                      className="bg-white border border-gray-200 rounded-lg shadow-sm"
                    >
                      <div className="p-3">
                        <div className="flex items-center justify-between">
                          <div className="min-w-0 flex-1">
                            <div className="font-medium text-gray-800 truncate">
                              {ap.name}
                            </div>
                            <div className="text-xs text-gray-500 font-mono truncate">
                              {ap.mac}
                            </div>
                          </div>
                          <div className="flex items-center gap-3">
                            <SignalIndicator strength={ap.signal_strength} />
                            <span className="inline-flex items-center px-2 py-1 rounded-md text-xs font-medium border border-gray-300 text-gray-700">
                              {ap.proximity_percent}%
                            </span>
                          </div>
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

const DeviceRow = ({ device, clients }: DeviceRowProps) => {
  const [isExpanded, setIsExpanded] = useState(false);
  const [expandedClient, setExpandedClient] = useState<string | null>(null);

  const deviceClients = clients.filter(
    (client: Clients) => client.mac_address_ap === device.mac
  );

  return (
    <div className="bg-white border border-gray-200 rounded-lg shadow-sm hover:shadow-md transition-shadow mb-4 w-full ">
      <div
        className="flex items-center justify-between p-4 cursor-pointer hover:bg-gray-50/50"
        onClick={() => setIsExpanded(!isExpanded)}
      >
        <div className="flex items-center gap-3 min-w-0 flex-1">
          <div className="flex-shrink-0">
            <div className="w-12 h-12 bg-gradient-to-br from-blue-500 to-blue-600 rounded-lg flex items-center justify-center">
              <WifiHigh className="w-6 h-6 text-white" />
            </div>
          </div>
          <div className="min-w-0 flex-1">
            <div className="flex items-center gap-2 mb-1">
              <h3 className="font-semibold text-gray-900 truncate">
                {device.name}
              </h3>
              <div
                className={`w-2 h-2 rounded-full ${device.state === 1 ? "bg-green-500" : "bg-red-500"}`}
              />
            </div>
            <div className="text-sm text-gray-500 font-mono">{device.mac}</div>
          </div>
        </div>

        <div className="flex items-center gap-4">
          <div className="text-right">
            <div className="flex items-center gap-2">
              <Pulse className="w-4 h-4 text-gray-500" />
              <span className="font-medium text-gray-900">
                {deviceClients.length}
              </span>
              <span className="text-sm text-gray-600">
                client{deviceClients.length !== 1 ? "s" : ""}
              </span>
            </div>
          </div>
          {/* <div className="flex-shrink-0">
            {isExpanded ? (
              <CaretDown className="w-5 h-5 text-gray-400" />
            ) : (
              <CaretRight className="w-5 h-5 text-gray-400" />
            )}
          </div> */}
        </div>
      </div>

      <div className="border-t bg-gray-50/30">
        <div className="p-4">
          {deviceClients.length === 0 ? (
            <div className="text-center py-8">
              <div className="w-16 h-16 bg-gray-100 rounded-full flex items-center justify-center mx-auto mb-3">
                <User className="w-8 h-8 text-gray-400" />
              </div>
              <p className="text-gray-500 italic">
                No clients connected to this access point
              </p>
            </div>
          ) : (
            <div className="space-y-3">
              <h4 className="font-medium text-gray-800 mb-3">
                Connected Clients ({deviceClients.length})
              </h4>
              {deviceClients.map((client: Clients) => (
                <ClientRow
                  key={client.mac_address_client}
                  client={client}
                  isExpanded={expandedClient === client.mac_address_client}
                  onToggle={() =>
                    setExpandedClient(
                      expandedClient === client.mac_address_client
                        ? null
                        : client.mac_address_client
                    )
                  }
                />
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

const FloorSection = ({ floorNumber, devices, clients }: FloorSectionProps) => {
  return (
    <div className="mb-8 w-full">
      <div className="mb-4">
        <div className="flex items-center gap-3 mb-2">
          <h2 className="text-xl font-bold text-gray-900">
            Floor {floorNumber}
          </h2>
        </div>
        <div className="h-1 bg-gradient-to-r from-indigo-500 to-purple-600 rounded-full w-20"></div>
      </div>

      <div className="space-y-4">
        {devices.length === 0 ? (
          <div className="bg-white border border-gray-200 rounded-lg shadow-sm p-8 text-center">
            <div className="w-16 h-16 bg-gray-100 rounded-full flex items-center justify-center mx-auto mb-4">
              <WifiHigh className="w-8 h-8 text-gray-400" />
            </div>
            <p className="text-gray-500">
              No access points found on this floor
            </p>
          </div>
        ) : (
          devices.map((device: Device) => (
            <DeviceRow key={device.mac} device={device} clients={clients} />
          ))
        )}
      </div>
    </div>
  );
};

export default function ClientTable() {
  const { data: clients, isLoading, error } = useActiveClients();
  const { data: devices } = useActiveDevices();

  if (isLoading) {
    return (
      <div className="p-6">
        <div className="animate-pulse space-y-4">
          <div className="h-8 bg-gray-200 rounded w-1/4"></div>
          <div className="space-y-3">
            {[1, 2, 3].map((i) => (
              <div key={i} className="h-20 bg-gray-200 rounded-lg"></div>
            ))}
          </div>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="p-6">
        <div className="bg-red-50 border border-red-200 rounded-lg">
          <div className="p-6 text-center">
            <div className="w-16 h-16 bg-red-100 rounded-full flex items-center justify-center mx-auto mb-4">
              {/* <Activity className="w-8 h-8 text-red-500" /> */}
            </div>
            <h3 className="font-semibold text-red-800 mb-2">
              Error Loading Data
            </h3>
            <p className="text-red-600">
              Unable to fetch client information. Please try again.
            </p>
          </div>
        </div>
      </div>
    );
  }

  // Group devices by floor
  const floor1Devices =
    devices?.filter((device) => device.name.includes("EXA")) || [];

  const floor2Devices =
    devices?.filter(
      (device) =>
        device.name.includes("AFRIYAN") || device.name.includes("ARJI")
    ) || [];

  const floor3Devices =
    devices?.filter((device) => device.name.includes("AFRIYAN")) || [];

  return (
    <div className="mt-10 w-5/6 lg:w-full lg:min-w-max">
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900 mb-2">
          Network Overview
        </h1>
        <p className="text-gray-600">
          Monitor access points and connected clients across all floors
        </p>
      </div>

      <div className="space-y-8">
        <FloorSection
          floorNumber={1}
          devices={floor1Devices}
          clients={clients || []}
        />
        <FloorSection
          floorNumber={2}
          devices={floor2Devices}
          clients={clients || []}
        />
        <FloorSection
          floorNumber={3}
          devices={floor3Devices}
          clients={clients || []}
        />
      </div>
    </div>
  );
}
