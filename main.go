package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/core"
)

const (
	InstanceShape = "VM.Standard.A1.Flex"
)

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
}

func main() {
	godotenv.Load()

	config, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}

	configProvider := common.DefaultConfigProvider()
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
		CompartmentID:      os.Getenv("COMPARTMENT_ID"),
		SubnetID:           os.Getenv("SUBNET_ID"),
		ImageID:            os.Getenv("IMAGE_ID"),
		SSHPublicKey:       os.Getenv("SSH_PUBLIC_KEY"),
		AvailabilityDomain: os.Getenv("AVAILABILITY_DOMAIN"),
		DisplayName:        os.Getenv("DISPLAY_NAME"),
		OCPUs:              float32(ocpus),
		Memory:             float32(memory),
		RetryDelay:         retryDelay,
		CheckOnly:          checkOnly,
		PeekBeforeLaunch:   peekBeforeLaunch,
	}

	if config.CompartmentID == "" || config.SubnetID == "" || config.ImageID == "" || config.AvailabilityDomain == "" || config.DisplayName == "" {
		return Config{}, log.New(os.Stderr, "", 0).Output(2, "Missing required environment variables.")
	}

	return config, nil
}

func attempt(ctx context.Context, client core.ComputeClient, config Config) bool {
	if config.PeekBeforeLaunch || config.CheckOnly {
		log.Printf("Peeking capacity in %s...", config.AvailabilityDomain)
		available := checkCapacity(ctx, client, config)
		if !available {
			log.Println("No capacity reported via peek. Skipping launch attempt.")
			return false
		}
		if config.CheckOnly {
			log.Println("CHECK_ONLY mode enabled. Skipping launch attempt.")
			return false
		}
	}

	log.Printf("Attempting to create instance in %s...", config.AvailabilityDomain)
	request := core.LaunchInstanceRequest{
		LaunchInstanceDetails: core.LaunchInstanceDetails{
			CompartmentId:      common.String(config.CompartmentID),
			AvailabilityDomain: common.String(config.AvailabilityDomain),
			DisplayName:        common.String(config.DisplayName),
			Shape:              common.String(InstanceShape),
			ShapeConfig: &core.LaunchInstanceShapeConfigDetails{
				Ocpus:       common.Float32(config.OCPUs),
				MemoryInGBs: common.Float32(config.Memory),
			},
			SourceDetails: core.InstanceSourceViaImageDetails{
				ImageId: common.String(config.ImageID),
			},
			CreateVnicDetails: &core.CreateVnicDetails{
				SubnetId:       common.String(config.SubnetID),
				AssignPublicIp: common.Bool(true),
			},
			Metadata: map[string]string{
				"ssh_authorized_keys": config.SSHPublicKey,
			},
		},
	}

	response, err := client.LaunchInstance(ctx, request)
	if err != nil {
		if strings.Contains(err.Error(), "Out of host capacity") {
			log.Println("Out of capacity. Retrying later...")
		} else if strings.Contains(err.Error(), "LimitExceeded") {
			log.Printf("Limit exceeded: %v", err)
			return true
		} else {
			log.Printf("Failed to launch instance: %v", err)
		}
		return false
	}

	log.Printf("Successfully launched instance! OCID: %s", *response.Instance.Id)
	return true
}

func checkCapacity(ctx context.Context, client core.ComputeClient, config Config) bool {
	request := core.CreateComputeCapacityReportRequest{
		CreateComputeCapacityReportDetails: core.CreateComputeCapacityReportDetails{
			CompartmentId:      common.String(config.CompartmentID),
			AvailabilityDomain: common.String(config.AvailabilityDomain),
			ShapeAvailabilities: []core.CreateCapacityReportShapeAvailabilityDetails{
				{
					InstanceShape: common.String(InstanceShape),
					InstanceShapeConfig: &core.CapacityReportInstanceShapeConfig{
						Ocpus:       common.Float32(config.OCPUs),
						MemoryInGBs: common.Float32(config.Memory),
					},
				},
			},
		},
	}

	response, err := client.CreateComputeCapacityReport(ctx, request)
	if err != nil {
		log.Printf("Failed to check capacity: %v", err)
		return false
	}

	if len(response.ComputeCapacityReport.ShapeAvailabilities) == 0 {
		return false
	}

	status := response.ComputeCapacityReport.ShapeAvailabilities[0].AvailabilityStatus
	log.Printf("Capacity status: %s", status)
	return status == core.CapacityReportShapeAvailabilityAvailabilityStatusAvailable
}
