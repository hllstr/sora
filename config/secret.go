package config

import (
	"os"
	"strings"
)

type Secrets struct {
	Owner []string
}

func loadSecret() Secrets {
	ownerEnv := os.Getenv("OWNER")
	return Secrets{
		Owner: strings.Split(ownerEnv, ","),
	}
}
