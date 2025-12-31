package main

import (
	"context"
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/oracle/oci-go-sdk/v65/core"
)

const (
	InstanceShape = "VM.Standard.A1.Flex"
)

func main() {
	godotenv.Load()

	config, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}

	configProvider, err := getConfigurationProvider()
	if err != nil {
		log.Fatalf("Failed to create configuration provider: %v", err)
	}

	computeClient, err := core.NewComputeClientWithConfigurationProvider(configProvider)
	if err != nil {
		log.Fatalf("Failed to create compute client: %v", err)
	}

	for {
		success := attempt(context.Background(), computeClient, config)
		if success && !config.CheckOnly {
			log.Println("Instance created or limit reached. Exiting.")
			break
		}

		log.Printf("Waiting for %d seconds before next attempt...", config.RetryDelay)
		time.Sleep(time.Duration(config.RetryDelay) * time.Second)
	}
}
