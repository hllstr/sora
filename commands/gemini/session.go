package gemini

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"google.golang.org/genai"
)

var (
	GenaiClient    *genai.Client
	ActiveSessions = make(map[string]*genai.Chat)
)

func SaveSession(chatID string, history []*genai.Content) {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Printf("Error getting current directory: %v", err)
		return
	}

	data, err := json.MarshalIndent(history, "", "  ")
	if err != nil {
		log.Printf("Error marshaling history: %v", err)
		return
	}
	loc := fmt.Sprintf(currentDir+"/commands/gemini//history/%s.json", chatID)
	err = os.WriteFile(loc, data, 0644)
	if err != nil {
		log.Printf("Error writing history: %v, \n Current Directory: %s", err, currentDir)
	}
}

func LoadSession(chatID string) []*genai.Content {
	var history []*genai.Content
	currentDir, err := os.Getwd()
	if err != nil {
		log.Printf("Error getting current directory: %v", err)
		return nil
	}
	loc := fmt.Sprintf(currentDir+"/commands/gemini/history/%s.json", chatID)
	if _, err := os.Stat(loc); os.IsNotExist(err) {
		return nil
	}
	historyData, err := os.ReadFile(loc)
	if err == nil {
		json.Unmarshal(historyData, &history)
	} else {
		log.Printf("Error reading history: %v", err)
	}
	return history
}
