package gemini

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"google.golang.org/genai"
)

var (
	GenaiClient    *genai.Client
	ActiveSessions = make(map[string]*genai.Chat)
)

func SaveSession(chatID string, history []*genai.Content) {
	dir := filepath.Join("commands", "gemini", "history")
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Printf("Error creating directory: %v", err)
		return
	}

	data, err := json.MarshalIndent(history, "", "  ")
	if err != nil {
		log.Printf("Error marshaling history: %v", err)
		return
	}
	loc := filepath.Join(dir, chatID+".json")
	err = os.WriteFile(loc, data, 0644)
	if err != nil {
		log.Printf("Error writing history: %v", err)
	}
}

func LoadSession(chatID string) []*genai.Content {
	var history []*genai.Content
	loc := filepath.Join("commands", "gemini", "history", chatID+".json")
	if _, err := os.Stat(loc); os.IsNotExist(err) {
		return nil
	}
	historyData, err := os.ReadFile(loc)
	if err != nil {
		log.Printf("Error reading history: %v", err)
		return nil
	}
	err = json.Unmarshal(historyData, &history)
	if err != nil {
		log.Printf("Error unmarshaling history: %v", err)
		return nil
	}
	
	return history
}
