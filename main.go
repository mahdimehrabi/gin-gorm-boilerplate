package main

import (
	bootstrap "boilerplate/core"
	"log"
	"os"

	"github.com/getsentry/sentry-go"
	"github.com/joho/godotenv"
	"go.uber.org/fx"
)

func main() {
	godotenv.Load()

	err := sentry.Init(
		sentry.ClientOptions{
			Dsn:         os.Getenv("SentryDSN"),
			Environment: os.Getenv("Environment"),
			Release:     os.Getenv("AppName") + "@1.0.0",
			// Enable printing of SDK debug messages.
			// Useful when getting started or trying to figure something out.
			Debug: true,
		},
	)
	if err != nil {
		log.Fatal("failed to init sentry: ", err)
		return
	}
	fx.New(bootstrap.BootstrapModule).Run()
}
