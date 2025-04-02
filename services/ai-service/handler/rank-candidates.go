package handler

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	openai "github.com/sashabaranov/go-openai"
	"github.com/srishti13722/ai-hirehub/ai-service/config"
)

type CandidateRankingRequest struct {
	JobDescription string   `json:"job_description"`
	Resumes        []string `json:"resumes"` 
}

func RankCandidates(c *fiber.Ctx) error {
	var req CandidateRankingRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if req.JobDescription == "" || len(req.Resumes) == 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Job description and at least one resume are required"})
	}

	var sb strings.Builder
	for i, resume := range req.Resumes {
		sb.WriteString("Candidate ")
		sb.WriteString(string(rune('A' + i)))
		sb.WriteString(":\n")
		sb.WriteString(resume + "\n\n")
	}

	prompt := `
You are a technical recruiter. Given the job description and the following candidate resumes, rank them from best fit to least fit and explain your reasoning for each.

Job Description:
"""` + req.JobDescription + `"""

Candidate Resumes:
"""` + sb.String() + `"""

Output in this format:
1. Candidate A - Excellent fit because...
2. Candidate B - Good fit...
3. Candidate C - Lacks...
`

	resp, err := config.OpenAIClient.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4,
			Messages: []openai.ChatCompletionMessage{
				{Role: "user", Content: prompt},
			},
		},
	)
	if err != nil {
		log.Println("OpenAI Error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "OpenAI failed: " + err.Error()})
	}

	return c.JSON(fiber.Map{
		"ranking": resp.Choices[0].Message.Content,
	})
}
