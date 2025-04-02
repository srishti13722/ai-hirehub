package handler

import (
	"context"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	openai "github.com/sashabaranov/go-openai"
	"github.com/srishti13722/ai-hirehub/ai-service/config"
)

type ResumeScoreRequest struct {
	ResumeText string `json:"resume_text"`
}

func ScoreResume(c *fiber.Ctx) error {
	var req ResumeScoreRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if req.ResumeText == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Resume text is required"})
	}

	prompt := `
You're a professional resume reviewer. Analyze the following resume and provide:

1. A score out of 100 based on relevance, formatting, clarity, and impact
2. What’s good about the resume
3. What can be improved (specific suggestions)
4. Missing keywords or skills for tech roles

Resume:
"""` + req.ResumeText + `"""

Output in this JSON format:
{
  "score": 0-100,
  "strengths": "...",
  "weaknesses": "...",
  "suggestions": "...",
  "missing_keywords": []
}
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
		"resume_review": resp.Choices[0].Message.Content,
	})
}
