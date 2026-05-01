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
	json.NewEncoder(w).Encode(result)
}

func callClaude(query string) (*MediaResult, error) {
	identifyMediaTool := anthropic.ToolParam{
		Name:        "identify_media",
		Description: anthropic.String("Identify the movie or TV show the user is describing"),
		InputSchema: anthropic.ToolInputSchemaParam{
			Properties: map[string]any{
				"title":   map[string]any{"type": "string", "description": "The movie or TV show title"},
				"year":    map[string]any{"type": "string", "description": "Release year"},
				"actors":  map[string]any{"type": "string", "description": "the main cast as a comma separated list e.g. 'Ben Stiller, Robert De Niro'"},
				"type":    map[string]any{"type": "string", "enum": []string{"Movie", "TV Show"}},
				"genre":   map[string]any{"type": "string", "description": "Comma-separated genres e.g. Comedy, Romance"},
				"summary": map[string]any{"type": "string", "description": "2-3 sentence description"},
				"rating":  map[string]any{"type": "number", "description": "IMDB rating as a number e.g. 8.6. Always include this field — use 0 if unknown"},
				},
				ExtraFields: map[string]any{
					"required": []string{"title", "year", "actors", "type", "genre", "summary", "rating"},
				},
			},
		}

	message, err := getClaudeClient().Messages.New(context.TODO(), anthropic.MessageNewParams{
		Model:      anthropic.ModelClaudeSonnet4_5_20250929,
		MaxTokens:  1024,
		Tools:      []anthropic.ToolUnionParam{{OfTool: &identifyMediaTool}},
		ToolChoice: anthropic.ToolChoiceParamOfTool("identify_media"),
		System: []anthropic.TextBlockParam{
			{Text: "You are a movie and TV show expert. Identify what the user is describing and call the identify_media tool with the result."},
		},
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(query)),
		},
	})
	if err != nil {
		return nil, err
	}

	for _, block := range message.Content {
		if toolUse, ok := block.AsAny().(anthropic.ToolUseBlock); ok {
			var result MediaResult
			if err := json.Unmarshal([]byte(toolUse.JSON.Input.Raw()), &result); err != nil {
				return nil, fmt.Errorf("failed to unmarshal tool result: %v", err)
			}
			return &result, nil
		}
	}

	return nil, fmt.Errorf("no tool use block in Claude response")
}
