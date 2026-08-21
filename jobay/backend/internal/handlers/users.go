package handlers

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"jobay/internal/agent"
	"jobay/internal/database"
	"jobay/internal/models"
	"jobay/internal/utils"

	"github.com/gin-gonic/gin"
)

func UploadCV(c *gin.Context) {
	name := c.PostForm("name")
	if name == "" {
		utils.Error(c, http.StatusBadRequest, "Name is required")
		return
	}

	email := c.PostForm("email")
	phone := c.PostForm("phone")

	file, err := c.FormFile("cv")
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "CV file is required")
		return
	}

	ext := filepath.Ext(file.Filename)
	if ext != ".pdf" && ext != ".docx" && ext != ".doc" {
		utils.Error(c, http.StatusBadRequest, "Only PDF, DOCX, DOC allowed")
		return
	}

	slug := utils.Slugify(name)
	if slug == "" {
		slug = "user-" + utils.RandID()[:8]
	}

	db := database.GetDB()
	for _, u := range db.Users {
		if u.Slug == slug {
			utils.Error(c, http.StatusBadRequest, "User already exists")
			return
		}
	}

	uploadsDir := "/app/data/cvs"
	if err := os.MkdirAll(uploadsDir, 0755); err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to create uploads dir")
		return
	}

	filename := slug + ext
	filepath := filepath.Join(uploadsDir, filename)
	if err := c.SaveUploadedFile(file, filepath); err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to save file")
		return
	}

	user := models.User{
		Slug:      slug,
		Name:      name,
		Email:     email,
		Phone:     phone,
		CVPath:    "/uploads/" + filename,
		CreatedAt: database.Now(),
	}
	db.Users = append(db.Users, user)
	database.Save()

	// Trigger agent pipeline in background
	go func() {
		runner := agent.NewRunner()
		if err := runner.Run(slug); err != nil {
			log.Printf("[Agent] Run failed for %s: %v", slug, err)
		}
	}()

	utils.JSON(c, http.StatusCreated, gin.H{
		"slug": slug,
		"name": name,
		"url":  "/" + slug,
	})
}

func GetUser(c *gin.Context) {
	slug := c.Param("slug")
	db := database.GetDB()
	for _, u := range db.Users {
		if u.Slug == slug {
			utils.JSON(c, http.StatusOK, u)
			return
		}
	}
	utils.Error(c, http.StatusNotFound, "User not found")
}

func ListUserJobs(c *gin.Context) {
	slug := c.Param("slug")
	db := database.GetDB()
	jobs := make([]models.Job, 0)
	for _, j := range db.Jobs {
		if j.UserSlug == slug {
			jobs = append(jobs, j)
		}
	}
	utils.JSON(c, http.StatusOK, gin.H{"jobs": jobs, "total": len(jobs)})
}

func ListUserActions(c *gin.Context) {
	slug := c.Param("slug")
	db := database.GetDB()
	actions := make([]models.Action, 0)
	for _, a := range db.Actions {
		if a.UserSlug == slug {
			actions = append(actions, a)
		}
	}
	// Sort desc
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

func GetUserStats(c *gin.Context) {
	slug := c.Param("slug")
	db := database.GetDB()
	stats := computeStats(db, slug)
	utils.JSON(c, http.StatusOK, stats)
}
