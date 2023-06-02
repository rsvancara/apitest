package util

import (
	"fmt"

	"rainier/internal/config"

	"github.com/segmentio/ksuid"
)

// Generate a unique identifier
func GenUUID() string {
	id := ksuid.New()
	return id.String()
}

// SiteTemplate loads the correct template directory for the site
func SiteTemplate(path string) (string, error) {

	cfg, err := config.GetConfig()
	if err != nil {
		return "", fmt.Errorf("error loading template directory %s", err)
	}

	return "sites/" + cfg.Site + path, nil

}
