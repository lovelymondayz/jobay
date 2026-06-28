package service

import (
	"fmt"
	"time"

	"wifi-tracker-be/src/dtos"
	"wifi-tracker-be/src/models"

	"github.com/google/uuid"
	"github.com/unpoller/unifi"
)

type UnifiIntegrationServiceInterface interface {
	// GetClients() ([]models.Clients, error)
	GetClients() ([]dtos.ActiveClient, error)
	GetDevices() ([]models.Devices, error)
}

type UnifiIntegrationService struct {
	unifiClient *unifi.Unifi
}

func NewUnifiIntegrationService(unifiClient *unifi.Unifi) UnifiIntegrationServiceInterface {
	return &UnifiIntegrationService{
		unifiClient: unifiClient,
	}
}

/*
* Old

	func (s *UnifiIntegrationService) GetClients() ([]models.Clients, error) {
		sites, err := s.unifiClient.GetSites()
		if err != nil {
			return nil, err
		}

		if len(sites) == 0 {
			return nil, fmt.Errorf("no sites available")
		}

		sitesSlice := []*unifi.Site{sites[0]}
		clients, err := s.unifiClient.GetClients(sitesSlice)
		if err != nil {
			return nil, err
		}

		layout := "2006-01-02 15:04:05"
		var connections []models.Clients

		for _, client := range clients {
			connectedAt := time.Unix(int64(client.FirstSeen.Val), 0)
			lastSeen := time.Unix(int64(client.LastSeen.Val), 0)

			var disconnectedAt *string
			// isActive := true
			if time.Since(lastSeen) > 5*time.Minute {
				// isActive = false
				formatted := lastSeen.Format(layout)
				disconnectedAt = &formatted
			}

			connection := models.Clients{
				// ID:             i + 1,
				MACAddress:     client.Mac,
				APName:         client.ApName,
				TxBytes:        int64(client.TxBytes.Val),
				RxBytes:        int64(client.RxBytes.Val),
				ConnectedAt:    connectedAt.Format(layout),
				DisconnectedAt: disconnectedAt,
				// IsActive:       isActive,
			}

			connections = append(connections, connection)
		}
		return connections, nil
	}
*/
func (s *UnifiIntegrationService) GetClients() ([]dtos.ActiveClient, error) {
	sites, err := s.unifiClient.GetSites()
	if err != nil {
		return nil, err
	}

	var targetSite *unifi.Site
	for _, site := range sites {
		if site.Name == "default" || site.Desc == "default" {
			targetSite = site
			break
		}
	}
	if targetSite == nil && len(sites) > 0 {
		targetSite = sites[0]
	}
	if targetSite == nil {
		return nil, fmt.Errorf("no sites available")
	}

	sitesSlice := []*unifi.Site{targetSite}

	clients, err := s.unifiClient.GetClients(sitesSlice)
	if err != nil {
		return nil, err
	}

	devices, err := s.unifiClient.GetDevices(sitesSlice)
	if err != nil {
		return nil, err
	}

	// Build AP map
	apMap := make(map[string]string)
	for _, uap := range devices.UAPs {
		apMap[uap.Mac] = uap.Name
	}

	layout := "2006-01-02 15:04:05"
	var connections []dtos.ActiveClient

	for _, client := range clients {
		connectedAt := time.Unix(int64(client.FirstSeen.Val), 0)
		lastSeenTime := time.Unix(int64(client.LastSeen.Val), 0)
		lastSeen := lastSeenTime.Format(layout)

		var disconnectedAt *string
		isActive := true
		if time.Since(lastSeenTime) > 5*time.Minute {
			isActive = false
			formatted := lastSeenTime.Format(layout)
			disconnectedAt = &formatted
		}

		apName := "Unknown"
		if client.ApMac != "" {
			if name, exists := apMap[client.ApMac]; exists {
				apName = name
			}
		}

		ap, ok := findAPByMAC(devices.UAPs, client.ApMac)
		if !ok {
			continue // skip if no AP found
		}

		nearbyAPs, err := s.getNearbyAPs(ap, apMap)
		if err != nil {
			continue
		}

		connection := dtos.ActiveClient{
			MACAddressClient: client.Mac,
			MACAddressAP:     client.ApMac,
			APName:           apName,
			Hostname:         client.Hostname,
			TxBytes:          int64(client.TxBytes.Val),
			RxBytes:          int64(client.RxBytes.Val),
			ConnectedAt:      connectedAt.Format(layout),
			DisconnectedAt:   disconnectedAt,
			LastSeen:         lastSeen,
			IsActive:         isActive,
			NearbyAPs:        nearbyAPs,
			SignalStrength:   int(client.Signal.Val), // 👈 Add this line
		}

		connections = append(connections, connection)
	}

	return connections, nil
}

func findAPByMAC(aps []*unifi.UAP, mac string) (*unifi.UAP, bool) {
	for _, ap := range aps {
		if ap.Mac == mac {
			return ap, true
		}
	}
	return nil, false
}

func (s *UnifiIntegrationService) getNearbyAPs(currentAP *unifi.UAP, apMap map[string]string) ([]dtos.NearbyAP, error) {
	var nearby []dtos.NearbyAP

	// Get all APs
	allSites, _ := s.unifiClient.GetSites()
	devices, err := s.unifiClient.GetDevices(allSites)
	if err != nil {
		return nearby, err
	}

	for _, ap := range devices.UAPs {
		if ap.Mac == currentAP.Mac {
			continue
		}

		nearby = append(nearby, dtos.NearbyAP{
			MAC:  ap.Mac,
			Name: apMap[ap.Mac],
		})

		if len(nearby) == 2 {
			break // limit to 2 for now
		}
	}

	return nearby, nil
}

func (s *UnifiIntegrationService) GetDevices() ([]models.Devices, error) {
	sites, err := s.unifiClient.GetSites()
	if err != nil {
		return nil, err
	}

	var targetSite *unifi.Site
	for _, site := range sites {
		if site.Name == "default" || site.Desc == "default" {
			targetSite = site
			break
		}
	}
	if targetSite == nil && len(sites) > 0 {
		targetSite = sites[0]
	}
	if targetSite == nil {
		return nil, fmt.Errorf("no sites available")
	}

	sitesSlice := []*unifi.Site{targetSite}
	device, err := s.unifiClient.GetDevices(sitesSlice)
	if err != nil {
		return nil, err
	}
	var connections []models.Devices
	for _, uap := range device.UAPs {
		// Generate new UUIDs for each device
		deviceID := uuid.New()
		userID := uuid.New() // Replace with real user ID if needed
		connectedAt := time.Unix(int64(uap.ConnectedAt.Val), 0).Format("2006:01:02 15:04:05")

		connection := models.Devices{
			DeviceID:    deviceID,
			UserID:      userID,
			Name:        uap.Name,
			MacAddress:  uap.Mac,
			DeviceType:  uap.Type,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			ConnectedAt: connectedAt, // now a string
			CreatedBy:   "system",
			UpdatedBy:   "system",
		}
		connections = append(connections, connection)
	}

	return connections, nil
}
