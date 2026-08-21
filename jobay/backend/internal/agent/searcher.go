package agent

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// JobListing represents a job found from any source
type JobListing struct {
	Company     string `json:"company"`
	Role        string `json:"role"`
	URL         string `json:"url"`
	Location    string `json:"location"`
	Source      string `json:"source"`
	Description string `json:"description,omitempty"`
}

// Searcher handles job searching across multiple sources
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

// Search searches for jobs via available sources
func (s *Searcher) Search(profile *UserProfile, maxResults int) []JobListing {
	allJobs := []JobListing{}

	queries := s.generateQueries(profile)

	for _, query := range queries[:min(len(queries), 5)] {
		// Firecrawl search (if key is set)
		if s.apiKey != "" {
			jobs := s.searchFirecrawl(query, profile.PreferredLocations)
			allJobs = append(allJobs, jobs...)
		}

		// RemoteOK search (always available, no API key needed)
		jobs := s.searchRemoteOK(query, profile.PreferredLocations)
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
		"query": query,
		"limit": 10,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return jobs
	}

	req, err := http.NewRequest("POST", s.baseURL+"/search", strings.NewReader(string(body)))
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
		job := remoteOKToJobListing(item)
		if job != nil && s.matchesLocation(job.Location, locations) {
			jobs = append(jobs, *job)
		}
	}

	return jobs
}

func remoteOKToJobListing(item map[string]interface{}) *JobListing {
	company, _ := item["company"].(string)
	role, _ := item["position"].(string)
	jobURL, _ := item["url"].(string)
	location, _ := item["location"].(string)
	if location == "" {
		location = "Remote"
	}
	source, _ := item["source"].(string)

	return &JobListing{
		Company:  company,
		Role:     role,
		URL:      jobURL,
		Location: location,
		Source:   source,
	}
}

func extractCompanyFromTitle(title, pageURL string) string {
	lowerTitle := strings.ToLower(title)

	separators := []string{" - ", " | ", " · ", " @ "}
	for _, sep := range separators {
		if idx := strings.Index(lowerTitle, sep); idx > 0 {
			return strings.TrimSpace(title[:idx])
		}
	}

	atPatterns := []string{" at ", " with ", " for ", " hiring "}
	for _, pat := range atPatterns {
		if idx := strings.Index(lowerTitle, pat); idx > 0 {
			rest := title[idx+len(pat):]
			parts := strings.Fields(rest)
			if len(parts) >= 2 {
				return strings.Join(parts[:min(3, len(parts))], " ")
			}
			return rest
		}
	}

	u, err := url.Parse(pageURL)
	if err == nil {
		host := u.Hostname()
		host = strings.TrimPrefix(host, "www.")
		if idx := strings.Index(host, "."); idx > 0 {
			host = host[:idx]
		}
		return strings.Title(host)
	}

	return "Unknown"
}

func extractRoleFromTitle(title string) string {
	lowerTitle := strings.ToLower(title)
	separators := []string{" - ", " | ", " · ", " @ "}
	for _, sep := range separators {
		if idx := strings.Index(lowerTitle, sep); idx > 0 {
			return strings.TrimSpace(title[idx+len(sep):])
		}
	}

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
		"jakarta":    {"jakarta", "jkt"},
		"bandung":    {"bandung", "bdg"},
		"surabaya":   {"surabaya", "sby"},
		"bali":       {"bali", "denpasar"},
		"yogyakarta": {"yogyakarta", "jogja", "yogya"},
		"singapore":  {"singapore", "sg"},
		"malaysia":   {"malaysia", "kuala lumpur", "kl"},
		"remote":     {"remote", "work from home", "wfh", "anywhere"},
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
