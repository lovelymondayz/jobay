package dtos

import (
	"time"

	"github.com/google/uuid"
)

type ActiveClient struct {
	DeviceID         uuid.UUID  `json:"device_id"`
	MACAddressClient string     `json:"mac_address_client"`
	APID             uuid.UUID  `json:"ap_id"`
	MACAddressAP     string     `json:"mac_address_ap"`
	APName           string     `json:"ap_name"`
	UserName         string     `json:"user_name"`
	Hostname         string     `json:"hostname"`
	TxBytes          int64      `json:"tx_bytes"`
	RxBytes          int64      `json:"rx_bytes"`
	ConnectedAt      string     `json:"connected_at"`
	DisconnectedAt   *string    `json:"disconnected_at"`
	LastSeen         string     `json:"last_seen"`
	IsActive         bool       `json:"is_active"`
	SignalStrength   int        `json:"signal_strength"`
	NearbyAPs        []NearbyAP `json:"nearby_aps"`
	// CurrentAPCoord   Coordinate `json:"current_ap_coordinates"`
}

type Coordinate struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type APDistance struct {
	MAC      string
	Name     string
	Distance float64
	RSSI     int
	Coord    Coordinate
}
type NearbyAP struct {
	MAC  string `json:"mac"`
	Name string `json:"name"`
	// Distance         float64 `json:"distance_meters"`
	ProximityPercent float64 `json:"proximity_percent"`
	RSSI             int     `json:"rssi,omitempty"`
	SignalStrength   int     `json:"signal_strength"`
}

// APInfo berisi informasi AP
type APInfo struct {
	Name  string
	Model string
	X     float64 // Koordinat X jika tersedia
	Y     float64 // Koordinat Y jika tersedia
}
type APInfoResponse struct {
	Mac   string `json:"mac"`
	Name  string `json:"name"`
	Model string `json:"model"`
	State int    `json:"state"` // Status AP (1=connected, 0=disconnected)
}

// StaResponse matches /stat/sta/{mac} UniFi private API
type StaResponse struct {
	Mac       string `json:"mac"`
	ApMac     string `json:"ap_mac"`
	Signal    int    `json:"signal"`
	Hostname  string `json:"hostname"`
	TxBytes   int64  `json:"tx_bytes"`
	RxBytes   int64  `json:"rx_bytes"`
	FirstSeen int64  `json:"first_seen"`
	LastSeen  int64  `json:"last_seen"`
}

type ConnectionHistoriesResponse struct {
	ConnectionHistoryID uuid.UUID `json:"connection_history_id"`
	UserName            string    `json:"user_name"`
	DeviceName          string    `json:"device_name"`
	ToAps               string    `json:"to_aps"`
	FromAps             string    `json:"from_aps"`
	CreatedAt           time.Time `json:"created_at"`
}

type ConnectionHistoriesDataExist struct {
	ConnectionHistoryID uuid.UUID `json:"connection_history_id"`
	ToAps               uuid.UUID `json:"to_aps"`
	ToMacAddressDevice         string    `json:"to_mac_address_device"`
}
