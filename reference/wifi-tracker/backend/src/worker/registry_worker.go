package worker

import (
	"context"
	"wifi-tracker-be/src/config"
	"wifi-tracker-be/src/repository"
	"wifi-tracker-be/src/service"
)

// StartAllWorkers menjalankan semua background worker
func StartAllWorkers(ctx context.Context) {
	// Init repository & service
	userRepo := repository.NewUserRepository()
	apRepo := repository.NewAccessPointRepository()
	deviceRepo := repository.NewDeviceRepository()
	historyRepo := repository.NewConnectionHistoriesRepository()

	unifiService := service.NewUnifiIntegrationServiceWoUnifiPoller(
		config.DB,
		deviceRepo,
		apRepo,
		userRepo,
		historyRepo,
		config.UnifiClient2, // ← Tambah ini juga
	)

	connWorker := NewConnectionWorker(unifiService)
	connWorker.Start(ctx)
}
