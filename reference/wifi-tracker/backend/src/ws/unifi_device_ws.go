package ws

import (
	"log"
	"time"

	"wifi-tracker-be/src/dtos"
	"wifi-tracker-be/src/service"

	"github.com/gofiber/websocket/v2"
)

type UniFiDeviceWSInterface interface {
	DeviceActiveHandler(c *websocket.Conn)
	StartDeviceWatcher()
	HandleDeviceBroadcast()
}

type UniFiDeviceWS struct {
	unifiService service.UnifiIntegrationServiceWoUnifiPollerInterface
}

func NewUniFiDeviceWS(unifiService service.UnifiIntegrationServiceWoUnifiPollerInterface) UniFiDeviceWSInterface {
	return &UniFiDeviceWS{
		unifiService: unifiService,
	}
}

var (
	deviceClients     = make(map[*websocket.Conn]bool)
	deviceBroadcast   = make(chan []dtos.APInfoResponse)
	deviceStatusCache = make(map[string]dtos.APInfoResponse)
)

func (s UniFiDeviceWS) DeviceActiveHandler(c *websocket.Conn) {
	deviceClients[c] = true
	log.Printf("[WS][DEVICE] 🔌 Client terhubung (%s)\n", c.RemoteAddr())

	initial := make([]dtos.APInfoResponse, 0, len(deviceStatusCache))
	for _, ap := range deviceStatusCache {
		initial = append(initial, ap)
	}
	if len(initial) > 0 {
		if err := c.WriteJSON(initial); err != nil {
			log.Printf("[WS][DEVICE] Gagal kirim snapshot awal ke client (%s): %v\n", c.RemoteAddr(), err)
		} else {
			log.Printf("[WS][DEVICE] Kirim snapshot awal ke client (%s): %d data\n", c.RemoteAddr(), len(initial))
		}
	} else {
		log.Printf("[WS][DEVICE] Tidak ada snapshot awal untuk client (%s)\n", c.RemoteAddr())
	}

	defer func() {
		delete(deviceClients, c)
		log.Printf("[WS][DEVICE] Client terputus (%s)\n", c.RemoteAddr())
		c.Close()
	}()

	for {
		if _, _, err := c.ReadMessage(); err != nil {
			break
		}
	}
}

func (s UniFiDeviceWS) StartDeviceWatcher() {
	for {
		apData, err := s.unifiService.GetDevices()
		if err != nil {
			log.Printf("[WS][WATCHER] Gagal ambil data AP: %v\n", err)
			time.Sleep(5 * time.Second)
			continue
		}

		currentMap := make(map[string]dtos.APInfoResponse)
		changedDevices := make([]dtos.APInfoResponse, 0)

		// Loop hasil baru dan update cache
		for _, ap := range apData {
			currentMap[ap.Mac] = ap

			lastStatus, exists := deviceStatusCache[ap.Mac]
			if !exists || lastStatus.State != ap.State {
				changedDevices = append(changedDevices, ap)
			}
		}

		// Deteksi AP yang hilang dari polling
		for mac, oldAP := range deviceStatusCache {
			if _, found := currentMap[mac]; !found {
				log.Printf("[WS][WATCHER] AP %s hilang dari polling, tandai tidak aktif\n", mac)

				oldAP.State = 0 // Nonaktif
				currentMap[mac] = oldAP
				changedDevices = append(changedDevices, oldAP)
			}
		}

		// Perbarui cache
		deviceStatusCache = currentMap
		// if len(changedDevices) > 0 {
		// 	log.Printf("[WS][WATCHER] %d AP berubah status (termasuk hilang)\n", len(changedDevices))
		// 	deviceBroadcast <- changedDevices
		// } 

		if len(changedDevices) > 0 {
			// Ambil seluruh cache untuk dikirim (bukan hanya yang berubah)
			allDevices := make([]dtos.APInfoResponse, 0, len(deviceStatusCache))
			for _, ap := range deviceStatusCache {
				allDevices = append(allDevices, ap)
			}

			log.Printf("[WS][WATCHER] Ada %d AP berubah, kirim total %d AP ke client\n", len(changedDevices), len(allDevices))
			deviceBroadcast <- allDevices
		} else {
			log.Println("[WS][WATCHER] Tidak ada perubahan")
		}

		time.Sleep(5 * time.Second)
	}
}

func (s UniFiDeviceWS) HandleDeviceBroadcast() {
	for {
		msg := <-deviceBroadcast
		log.Printf("[WS][DEVICE] Broadcast %d AP ke %d client\n", len(msg), len(deviceClients))

		for conn := range deviceClients {
			err := conn.WriteJSON(msg)
			if err != nil {
				log.Printf("[WS][DEVICE] Gagal kirim ke client (%s): %v\n", conn.RemoteAddr(), err)
				conn.Close()
				delete(deviceClients, conn)
			} else {
				log.Printf("[WS][DEVICE] Kirim ke client (%s)\n", conn.RemoteAddr())
			}
		}
	}
}
