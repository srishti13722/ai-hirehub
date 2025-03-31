package config

import (
	"os"

	openai "github.com/sashabaranov/go-openai"
)

var OpenAIClient *openai.Client

func InitOpenAI() {
	OpenAIClient = openai.NewClient(os.Getenv("OPENAI_API_KEY"))
}
