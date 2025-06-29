package app

import (
	"context"
	"cosmos-server/pkg/auth"
	c "cosmos-server/pkg/config"
	"cosmos-server/pkg/routes"
	"cosmos-server/pkg/server"
	"cosmos-server/pkg/storage"
	"cosmos-server/pkg/user"
	"fmt"
	"net/http"
)

type App struct {
	config *c.Config
	routes *routes.HTTPRoutes
}

func NewApp(config *c.Config) (*App, error) {
	storageService, err := storage.NewMongoService(config.StorageConfig)
	if err != nil {
		return nil, err
	}

	authService := auth.NewAuthService(config.AuthConfig, storageService)
	userService := user.NewUserService(storageService)

	httpRoutes := routes.NewHTTPRoutes(authService, userService)

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

			if err := app.routes.UserService.RegisterAdminUser(context.Background(), adminUsername, adminEmail, adminPassword); err != nil {
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
