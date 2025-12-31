package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
)

// notifyDiscord sends an embed message to Discord via webhook
func notifyDiscord(config Config, instanceOCID string) {
	if config.DiscordWebhookURL == "" {
		return
	}

	embed := map[string]interface{}{
		"title":       "🎉 OCI Instance Launched!",
		"description": "インスタンスの作成に成功しました。",
		"color":       0x00FF00,
		"fields": []map[string]interface{}{
			{"name": "Name", "value": config.DisplayName, "inline": true},
			{"name": "Shape", "value": InstanceShape, "inline": true},
			{"name": "OCPUs", "value": strconv.FormatFloat(float64(config.OCPUs), 'f', 0, 32), "inline": true},
			{"name": "Memory (GB)", "value": strconv.FormatFloat(float64(config.Memory), 'f', 0, 32), "inline": true},
			{"name": "Availability Domain", "value": config.AvailabilityDomain, "inline": false},
			{"name": "OCID", "value": instanceOCID, "inline": false},
		},
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}

	payload := map[string]interface{}{
		"embeds": []interface{}{embed},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal Discord payload: %v", err)
		return
	}

	resp, err := http.Post(config.DiscordWebhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Failed to send Discord notification: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		log.Println("Discord notification sent successfully.")
	} else {
		log.Printf("Discord notification failed with status: %d", resp.StatusCode)
	}
}
