// import { HistoriesPage } from "@/pages/HistoriesPage";

import {
  House,
  ClockCounterClockwise,
  MapPinArea,
  MapPin,
  Truck,
  GearSix,
  FileText,
} from "@phosphor-icons/react";

import type { IconProps } from "@phosphor-icons/react";

type SidebarItem = {
  to: string;
  icon: React.ComponentType<IconProps>;
  label: string;
};

export const sidebarMenu: SidebarItem[] = [
  { to: "/", icon: House, label: "Home" },
  { to: "/histories", icon: ClockCounterClockwise, label: "Histories" },
  { to: "/miningmap", icon: MapPinArea, label: "Mining Map" },
  { to: "/miners-location", icon: MapPin, label: "Miners Location" },
  { to: "/dumping", icon: Truck, label: "Dumping" },
  { to: "/processing", icon: GearSix, label: "Processing" },
  // { to: "/report", icon: FileText, label: "Report" },
  { to: "/regist", icon: FileText, label: "Devices" },
  { to: "/regist-user", icon: FileText, label: "User" },
];
