package handler

import (
	"context"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	openai "github.com/sashabaranov/go-openai"
	"github.com/srishti13722/ai-hirehub/ai-service/config"
)

type CoverLetterRequest struct {
	ResumeText     string `json:"resume_text"`
	JobDescription string `json:"job_description"`
}

func GenerateCoverLetter(c *fiber.Ctx) error {
	var req CoverLetterRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if req.ResumeText == "" || req.JobDescription == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Both resume_text and job_description are required"})
	}

	prompt := `
Given the following resume and job description, write a personalized, professional cover letter tailored for this job. The tone should be enthusiastic, confident, and highlight matching skills and experiences.

Resume:
"""` + req.ResumeText + `"""

Job Description:
"""` + req.JobDescription + `"""

Output only the cover letter, no headers.
`

	resp, err := config.OpenAIClient.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    "user",
					Content: prompt,
				},
			},
		},
	)
	if err != nil {
		log.Println("OpenAI Error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "OpenAI failed: " + err.Error()})
	}

	return c.JSON(fiber.Map{
		"cover_letter": resp.Choices[0].Message.Content,
	})
}
