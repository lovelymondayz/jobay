package agent

import (
	"encoding/json"
	"fmt"
	"strings"
)

// UserProfile represents the structured profile extracted from a CV
type UserProfile struct {
	Name               string   `json:"name"`
	Skills             []string `json:"skills"`
	YearsExperience    int      `json:"years_experience"`
	CurrentRole        string   `json:"current_role"`
	PreferredRoles    []string `json:"preferred_roles"`
	Location           string   `json:"location"`
	RemotePreference   bool     `json:"remote_preference"`
	SalaryMin          int      `json:"salary_min"`
	SalaryMax          int      `json:"salary_max"`
	Industries         []string `json:"industries"`
	EducationLevel     string   `json:"education_level"`
	Languages          []string `json:"languages"`
	Summary            string   `json:"summary"`
	PreferredLocations []string `json:"preferred_locations"`
}

var extractionSystemPrompt = `You are an expert CV/resume parser. Analyze the given CV text and extract structured information.
Return ONLY valid JSON with the following schema:
{
  "name": "string - full name",
  "skills": ["skill1", "skill2", ...] - technical and soft skills,
  "years_experience": number - total years of professional work experience,
  "current_role": "string - current or most recent job title",
  "preferred_roles": ["role1", "role2", ...] - desired job titles,
  "location": "string - current location (city, country)",
  "remote_preference": boolean - whether they prefer remote work,
  "salary_min": number - minimum salary expectation (USD, 0 if unknown),
  "salary_max": number - maximum salary expectation (USD, 0 if unknown),
  "industries": ["industry1", "industry2", ...],
  "education_level": "string - high_school, bachelor, master, phd, or other",
  "languages": ["lang1", "lang2", ...],
  "summary": "string - 1-2 sentence professional summary",
  "preferred_locations": ["city1", "city2", ...] - preferred work locations
}

Rules:
- Return ONLY the JSON object, no markdown, no explanation
- If a field is unknown, use reasonable defaults
- Years experience should be a conservative estimate
- Skills should include both technical and notable soft skills
- Array items should be lowercase`

var extractionUserPromptTemplate = `Extract the professional profile from this CV text:

---CV START---
%s
---CV END---

Return the JSON profile:`

// ExtractProfile uses AI to extract a structured profile from CV text
func ExtractProfile(cvText string) (*UserProfile, error) {
	ai := NewAIClient()

	prompt := fmt.Sprintf(extractionUserPromptTemplate, cvText)

	var profile UserProfile
	if err := ai.ChatJSON(extractionSystemPrompt, prompt, &profile); err != nil {
		// If JSON parsing fails, try to extract from a more lenient response
		response, chatErr := ai.Chat(extractionSystemPrompt, prompt)
		if chatErr != nil {
			return nil, fmt.Errorf("AI extraction failed: %w", err)
		}
		// Clean up response
		response = strings.TrimSpace(response)
		if idx := strings.Index(response, "{"); idx != -1 {
			response = response[idx:]
		}
		if idx := strings.LastIndex(response, "}"); idx != -1 {
			response = response[:idx+1]
		}
		if jsonErr := json.Unmarshal([]byte(response), &profile); jsonErr != nil {
			return nil, fmt.Errorf("failed to parse profile: %w", jsonErr)
		}
	}

	return &profile, nil
}
