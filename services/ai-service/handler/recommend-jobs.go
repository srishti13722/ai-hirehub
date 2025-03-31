package handler

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	openai "github.com/sashabaranov/go-openai"
	"github.com/srishti13722/ai-hirehub/ai-service/config"
)

type RecommendRequest struct {
	ResumeText string   `json:"resume_text"` 
	Skills     []string `json:"skills"`      
}

func RecommendJobs(c *fiber.Ctx) error {
	var req RecommendRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Call job-service API
	jobServiceURL := os.Getenv("JOB_SERVICE_URL") 

	respJobs, err := http.Get(jobServiceURL)
	if err != nil || respJobs.StatusCode != http.StatusOK {
		log.Println("Failed to fetch jobs from job-service")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not fetch jobs from job-service"})
	}
	defer respJobs.Body.Close()

	body, _ := ioutil.ReadAll(respJobs.Body)

	var jobList []struct {
		JobTitle       string   `json:"job_title"`
		JobDescription string   `json:"job_description"`
		SkillsRequired []string `json:"skills_required"`
	}

	if err := json.Unmarshal(body, &jobList); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Error parsing job-service response"})
	}

	// Prepare job strings
	var jobLines []string
	for _, job := range jobList {
		line := job.JobTitle + " - " + job.JobDescription + " (" + strings.Join(job.SkillsRequired, ", ") + ")"
		jobLines = append(jobLines, line)
	}

	// Build prompt
	var base string
	if req.ResumeText != "" {
		base = "Resume:\n" + req.ResumeText
	} else {
		base = "Skills: " + strings.Join(req.Skills, ", ")
	}

	prompt := `
	Given the following candidate info:
	
	"""` + base + `"""
	
	And the following job openings:
	
	"""` + strings.Join(jobLines, "\n") + `"""
	
	Recommend the top 3 matching jobs and explain why they match in plain text.
	`

	// Ask OpenAI
	resp, err := config.OpenAIClient.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{Role: "user", Content: prompt},
			},
		},
	)
	if err != nil {
		log.Println("OpenAI Error:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "OpenAI failed: " + err.Error()})
	}

	return c.JSON(fiber.Map{
		"recommendations": resp.Choices[0].Message.Content,
	})
}
