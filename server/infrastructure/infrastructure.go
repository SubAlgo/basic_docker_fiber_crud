package infrastructure

import (
	"os"
	"strings"
)

type Infrastructure struct {
	AppPort string
}

func Get() Infrastructure {
	appPort := os.Getenv("APP_PORT")
	if strings.TrimSpace(appPort) == "" {
		appPort = "9000"
	}

	return Infrastructure{
		AppPort: appPort,
	}
}
