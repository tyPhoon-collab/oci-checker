package main

import (
	"context"
	"log"
	"strings"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/core"
)

// attempt tries to launch an instance. Returns true if successful or limit exceeded.
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
	notifyDiscord(config, *response.Instance.Id)
	return true
}

// checkCapacity checks if capacity is available for the specified shape
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
