package main

import (
	"log"
	"os"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/common/auth"
)

// getConfigurationProvider returns an OCI configuration provider.
// It tries Instance Principal first (for OCI instances), then falls back to ~/.oci/config.
func getConfigurationProvider() (common.ConfigurationProvider, error) {
	// 0. Skip Instance Principal if explicitly requested (useful for local development)
	if os.Getenv("OCI_SKIP_INSTANCE_PRINCIPAL") == "true" {
		log.Println("Skipping Instance Principal as requested via OCI_SKIP_INSTANCE_PRINCIPAL.")
	} else {
		// 1. Try Instance Principal (for OCI Instances)
		provider, err := auth.InstancePrincipalConfigurationProvider()
		if err == nil {
			log.Println("Using Instance Principal for authentication.")
			return provider, nil
		}
		log.Printf("Instance Principal not available: %v. Falling back to default config provider.", err)
	}

	// 2. Fallback to default config provider (~/.oci/config)
	log.Println("Using default config provider (~/.oci/config).")
	return common.DefaultConfigProvider(), nil
}
