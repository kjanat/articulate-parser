package services

import (
	"github.com/kjanat/articulate-parser/internal/config"
	"github.com/kjanat/articulate-parser/internal/interfaces"
)

// NewLoggerFromConfig creates a logger based on the configuration.
// Returns a JSON logger if cfg.LogFormat is "json", otherwise returns a text logger.
func NewLoggerFromConfig(cfg *config.Config) interfaces.Logger {
	if cfg.LogFormat == "json" {
		return NewSlogLogger(cfg.LogLevel)
	}
	return NewTextLogger(cfg.LogLevel)
}
