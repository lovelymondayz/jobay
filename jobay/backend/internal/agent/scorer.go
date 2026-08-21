package agent

import (
	"encoding/json"
	"fmt"
)

// ScoredJob represents a job with an AI match score
type ScoredJob struct {
	JobListing
	Score       int      `json:"score"`
	MatchReason string   `json:"match_reason"`
	MatchSkills []string `json:"match_skills"`
	Gaps        []string `json:"gaps"`
}

var scoringSystemPrompt = `You are a precise job-matching AI. Given a user profile and a job listing, score how well the job matches the user's profile.

Return ONLY a JSON object with:
{
  "score": number 0-100 (0=terrible fit, 100=perfect fit),
  "match_reason": "string - one sentence explaining why this matches",
  "match_skills": ["skill1", "skill2"] - user's skills that match this job,
  "gaps": ["gap1", "gap2"] - requirements the user doesn't meet
}

Scoring criteria:
- Skills match: 40% weight (how many required skills does the user have?)
- Experience level: 30% weight (is the user's experience appropriate?)
- Location fit: 15% weight (does location match preference?)
- Role alignment: 15% weight (is this the kind of role the user wants?)

Important: Be honest. If a job is clearly not a fit, give it a low score. Do NOT inflate scores.`
var scoringUserPromptTemplate = `Score this job against this user profile:

USER PROFILE:
%s

JOB LISTING:
- Company: %s
- Role: %s
- Location: %s
- Description: %s

Return the scoring JSON:`

// ScoreJobs scores each job against the user profile
func ScoreJobs(jobs []JobListing, profile *UserProfile) []ScoredJob {
	if len(jobs) == 0 {
		return []ScoredJob{}
	}

	ai := NewAIClient()
	scored := make([]ScoredJob, 0, len(jobs))

	profileJSON, err := json.MarshalIndent(profile, "", "  ")
	if err != nil {
		return scored
	}

	for _, job := range jobs {
		desc := job.Description
		if desc == "" {
			desc = job.Role + " at " + job.Company
		}

		prompt := fmt.Sprintf(scoringUserPromptTemplate, string(profileJSON), job.Company, job.Role, job.Location, desc)

		var result ScoredJob
		if err := ai.ChatJSON(scoringSystemPrompt, prompt, &result); err != nil {
			continue
		}

		result.JobListing = job
		scored = append(scored, result)
	}

	return scored
}
