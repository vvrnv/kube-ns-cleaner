package config

import (
	"os"
	"strings"

	log "github.com/rs/zerolog/log"

	"github.com/spf13/viper"
)

// config is global object that holds all application level variables.
var Config appConfig

type appConfig struct {
	ExcludedNamespaces []string `mapstructure:"excludedNamespaces"`
	ScalingLifeTime    int      `mapstructure:"scalingLifeTime"`
	DeleteingLifeTime  int      `mapstructure:"deletingLifeTime"`
	Cron               string   `mapstructure:"cron"`
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

	// excluded namespaces from config
	excludedNamespaces := v.GetStringSlice("excludedNamespaces")

	if os.Getenv("KUBE_NS_CLEANER_LOGS_DIR") == "" {
		os.Setenv("KUBE_NS_CLEANER_LOGS_DIR", "kube-ns-cleaner.json")
		log.Info().Msgf("ENV variable `KUBE_NS_CLEANER_LOGS_DIR` is not set. Default logs dir are %s", os.Getenv("KUBE_NS_CLEANER_LOGS_DIR"))
	}

	if len(excludedNamespaces) != 0 {
		log.Info().Msg("Excluded namespaces: " + strings.Join(excludedNamespaces, ", "))
	} else {
		log.Info().Msg("Excluded namespaces is not set in config file")
	}

	return v.Unmarshal(&Config)
}
