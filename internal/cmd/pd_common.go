package cmd

import (
	"github.com/engswee/flashpipe/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Helper functions for Partner Directory commands to support reading
// configuration from both command-line flags and nested config file keys

// getConfigStringWithFallback reads a string value from command flag,
// falling back to a nested config key if the flag wasn't explicitly set
func getConfigStringWithFallback(cmd *cobra.Command, flagName, configKey string) string {
	// Check if flag was explicitly set on command line
	if cmd.Flags().Changed(flagName) {
		return config.GetString(cmd, flagName)
	}

	// Try to get from nested config key
	if viper.IsSet(configKey) {
		return viper.GetString(configKey)
	}

	// Fall back to flag default
	return config.GetString(cmd, flagName)
}

// getConfigBoolWithFallback reads a bool value from command flag,
// falling back to a nested config key if the flag wasn't explicitly set
func getConfigBoolWithFallback(cmd *cobra.Command, flagName, configKey string) bool {
	// Check if flag was explicitly set on command line
	if cmd.Flags().Changed(flagName) {
		return config.GetBool(cmd, flagName)
	}

	// Try to get from nested config key
	if viper.IsSet(configKey) {
		return viper.GetBool(configKey)
	}

	// Fall back to flag default
	return config.GetBool(cmd, flagName)
}

// getConfigStringSliceWithFallback reads a string slice value from command flag,
// falling back to a nested config key if the flag wasn't explicitly set
func getConfigStringSliceWithFallback(cmd *cobra.Command, flagName, configKey string) []string {
	// Check if flag was explicitly set on command line
	if cmd.Flags().Changed(flagName) {
		return config.GetStringSlice(cmd, flagName)
	}

	// Try to get from nested config key
	if viper.IsSet(configKey) {
		return viper.GetStringSlice(configKey)
	}

	// Fall back to flag default
	return config.GetStringSlice(cmd, flagName)
}

// contains checks if a string slice contains a specific string
func contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}
