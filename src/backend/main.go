package main

import (
	"github.com/chiefwhitecloud/golf-app/service"
	"os"
)

func main() {

	var PORT string

	if PORT = os.Getenv("PORT"); PORT == "" {
		PORT = "3001"
	}

	a := service.App{}
	a.Initialize(
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"))

	a.Run(":" + PORT)

}
