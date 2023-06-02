package searchservice

import (
	"rainier/internal/config"
)

type PhoneService struct {
	hConfig *config.AppConfig
}

// Initialize Phone Service
func (ps *PhoneService) Initialize(config *config.AppConfig) {
	ps.hConfig = config
}
