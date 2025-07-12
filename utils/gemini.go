package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"regexp"

	"cloud.google.com/go/vertexai/genai"
	"google.golang.org/api/option"
)

func CallGemini(ctx context.Context, credsFile string, imageBytes []byte) (map[string]interface{}, error) {
	credsData, _ := os.ReadFile(credsFile)
	var creds struct {
		ProjectID string `json:"project_id"`
	}
	json.Unmarshal(credsData, &creds)

	client, err := genai.NewClient(ctx, creds.ProjectID, "us-central1", option.WithCredentialsFile(credsFile))
	if err != nil {
		return nil, err
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-2.0-flash-lite-001")
	prompt := `You are a visual moderation agent. Analyze the uploaded image and identify if it contains any of the following categories. For each category, reply Yes/No and explain briefly if Yes.

1. Nudity or Obscene content
2. Sexual acts or intent
3. QR codes (excluding on packaging)
4. Guns, weapons, or explosives
5. Violence or Gore
6. Dead bodies or corpses
7. Animal abuse or cruelty
8. Drugs, Alcohol, Tobacco
9. Hate symbols or middle finger
10. Gambling
11. PII (text containing phone numbers, emails, etc.)

Respond in this JSON format:
{
  "nudity": {{"flagged": true/false, "reason": "..." }},
  ...

  Note that flagged should be true/false no other values
}` // full prompt here
	resp, err := model.GenerateContent(ctx, genai.Text(prompt), genai.ImageData("image/jpeg", imageBytes))
	if err != nil {
		return nil, err
	}

	var result string
	for _, candidate := range resp.Candidates {
		for _, part := range candidate.Content.Parts {
			result += fmt.Sprintln(part)
		}
	}

	re := regexp.MustCompile(`(?s){.*}`)
	jsonMatch := re.FindString(result)
	if jsonMatch == "" {
		return nil, fmt.Errorf("No valid JSON found in Gemini output")
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(jsonMatch), &parsed); err != nil {
		return nil, fmt.Errorf("Failed to parse JSON: %v", err)
	}
	return parsed, nil

}
