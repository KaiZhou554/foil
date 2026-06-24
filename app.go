package main

import (
	"context"

	"lets_config/config"
)

// App struct
type App struct {
	ctx    context.Context
	Config *config.Manager
}

// NewApp creates a new App application struct.
// cfg is the config manager (created in main).
func NewApp(cfg *config.Manager) *App {
	return &App{Config: cfg}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods.
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// GetConfig returns the full application configuration.
func (a *App) GetConfig() (*config.Config, error) {
	return a.Config.Get(), nil
}

// SaveConfig persists the full application configuration.
func (a *App) SaveConfig(cfg config.Config) error {
	return a.Config.Update(func(c *config.Config) error {
		*c = cfg
		return nil
	})
}
