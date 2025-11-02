package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// ---- CHATGPT HANDLER ----
func chat(prompt string) (string, error) {
	reqBody, _ := json.Marshal(map[string]interface{}{
		"model": "gpt-4o-mini",
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
	})

	req, _ := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(reqBody))
	req.Header.Set("Authorization", "Bearer "+os.Getenv("OPENAI_API_KEY"))
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	var resp struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "No response from OpenAI.", nil
	}

	return resp.Choices[0].Message.Content, nil
}

// ---- OLLAMA (CODE) HANDLER ----
func code(prompt string) (string, error) {
	ollamaURL := os.Getenv("OLLAMA_BASE_URL") + "/api/generate"
	body, _ := json.Marshal(map[string]string{
		"model":  os.Getenv("OLLAMA_MODEL"),
		"prompt": prompt,
	})

	req, _ := http.NewRequest("POST", ollamaURL, bytes.NewBuffer(body))
	req.Header.Set("X-API-Key", os.Getenv("OLLAMA_API_KEY"))
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	b, _ := io.ReadAll(res.Body)
	return string(b), nil
}

// ---- SERVER ----
func main() {
	http.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		var body struct{ Prompt string `json:"prompt"` }
		_ = json.NewDecoder(r.Body).Decode(&body)

		if body.Prompt == "" {
			http.Error(w, "Missing prompt", http.StatusBadRequest)
			return
		}

		reply, err := chat(body.Prompt)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{"response": reply})
	})

	http.HandleFunc("/code", func(w http.ResponseWriter, r *http.Request) {
		var body struct{ Prompt string `json:"prompt"` }
		_ = json.NewDecoder(r.Body).Decode(&body)

		if body.Prompt == "" {
			http.Error(w, "Missing prompt", http.StatusBadRequest)
			return
		}

		reply, err := code(body.Prompt)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{"response": reply})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("🚀 Multi-AI API running on port " + port)
	http.ListenAndServe(":"+port, nil)
}
