package app

import (
	"cosmos-server/pkg/auth"
	c "cosmos-server/pkg/config"
	"cosmos-server/pkg/routes"
	"cosmos-server/pkg/server"
	"cosmos-server/pkg/storage"
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

	authService := auth.NewStorageAuthService(storageService)

	httpRoutes := routes.NewHTTPRoutes(authService)

	return &App{
		config: config,
		routes: httpRoutes,
	}, nil
}

func (app *App) RunServer() error {
	s := &http.Server{
		Addr:    app.config.ServerConfig.Port,
		Handler: server.NewGinHandler(app.routes),
	}

	return server.StartServer(s)
}
