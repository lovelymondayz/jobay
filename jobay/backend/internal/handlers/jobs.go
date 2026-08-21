package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"jobay/internal/database"
	"jobay/internal/models"
	"jobay/internal/utils"

	"github.com/gin-gonic/gin"
)

func computeStats(db *models.Database, slug string) models.Stats {
	var s models.Stats
	for _, j := range db.Jobs {
		if slug != "" && j.UserSlug != slug {
			continue
		}
		switch j.Status {
		case models.StatusDiscovered:
			s.Discovered++
		case models.StatusQualified:
			s.Qualified++
		case models.StatusReview:
			s.Review++
		case models.StatusApplied:
			s.Applied++
		case models.StatusOutcomeInterview:
			s.OutcomeInterview++
		case models.StatusOutcomeRejected:
			s.OutcomeRejected++
		}
		s.Total++
	}
	return s
}

func GetStatus(c *gin.Context) {
	db := database.GetDB()
	slug := c.Query("user_slug")
	stats := computeStats(db, slug)
	utils.JSON(c, http.StatusOK, gin.H{
		"jobs":    db.Jobs,
		"actions": db.Actions,
		"agent":   db.Agent,
		"runs":    db.Runs,
		"stats":   stats,
	})
}

func ListJobs(c *gin.Context) {
	db := database.GetDB()
	status := c.Query("status")
	search := c.Query("search")
	userSlug := c.Query("user_slug")

	jobs := make([]models.Job, 0)
	for _, j := range db.Jobs {
		if status != "" && string(j.Status) != status {
			continue
		}
		if userSlug != "" && j.UserSlug != userSlug {
			continue
		}
		if search != "" {
			s := strings.ToLower(search)
			if !strings.Contains(strings.ToLower(j.Company), s) && !strings.Contains(strings.ToLower(j.Role), s) {
				continue
			}
		}
		jobs = append(jobs, j)
	}
	utils.JSON(c, http.StatusOK, gin.H{"jobs": jobs, "total": len(jobs)})
}

func GetJob(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid job ID")
		return
	}
	db := database.GetDB()
	for _, j := range db.Jobs {
		if j.ID == id {
			utils.JSON(c, http.StatusOK, j)
			return
		}
	}
	utils.Error(c, http.StatusNotFound, "Not found")
}

func CreateJob(c *gin.Context) {
	var req struct {
		Company  string `json:"company" binding:"required"`
		Role     string `json:"role" binding:"required"`
		URL      string `json:"url"`
		Status   string `json:"status"`
		Score    *int   `json:"score"`
		Notes    string `json:"notes"`
		UserSlug string `json:"user_slug"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	db := database.GetDB()
	job := models.Job{
		ID:        nextID(db.Jobs),
		Company:   req.Company,
		Role:      req.Role,
		URL:       req.URL,
		Status:    models.JobStatus(req.Status),
		Score:     req.Score,
		Notes:     req.Notes,
		UserSlug:  req.UserSlug,
		CreatedAt: database.Now(),
	}
	if job.Status == "" {
		job.Status = models.StatusDiscovered
	}
	db.Jobs = append(db.Jobs, job)
	database.Save()
	utils.JSON(c, http.StatusCreated, job)
}

func UpdateJob(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid job ID")
		return
	}
	var req struct {
		Company string `json:"company"`
		Role    string `json:"role"`
		URL     string `json:"url"`
		Status  string `json:"status"`
		Score   *int   `json:"score"`
		Notes   string `json:"notes"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	db := database.GetDB()
	for i := range db.Jobs {
		if db.Jobs[i].ID == id {
			if req.Company != "" {
				db.Jobs[i].Company = req.Company
			}
			if req.Role != "" {
				db.Jobs[i].Role = req.Role
			}
			if req.URL != "" {
				db.Jobs[i].URL = req.URL
			}
			if req.Status != "" {
				db.Jobs[i].Status = models.JobStatus(req.Status)
			}
			if req.Score != nil {
				db.Jobs[i].Score = req.Score
			}
			if req.Notes != "" {
				db.Jobs[i].Notes = req.Notes
			}
			database.Save()
			utils.JSON(c, http.StatusOK, db.Jobs[i])
			return
		}
	}
	utils.Error(c, http.StatusNotFound, "Not found")
}

func DeleteJob(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid job ID")
		return
	}
	db := database.GetDB()
	for i := range db.Jobs {
		if db.Jobs[i].ID == id {
			db.Jobs = append(db.Jobs[:i], db.Jobs[i+1:]...)
			database.Save()
			utils.JSON(c, http.StatusOK, gin.H{"ok": true})
			return
		}
	}
	utils.Error(c, http.StatusNotFound, "Not found")
}
