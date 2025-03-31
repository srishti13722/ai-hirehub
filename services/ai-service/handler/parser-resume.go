package handler

import (
	"context"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	openai "github.com/sashabaranov/go-openai"
	"github.com/srishti13722/ai-hirehub/ai-service/config"
)

type ResumeInput struct {
	ResumeText string `json:"resume_text"`
}

func ParseResume(c *fiber.Ctx) error {
	var input ResumeInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	prompt := `
Given the following resume text, extract the following information in JSON format:
- Full Name
- Contact Info
- Skills (as array)
- Education (as array)
- Work Experience (as array)
- Summary (1-2 lines)

Resume:
"""` + input.ResumeText + `"""

Return JSON like:
{
  "name": "",
  "contact_info": "",
  "skills": [],
  "education": [],
  "experience": [],
  "summary": ""
}
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
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to process resume"})
	}

	// Return raw output as JSON
	return c.JSON(fiber.Map{"parsed_resume": resp.Choices[0].Message.Content})
}
