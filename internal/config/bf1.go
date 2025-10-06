// Package config provides configuration management for the grompt
package config

import (
	"os"
	"strconv"
	"strings"
)

// BF1Config holds configuration for BF1 (loop v1) mode
type BF1Config struct {
	Enabled    bool
	WIPCap     int
	CooldownH  int
	CanaryOnly bool
}

// GetBF1Config returns BF1 mode configuration from environment
func GetBF1Config() BF1Config {
	config := BF1Config{
		Enabled:    false,
		WIPCap:     1,
		CooldownH:  24,
		CanaryOnly: true,
	}

	// Check if BF1_MODE is enabled
	if bf1Mode := os.Getenv("BF1_MODE"); bf1Mode != "" {
		config.Enabled = strings.ToLower(bf1Mode) == "true"
	}

	// Override WIP cap if specified
	if wipCap := os.Getenv("BF1_WIP_CAP"); wipCap != "" {
		if cap, err := strconv.Atoi(wipCap); err == nil {
			config.WIPCap = cap
		}
	}

	// Override cooldown if specified
	if cooldown := os.Getenv("BF1_COOLDOWN_HOURS"); cooldown != "" {
		if hours, err := strconv.Atoi(cooldown); err == nil {
			config.CooldownH = hours
		}
	}

	// Override canary mode if specified
	if canary := os.Getenv("BF1_CANARY_ONLY"); canary != "" {
		config.CanaryOnly = strings.ToLower(canary) == "true"
	}

	return config
}

// IsBF1Mode returns true if BF1 mode is enabled
func IsBF1Mode() bool {
	return GetBF1Config().Enabled
}
