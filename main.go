package main

import (
	"context"
	"cosmos-server/pkg/app"
	"cosmos-server/pkg/config"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
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

	err = application.SetUpDatabase()
	if err != nil {
		fmt.Printf("Error setting up database: %v\n", err)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	application.StartSentinel(ctx)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	serverErrors := make(chan error, 1)
	go func() {
		if err := application.RunServer(); err != nil {
			serverErrors <- err
		}
	}()

	select {
	case <-quit:
		fmt.Println("Received shutdown signal...")
	case err := <-serverErrors:
		fmt.Printf("Server error: %v\n", err)
	}

	// We cancel the context of the sentinels
	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := application.Shutdown(shutdownCtx); err != nil {
		fmt.Printf("Server forced to shutdown: %v\n", err)
		return
	}

	// I wait a little to let the workers finish
	time.Sleep(2 * time.Second)

	fmt.Println("Server exited gracefully")
}
