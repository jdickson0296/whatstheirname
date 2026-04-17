package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

type SearchRequest struct {
	Query string `json:"query"`
}

type SearchResponse struct {
	Result string `json:"result"`
}

var claudeClient *anthropic.Client

func getClaudeClient() *anthropic.Client {
	if claudeClient == nil {
		c := anthropic.NewClient(
			option.WithAPIKey(os.Getenv("ANTHROPIC_API_KEY")),
		)
		claudeClient = &c
	}
	return claudeClient
}

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req SearchRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Query == "" {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	result, err := callClaude(req.Query)
	if err != nil {
		log.Printf("failed to call claude: %v", err)
		http.Error(w, fmt.Sprintf("failed to call Claude: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(SearchResponse{Result: result}); err != nil {
		log.Printf("failed to encode response: %v", err)
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func callClaude(query string) (string, error) {
	message, err := getClaudeClient().Messages.New(context.TODO(), anthropic.MessageNewParams{
		Model:     anthropic.ModelClaudeSonnet4_5_20250929,
		MaxTokens: 1024,
		System: []anthropic.TextBlockParam{
			{Text: "You are a movie, TV show, and actor expert. Answer questions concisely."},
		},
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(query)),
		},
	})

	if err != nil {
		return "", err
	}

	return message.Content[0].Text, nil
}
