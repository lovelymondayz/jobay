package main

import (
	"log"
	"net/http"
	"path/filepath"

	"jobay/internal/config"
	"jobay/internal/database"
	"jobay/internal/handlers"
	"jobay/internal/middleware"
	ws "jobay/internal/ws"

	"github.com/gin-gonic/gin"
)

func main() {
	dbPath := config.Get("DB_PATH", "/app/data/jobay.json")
	database.InitDB(dbPath)

	hub := ws.NewHub()

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(middleware.CORSMiddleware())

	uploadsDir := config.Get("UPLOADS_DIR", "/app/data/cvs")

	// Static files — serve assets and uploads under specific paths
	r.Static("/assets", filepath.Join(".", "frontend", "dist", "assets"))
	r.Static("/uploads", uploadsDir)

	// Health check
	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// ═══════════════════════════════════════════
	// PUBLIC API (no auth)
	// ═══════════════════════════════════════════
	public := r.Group("/api")
	{
		// Upload CV + create user
		public.POST("/upload", handlers.UploadCV)

		// User-scoped routes
		public.GET("/users/:slug", handlers.GetUser)
		public.GET("/users/:slug/jobs", handlers.ListUserJobs)
		public.GET("/users/:slug/actions", handlers.ListUserActions)
		public.GET("/users/:slug/stats", handlers.GetUserStats)

		// Job CRUD
		public.GET("/jobs", handlers.ListJobs)
		public.GET("/jobs/:id", handlers.GetJob)
		public.POST("/jobs", handlers.CreateJob)
		public.PATCH("/jobs/:id", handlers.UpdateJob)
		public.DELETE("/jobs/:id", handlers.DeleteJob)

		// Actions
		public.GET("/actions", handlers.ListActions)
		public.POST("/actions", handlers.CreateAction)

		// Agent
		public.GET("/agent", handlers.GetAgent)
		public.POST("/agent/toggle", handlers.ToggleAgent)
		public.POST("/agent/mode", handlers.SetAgentMode)
		public.POST("/agent/run", handlers.RunAgent)

		// Runs
		public.GET("/runs", handlers.ListRuns)

		// Stats
		public.GET("/stats", handlers.GetStats)

		// Status (full snapshot)
		public.GET("/status", handlers.GetStatus)
	}

	// WebSocket
	r.GET("/ws", gin.WrapH(http.HandlerFunc(hub.HandleWS)))
	r.GET("/ws/:slug", func(c *gin.Context) {
		hub.HandleWS(c.Writer, c.Request)
	})

	// SPA fallback
	r.NoRoute(func(c *gin.Context) {
		c.File(filepath.Join(".", "frontend", "dist", "index.html"))
	})

	port := config.Get("PORT", "8080")
	log.Printf("Jobay API starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start:", err)
	}
}
