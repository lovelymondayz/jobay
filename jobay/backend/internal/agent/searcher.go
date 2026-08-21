package agent

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// JobListing represents a job found from any source
type JobListing struct {
	Company    string `json:"company"`
	Role       string `json:"role"`
	URL        string `json:"url"`
	Location   string `json:"location"`
	Source     string `json:"source"`
	Salary     string `json:"salary"`
	JobType    string `json:"job_type"`
	PostedDate string `json:"posted_date"`
	Description string `json:"description,omitempty"`
}

// Searcher handles job searching across multiple sources
type Searcher struct {
	client *http.Client
}

func NewSearcher() *Searcher {
	return &Searcher{
		client: &http.Client{Timeout: 20 * time.Second},
	}
}

// Search searches all sources for matching jobs
func (s *Searcher) Search(profile *UserProfile, maxResults int) []JobListing {
	allJobs := []JobListing{}

	// Generate search queries
	queries := s.generateQueries(profile)

	// Search each source
	for _, query := range queries[:min(len(queries), 5)] {
		jobs := s.searchRemoteOK(query, profile.PreferredLocations)
		allJobs = append(allJobs, jobs...)

		jobs = s.searchGoogleJobs(query, profile.PreferredLocations)
		allJobs = append(allJobs, jobs...)

		jobs = s.searchJobStreet(query, profile.PreferredLocations)
		allJobs = append(allJobs, jobs...)

		jobs = s.searchLinkedIn(query, profile.PreferredLocations)
		allJobs = append(allJobs, jobs...)
	}

	// Deduplicate by URL
	seen := make(map[string]bool)
	unique := []JobListing{}
	for _, job := range allJobs {
		if job.URL == "" || seen[job.URL] {
			continue
		}
		seen[job.URL] = true
		unique = append(unique, job)
		if len(unique) >= maxResults {
			break
		}
	}

	return unique
}

func (s *Searcher) generateQueries(profile *UserProfile) []string {
	queries := []string{}

	// Use preferred roles if available, otherwise current role
	roles := profile.PreferredRoles
	if len(roles) == 0 && profile.CurrentRole != "" {
		roles = []string{profile.CurrentRole}
	}

	locations := profile.PreferredLocations
	if len(locations) == 0 && profile.Location != "" {
		locations = []string{profile.Location}
	}

	// Generate role x location combos
	for _, role := range roles[:min(len(roles), 3)] {
		for _, loc := range locations[:min(len(locations), 3)] {
			queries = append(queries, fmt.Sprintf("%s %s", role, loc))
		}
		// Also add remote
		if profile.RemotePreference {
			queries = append(queries, fmt.Sprintf("%s remote", role))
		}
	}

	// Add skill-based queries
	if len(profile.Skills) > 0 {
		queries = append(queries, fmt.Sprintf("%s jobs", strings.Join(profile.Skills[:min(len(profile.Skills), 3)], " ")))
	}

	// Ensure at least some queries
	if len(queries) == 0 {
		queries = []string{"software engineer jakarta", "developer indonesia", "remote developer"}
	}

	return queries
}

// RemoteOK has a public JSON API
func (s *Searcher) searchRemoteOK(query string, locations []string) []JobListing {
	jobs := []JobListing{}
	apiURL := fmt.Sprintf("https://remoteok.com/api?search=%s", url.QueryEscape(query))

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return jobs
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; Jobay/1.0)")

	resp, err := s.client.Do(req)
	if err != nil {
		return jobs
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return jobs
	}

	// RemoteOK returns array with first element as metadata
	var results []map[string]interface{}
	if err := json.Unmarshal(body, &results); err != nil {
		return jobs
	}

	for _, item := range results {
		job := RemoteOKToJobListing(item)
		if job != nil && s.matchesLocation(job.Location, locations) {
			jobs = append(jobs, *job)
		}
	}

	return jobs
}

func RemoteOKToJobListing(item map[string]interface{}) *JobListing {
	company, _ := item["company"].(string)
	role, _ := item["position"].(string)
	jobURL, _ := item["url"].(string)
	location, _ := item["location"].(string)
	if location == "" {
		location = "Remote"
	}
	source, _ := item["source"].(string)
	posted, _ := item["date"].(string)

	return &JobListing{
		Company:    company,
		Role:       role,
		URL:        jobURL,
		Location:   location,
		Source:     source,
		PostedDate: posted,
	}
}

// Google Jobs search via web search
func (s *Searcher) searchGoogleJobs(query string, locations []string) []JobListing {
	jobs := []JobListing{}

	searchURL := fmt.Sprintf("https://www.google.com/search?q=%s+jobs&ibp=htl;jobs", url.QueryEscape(query))
	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		return jobs
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	resp, err := s.client.Do(req)
	if err != nil {
		return jobs
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return jobs
	}

	// Google Jobs data is embedded in page as JSON
	content := string(body)

	// Simple extraction - look for Google Jobs JSON patterns
	// Google uses special markup for job listings
	// This is a basic scraper - Google changes patterns frequently
	if strings.Contains(content, "job-searches") || strings.Contains(content, "data-jk") {
		// We can't reliably parse Google's JS-rendered content
		// The frontend should use a proper Google Jobs API if needed
		// Return empty for now, rely on other sources
	}

	return jobs
}

// JobStreet search
func (s *Searcher) searchJobStreet(query string, locations []string) []JobListing {
	jobs := []JobListing{}
	// JobStreet doesn't have a public API
	// Web scraping is fragile and often blocked
	// We'll skip it and rely on RemoteOK which has a real API
	return jobs
}

// LinkedIn search
func (s *Searcher) searchLinkedIn(query string, locations []string) []JobListing {
	jobs := []JobListing{}
	// LinkedIn has aggressive bot protection
	// Public job search requires authentication
	// We'll skip this and rely on RemoteOK
	return jobs
}

func (s *Searcher) matchesLocation(jobLocation string, preferredLocations []string) bool {
	if len(preferredLocations) == 0 || jobLocation == "" {
		return true
	}

	jobLocLower := strings.ToLower(jobLocation)
	for _, pref := range preferredLocations {
		if strings.Contains(jobLocLower, strings.ToLower(pref)) {
			return true
		}
	}

	// Always include "remote" jobs
	return strings.Contains(jobLocLower, "remote")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
