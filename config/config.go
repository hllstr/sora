package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

/*
	Nanti kalo mau nambahin Environtment Variables
	lain kaya API Key dsb., Bisa ditambahin dibagian sini aja.
*/

type Configuration struct {
	Owner  string
	Prefix []string
	Mode   string
}

func LoadConf() *Configuration {
	if err := godotenv.Load(); err != nil {
		log.Println("WARN: Environtment variables (.env) file not found.")
	}
	envPrefixes := os.Getenv("PREFIXES")
	return &Configuration{
		Owner:  os.Getenv("OWNER"),
		Prefix: strings.Split(envPrefixes, ","),
		Mode:   os.Getenv("MODE"),
	}
}
