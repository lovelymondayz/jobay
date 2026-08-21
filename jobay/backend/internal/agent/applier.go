package agent

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Applier handles auto-applying to jobs
type Applier struct {
	client *http.Client
}

func NewApplier() *Applier {
	return &Applier{
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

// ApplyAttempt attempts to apply to a job
func (a *Applier) ApplyAttempt(job JobListing, profile *UserProfile) (success bool, message string) {
	if job.URL == "" {
		return false, "No application URL"
	}

	req, err := http.NewRequest("GET", job.URL, nil)
	if err != nil {
		return false, fmt.Sprintf("Failed to open job page: %v", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")

	resp, err := a.client.Do(req)
	if err != nil {
		return false, fmt.Sprintf("Failed to fetch job page: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Sprintf("Failed to read response: %v", err)
	}

	content := string(body)

	if strings.Contains(content, "greenhouse") ||
		strings.Contains(content, "lever.co") ||
		strings.Contains(content, "workday") {
		return false, "Third-party ATS detected - manual apply required"
	}

	if strings.Contains(content, "<form") &&
		(strings.Contains(content, "apply") || strings.Contains(content, "application")) {
		return false, "Application form detected but auto-fill not yet implemented"
	}

	if strings.Contains(content, "http") && strings.Contains(content, "apply") {
		return true, "External application link found"
	}

	return true, "Application page accessible"
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
