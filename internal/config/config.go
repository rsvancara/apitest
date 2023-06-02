package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

//AppConfig Application Configuration
type AppConfig struct {
	Env            string `envconfig:"ENV"`            //PROD,DEV
	Site           string `envconfig:"SITE"`           // defines site name and location of template directories etc...
	SearchLayer    string `envconfig:"SEARCHLAYER"`    //Search Layer URI  Example http://localhost:9088
	LocationHelper string `envconfig:"LOCATIONHELPER"` // Location helper service used for autocomplete in city/state searches
	NameHelper     string `envconfig:"NAMEHELPER"`     // Name helper service used for creating sane names to pass into search leayer
}

// GetConfig get the current configuration from the environment
func GetConfig() (AppConfig, error) {
	var cfg AppConfig
	err := envconfig.Process("", &cfg)
	if err != nil {
		fmt.Println(err)
	}
	return cfg, nil
}
