package main

import (
	"cosmos-server/pkg/app"
	"cosmos-server/pkg/config"
	"fmt"
)

func main() {
	conf, err := config.NewConfiguration()
	if err != nil {
		fmt.Printf("Error reading config: %v\n", err)
		return
	}

	application, err := app.NewApp(conf)
	if err != nil {
		fmt.Printf("Error initializing application: %v\n", err)
		return
	}

	if err := application.RunServer(); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
		return
	}
}
