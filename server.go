package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/SnowAI2003/multi-ai-cmd/internal/config"
	"github.com/SnowAI2003/multi-ai-cmd/internal/llm"
)

func main() {
	cfg := config.Load()
	openaiClient := llm.NewOpenAIClient(cfg)
	ollamaClient := llm.NewOllamaClient(cfg)

	http.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		var body struct{ Prompt string `json:"prompt"` }
		_ = json.NewDecoder(r.Body).Decode(&body)
		if strings.TrimSpace(body.Prompt) == "" {
			http.Error(w, "missing prompt", http.StatusBadRequest)
			return
		}
		reply, err := openaiClient.Chat(r.Context(), body.Prompt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"response": reply})
	})

	http.HandleFunc("/code", func(w http.ResponseWriter, r *http.Request) {
		var body struct{ Prompt string `json:"prompt"` }
		_ = json.NewDecoder(r.Body).Decode(&body)
		if strings.TrimSpace(body.Prompt) == "" {
			http.Error(w, "missing prompt", http.StatusBadRequest)
			return
		}
		reply, err := ollamaClient.Generate(r.Context(), body.Prompt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"response": reply})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Println("Server running on port " + port)
	http.ListenAndServe(":"+port, nil)
}
