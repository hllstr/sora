package config

import (
	"log"

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
	secret := loadSecret()
	setting := loadSetting()

	return &Configuration{
		Owner:  secret.Owner,
		Prefix: setting.Prefix,
		Mode:   setting.Mode,
	}
}
