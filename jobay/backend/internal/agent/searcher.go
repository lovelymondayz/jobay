package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// Firecrawl job listing from search results
type FirecrawlJobListing struct {
	Title       string `json:"title"`
	URL         string `json:"url"`
	Description string `json:"description"`
	Company     string `json:"companyName,omitempty"`
	Location    string `json:"location,omitempty"`
	Source      string `json:"source,omitempty"`
}

// JobListing represents a job found from any source
type JobListing struct {
	Company     string `json:"company"`
	Role        string `json:"role"`
	URL         string `json:"url"`
	Location    string `json:"location"`
	Source      string `json:"source"`
	Salary      string `json:"salary"`
	JobType     string `json:"job_type"`
	PostedDate  string `json:"posted_date"`
	Description string `json:"description,omitempty"`
}

// Searcher handles job searching via Firecrawl
type Searcher struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

func NewSearcher() *Searcher {
	apiKey := os.Getenv("FIRECRAWL_API_KEY")
	baseURL := "https://api.firecrawl.dev/v1"
	
	return &Searcher{
		apiKey:  apiKey,
		baseURL: baseURL,
		client:  &http.Client{Timeout: 30 * time.Second},
	}
}

// Search searches for jobs via Firecrawl
func (s *Searcher) Search(profile *UserProfile, maxResults int) []JobListing {
	allJobs := []JobListing{}
	
	queries := s.generateQueries(profile)
	
	for _, query := range queries[:min(len(queries), 5)] {
		jobs := s.searchFirecrawl(query, profile.PreferredLocations)
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
	
	roles := profile.PreferredRoles
	if len(roles) == 0 && profile.CurrentRole != "" {
		roles = []string{profile.CurrentRole}
	}
	
	locations := profile.PreferredLocations
	if len(locations) == 0 && profile.Location != "" {
		locations = []string{profile.Location}
	}
	
	for _, role := range roles[:min(len(roles), 3)] {
		for _, loc := range locations[:min(len(locations), 3)] {
			queries = append(queries, fmt.Sprintf("%s %s jobs", role, loc))
		}
		if profile.RemotePreference {
			queries = append(queries, fmt.Sprintf("%s remote jobs", role))
		}
	}
	
	if len(queries) == 0 {
		queries = []string{"software engineer jakarta", "developer indonesia"}
	}
	
	return queries
}

func (s *Searcher) searchFirecrawl(query string, locations []string) []JobListing {
	jobs := []JobListing{}
	
	reqBody := map[string]interface{}{
		"query":    query,
		"limit":    10,
		"scrapeOptions": map[string]interface{}{
			"formats": []string{"markdown"},
		},
	}
	
	body, err := json.Marshal(reqBody)
	if err != nil {
		return jobs
	}
	
	req, err := http.NewRequest("POST", s.baseURL+"/search", bytes.NewBuffer(body))
	if err != nil {
		return jobs
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.apiKey)
	
	resp, err := s.client.Do(req)
	if err != nil {
		return jobs
	}
	defer resp.Body.Close()
	
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return jobs
	}
	
	var result struct {
		Success bool `json:"success"`
		Data    []struct {
			URL         string `json:"url"`
			Title       string `json:"title"`
			Description string `json:"description"`
			Markdown    string `json:"markdown"`
			Metadata    struct {
				Title       string `json:"title"`
				Description string `json:"description"`
			} `json:"metadata"`
		} `json:"data"`
	}
	
	if err := json.Unmarshal(respBody, &result); err != nil {
		return jobs
	}
	
	if !result.Success {
		return jobs
	}
	
	for _, item := range result.Data {
		job := s.parseSearchResult(item.Title, item.URL, item.Description)
		if job != nil && s.matchesLocation(job.Location, locations) {
			jobs = append(jobs, *job)
		}
	}
	
	return jobs
}

func (s *Searcher) parseSearchResult(title, pageURL, description string) *JobListing {
	if pageURL == "" {
		return nil
	}
	
	// Extract company from title or URL
	company := extractCompanyFromTitle(title, pageURL)
	location := extractLocationFromText(description)
	
	role := extractRoleFromTitle(title)
	if role == "" {
		role = "Software Engineer"
	}
	
	return &JobListing{
		Company:     company,
		Role:        role,
		URL:         pageURL,
		Location:    location,
		Source:      "Firecrawl",
		Description: description,
	}
}

func extractCompanyFromTitle(title, pageURL string) string {
	// Try common patterns
	lowerTitle := strings.ToLower(title)
	
	// Pattern: "Company - Role" or "Company | Role"
	separators := []string{" - ", " | ", " · ", " @ "}
	for _, sep := range separators {
		if idx := strings.Index(lowerTitle, sep); idx > 0 {
			return strings.TrimSpace(title[:idx])
		}
	}
	
	// Pattern: "at Company" or "with Company"
	atPatterns := []string{" at ", " with ", " for ", " hiring "}
	for _, pat := range atPatterns {
		if idx := strings.Index(lowerTitle, pat); idx > 0 {
			rest := title[idx+len(pat):]
			// Take the company name (first 2-3 words after "at")
			parts := strings.Fields(rest)
			if len(parts) >= 2 {
				return strings.Join(parts[:min(3, len(parts))], " ")
			}
			return rest
		}
	}
	
	// Fallback: use domain name
	u, err := url.Parse(pageURL)
	if err == nil {
		host := u.Hostname()
		// Remove www. and .com/.co.id etc
		host = strings.TrimPrefix(host, "www.")
		if idx := strings.Index(host, "."); idx > 0 {
			host = host[:idx]
		}
		return strings.Title(host)
	}
	
	return "Unknown"
}

func extractRoleFromTitle(title string) string {
	// Remove company prefix
	lowerTitle := strings.ToLower(title)
	separators := []string{" - ", " | ", " · ", " @ "}
	for _, sep := range separators {
		if idx := strings.Index(lowerTitle, sep); idx > 0 {
			return strings.TrimSpace(title[idx+len(sep):])
		}
	}
	
	// Remove "hiring at Company" prefix
	atPatterns := []string{" is hiring ", " is looking for ", " is searching for "}
	for _, pat := range atPatterns {
		if idx := strings.Index(lowerTitle, pat); idx > 0 {
			return strings.TrimSpace(title[idx+len(pat):])
		}
	}
	
	return title
}

func extractLocationFromText(text string) string {
	if text == "" {
		return ""
	}
	
	lower := strings.ToLower(text)
	
	locations := map[string][]string{
		"jakarta": {"jakarta", "jkt"},
		"bandung": {"bandung", "bdg"},
		"surabaya": {"surabaya", "sby"},
		"bali": {"bali", "denpasar"},
		"yogyakarta": {"yogyakarta", "jogja", "yogya"},
		"singapore": {"singapore", "sg"},
		"malaysia": {"malaysia", "kuala lumpur", "kl"},
		"remote": {"remote", "work from home", "wfh", "anywhere"},
	}
	
	for loc, keywords := range locations {
		for _, kw := range keywords {
			if strings.Contains(lower, kw) {
				return strings.Title(loc)
			}
		}
	}
	
	return ""
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
	
	return strings.Contains(jobLocLower, "remote")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
