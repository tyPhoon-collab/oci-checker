package main

import (
	"log"
	"os"
	"strconv"
	"strings"
)

// Config holds all configuration for the OCI checker
type Config struct {
	CompartmentID      string
	SubnetID           string
	ImageID            string
	SSHPublicKey       string
	AvailabilityDomain string
	DisplayName        string
	OCPUs              float32
	Memory             float32
	RetryDelay         int
	CheckOnly          bool
	PeekBeforeLaunch   bool
	DiscordWebhookURL  string
}

// loadConfig reads configuration from environment variables
func loadConfig() (Config, error) {
	ocpusStr := os.Getenv("OCPUS")
	memoryStr := os.Getenv("MEMORY_IN_GBS")
	retryDelayStr := os.Getenv("RETRY_DELAY")

	ocpus, _ := strconv.ParseFloat(ocpusStr, 32)
	memory, _ := strconv.ParseFloat(memoryStr, 32)
	retryDelay, _ := strconv.Atoi(retryDelayStr)
	if retryDelay == 0 {
		retryDelay = 60
	}

	checkOnly := strings.ToLower(os.Getenv("CHECK_ONLY")) == "true"
	peekBeforeLaunch := strings.ToLower(os.Getenv("PEEK_BEFORE_LAUNCH")) == "true"

	config := Config{
		CompartmentID:      os.Getenv("OCI_COMPARTMENT_ID"),
		SubnetID:           os.Getenv("OCI_SUBNET_ID"),
		ImageID:            os.Getenv("OCI_IMAGE_ID"),
		SSHPublicKey:       os.Getenv("OCI_SSH_PUBLIC_KEY"),
		AvailabilityDomain: os.Getenv("OCI_AVAILABILITY_DOMAIN"),
		DisplayName:        os.Getenv("OCI_DISPLAY_NAME"),
		OCPUs:              float32(ocpus),
		Memory:             float32(memory),
		RetryDelay:         retryDelay,
		CheckOnly:          checkOnly,
		PeekBeforeLaunch:   peekBeforeLaunch,
		DiscordWebhookURL:  os.Getenv("DISCORD_WEBHOOK_URL"),
	}

	if config.CompartmentID == "" || config.SubnetID == "" || config.ImageID == "" || config.AvailabilityDomain == "" || config.DisplayName == "" {
		return Config{}, log.New(os.Stderr, "", 0).Output(2, "Missing required environment variables.")
	}

	return config, nil
}
