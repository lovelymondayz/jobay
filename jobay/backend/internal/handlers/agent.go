package handlers

import (
	"log"
	"net/http"
	"time"

	"jobay/internal/agent"
	"jobay/internal/database"
	"jobay/internal/models"
	"jobay/internal/utils"

	"github.com/gin-gonic/gin"
)

func ListActions(c *gin.Context) {
	db := database.GetDB()
	slug := c.Query("user_slug")
	actions := make([]models.Action, 0)
	for _, a := range db.Actions {
		if slug != "" && a.UserSlug != slug {
			continue
		}
		actions = append(actions, a)
	}
	// Sort by created_at desc
	for i := range actions {
		for j := i + 1; j < len(actions); j++ {
			if actions[i].CreatedAt < actions[j].CreatedAt {
				actions[i], actions[j] = actions[j], actions[i]
			}
		}
	}
	if len(actions) > 50 {
		actions = actions[:50]
	}
	utils.JSON(c, http.StatusOK, actions)
}

func CreateAction(c *gin.Context) {
	var req struct {
		Type     string `json:"type" binding:"required"`
		Message  string `json:"message" binding:"required"`
		JobID    *int   `json:"job_id"`
		UserSlug string `json:"user_slug"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	db := database.GetDB()
	action := models.Action{
		ID:        nextID(db.Actions),
		Type:      req.Type,
		Message:   req.Message,
		JobID:     req.JobID,
		UserSlug:  req.UserSlug,
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
	db.Actions = append(db.Actions, action)
	database.Save()
	utils.JSON(c, http.StatusCreated, action)
}

func GetAgent(c *gin.Context) {
	db := database.GetDB()
	utils.JSON(c, http.StatusOK, db.Agent)
}

func ToggleAgent(c *gin.Context) {
	db := database.GetDB()
	if db.Agent.IsRunning == 1 {
		db.Agent.IsRunning = 0
	} else {
		db.Agent.IsRunning = 1
	}
	db.Agent.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	database.Save()
	utils.JSON(c, http.StatusOK, db.Agent)
}

func SetAgentMode(c *gin.Context) {
	var req struct {
		Mode string `json:"mode"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	if req.Mode != "review-each" && req.Mode != "routine-auto" {
		utils.Error(c, http.StatusBadRequest, "Invalid mode")
		return
	}
	db := database.GetDB()
	db.Agent.Mode = req.Mode
	db.Agent.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	database.Save()
	utils.JSON(c, http.StatusOK, db.Agent)
}

func ListRuns(c *gin.Context) {
	db := database.GetDB()
	utils.JSON(c, http.StatusOK, db.Runs)
}

func RunAgent(c *gin.Context) {
	var req struct {
		UserSlug string `json:"user_slug"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	if req.UserSlug == "" {
		utils.Error(c, http.StatusBadRequest, "user_slug is required")
		return
	}

	go func() {
		runner := agent.NewRunner()
		if err := runner.Run(req.UserSlug); err != nil {
			log.Printf("[Agent] Manual run failed for %s: %v", req.UserSlug, err)
		}
	}()

	utils.JSON(c, http.StatusAccepted, gin.H{"status": "started", "user_slug": req.UserSlug})
}

func GetStats(c *gin.Context) {
	db := database.GetDB()
	slug := c.Query("user_slug")
	stats := computeStats(db, slug)
	utils.JSON(c, http.StatusOK, stats)
}
