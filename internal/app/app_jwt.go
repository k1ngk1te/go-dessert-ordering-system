package app

import "os"

type ApplicationJwt struct {
	SECRET string
}

func NewApplicationJwt(loggers *ApplicationLoggers) *ApplicationJwt {
	secret := os.Getenv("JWT_SECRET_KEY")

	if secret == "" {
		loggers.Error.Fatalf("error: JWT_SECRET_KEY is not valid.")
	}

	return &ApplicationJwt{
		SECRET: secret,
	}
}
