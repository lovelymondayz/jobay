package routes

import (
	"wifi-tracker-be/src/config"
	"wifi-tracker-be/src/handlers"
	"wifi-tracker-be/src/middleware"
	"wifi-tracker-be/src/repository"
	"wifi-tracker-be/src/service"
	"wifi-tracker-be/src/ws"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func SetupRoutes(app *fiber.App) {
	
	// Dependencies
	userRepo := repository.NewUserRepository()
	apRepo := repository.NewAccessPointRepository()
	deviceRepo := repository.NewDeviceRepository()
	historyRepo := repository.NewConnectionHistoriesRepository()
	siteRepo := repository.NewSiteRepository()
	roleRepo := repository.NewRoleRepository()
	clientRepo := repository.NewClientRepository()

	//services
	authService := service.NewAuthService(config.DB, userRepo)
	unifiService := service.NewUnifiIntegrationService(config.UnifiClient)
	unifiService2 := service.NewUnifiIntegrationServiceWoUnifiPoller(config.DB,deviceRepo,apRepo,userRepo,historyRepo,config.UnifiClient2)
	siteService := service.NewSiteService(siteRepo, config.DB)
	roleService := service.NewRoleService(roleRepo, config.DB)
	deviceService := service.NewDevicesService(deviceRepo, config.DB)
	clientDevice := service.NewClientService(clientRepo, config.DB)
	userService := service.NewUserService(userRepo, config.DB)
	apService := service.NewAccessPointService(apRepo, config.DB)

	//Handlers
	unifiHandler := handlers.NewUnifiIntegrationHandler(unifiService, unifiService2)
	authHandler := handlers.NewAuthHandler(authService)
	healthHandler := handlers.NewHealthHandler()
	siteHandler := handlers.NewSitesHandler(siteService)
	roleHandler := handlers.NewRoleHandler(roleService)
	deviceHandler := handlers.NewDevicesHandler(deviceService)
	clientHandler := handlers.NewClientsHandler(clientDevice)
	userHandler := handlers.NewUserHandler(userService)
	apHandler := handlers.NewAccessPointHandler(apService)

	//workerIntialize
	// connWorker := worker.NewConnectionWorker(unifiService2)
	
	// Routes REST API
	api := app.Group("/api/v1")

	// Health check endpoint
	api.Get("/healthCheck", middleware.JWTProtected(),healthHandler.HealthCheck)

	// Authentication endpoints
	// api.Post("/register", authHandler.Register)
	api.Post("/auth/login", authHandler.Login)

	// User profile endpoint (protected)
	// api.Get("/profile", middleware.JWTProtected(), userHandler.Profile)

	// api.Get("/clients/active", unifiHandler.GetClients)
	// uniFi Route
	api.Get("/devices", unifiHandler.GetDevices2)
	api.Get("/clients/active", unifiHandler.GetClients2)
	api.Get("/devices/active", unifiHandler.GetDevices)
	api.Get("/clients/active/:macAP", unifiHandler.GetClientByMACAP)
	api.Get("/histories", unifiHandler.FindAllHistoryConnection)

	// site Route
	api.Get("/sites", siteHandler.FindAllSites)
	api.Get("/sites/:id", siteHandler.FindSiteById)
	api.Post("/sites", siteHandler.CreateSite)
	api.Put("/sites/:id", siteHandler.UpdateSite)
	api.Delete("/sites/:id", siteHandler.DeleteSite)

	// role Route
	api.Get("/roles", roleHandler.FindAll)
	api.Get("/roles/:id", roleHandler.FindByID)
	api.Post("/roles", roleHandler.Create)
	api.Put("/roles/:id", roleHandler.Update)
	api.Delete("/roles/:id", roleHandler.Delete)

	// device Route
	api.Get("/user-devices", deviceHandler.GetAll)
	api.Get("/user-devices/:id", deviceHandler.GetByID)
	api.Post("/user-devices", deviceHandler.Create)
	api.Put("/user-devices/:id", deviceHandler.Update)
	api.Delete("/user-devices/:id", deviceHandler.Delete)

	// client Route
	api.Get("/clients", clientHandler.GetAll)
	api.Get("/clients/:id", clientHandler.GetByID)
	api.Post("/clients", clientHandler.Create)
	api.Put("/clients/:id", clientHandler.Update)
	api.Delete("/clients/:id", clientHandler.Delete)

	// user Route
	api.Get("/users", userHandler.GetAll)
	api.Get("/users/:id", userHandler.GetByID)
	api.Post("/users", userHandler.Create)
	api.Put("/users/:id", userHandler.Update)
	api.Delete("/users/:id", userHandler.Delete)

	// ap Route
	// api.Get("/aps", apHandler.GetAll)
	api.Get("/aps/:id", apHandler.GetByID)
	api.Post("/aps", apHandler.Create)
	api.Put("/aps/:id", apHandler.Update)
	api.Delete("/aps/:id", apHandler.Delete)

	// WS
	deviceWsService := ws.NewUniFiDeviceWS(unifiService2)

	// Route WebSocket
	wsRoute := app.Group("/ws/v1")
	
	// WS endpoint
	wsRoute.Get("/devices", websocket.New(deviceWsService.DeviceActiveHandler))
	
	// worker
	go deviceWsService.StartDeviceWatcher()
	go deviceWsService.HandleDeviceBroadcast()

	// Jalankan worker dengan context
	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()

	// connWorker.Start(ctx)
}