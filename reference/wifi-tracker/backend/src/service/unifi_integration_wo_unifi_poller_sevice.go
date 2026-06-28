package service

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"wifi-tracker-be/src/config"
	"wifi-tracker-be/src/dtos"
	"wifi-tracker-be/src/models"
	"wifi-tracker-be/src/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UnifiIntegrationServiceWoUnifiPollerInterface interface {
	GetClientsWithNearbyAPs() ([]dtos.ActiveClient, error)
	GetDevices() ([]dtos.APInfoResponse, error)
	GetClientByMACAP(mac string) (*dtos.ActiveClient, error)
	SaveConnectionHistory(client dtos.ActiveClient, deviceID,fromAP, toAP uuid.UUID) error
	CheckIfConnectionLoggedToday(deviceID uuid.UUID, baypass bool,date time.Time) (bool, *dtos.ConnectionHistoriesDataExist, error)
	FindAllHistoryConnection(
		search string, 
		dateFrom string, 
		dateTo string,
		page int,
		size int,
	) (map[string]any, error)
}

type UnifiIntegrationServiceWoUnifiPoller struct {
	DB          *gorm.DB
	deviceRepo  repository.DeviceRepositoryInterface
	apRepo      repository.AccessPointRepositoryInterface
	userRepo    repository.UserRepositoryInterface
	historyRepo repository.ConnectionHistoriesRepositoryInterface
	httpClient  *http.Client
	baseURL     string
	site        string
	username    string
	password    string
}

func NewUnifiIntegrationServiceWoUnifiPoller(
	db *gorm.DB,
	deviceRepo repository.DeviceRepositoryInterface,
	apRepo repository.AccessPointRepositoryInterface,
	userRepo repository.UserRepositoryInterface,
	historyRepo repository.ConnectionHistoriesRepositoryInterface,
	client *http.Client,
) UnifiIntegrationServiceWoUnifiPollerInterface {
	return &UnifiIntegrationServiceWoUnifiPoller{
		DB:          db,
		deviceRepo:  deviceRepo,
		apRepo:      apRepo,
		userRepo:    userRepo,
		historyRepo: historyRepo,
		httpClient:  client,
		baseURL:     os.Getenv("UNIFI_CONTROLLER_URL"),
		site:        "default",
		username:    os.Getenv("UNIFI_USERNAME"),
		password:    os.Getenv("UNIFI_PASSWORD"),
	}
}

func (s *UnifiIntegrationServiceWoUnifiPoller) SaveConnectionHistory(client dtos.ActiveClient, deviceID,fromAP, toAP uuid.UUID) error {
	user, _ := s.userRepo.GetUserIdByDeviceMAC(s.DB, []string{client.MACAddressClient})
	apFrom, _ := s.apRepo.FindByID(s.DB, fromAP)
	apTo, _ := s.apRepo.FindByID(s.DB, toAP)

	var userID uuid.UUID

	if id, ok := user[client.MACAddressClient]; ok {
		userID = id
	}

	history := &models.ConnectionHistories{
		UserID:    userID,
		DeviceID:  deviceID,
		FromAPs:   apFrom.AccessPointID,
		ToAPs:     apTo.AccessPointID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		CreatedBy: "System",
		UpdatedBy: "System",
	}

	err := s.historyRepo.Create(s.DB, history)
	if err != nil {
		return err
	}

	return nil
}

func (s *UnifiIntegrationServiceWoUnifiPoller) CheckIfConnectionLoggedToday(deviceID uuid.UUID, baypass bool,date time.Time) (bool, *dtos.ConnectionHistoriesDataExist, error){
	if baypass {
		return false, &dtos.ConnectionHistoriesDataExist{}, nil
	}

	datas, err := s.historyRepo.ExistsToday(s.DB, deviceID, date)
	if err != nil {
		return true, &dtos.ConnectionHistoriesDataExist{}, err
	}

	return false, datas, nil
}

func (s *UnifiIntegrationServiceWoUnifiPoller) GetDevices() ([]dtos.APInfoResponse, error) {
	apsURL := fmt.Sprintf("%s/api/s/%s/stat/device", s.baseURL, s.site)
	req, err := http.NewRequest("GET", apsURL, nil)
	if err != nil {
		return nil, err
	}
	resp, err := s.doRequestWithAuthRetry(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to get access points, status: %d", resp.StatusCode)
	}

	apBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var apRaw struct {
		Data []dtos.APInfoResponse `json:"data"`
	}
	if err := json.Unmarshal(apBody, &apRaw); err != nil {
		return nil, fmt.Errorf("failed to unmarshal APs response: %w", err)
	}

	return apRaw.Data, nil
}

func (s *UnifiIntegrationServiceWoUnifiPoller) GetClientsWithNearbyAPs() ([]dtos.ActiveClient, error) {
	// Get all clients
	clients, err := s.getClients()
	if err != nil {
		return nil, fmt.Errorf("failed to get clients: %w", err)
	}

	// Get all APs with their signal information
	aps, err := s.getAccessPoints()
	if err != nil {
		return nil, fmt.Errorf("failed to get access points: %w", err)
	}

	// Get detailed signal information for each client from each AP
	clientSignals, err := s.getClientSignalData()
	if err != nil {
		return nil, fmt.Errorf("failed to get client signal data: %w", err)
	}

	var result []dtos.ActiveClient
	layout := "2006-01-02 15:04:05"

	userMap, apMap := s.getMappingAPtoDB(clients)

	deviceMap := s.getMappingDeviceToDB(clients)

	apIdMap := s.getAPIDFromDB(aps)

	for _, sta := range clients {
		connectedAt := time.Unix(sta.FirstSeen, 0).Format(layout)
		lastSeen := time.Unix(sta.LastSeen, 0)

		var disconnectedAt *string
		if time.Since(lastSeen) > 5*time.Minute {
			dis := lastSeen.Format(layout)
			disconnectedAt = &dis
		}

		// Find nearby APs based on signal strength
		nearby := s.findNearbyAPs(sta.Mac, sta.ApMac, sta.Signal, aps, clientSignals)

		// Get connected AP name

		userName := "user/device unknown not registered"
		if name, exists := userMap[sta.Mac]; exists {
			userName = name
		}

		aPName := "AP unknown not registered"
		if name, exists := apMap[sta.ApMac]; exists {
			aPName = name
		}

		// Get device ID
		deviceID := uuid.Nil // default: device unknown not registered
		if id, exists := deviceMap[sta.Mac]; exists {
			deviceID = id
		}

		// Get ap ID
		apID := uuid.Nil // default: device unknown not registered
		if id, exists := apIdMap[sta.ApMac]; exists {
			apID = id
		}

		result = append(result, dtos.ActiveClient{
			DeviceID:         deviceID,
			MACAddressClient: sta.Mac,
			APID:             apID,
			MACAddressAP:     sta.ApMac,
			APName:           aPName,
			UserName:         userName,
			Hostname:         sta.Hostname,
			TxBytes:          sta.TxBytes,
			RxBytes:          sta.RxBytes,
			SignalStrength:   sta.Signal,
			ConnectedAt:      connectedAt,
			DisconnectedAt:   disconnectedAt,
			LastSeen:         lastSeen.Format(layout),
			IsActive:         disconnectedAt == nil,
			NearbyAPs:        nearby,
		})
	}

	return result, nil
}

func (s *UnifiIntegrationServiceWoUnifiPoller) GetClientByMACAP(mac string) (*dtos.ActiveClient, error) {
	userName := "user/device unknown not registered"
	aPName := "AP unknown not registered"
	deviceID := uuid.Nil
	apID := uuid.Nil

	// Get all clients
	clients, err := s.getClients()
	if err != nil {
		return nil, fmt.Errorf("failed to get clients: %w", err)
	}

	userMap, apMap := s.getMappingAPtoDB(clients)

	deviceMap := s.getMappingDeviceToDB(clients)

	apIdMap := s.getAPIDFromDB(map[string]dtos.APInfo{
		mac: {},
	})

	// Cari client yang cocok dengan MAC
	var target *dtos.StaResponse
	for _, sta := range clients {
		if strings.EqualFold(sta.ApMac, mac) {
			target = &sta
			break
		}
		if name, exists := userMap[sta.Mac]; exists {
			userName = name
		}
		if name, exists := apMap[sta.ApMac]; exists {
			aPName = name
		}

		if id, exists := deviceMap[sta.Mac]; exists {
			deviceID = id
		}

		// Get ap ID
		if id, exists := apIdMap[sta.ApMac]; exists {
			apID = id
		}
	}
	if target == nil {
		return nil, fmt.Errorf("client connected to MAC %s not found", mac)
	}

	// Get all APs
	aps, err := s.getAccessPoints()
	if err != nil {
		return nil, fmt.Errorf("failed to get access points: %w", err)
	}

	// Get signal data dari AP lain
	clientSignals, err := s.getClientSignalData()
	if err != nil {
		return nil, fmt.Errorf("failed to get client signal data: %w", err)
	}

	// Format tanggal
	layout := "2006-01-02 15:04:05"
	connectedAt := time.Unix(target.FirstSeen, 0).Format(layout)
	lastSeen := time.Unix(target.LastSeen, 0)

	var disconnectedAt *string
	if time.Since(lastSeen) > 5*time.Minute {
		dis := lastSeen.Format(layout)
		disconnectedAt = &dis
	}

	// Dapatkan nearby APs
	nearby := s.findNearbyAPs(target.Mac, target.ApMac, target.Signal, aps, clientSignals)

	// Nama AP terkoneksi
	// connectedAPName := "Unknown"
	// if ap, exists := aps[target.ApMac]; exists {
	// 	connectedAPName = ap.Name
	// }

	// Get device ID


	return &dtos.ActiveClient{
		DeviceID:         deviceID,
		MACAddressClient: target.Mac,
		MACAddressAP:     target.ApMac,
		APID:             apID,
		APName:           aPName,
		Hostname:         target.Hostname,
		UserName:         userName,
		TxBytes:          target.TxBytes,
		RxBytes:          target.RxBytes,
		SignalStrength:   target.Signal,
		ConnectedAt:      connectedAt,
		DisconnectedAt:   disconnectedAt,
		LastSeen:         lastSeen.Format(layout),
		IsActive:         disconnectedAt == nil,
		NearbyAPs:        nearby,
	}, nil
}

func (s *UnifiIntegrationServiceWoUnifiPoller) FindAllHistoryConnection(
		search string, 
		dateFrom string, 
		dateTo string,
		page int,
		size int,
	) (map[string]any, error) {
	datas, totals, maxPages, err := s.historyRepo.FindAll(s.DB, search, dateFrom, dateTo, page, size)
	if err != nil {
		return nil, err
	}


	return map[string]any{
		"items": datas,
		"total_data": totals,
		"max_pages": maxPages,
		"page": page,
		"visible": len(datas),
		"size": size,
	}, nil
}

func (s *UnifiIntegrationServiceWoUnifiPoller) getMappingAPtoDB(clients []dtos.StaResponse) (map[string]string, map[string]string) {
	macClients := make([]string, 0)
	macAPs := make([]string, 0)

	for _, sta := range clients {
		macClients = append(macClients, sta.Mac)
		macAPs = append(macAPs, sta.ApMac)
	}

	// Ambil user & AP dari DB sekali saja
	userMap, _ := s.userRepo.GetUserByDeviceMAC(s.DB, macClients)
	apMap, _ := s.apRepo.GetAccessPointByMAC(s.DB, macAPs)

	return userMap, apMap
}

func (s *UnifiIntegrationServiceWoUnifiPoller) getMappingDeviceToDB(clients []dtos.StaResponse) map[string]uuid.UUID {
	macClients := make([]string, 0)

	for _, sta := range clients {
		macClients = append(macClients, sta.Mac)
	}

	// Ambil user & AP dari DB sekali saja
	deviceMap, _ := s.deviceRepo.GetUDeviceIDByMAC(s.DB, macClients)

	return deviceMap
}

// func (s *UnifiIntegrationServiceWoUnifiPoller) GetMappingDeviceIdToDBPub(clients []dtos.StaResponse) (map[string]uuid.UUID) {
// 	macClients := make([]string, 0)

// 	for _, sta := range clients {
// 		macClients = append(macClients, sta.Mac)
// 	}

// 	// Ambil user & AP dari DB sekali saja
// 	deviceMap, _ := s.deviceRepo.GetUDeviceIDByMAC(s.DB, macClients)

// 	return deviceMap
// }

func (s *UnifiIntegrationServiceWoUnifiPoller) getAPNamesFromDB(aps map[string]dtos.APInfo) map[string]string {
	macAPs := make([]string, 0, len(aps))
	for mac := range aps {
		macAPs = append(macAPs, mac)
	}

	apMap, _ := s.apRepo.GetAccessPointByMAC(s.DB, macAPs)
	return apMap
}
func (s *UnifiIntegrationServiceWoUnifiPoller) getAPIDFromDB(aps map[string]dtos.APInfo) map[string]uuid.UUID {
	macAPs := make([]string, 0, len(aps))
	for mac := range aps {
		macAPs = append(macAPs, mac)
	}

	apMap, _ := s.apRepo.GetAccessPointIdByMAC(s.DB, macAPs)
	return apMap
}

func (s *UnifiIntegrationServiceWoUnifiPoller) getClients() ([]dtos.StaResponse, error) {
	clientsURL := fmt.Sprintf("%s/api/s/%s/stat/sta", s.baseURL, s.site)
	req, err := http.NewRequest("GET", clientsURL, nil)
	if err != nil {
		return nil, err
	}
	resp, err := s.doRequestWithAuthRetry(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to get clients, status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var raw struct {
		Data []dtos.StaResponse `json:"data"`
	}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("failed to unmarshal clients response: %w", err)
	}

	return raw.Data, nil
}

func (s *UnifiIntegrationServiceWoUnifiPoller) getAccessPoints() (map[string]dtos.APInfo, error) {
	apsURL := fmt.Sprintf("%s/api/s/%s/stat/device", s.baseURL, s.site)
	req, err := http.NewRequest("GET", apsURL, nil)
	if err != nil {
		return nil, err
	}
	resp, err := s.doRequestWithAuthRetry(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to get access points, status: %d", resp.StatusCode)
	}

	apBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var apRaw struct {
		Data []struct {
			Mac   string  `json:"mac"`
			Name  string  `json:"name"`
			Model string  `json:"model"`
			X     float64 `json:"x,omitempty"` // Koordinat X jika tersedia
			Y     float64 `json:"y,omitempty"` // Koordinat Y jika tersedia
			State int     `json:"state"`       // Status AP (1=connected, 0=disconnected)
		} `json:"data"`
	}
	if err := json.Unmarshal(apBody, &apRaw); err != nil {
		return nil, fmt.Errorf("failed to unmarshal APs response: %w", err)
	}

	apMap := make(map[string]dtos.APInfo)
	for _, ap := range apRaw.Data {
		// Hanya masukkan AP yang aktif
		if ap.State == 1 {
			apMap[ap.Mac] = dtos.APInfo{
				Name:  ap.Name,
				Model: ap.Model,
				X:     ap.X,
				Y:     ap.Y,
			}
		}
	}

	return apMap, nil
}

func (s *UnifiIntegrationServiceWoUnifiPoller) getClientSignalData() (map[string]map[string]int, error) {
	// Coba mendapatkan data sinyal yang lebih akurat dari endpoint lain
	clientSignals := make(map[string]map[string]int)

	// Endpoint alternatif untuk mendapatkan scan results
	scanURL := fmt.Sprintf("%s/api/s/%s/stat/sta", s.baseURL, s.site)
	req, err := http.NewRequest("GET", scanURL, nil)
	if err != nil {
		return nil, err
	}
	resp, err := s.doRequestWithAuthRetry(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return clientSignals, nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return clientSignals, nil
	}

	var raw struct {
		Data []struct {
			Mac        string `json:"mac"`
			ApMac      string `json:"ap_mac"`
			Signal     int    `json:"signal"`
			NoiseFloor int    `json:"noise,omitempty"`
			Rssi       int    `json:"rssi,omitempty"`
			// Scan results jika tersedia
			ScanTable []struct {
				Bssid  string `json:"bssid"`
				Rssi   int    `json:"rssi"`
				Signal int    `json:"signal"`
			} `json:"scan_table,omitempty"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &raw); err != nil {
		return clientSignals, nil
	}

	// Parse scan results jika tersedia
	for _, sta := range raw.Data {
		if len(sta.ScanTable) > 0 {
			if clientSignals[sta.Mac] == nil {
				clientSignals[sta.Mac] = make(map[string]int)
			}

			for _, scan := range sta.ScanTable {
				// Gunakan RSSI jika tersedia, jika tidak gunakan signal
				signalValue := scan.Rssi
				if signalValue == 0 {
					signalValue = scan.Signal
				}
				clientSignals[sta.Mac][scan.Bssid] = signalValue
			}
		}
	}

	return clientSignals, nil
}

/**
func (s *UnifiIntegrationServiceWoUnifiPoller) findNearbyAPs(
	clientMAC, connectedAPMAC string,
	clientSignal int,
	aps map[string]dtos.APInfo,
	clientSignals map[string]map[string]int,
) []dtos.NearbyAP {

	var nearby []dtos.NearbyAP
	apNameMap := s.getAPNamesFromDB(aps)

	if signals, exists := clientSignals[clientMAC]; exists {
		for apMAC, signal := range signals {
			if apMAC == connectedAPMAC {
				continue
			}
			if _, ok := aps[apMAC]; ok {
				nearby = append(nearby, dtos.NearbyAP{
					MAC:              apMAC,
					Name:             s.fallbackName(apNameMap, apMAC),
					SignalStrength:   signal,
					ProximityPercent: s.calculateProximityFromSignal(signal, clientSignal),
				})
			}
		}
	} else {
		for apMAC := range aps {
			if apMAC == connectedAPMAC {
				continue
			}
			estimatedSignal := s.estimateSignalStrength(clientSignal, apMAC)
			proximity := s.calculateProximityFromSignal(estimatedSignal, clientSignal)

			if proximity > 10.0 {
				nearby = append(nearby, dtos.NearbyAP{
					MAC:              apMAC,
					Name:             s.fallbackName(apNameMap, apMAC),
					SignalStrength:   estimatedSignal,
					ProximityPercent: proximity,
				})
			}
		}
	}

	sort.Slice(nearby, func(i, j int) bool {
		return nearby[i].ProximityPercent > nearby[j].ProximityPercent
	})

	if len(nearby) > 2 {
		return nearby[:2]
	}
	return nearby
}
*/

func (s *UnifiIntegrationServiceWoUnifiPoller) findNearbyAPs(
	clientMAC, connectedAPMAC string,
	clientSignal int,
	aps map[string]dtos.APInfo,
	clientSignals map[string]map[string]int,
) []dtos.NearbyAP {
	var nearby []dtos.NearbyAP

	// Normalisasi MAC
	clientMAC = strings.ToLower(clientMAC)
	connectedAPMAC = strings.ToLower(connectedAPMAC)

	log.Printf("Finding nearby APs for client %s connected to %s with signal %d", 
		clientMAC, connectedAPMAC, clientSignal)

	// Ambil nama AP dari DB
	apNameMap := s.getAPNamesFromDB(aps)

	// Ambil peta mesh
	uplinkMap, downlinkMap, err := s.getMeshUplinkAndDownlinkMap()
	if err != nil {
		log.Printf("Error getting mesh uplink/downlink map: %v", err)
		// Fallback: coba cari nearby AP berdasarkan signal saja
		return s.findNearbyAPsBySignalOnly(clientMAC, connectedAPMAC, clientSignal, aps, clientSignals, apNameMap)
	}

	// Buat set AP mesh yang terhubung dengan connectedAP
	meshConnectedAPs := make(map[string]bool)

	// 1. Jika connectedAP adalah mesh uplink device, tambahkan parent-nya
	if parentMAC, hasParent := uplinkMap[connectedAPMAC]; hasParent {
		parentMAC = strings.ToLower(parentMAC)
		meshConnectedAPs[parentMAC] = true
		log.Printf("Added parent AP: %s for mesh uplink: %s", parentMAC, connectedAPMAC)
		
		// Tambahkan juga sibling dari parent yang sama
		if siblings, hasSiblings := downlinkMap[parentMAC]; hasSiblings {
			for _, sibling := range siblings {
				sibling = strings.ToLower(sibling)
				if sibling != connectedAPMAC {
					meshConnectedAPs[sibling] = true
					log.Printf("Added sibling AP: %s (same parent: %s)", sibling, parentMAC)
				}
			}
		}
	}

	// 2. Jika connectedAP adalah parent dari mesh downlink devices, tambahkan children-nya
	if children, hasChildren := downlinkMap[connectedAPMAC]; hasChildren {
		for _, child := range children {
			child = strings.ToLower(child)
			meshConnectedAPs[child] = true
			log.Printf("Added child AP: %s for mesh parent: %s", child, connectedAPMAC)
		}
	}

	// 3. Jika connectedAP adalah mesh device, tambahkan sibling-nya (yang punya parent yang sama)
	if parentMAC, hasParent := uplinkMap[connectedAPMAC]; hasParent {
		parentMAC = strings.ToLower(parentMAC)
		if siblings, hasSiblings := downlinkMap[parentMAC]; hasSiblings {
			for _, sibling := range siblings {
				sibling = strings.ToLower(sibling)
				if sibling != connectedAPMAC {
					meshConnectedAPs[sibling] = true
					log.Printf("Added sibling AP: %s for mesh device: %s", sibling, connectedAPMAC)
				}
			}
		}
	}

	log.Printf("Found %d mesh-connected APs for %s", len(meshConnectedAPs), connectedAPMAC)

	// Jika tidak ada mesh AP yang ditemukan, fallback ke pencarian berdasarkan signal
	if len(meshConnectedAPs) == 0 {
		log.Printf("No mesh-connected APs found, using signal-based fallback")
		return s.findNearbyAPsBySignalOnly(clientMAC, connectedAPMAC, clientSignal, aps, clientSignals, apNameMap)
	}

	// Proses AP yang terhubung secara mesh
	if signals, exists := clientSignals[clientMAC]; exists {
		// Gunakan sinyal aktual jika tersedia
		for meshAP := range meshConnectedAPs {
			if signal, hasSignal := signals[meshAP]; hasSignal {
				if _, apExists := aps[meshAP]; apExists {
					proximity := s.calculateProximityFromSignal(signal, clientSignal)
					nearby = append(nearby, dtos.NearbyAP{
						MAC:              meshAP,
						Name:             s.fallbackName(apNameMap, meshAP),
						SignalStrength:   signal,
						ProximityPercent: proximity,
					})
					log.Printf("Added nearby AP with actual signal: %s (signal: %d, proximity: %.1f%%)", 
						meshAP, signal, proximity)
				} else {
					log.Printf("Mesh AP %s not found in aps map", meshAP)
				}
			} else {
				log.Printf("No signal data for mesh AP %s", meshAP)
			}
		}
	}

	// Jika tidak cukup AP dengan sinyal aktual, coba estimasi
	if len(nearby) < 2 {
		log.Printf("Not enough APs with actual signals (%d), trying estimation", len(nearby))
		
		for meshAP := range meshConnectedAPs {
			// Skip jika sudah ditambahkan dengan sinyal aktual
			alreadyAdded := false
			for _, existingAP := range nearby {
				if existingAP.MAC == meshAP {
					alreadyAdded = true
					break
				}
			}
			
			if !alreadyAdded {
				if _, apExists := aps[meshAP]; apExists {
					estSignal := s.estimateSignalStrength(clientSignal, meshAP)
					proximity := s.calculateProximityFromSignal(estSignal, clientSignal)
					
					// Hanya tambahkan jika proximity cukup tinggi
					if proximity > 15.0 {
						nearby = append(nearby, dtos.NearbyAP{
							MAC:              meshAP,
							Name:             s.fallbackName(apNameMap, meshAP),
							SignalStrength:   estSignal,
							ProximityPercent: proximity,
						})
						log.Printf("Added nearby AP with estimated signal: %s (signal: %d, proximity: %.1f%%)", 
							meshAP, estSignal, proximity)
					} else {
						log.Printf("Skipped mesh AP %s due to low proximity: %.1f%%", meshAP, proximity)
					}
				}
			}
		}
	}

	// Urutkan berdasarkan proximity tertinggi
	sort.Slice(nearby, func(i, j int) bool {
		return nearby[i].ProximityPercent > nearby[j].ProximityPercent
	})

	// Batasi maksimum 2 AP
	if len(nearby) > 2 {
		nearby = nearby[:2]
	}

	log.Printf("Final nearby APs count: %d", len(nearby))
	for i, ap := range nearby {
		log.Printf("Nearby AP %d: %s (%s) - Signal: %d, Proximity: %.1f%%", 
			i+1, ap.MAC, ap.Name, ap.SignalStrength, ap.ProximityPercent)
	}

	return nearby
}

// Fallback function untuk mencari nearby AP berdasarkan signal saja
func (s *UnifiIntegrationServiceWoUnifiPoller) findNearbyAPsBySignalOnly(
	clientMAC, connectedAPMAC string,
	clientSignal int,
	aps map[string]dtos.APInfo,
	clientSignals map[string]map[string]int,
	apNameMap map[string]string,
) []dtos.NearbyAP {
	var nearby []dtos.NearbyAP

	log.Printf("Using signal-only fallback for client %s", clientMAC)

	// Cek apakah ada data sinyal untuk client ini
	if signals, exists := clientSignals[clientMAC]; exists {
		for apMAC, signal := range signals {
			apMAC = strings.ToLower(apMAC)
			
			// Skip AP yang sedang terhubung
			if apMAC == connectedAPMAC {
				continue
			}
			
			// Pastikan AP ada dalam daftar
			if _, apExists := aps[apMAC]; apExists {
				proximity := s.calculateProximityFromSignal(signal, clientSignal)
				
				// Hanya tambahkan jika proximity cukup tinggi
				if proximity > 10.0 {
					nearby = append(nearby, dtos.NearbyAP{
						MAC:              apMAC,
						Name:             s.fallbackName(apNameMap, apMAC),
						SignalStrength:   signal,
						ProximityPercent: proximity,
					})
					log.Printf("Added fallback AP: %s (signal: %d, proximity: %.1f%%)", 
						apMAC, signal, proximity)
				}
			}
		}
	}

	// Urutkan berdasarkan proximity tertinggi
	sort.Slice(nearby, func(i, j int) bool {
		return nearby[i].ProximityPercent > nearby[j].ProximityPercent
	})

	// Batasi maksimum 2 AP
	if len(nearby) > 2 {
		nearby = nearby[:2]
	}

	log.Printf("Fallback found %d nearby APs", len(nearby))
	return nearby
}


func (s *UnifiIntegrationServiceWoUnifiPoller) calculateProximityFromSignal(apSignal, clientSignal int) float64 {
	// Hitung selisih absolut sinyal
	diff := math.Abs(float64(apSignal - clientSignal))

	// Skala proximity berdasarkan selisih sinyal
	var proximity float64

	switch {
	case diff <= 2:
		proximity = 85.0 + (2.0-diff)*5.0 // 85-95% untuk selisih 0-2 dBm
	case diff <= 5:
		proximity = 70.0 + (5.0-diff)*5.0 // 70-85% untuk selisih 2-5 dBm
	case diff <= 10:
		proximity = 50.0 + (10.0-diff)*4.0 // 50-70% untuk selisih 5-10 dBm
	case diff <= 20:
		proximity = 25.0 + (20.0-diff)*2.5 // 25-50% untuk selisih 10-20 dBm
	case diff <= 35:
		proximity = 5.0 + (35.0-diff)*1.33 // 5-25% untuk selisih 20-35 dBm
	default:
		proximity = 0.0 // > 35 dBm selisih = tidak nearby
	}

	// Faktor koreksi berdasarkan kekuatan sinyal AP
	strengthFactor := 1.0
	if apSignal > -50 {
		strengthFactor = 1.1 // Bonus untuk sinyal kuat
	} else if apSignal < -80 {
		strengthFactor = 0.8 // Penalty untuk sinyal lemah
	}

	proximity *= strengthFactor

	// Batasi range 0-95%
	if proximity > 95.0 {
		proximity = 95.0
	}
	if proximity < 0 {
		proximity = 0
	}

	return math.Round(proximity*100) / 100
}

func (s *UnifiIntegrationServiceWoUnifiPoller) estimateSignalStrength(baseSignal int, apMAC string) int {
	// Buat hash sederhana dari MAC address untuk variasi yang konsisten
	hash := 0
	for i, b := range apMAC {
		hash += int(b) * (i + 1) // Berikan bobot berbeda untuk setiap posisi
	}

	// Buat variasi yang lebih beragam berdasarkan hash MAC
	variation1 := (hash % 25) - 12       // Variasi -12 sampai +12 dBm
	variation2 := ((hash >> 3) % 15) - 7 // Variasi tambahan -7 sampai +7 dBm

	totalVariation := variation1 + variation2

	// Degradasi jarak yang berbeda untuk setiap AP
	// AP yang berbeda memiliki jarak yang berbeda dari client
	distanceFactor := 5 + (hash % 20) // 5-25 dBm degradasi

	estimatedSignal := baseSignal - distanceFactor + totalVariation

	// Pastikan AP nearby selalu lebih lemah dari yang terkoneksi
	minWeakness := 3 + (hash % 10) // Minimal 3-13 dBm lebih lemah
	if estimatedSignal > baseSignal-minWeakness {
		estimatedSignal = baseSignal - minWeakness
	}

	// Batasi dalam range yang masuk akal
	if estimatedSignal > -35 {
		estimatedSignal = -35
	}
	if estimatedSignal < -90 {
		estimatedSignal = -90
	}

	return estimatedSignal
}

func (s *UnifiIntegrationServiceWoUnifiPoller) doRequestWithAuthRetry(req *http.Request) (*http.Response, error) {
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusUnauthorized {
		// Token mungkin expired, login ulang
		resp.Body.Close() // Penting!
		config.InitUniFiClient2()
		// Retry request setelah login ulang
		resp, err = s.httpClient.Do(req)
	}

	return resp, err
}

func (s *UnifiIntegrationServiceWoUnifiPoller) fallbackName(m map[string]string, key string) string {
	if name, exists := m[key]; exists {
		return name
	}
	return "AP unknown not registered"
}

func (s *UnifiIntegrationServiceWoUnifiPoller) getMeshUplinkAndDownlinkMap() (
	uplinkMap map[string]string,
	downlinkMap map[string][]string,
	err error,
) {
	uplinkMap = make(map[string]string)
	downlinkMap = make(map[string][]string)

	url := fmt.Sprintf("%s/api/s/%s/stat/device", s.baseURL, s.site)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.doRequestWithAuthRetry(req)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	// Struktur yang benar sesuai data JSON yang diberikan
	var raw struct {
		Data []struct {
			// Fields dari device utama
			Serialno string `json:"serialno"` // MAC address device utama
			Mac      string `json:"mac"`      // MAC address device utama
			Name     string `json:"name"`     // Nama device
			Type     string `json:"type"`     // Tipe device
			
			// Uplink table berada di dalam setiap device
			UplinkTable []struct {
				Serialno       string `json:"serialno"`        // MAC address mesh device
				UplinkMac      string `json:"uplink_mac"`      // MAC address parent
				Type           string `json:"type"`            // connection type (wireless)
				Connected      bool   `json:"connected"`       // status koneksi
				ApMac          string `json:"ap_mac"`          // MAC address AP
				UplinkSource   string `json:"uplink_source"`   // sumber uplink
				Radio          string `json:"radio"`           // radio yang digunakan
				Essid          string `json:"essid"`           // SSID mesh
				ApConnected    bool   `json:"ap_connected"`    // status AP connected
				Available      bool   `json:"available"`       // status available
				UplinkDeviceName string `json:"uplink_device_name"` // nama device parent
			} `json:"uplink_table"`
			
			// Downlink table berada di dalam setiap device
			DownlinkTable []struct {
				Serialno string `json:"serialno"` // MAC address mesh device
				ApMac    string `json:"ap_mac"`   // MAC address AP yang melayani
				Mac      string `json:"mac"`      // MAC address device
				Type     string `json:"type"`     // connection type
				Radio    string `json:"radio"`    // radio yang digunakan
			} `json:"downlink_table"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, nil, fmt.Errorf("failed to decode respons: %w", err)
	}

	totalUplinks := 0
	totalDownlinks := 0

	// Iterasi melalui setiap device dalam data
	for _, device := range raw.Data {

		// Process uplink table untuk device ini
		for _, uplink := range device.UplinkTable {
			// PERBAIKAN: Gunakan MAC address device utama sebagai child device
			meshDeviceMAC := strings.ToLower(device.Mac)
			if meshDeviceMAC == "" {
				meshDeviceMAC = strings.ToLower(device.Serialno)
			}
			
			// Gunakan uplink_mac sebagai parent, fallback ke ap_mac
			parentMAC := strings.ToLower(uplink.UplinkMac)
			if parentMAC == "" {
				parentMAC = strings.ToLower(uplink.ApMac)
			}

			// Kondisi untuk mesh uplink:
			// 1. Type harus wireless ATAU uplink_source mengandung "lldp_uplink" atau "wireless"
			// 2. Connected harus true
			// 3. Parent MAC tidak boleh sama dengan device MAC (bukan self-reference)
			// 4. Parent MAC tidak boleh kosong
			isWirelessMesh := uplink.Type == "wireless" || 
							strings.Contains(uplink.UplinkSource, "lldp_uplink") ||
							strings.Contains(uplink.UplinkSource, "wireless")
			
			if isWirelessMesh && uplink.Connected && 
			   parentMAC != "" && meshDeviceMAC != parentMAC {
				uplinkMap[meshDeviceMAC] = parentMAC
				
				// Tambahkan ke downlink map
				if !contains(downlinkMap[parentMAC], meshDeviceMAC) {
					downlinkMap[parentMAC] = append(downlinkMap[parentMAC], meshDeviceMAC)
				}
				totalUplinks++
			} else {
				log.Printf("Skipped uplink: device=%s, parent=%s, reason: type=%s, connected=%t, same_mac=%t, empty_parent=%t", 
					meshDeviceMAC, parentMAC, uplink.Type, uplink.Connected, 
					meshDeviceMAC == parentMAC, parentMAC == "")
			}
		}

		// Process downlink table untuk device ini
		for _, downlink := range device.DownlinkTable {
			// PERBAIKAN: Gunakan MAC address dari downlink entry sebagai child device
			meshDeviceMAC := strings.ToLower(downlink.Mac)
			if meshDeviceMAC == "" {
				meshDeviceMAC = strings.ToLower(downlink.Serialno)
			}
			
			// Gunakan MAC address device utama sebagai parent (karena downlink menunjukkan child dari device ini)
			parentMAC := strings.ToLower(device.Mac)
			if parentMAC == "" {
				parentMAC = strings.ToLower(device.Serialno)
			}

			// Kondisi untuk mesh downlink:
			// 1. Type bisa wireless atau kosong (karena beberapa data punya type kosong)
			// 2. Radio harus ada (na/ng untuk wireless)
			// 3. Device MAC tidak boleh sama dengan parent MAC
			// 4. Parent MAC tidak boleh kosong
			isWirelessMesh := downlink.Type == "wireless" || 
							(downlink.Type == "" && downlink.Radio != "")
			
			if isWirelessMesh && meshDeviceMAC != "" && parentMAC != "" && 
			   meshDeviceMAC != parentMAC {
				// Downlink menunjukkan bahwa meshDeviceMAC adalah child dari parentMAC
				uplinkMap[meshDeviceMAC] = parentMAC
				
				// Tambahkan ke downlink map
				if !contains(downlinkMap[parentMAC], meshDeviceMAC) {
					downlinkMap[parentMAC] = append(downlinkMap[parentMAC], meshDeviceMAC)
				}
				totalDownlinks++
			} else {
				log.Printf("Skipped downlink: device=%s, parent=%s, reason: type=%s, radio=%s, same_mac=%t, empty_fields=%t", 
					meshDeviceMAC, parentMAC, downlink.Type, downlink.Radio,
					meshDeviceMAC == parentMAC, meshDeviceMAC == "" || parentMAC == "")
			}
		}
	}

	log.Printf("Processed %d devices with %d uplink entries and %d downlink entries", 
		len(raw.Data), totalUplinks, totalDownlinks)
	log.Printf("Final Uplink map: %v", uplinkMap)
	log.Printf("Final Downlink map: %v", downlinkMap)

	return uplinkMap, downlinkMap, nil
}

// Helper function untuk mengecek apakah slice contains string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// TAMBAHAN: Function untuk menggabungkan semua device yang terhubung dalam mesh
// func (s *UnifiIntegrationServiceWoUnifiPoller) getMeshedDevices() (map[string][]string, error) {
// 	uplinkMap, downlinkMap, err := s.getMeshUplinkAndDownlinkMap()
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Gabungkan semua device yang terhubung dalam mesh
// 	meshedDevices := make(map[string][]string)
	
// 	// Tambahkan semua device yang memiliki uplink (child devices)
// 	for childMAC, parentMAC := range uplinkMap {
// 		// Tambahkan ke parent's list
// 		if !contains(meshedDevices[parentMAC], childMAC) {
// 			meshedDevices[parentMAC] = append(meshedDevices[parentMAC], childMAC)
// 		}
		
// 		// Tambahkan ke child's list (bidirectional relationship)
// 		if !contains(meshedDevices[childMAC], parentMAC) {
// 			meshedDevices[childMAC] = append(meshedDevices[childMAC], parentMAC)
// 		}
// 	}
	
// 	// Tambahkan dari downlink map
// 	for parentMAC, childrenMACs := range downlinkMap {
// 		for _, childMAC := range childrenMACs {
// 			// Tambahkan ke parent's list
// 			if !contains(meshedDevices[parentMAC], childMAC) {
// 				meshedDevices[parentMAC] = append(meshedDevices[parentMAC], childMAC)
// 			}
			
// 			// Tambahkan ke child's list (bidirectional relationship)
// 			if !contains(meshedDevices[childMAC], parentMAC) {
// 				meshedDevices[childMAC] = append(meshedDevices[childMAC], parentMAC)
// 			}
// 		}
// 	}
	
// 	log.Printf("Meshed devices map: %v", meshedDevices)
// 	return meshedDevices, nil
// }

// Helper function untuk normalisasi MAC address comparison
// func (s *UnifiIntegrationServiceWoUnifiPoller) normalizeMACForComparison(mac string) string {
// 	return strings.ToLower(strings.ReplaceAll(mac, ":", ""))
// }

// Improved version with better MAC address handling
// func (s *UnifiIntegrationServiceWoUnifiPoller) findNearbyAPsImproved(
// 	clientMAC, connectedAPMAC string,
// 	clientSignal int,
// 	aps map[string]dtos.APInfo,
// 	clientSignals map[string]map[string]int,
// ) []dtos.NearbyAP {
// 	var nearby []dtos.NearbyAP

// 	// Normalize MAC addresses for consistent comparison
// 	clientMAC = strings.ToLower(clientMAC)
// 	connectedAPMAC = strings.ToLower(connectedAPMAC)

// 	// Ambil nama AP dari database
// 	apNameMap := s.getAPNamesFromDB(aps)

// 	// Ambil peta mesh uplink
// 	meshUplinkMap, err := s.getMeshUplinkMap()
// 	if err != nil {
// 		log.Printf("Error getting mesh uplink map: %v", err)
// 		return nearby
// 	}

// 	// Buat set AP yang berhubungan mesh dengan connected AP
// 	meshRelatedAPs := make(map[string]bool)
	
// 	// 1. Cari AP mesh children yang uplink ke connected AP
// 	for meshMAC, uplinkMAC := range meshUplinkMap {
// 		if strings.EqualFold(uplinkMAC, connectedAPMAC) {
// 			meshRelatedAPs[meshMAC] = true
// 		}
// 	}

// 	// 2. Cari AP parent (root) jika connected AP adalah mesh child
// 	if parentMAC, isChild := meshUplinkMap[connectedAPMAC]; isChild {
// 		meshRelatedAPs[parentMAC] = true
		
// 		// 3. Cari sibling AP (AP lain yang uplink ke parent yang sama)
// 		for meshMAC, uplinkMAC := range meshUplinkMap {
// 			if strings.EqualFold(uplinkMAC, parentMAC) && !strings.EqualFold(meshMAC, connectedAPMAC) {
// 				meshRelatedAPs[meshMAC] = true
// 			}
// 		}
// 	}

// 	// 4. Proses hanya AP yang berhubungan mesh
// 	if signals, exists := clientSignals[clientMAC]; exists {
// 		// Gunakan data sinyal aktual jika tersedia
// 		for meshMAC := range meshRelatedAPs {
// 			if signal, hasSignal := signals[meshMAC]; hasSignal {
// 				// Cari AP info dengan case-insensitive comparison
// 				// var apInfo dtos.APInfo
// 				var apExists bool
// 				for apMAC, info := range aps {
// 					if strings.EqualFold(apMAC, meshMAC) {
// 						log.Printf("AP info: %v", info)
// 						// apInfo = info
// 						apExists = true
// 						break
// 					}
// 				}
				
// 				if apExists {
// 					nearby = append(nearby, dtos.NearbyAP{
// 						MAC:              meshMAC,
// 						Name:             s.fallbackName(apNameMap, meshMAC),
// 						SignalStrength:   signal,
// 						ProximityPercent: s.calculateProximityFromSignal(signal, clientSignal),
// 					})
// 				}
// 			}
// 		}
// 	} else {
// 		// Fallback: estimasi sinyal untuk AP mesh yang related
// 		for meshMAC := range meshRelatedAPs {
// 			// Cari AP info dengan case-insensitive comparison
// 			var apExists bool
// 			for apMAC := range aps {
// 				if strings.EqualFold(apMAC, meshMAC) {
// 					apExists = true
// 					break
// 				}
// 			}
			
// 			if apExists {
// 				estimatedSignal := s.estimateSignalStrength(clientSignal, meshMAC)
// 				proximity := s.calculateProximityFromSignal(estimatedSignal, clientSignal)

// 				// Hanya tampilkan jika proximity cukup tinggi
// 				if proximity > 15.0 {
// 					nearby = append(nearby, dtos.NearbyAP{
// 						MAC:              meshMAC,
// 						Name:             s.fallbackName(apNameMap, meshMAC),
// 						SignalStrength:   estimatedSignal,
// 						ProximityPercent: proximity,
// 					})
// 				}
// 			}
// 		}
// 	}

// 	// Urutkan berdasarkan proximity tertinggi
// 	sort.Slice(nearby, func(i, j int) bool {
// 		return nearby[i].ProximityPercent > nearby[j].ProximityPercent
// 	})

// 	// Batasi maksimal 2 AP nearby
// 	if len(nearby) > 2 {
// 		return nearby[:2]
// 	}
// 	return nearby
// }

