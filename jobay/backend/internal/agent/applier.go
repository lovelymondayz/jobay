package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

// Applier handles auto-applying to jobs via Playwright worker
type Applier struct {
	client     *http.Client
	workerURL  string
}

func NewApplier() *Applier {
	workerURL := os.Getenv("WORKER_URL")
	if workerURL == "" {
		workerURL = "http://localhost:3011"
	}

	return &Applier{
		client:    &http.Client{Timeout: 60 * time.Second},
		workerURL: workerURL,
	}
}

// ApplyAttempt attempts to apply to a job via Playwright worker
func (a *Applier) ApplyAttempt(job JobListing, profile *UserProfile) (success bool, message string) {
	if job.URL == "" {
		return false, "No application URL"
	}

	// Call Playwright worker
	reqBody := map[string]interface{}{
		"jobUrl":     job.URL,
		"userProfile": profile,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return false, fmt.Sprintf("Failed to marshal request: %v", err)
	}

	req, err := http.NewRequest("POST", a.workerURL+"/apply", bytes.NewBuffer(body))
	if err != nil {
		return false, fmt.Sprintf("Failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return false, fmt.Sprintf("Worker request failed: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Sprintf("Failed to read worker response: %v", err)
	}

	var result struct {
		Ok          bool     `json:"ok"`
		Filled      []string `json:"filled"`
		Submitted   bool     `json:"submitted"`
		Screenshot  string   `json:"screenshot"`
		Details     []string `json:"details"`
		Error       string   `json:"error"`
	}

	if err := json.Unmarshal(respBody, &result); err != nil {
		return false, fmt.Sprintf("Failed to parse worker response: %v", err)
	}

	if !result.Ok {
		return false, result.Error
	}

	if result.Submitted {
		return true, fmt.Sprintf("Applied! Filled: %v", result.Filled)
	}

	if len(result.Filled) > 0 {
		return false, fmt.Sprintf("Filled %v but didn't submit", result.Filled)
	}

	return false, "No form fields found"
}

func (a *Applier) DetectApplicationURL(pageURL string) string {
	return pageURL
}

func (a *Applier) ExtractFormFields(htmlContent string) map[string]string {
	return make(map[string]string)
}

func (a *Applier) FillForm(formURL string, fields map[string]string, profile *UserProfile) error {
	data := url.Values{}
	data.Set("name", profile.Name)
	data.Set("email", "")
	data.Set("phone", "")

	req, err := http.NewRequest("POST", formURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36")

	resp, err := a.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("form submission failed with status %d", resp.StatusCode)
	}

	return nil
}
