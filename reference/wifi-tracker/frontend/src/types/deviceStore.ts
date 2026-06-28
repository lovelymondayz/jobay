import type { Clients } from "@/types/clientsStore";

export interface Device {
  device_id: string;
  user_id: string;
  name: string;
  mac_address: string;
  ConnectionHistories: string;
  device_type: string;
  connection_at: string;
  CreatedAt: string;
  updated_at: string;
  DeletedAt: string;
  created_by: string;
  updated_by: string;

  mac: string;
  model: string;
  state: number;

  // device_type: "mobile" | "laptop" | "desktop" | "tablet";
  is_active: number; // 0 or 1
  is_online: number; // 0 or 1
  status: "active" | "inactive" | "suspended";
  connected_at: string; // ISO date string
  created_at?: string; // Optional fields that might come from API
}
export interface DeviceRowProps {
  device: Device;
  clients: Clients[];
}
export interface FloorSectionProps {
  floorNumber: number;
  devices: Device[];
  clients: Clients[];
}

// types/deviceStore.ts

export interface DeviceFormProps {
  device?: Device | null;
  onSubmit: (data: Device) => void;
  onCancel: () => void;
}
