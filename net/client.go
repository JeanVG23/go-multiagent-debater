package net

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// ChatMessage représente un message pour l'API Groq (role: system, user, assistant)
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest correspond au payload envoyé à /openai/v1/chat/completions
type ChatRequest struct {
	Model    string        `json:"model"`
	Messages []ChatMessage `json:"messages"`
}

// ChatResponse correspond au format de réponse de Groq
type ChatResponse struct {
	Choices []struct {
		Message ChatMessage `json:"message"`
	} `json:"choices"`
}

// ClientGroq encapsule la clé d'API et le client HTTP pour échanger avec Groq
type ClientGroq struct {
	apiKey     string
	httpClient *http.Client
}

// NewClientGroq initialise le client en chargeant la clé d'API depuis l'environnement ou .env
func NewClientGroq() (*ClientGroq, error) {
	// Charge le .env s'il existe (dans le dossier courant ou un des dossiers parents)
	_ = godotenv.Load()
	_ = godotenv.Load("../.env")
	_ = godotenv.Load("../../.env")
	_ = godotenv.Load("../../../.env")

	apiKey := os.Getenv("GROQ_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("la variable GROQ_API_KEY n'est pas définie dans le .env ou l'environnement")
	}

	return &ClientGroq{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 45 * time.Second,
		},
	}, nil
}

// Chat envoie une liste ordonnée de messages (historique) à un modèle Groq donné
func (c *ClientGroq) Chat(model string, messages []ChatMessage) (string, error) {
	payload := ChatRequest{
		Model:    model,
		Messages: messages,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("erreur marshal JSON : %w", err)
	}

	url := "https://api.groq.com/openai/v1/chat/completions"
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("erreur création de requête : %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("erreur envoi requête HTTP : %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("erreur API Groq (status %d) : %s", resp.StatusCode, string(body))
	}

	var result ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("erreur décodage réponse JSON : %w", err)
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("aucun choix renvoyé par Groq")
	}

	return result.Choices[0].Message.Content, nil
}

// Query est un raccourci pour envoyer un prompt utilisateur avec éventuellement un prompt système
func (c *ClientGroq) Query(model string, systemPrompt string, userPrompt string) (string, error) {
	var messages []ChatMessage
	if systemPrompt != "" {
		messages = append(messages, ChatMessage{
			Role:    "system",
			Content: systemPrompt,
		})
	}
	messages = append(messages, ChatMessage{
		Role:    "user",
		Content: userPrompt,
	})
	return c.Chat(model, messages)
}
