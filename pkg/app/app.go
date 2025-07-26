package app

import (
	"context"
	c "cosmos-server/pkg/config"
	"cosmos-server/pkg/log"
	"cosmos-server/pkg/routes"
	"cosmos-server/pkg/server"
	"cosmos-server/pkg/services/auth"
	"cosmos-server/pkg/services/team"
	"cosmos-server/pkg/services/user"
	"cosmos-server/pkg/storage"
	"fmt"
	"net/http"
)

type App struct {
	config *c.Config
	routes *routes.HTTPRoutes
}

func NewApp(config *c.Config) (*App, error) {
	logger, err := log.NewLogger(config.LogConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %v", err)
	}

	storageService, err := storage.NewPostgresService(config.StorageConfig, logger)
	if err != nil {
		return nil, err
	}

	authService := auth.NewAuthService(config.AuthConfig, storageService, auth.NewTranslator(), logger)
	userService := user.NewUserService(storageService, logger)
	teamService := team.NewTeamService(storageService)

	httpRoutes := routes.NewHTTPRoutes(authService, userService, teamService, logger)

	return &App{
		config: config,
		routes: httpRoutes,
	}, nil
}

func (app *App) SetUpDatabase() error {
	if adminPresent, err := app.routes.UserService.AdminUserPresent(context.Background()); err != nil {
		return fmt.Errorf("failed to check for admin user: %v", err)
	} else {
		if !adminPresent {
			adminUsername := app.config.SystemConfig.DefaultAdmin.Username
			adminEmail := app.config.SystemConfig.DefaultAdmin.Email
			adminPassword := app.config.SystemConfig.DefaultAdmin.Password

			if err := app.routes.UserService.RegisterUser(context.Background(), adminUsername, adminEmail, adminPassword, user.AdminUserRole); err != nil {
				return fmt.Errorf("failed to register admin user: %v", err)
			}
		}
	}

	return nil
}

func (app *App) RunServer() error {
	address := app.config.ServerConfig.Host + ":" + app.config.ServerConfig.Port
	s := &http.Server{
		Addr:    address,
		Handler: server.NewGinHandler(app.routes),
	}

	return server.StartServer(s)
}
