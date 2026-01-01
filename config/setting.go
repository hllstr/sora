package config

import (
	"encoding/json"
	"os"
)

type Settings struct {
	Prefix []string `json:"prefixes"`
	Mode   string   `json:"mode"`
}

func loadSetting() Settings {
	var set Settings
	file, err := os.ReadFile("settings.json")
	if err == nil {
		json.Unmarshal(file, &set)
	}
	return set
}
