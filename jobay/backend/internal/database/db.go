package database

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"jobay/internal/models"
)

var (
	db     *models.Database
	mu     sync.RWMutex
	dbPath string
)

func InitDB(path string) {
	dbPath = path
	db = &models.Database{
		Jobs:    []models.Job{},
		Users:   []models.User{},
		Actions: []models.Action{},
		Agent: models.AgentStatus{
			ID:         1,
			Mode:       "review-each",
			AIProvider: "9router",
			IsRunning:  0,
			UpdatedAt:  Now(),
		},
		Runs: []models.Run{},
	}
	load()
}

func GetDB() *models.Database {
	return db
}

func load() {
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Printf("DB dir error: %v", err)
		return
	}
	data, err := os.ReadFile(dbPath)
	if err != nil {
		log.Printf("DB read error (first run?): %v", err)
		seed()
		return
	}
	if err := json.Unmarshal(data, db); err != nil {
		log.Printf("DB parse error: %v", err)
		seed()
	}
}

func Save() {
	mu.Lock()
	defer mu.Unlock()
	save()
}

func save() {
	data, err := json.MarshalIndent(db, "", "  ")
	if err != nil {
		log.Printf("DB marshal error: %v", err)
		return
	}
	if err := os.WriteFile(dbPath, data, 0644); err != nil {
		log.Printf("DB write error: %v", err)
	}
}

func seed() {
	log.Println("Seeding dummy data")
	db.Jobs = []models.Job{
		{ID: 1, Company: "TechCorp", Role: "Senior Backend Engineer", URL: "https://techcorp.com/jobs/123", Status: models.StatusApplied, Score: intPtr(88), UserSlug: "admin", CreatedAt: Now()},
		{ID: 2, Company: "StartupXYZ", Role: "Full Stack Developer", URL: "https://startupxyz.com/jobs/456", Status: models.StatusReview, Score: intPtr(75), UserSlug: "admin", CreatedAt: Now()},
		{ID: 3, Company: "BigTech Inc", Role: "Staff Engineer", URL: "https://bigtech.com/jobs/789", Status: models.StatusDiscovered, Score: intPtr(62), UserSlug: "admin", CreatedAt: Now()},
		{ID: 4, Company: "RemoteCo", Role: "React Developer", URL: "https://remote.co/jobs/101", Status: models.StatusQualified, Score: intPtr(82), UserSlug: "admin", CreatedAt: Now()},
		{ID: 5, Company: "AI Labs", Role: "ML Engineer", URL: "https://ailabs.com/jobs/202", Status: models.StatusOutcomeRejected, Score: intPtr(70), UserSlug: "admin", CreatedAt: Now()},
		{ID: 6, Company: "DataDriven", Role: "Data Engineer", URL: "https://datadriven.com/jobs/303", Status: models.StatusOutcomeInterview, Score: intPtr(85), UserSlug: "admin", CreatedAt: Now()},
		{ID: 7, Company: "CloudFirst", Role: "DevOps Engineer", URL: "https://cloudfirst.com/jobs/404", Status: models.StatusDiscovered, Score: intPtr(78), UserSlug: "admin", CreatedAt: Now()},
		{ID: 8, Company: "FinTech Pro", Role: "Backend Engineer", URL: "https://fintechpro.com/jobs/505", Status: models.StatusApplied, Score: intPtr(91), UserSlug: "admin", CreatedAt: Now()},
	}
	db.Actions = []models.Action{
		{ID: 1, Type: "discover", Message: "Found 3 new roles from career pages", UserSlug: "admin", CreatedAt: Now()},
		{ID: 2, Type: "score", Message: "TechCorp Senior Backend Engineer scored 88/100", JobID: intPtr(1), UserSlug: "admin", CreatedAt: Now()},
		{ID: 3, Type: "apply", Message: "Application submitted to TechCorp", JobID: intPtr(1), UserSlug: "admin", CreatedAt: Now()},
		{ID: 4, Type: "score", Message: "StartupXYZ Full Stack Developer scored 75/100", JobID: intPtr(2), UserSlug: "admin", CreatedAt: Now()},
		{ID: 5, Type: "review", Message: "Flagged for manual review: StartupXYZ", JobID: intPtr(2), UserSlug: "admin", CreatedAt: Now()},
		{ID: 6, Type: "score", Message: "RemoteCo React Developer scored 82/100", JobID: intPtr(4), UserSlug: "admin", CreatedAt: Now()},
		{ID: 7, Type: "outcome", Message: "Interview scheduled with DataDriven", JobID: intPtr(6), UserSlug: "admin", CreatedAt: Now()},
	}
	db.Runs = []models.Run{
		{ID: "demo-run-1", StartedAt: Now(), EndedAt: Now(), TotalFound: 8, TotalQualified: 5, TotalApplied: 2, TotalSkipped: 1},
	}
	save()
}

func intPtr(i int) *int {
	return &i
}

func Now() string {
	return time.Now().UTC().Format(time.RFC3339)
}

func nextID(items []models.Job) int {
	max := 0
	for _, item := range items {
		if item.ID > max {
			max = item.ID
		}
	}
	return max + 1
}

func nextIDActions(items []models.Action) int {
	max := 0
	for _, item := range items {
		if item.ID > max {
			max = item.ID
		}
	}
	return max + 1
}

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
