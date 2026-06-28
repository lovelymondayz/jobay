export interface NearbyAP {
  mac: string;
  name: string;
  proximity_percent: number;
  signal_strength: number;
}

export interface Clients {
  mac_address_client: string;
  mac_address_ap: string;
  ap_name: string;
  hostname: string;
  tx_bytes: number;
  rx_bytes: number;
  connected_at: string;
  disconnected_at: string | null;
  last_seen: string;
  is_active: boolean;
  signal_strength: number;
  site_id: string;
  nearby_aps: NearbyAP[];
}

export interface ClientRowProps {
  client: Clients;
  isExpanded: boolean;
  onToggle: () => void;
}
