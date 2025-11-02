package app

import (
	"context"
	c "cosmos-server/pkg/config"
	"cosmos-server/pkg/log"
	"cosmos-server/pkg/model"
	"cosmos-server/pkg/routes"
	"cosmos-server/pkg/sentinel"
	"cosmos-server/pkg/server"
	"cosmos-server/pkg/services/application"
	"cosmos-server/pkg/services/auth"
	"cosmos-server/pkg/services/mail"
	"cosmos-server/pkg/services/monitoring"
	"cosmos-server/pkg/services/team"
	"cosmos-server/pkg/services/token"
	"cosmos-server/pkg/services/user"
	"cosmos-server/pkg/storage"
	"fmt"
	"net/http"
)

type App struct {
	config *c.Config
	routes *routes.HTTPRoutes
	server *http.Server
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

	encryptor, err := token.NewAESEncryptor(config.TokenConfig.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create encryptor: %v", err)
	}

	mailService, err := mail.NewMailService(config.MailConfig, storageService, logger)

	authService := auth.NewAuthService(config.AuthConfig, storageService, auth.NewTranslator(), logger)
	userService := user.NewUserService(storageService, user.NewTranslator(), logger)
	teamService := team.NewTeamService(storageService, team.NewTranslator())
	applicationService := application.NewApplicationService(storageService, application.NewTranslator(), logger)
	monitoringService := monitoring.NewMonitoringService(storageService, monitoring.NewGithubService(), monitoring.NewOpenApiService(), mailService, config.SentinelConfig.MaxIntervalSeconds, config.SentinelConfig.MinIntervalSeconds, encryptor, monitoring.NewTranslator(), logger)

	tokenService := token.NewTokenService(encryptor, storageService, token.NewTranslator(), logger)

	httpRoutes := routes.NewHTTPRoutes(authService, userService, teamService, applicationService, monitoringService, tokenService, logger)

	return &App{
		config: config,
		routes: httpRoutes,
	}, nil
}

func (app *App) SetUpDatabase() error {
	err := app.setUpAdminUser()
	if err != nil {
		return err
	}

	err = app.setUpSentinelSettings()
	if err != nil {
		return err
	}

	return nil
}

func (app *App) setUpAdminUser() error {
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

func (app *App) setUpSentinelSettings() error {
	if sentinelSettingsPresent, err := app.routes.MonitoringService.SentinelSettingsPresent(context.Background()); err != nil {
		return fmt.Errorf("failed to check for sentinel settings: %v", err)
	} else {
		if !sentinelSettingsPresent {
			defaultSentinelInterval := app.config.SentinelConfig.DefaultIntervalSeconds
			defaultSentinelEnabled := app.config.SentinelConfig.DefaultEnabled
			err := app.routes.MonitoringService.InsertSentinelIntervalSetting(context.Background(), defaultSentinelInterval, defaultSentinelEnabled)
			if err != nil {
				return fmt.Errorf("failed to insert default sentinel interval setting: %v", err)
			}
		}
	}

	return nil
}

func (app *App) RunServer() error {
	address := ":" + app.config.ServerConfig.Port
	app.server = &http.Server{
		Addr:    address,
		Handler: server.NewGinHandler(app.routes),
	}

	return server.StartServer(app.server)
}

func (app *App) StartSentinel(ctx context.Context) {
	newSettingsChannel := make(chan model.SentinelSettings, 3) // I think 1 would suffice, but just in case

	sentinel := sentinel.NewSentinel(app.routes.Logger, app.routes.ApplicationService, app.routes.MonitoringService, newSettingsChannel, *app.config.SentinelConfig.SentinelWorkers)
	fallbackSettings := &model.SentinelSettings{
		Interval: app.config.SentinelConfig.DefaultIntervalSeconds,
		Enabled:  app.config.SentinelConfig.DefaultEnabled,
	}
	go sentinel.Start(ctx, fallbackSettings)

	app.routes.MonitoringService.StoreSentinelChannel(newSettingsChannel)
}

func (app *App) Shutdown(ctx context.Context) error {
	if app.server != nil {
		return app.server.Shutdown(ctx)
	}
	return nil
}
