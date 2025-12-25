package main

import (
	"os"

	"github.com/RivLawrient/login-system/backend/config"
	"github.com/RivLawrient/login-system/backend/internal"
)

func main() {
	config.LoadEnv()
	db := config.GetConnection()
	app := config.NewGin()
	validate := config.NewValidator()

	internal.Apps(&internal.AppsConfig{
		DB:       db,
		App:      app,
		Validate: validate,
	})

	port := os.Getenv("BE_PORT")
	if port == "" {
		port = "8080"
	}
	app.Run(":" + port)
}
