// main.go
package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// incoming message structure
type ChatInput struct {
	Message string `json:"message"`
}

// outgoing reply structure
type ChatOutput struct {
	Reply string `json:"reply"`
}

type MistralRequest struct {
	Messages    []Message `json:"messages"`
	Temperature float32   `json:"temperature"`
	Stream      bool      `json:"stream"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type MistralResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found; using system env vars")
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("chatbot API running on :" + port)

	router := gin.Default()
	router.POST("/chat", handleChat)
	router.Run(":" + port)

}

func handleChat(c *gin.Context) {
	var input ChatInput
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	reply, err := callMistral(input.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Mistral error"})
		return
	}

	c.JSON(http.StatusOK, ChatOutput{Reply: reply})
}

func callMistral(userMessage string) (string, error) {
	mistralURL := os.Getenv("MISTRAL_URL")
	bearerToken := os.Getenv("MISTRAL_TOKEN") // this is safer than hardcoding

	if mistralURL == "" || bearerToken == "" {
		log.Fatal("Missing MISTRAL_URL or MISTRAL_TOKEN in .env")
	}

	payload := MistralRequest{
		Messages: []Message{
			{Role: "user", Content: userMessage},
		},
		Temperature: 0.7,
		Stream:      false,
	}

	jsonData, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", mistralURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+bearerToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	log.Printf("Mistral response body: %s", body)

	var result MistralResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	if len(result.Choices) > 0 {
		return result.Choices[0].Message.Content, nil
	}

	return "{no reply}", nil
}
