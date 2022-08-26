//go:build test && !integration
// +build test,!integration

package config

import "github.com/RHEnVision/provisioning-backend/internal/config/parser"

func init() {
	parser.Viper.Set("featureFlags.environment", "test")
}
