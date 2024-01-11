package config

import (
	"strings"

	log "github.com/rs/zerolog/log"

	"github.com/spf13/viper"
)

// config is global object that holds all application level variables.
var Config appConfig

type appConfig struct {
	DefaultExcludedNamespaces []string `mapstructure:"defaultExcludedNamespaces"`
	ExcludedNamespaces        []string `mapstructure:"excludedNamespaces"`
	ScalingLifeTime           int      `mapstructure:"scalingLifeTime"`
	DeleteingLifeTime         int      `mapstructure:"deletingLifeTime"`
	Cron                      string   `mapstructure:"cron"`
}

// LoadConfig loads config from files
func LoadConfig(configPaths ...string) error {
	log.Info().Msg(("Check config"))
	v := viper.New()

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("/opt/app/")
	err := v.ReadInConfig()
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	// default excluded namespaces
	defaultExcludedNamespaces := v.GetStringSlice("defaultExcludedNamespaces")
	if len(defaultExcludedNamespaces) == 0 {
		log.Fatal().Msg("Default excluded namespaces is not set in config file")
	} else {
		log.Info().Msg("Default excluded namespaces: " + strings.Join(v.GetStringSlice("defaultExcludedNamespaces"), ", "))
	}

	// excluded namespaces from config
	excludedNamespaces := v.GetStringSlice("excludedNamespaces")

	if len(excludedNamespaces) != 0 {
		log.Info().Msg("Excluded namespaces: " + strings.Join(excludedNamespaces, ", "))
	} else {
		log.Info().Msg("Excluded namespaces is not set in config file")
	}

	return v.Unmarshal(&Config)
}
