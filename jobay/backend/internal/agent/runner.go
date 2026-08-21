package agent

import (
	"fmt"
	"log"
	"time"

	"jobay/internal/database"
	"jobay/internal/models"
)

// Runner orchestrates the full agent pipeline
type Runner struct {
	parser   *Parser
	searcher *Searcher
	applier  *Applier
}

func NewRunner() *Runner {
	return &Runner{
		parser:   NewParser(),
		searcher: NewSearcher(),
		applier:  NewApplier(),
	}
}

// Run executes the full pipeline for a user
func (r *Runner) Run(userSlug string) error {
	db := database.GetDB()

	// Get user
	var user *models.User
	for i := range db.Users {
		if db.Users[i].Slug == userSlug {
			user = &db.Users[i]
			break
		}
	}
	if user == nil {
		return fmt.Errorf("user not found: %s", userSlug)
	}

	// Get CV path
	cvPath := "/app/data/cvs" + user.CVPath[len("/uploads"):]
	if user.CVPath == "" {
		return fmt.Errorf("user has no CV: %s", userSlug)
	}

	// 1. Parse CV
	log.Printf("[Agent] Parsing CV for %s...", userSlug)
	cvText, err := r.parser.ExtractText(cvPath)
	if err != nil {
		return fmt.Errorf("CV parse failed: %w", err)
	}
	if cvText == "" {
		return fmt.Errorf("empty CV text")
	}

	// 2. Extract profile
	log.Printf("[Agent] Extracting profile for %s...", userSlug)
	profile, err := ExtractProfile(cvText)
	if err != nil {
		return fmt.Errorf("profile extraction failed: %w", err)
	}

	// 3. Search jobs
	log.Printf("[Agent] Searching jobs for %s...", userSlug)
	maxResults := configMaxResults()
	searchedJobs := r.searcher.Search(profile, maxResults)

	if len(searchedJobs) == 0 {
		log.Printf("[Agent] No jobs found for %s", userSlug)
		r.addAction(userSlug, "discover", "No jobs found for your profile", nil)
		return nil
	}

	// 4. Score jobs
	log.Printf("[Agent] Scoring %d jobs for %s...", len(searchedJobs), userSlug)
	scoredJobs := ScoreJobs(searchedJobs, profile)

	// 5. Apply threshold
	threshold := configThreshold()
	qualified := []ScoredJob{}
	for _, sj := range scoredJobs {
		if sj.Score >= threshold {
			qualified = append(qualified, sj)
		}
	}

	if len(qualified) == 0 {
		log.Printf("[Agent] No jobs above threshold for %s", userSlug)
		r.addAction(userSlug, "discover", fmt.Sprintf("Found %d jobs but none above score threshold", len(scoredJobs)), nil)
		return nil
	}

	// 6. Auto-apply
	totalApplied := 0
	for _, sj := range qualified {
		ok, msg := r.applier.ApplyAttempt(sj.JobListing, profile)
		status := models.StatusDiscovered
		if ok {
			status = models.StatusApplied
			totalApplied++
		}

		// Save job to DB
		job := models.Job{
			ID:           nextJobID(),
			Company:      sj.Company,
			Role:         sj.Role,
			URL:          sj.URL,
			Status:       status,
			Score:        &sj.Score,
			Notes:        sj.MatchReason,
			UserSlug:     userSlug,
			Location:     sj.Location,
			Source:       sj.Source,
			ApplyStatus:  map[bool]string{true: "applied", false: "failed"}[ok],
			MatchReasons: formatMatchReasons(sj),
			CreatedAt:    time.Now().UTC().Format(time.RFC3339),
		}

		db.Jobs = append(db.Jobs, job)

		// Add action
		actionType := "discover"
		if ok {
			actionType = "apply"
		}
		r.addAction(userSlug, actionType,
			fmt.Sprintf("%s at %s (score: %d) - %s", sj.Role, sj.Company, sj.Score, msg),
			&job.ID,
		)
	}

	// Save DB
	database.Save()

	log.Printf("[Agent] Completed for %s: found %d, qualified %d, applied %d",
		userSlug, len(scoredJobs), len(qualified), totalApplied)

	return nil
}

func (r *Runner) addAction(userSlug, actionType string, message string, jobID *int) {
	db := database.GetDB()
	action := models.Action{
		ID:        nextActionID(),
		Type:      actionType,
		Message:   message,
		JobID:     jobID,
		UserSlug:  userSlug,
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
	db.Actions = append(db.Actions, action)
}

func nextJobID() int {
	db := database.GetDB()
	maxID := 0
	for _, j := range db.Jobs {
		if j.ID > maxID {
			maxID = j.ID
		}
	}
	return maxID + 1
}

func nextActionID() int {
	db := database.GetDB()
	maxID := 0
	for _, a := range db.Actions {
		if a.ID > maxID {
			maxID = a.ID
		}
	}
	return maxID + 1
}

func configThreshold() int {
	// Could be read from env/config
	return 70
}

func configMaxResults() int {
	return 50
}

func formatMatchReasons(sj ScoredJob) string {
	return fmt.Sprintf("Match skills: %v", sj.MatchSkills)
}
