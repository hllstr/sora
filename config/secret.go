package config

import (
	"os"
)

type Secrets struct {
	Owner string
}

func loadSecret() Secrets {
	return Secrets{
		Owner: os.Getenv("OWNER"),
	}
}
