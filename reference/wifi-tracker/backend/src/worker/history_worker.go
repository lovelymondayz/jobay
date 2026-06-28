package worker

import (
	"context"
	"log"
	"time"
	"wifi-tracker-be/src/dtos"
	"wifi-tracker-be/src/service"

	"github.com/google/uuid"
)

type ConnectionWorker struct {
	unifiService service.UnifiIntegrationServiceWoUnifiPollerInterface
	lastSeen     map[string]connectionState
}

type connectionState struct {
	LastIdAP   uuid.UUID
	LastLogged time.Time
}

func NewConnectionWorker(svc service.UnifiIntegrationServiceWoUnifiPollerInterface) *ConnectionWorker {
	return &ConnectionWorker{
		unifiService: svc,
		lastSeen:     make(map[string]connectionState),
	}
}

func (w *ConnectionWorker) Start(ctx context.Context) {
	go func() {
		// Tunggu sampai ke menit penuh
		now := time.Now()
		next := now.Truncate(time.Minute).Add(time.Minute) // ke menit berikutnya
		time.Sleep(time.Until(next))

		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				w.run()
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (w *ConnectionWorker) run() {
	log.Println("[WORKER][CONNECTION] Running")
	clients, err := w.unifiService.GetClientsWithNearbyAPs()
	if err != nil {
		// log error
		return
	}

	if len(clients) == 0 {
		log.Println("[WORKER][CONNECTION] No clients found")
	}

	now := time.Now()

	log.Print(w.lastSeen)

	for _, client := range clients {
		if client.APID == uuid.Nil || client.DeviceID == uuid.Nil {
			continue
		}
		mac := client.MACAddressClient
		currAP := client.APID

		state, exists := w.lastSeen[mac]

		if !exists || state.LastIdAP != currAP || !w.IsSameDay(now, state.LastLogged, time.Local) {
			// Device baru atau pindah AP
			if state.LastIdAP == uuid.Nil {
				w.logConnection(client, false, client.DeviceID,  currAP)
			} else {
				w.logConnection(client, false, client.DeviceID, currAP)
				w.lastSeen[mac] = connectionState{
					LastIdAP:   currAP,
					LastLogged: now,
				}
			}
		} else {
			slotTime := w.truncateTo30Min(now)
			if slotTime.After(w.truncateTo30Min(state.LastLogged)) {
				w.logConnection(client, true, client.DeviceID,  currAP)
				w.lastSeen[mac] = connectionState{
					LastIdAP:   currAP,
					LastLogged: slotTime,
				}
			}
		}
	}
}

func (w *ConnectionWorker) logConnection(client dtos.ActiveClient, isThirtyMinutes bool, deviceID, toAP uuid.UUID) {
	// Jika DeviceID kosong, abaikan
	if client.DeviceID == uuid.Nil {
		log.Printf("[WORKER][CONNECTION] empty DeviceID: %s\n", client.MACAddressClient)
		return
	}

	var fromAP uuid.UUID

	isError, datas, err := w.unifiService.CheckIfConnectionLoggedToday(client.DeviceID, isThirtyMinutes, time.Now())
	if err != nil {
		log.Printf("[WORKER][CONNECTION] Cannot check history: %v\n", err)
		return
	}

	if isError {
		log.Printf("[WORKER][CONNECTION][ERROR]: %s -> %s\n", client.MACAddressClient, toAP)
		return
	}

	if isThirtyMinutes {
		err = w.unifiService.SaveConnectionHistory(client, deviceID, toAP, toAP)
		if err != nil {
			log.Printf("[WORKER][CONNECTION] Cannot save history: %v\n", err)
		}

		log.Printf("[WORKER][CONNECTION] Logged: %s -> %s\n", client.MACAddressClient, toAP)
		return
	}

	if datas.ToAps == toAP {
		w.lastSeen[datas.ToMacAddressDevice] = connectionState{
			LastIdAP:   datas.ToAps,
			LastLogged: time.Now(),
		}
		log.Printf("[WORKER][CONNECTION] Already logged: %s -> %s\n", client.MACAddressClient, toAP)
		return
	}

	if datas.ToAps != uuid.Nil{
		fromAP = datas.ToAps
	}else {
		fromAP = toAP
	}

	err = w.unifiService.SaveConnectionHistory(client, deviceID, fromAP, toAP)
	if err != nil {
		log.Printf("[WORKER][CONNECTION] Cannot save history: %v\n", err)
	}
}

func (w *ConnectionWorker) truncateTo30Min(t time.Time) time.Time {
	minute := t.Minute()
	roundedMin := 0
	if minute >= 30 {
		roundedMin = 30
	}
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), roundedMin, 0, 0, t.Location())
}

func (w *ConnectionWorker) IsSameDay(t1, t2 time.Time, loc *time.Location) bool {
	y1, m1, d1 := t1.In(loc).Date()
	y2, m2, d2 := t2.In(loc).Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}