package config

import (
	"context"
	"strings"

	"github.com/jainabhishek5986/employee-records/pkg/zaplogger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// GlobalConfig stores the config instance for global use
var GlobalConfig *Config

// Load loads config from command instance to predefined config variables
func Load(cmd *cobra.Command) (*Config, error) {
	err := viper.BindPFlags(cmd.Flags())
	if err != nil {
		return nil, err
	}

	// default viper configs
	viper.SetEnvPrefix("EXM")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// set default configs
	setDefaultConfig()

	if configFile, _ := cmd.Flags().GetString("config"); configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		viper.SetConfigName(".employee-records-service")
		viper.AddConfigPath("./")
		viper.AddConfigPath("$HOME/.employee-records-service")
	}

	if err := viper.ReadInConfig(); err != nil {
		zaplogger.Warn(context.Background(), "No configuration file found. Proceeding with defaults")
	}

	return populateConfig(new(Config))
}
