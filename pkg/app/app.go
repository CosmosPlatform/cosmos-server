package app

import (
	c "cosmos-server/pkg/config"
	"cosmos-server/pkg/server"
	"net/http"
)

type App struct {
	config c.Config
}

func NewApp(config c.Config) App {
	return App{
		config: config,
	}
}

func (app *App) RunServer() error {
	s := &http.Server{
		Addr:    app.config.Port,
		Handler: server.NewGinHandler(),
	}

	return server.StartServer(s)
}
