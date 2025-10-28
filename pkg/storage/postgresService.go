package storage

import (
	"context"
	"cosmos-server/pkg/config"
	"cosmos-server/pkg/log"
	"cosmos-server/pkg/model"
	"cosmos-server/pkg/storage/obj"
	errorUtils "errors"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresService struct {
	db     *gorm.DB
	logger log.Logger
}

func NewPostgresService(config config.StorageConfig, logger log.Logger) (*PostgresService, error) {
	db, err := gorm.Open(postgres.Open(config.DatabaseURL), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %v", err)
	}

	return &PostgresService{
		db:     db,
		logger: logger,
	}, nil
}

func (s *PostgresService) InsertUser(ctx context.Context, user *obj.User) error {
	err := gorm.G[obj.User](s.db).Create(ctx, user)
	if err != nil {
		if errorUtils.Is(err, gorm.ErrDuplicatedKey) {
			return ErrAlreadyExists
		}
		return fmt.Errorf("failed to insert user: %v", err)
	}

	return nil
}

func (s *PostgresService) GetUserWithEmail(ctx context.Context, email string) (*obj.User, error) {
	user, err := gorm.G[obj.User](s.db).Preload("Team", nil).Where("email = ?", email).First(ctx)

	if err != nil {
		if errorUtils.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get user with email %s: %v", email, err)
	}

	return &user, nil
}

func (s *PostgresService) GetUserWithRole(ctx context.Context, role string) (*obj.User, error) {
	user, err := gorm.G[*obj.User](s.db).Where("role = ?", role).First(ctx)
	if err != nil {
		if errorUtils.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get user with role %s: %v", role, err)
	}

	return user, nil
}

func (s *PostgresService) GetUsersWithFilter(ctx context.Context, filter string) ([]*obj.User, error) {
	users, err := gorm.G[*obj.User](s.db).Preload("Team", nil).Where("username ILIKE ? OR email ILIKE ?", "%"+filter+"%", "%"+filter+"%").Order("username").Find(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get users with filter '%s': %v", filter, err)
	}

	s.logger.Infow("got users with filter", "filter", filter, "users", users)

	return users, nil
}

func (s *PostgresService) InsertTeam(ctx context.Context, team *obj.Team) error {
	err := gorm.G[obj.Team](s.db).Create(ctx, team)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") ||
			strings.Contains(err.Error(), "violates unique constraint") ||
			strings.Contains(err.Error(), "23505") {
			return ErrAlreadyExists
		}
		return fmt.Errorf("failed to insert team: %v", err)
	}

	return nil
}

func (s *PostgresService) GetTeamsWithFilter(ctx context.Context, filter string) ([]*obj.Team, error) {
	teams, err := gorm.G[*obj.Team](s.db).Where("name ILIKE ?", "%"+filter+"%").Order("name").Find(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get teams with filter '%s': %v", filter, err)
	}

	return teams, nil
}

func (s *PostgresService) DeleteTeam(ctx context.Context, name string) error {
	rowsAffected, err := gorm.G[obj.Team](s.db).Where("name = ?", name).Delete(ctx)
	if err != nil {
		return fmt.Errorf("failed to get rows affected for team %s: %v", name, err)
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *PostgresService) DeleteUser(ctx context.Context, email string) error {
	rowsAffected, err := gorm.G[obj.User](s.db).Where("email = ?", email).Delete(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete user with email %s: %v", email, err)
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *PostgresService) AddUserToTeam(ctx context.Context, userEmail, teamName string) error {
	team, err := gorm.G[obj.Team](s.db).Where("name = ?", teamName).First(ctx)
	if err != nil {
		if errorUtils.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return fmt.Errorf("failed to get team with name %s: %v", teamName, err)
	}

	rowsAffected, err := gorm.G[obj.User](s.db).Where("email = ?", userEmail).Update(ctx, "team_id", team.ID)
	if err != nil {
		return fmt.Errorf("failed to add user %s to team %s: %v", userEmail, teamName, err)
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *PostgresService) RemoveUserFromTeam(ctx context.Context, userEmail string) error {
	rowsAffected, err := gorm.G[obj.User](s.db).Where("email = ?", userEmail).Update(ctx, "team_id", nil)
	if err != nil {
		return fmt.Errorf("failed to remove user %s from team: %v", userEmail, err)
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *PostgresService) GetTeamWithName(ctx context.Context, name string) (*obj.Team, error) {
	team, err := gorm.G[*obj.Team](s.db).Where("name = ?", name).First(ctx)
	if err != nil {
		if errorUtils.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get team with name %s: %v", name, err)
	}

	return team, nil
}

func (s *PostgresService) InsertApplication(ctx context.Context, application *obj.Application) error {
	existing, err := gorm.G[*obj.Application](s.db).Where("LOWER(name) = LOWER(?)", application.Name).First(ctx)
	if err == nil && existing != nil {
		return ErrAlreadyExists
	}
	if err != nil && !errorUtils.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("failed to check for existing application: %v", err)
	}

	err = gorm.G[obj.Application](s.db).Create(ctx, application)
	if err != nil {
		if errorUtils.Is(err, gorm.ErrDuplicatedKey) {
			return ErrAlreadyExists
		}
		return fmt.Errorf("failed to insert application: %v", err)
	}

	return nil
}

func (s *PostgresService) GetApplicationWithName(ctx context.Context, name string) (*obj.Application, error) {
	application, err := gorm.G[*obj.Application](s.db).Preload("Team", nil).Preload("Token", nil).Where("LOWER(name) = LOWER(?)", name).First(ctx)
	if err != nil {
		if errorUtils.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get application with name %s: %v", name, err)
	}

	return application, nil
}

func (s *PostgresService) GetApplicationsWithFilter(ctx context.Context, filter string) ([]*obj.Application, error) {
	applications, err := gorm.G[*obj.Application](s.db).Preload("Team", nil).Preload("Token", nil).Where("name ILIKE ?", "%"+filter+"%").Order("name").Find(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get applications with filter '%s': %v", filter, err)
	}

	return applications, nil
}

func (s *PostgresService) DeleteApplicationWithName(ctx context.Context, name string) error {
	rowsAffected, err := gorm.G[obj.Application](s.db).Where("name = ?", name).Delete(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete application with name %s: %v", name, err)
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *PostgresService) GetApplicationsByTeam(ctx context.Context, team string) ([]*obj.Application, error) {
	teamObj, err := s.GetTeamWithName(ctx, team)
	if err != nil {
		return nil, err
	}

	applications, err := gorm.G[*obj.Application](s.db).Preload("Team", nil).Preload("Token", nil).Where("team_id = ?", teamObj.ID).Order("name").Find(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get applications for team %s: %v", team, err)
	}

	return applications, nil
}

func (s *PostgresService) UpdateApplication(ctx context.Context, application *obj.Application) error {
	rowsAffected, err := gorm.G[*obj.Application](s.db).Where("id = ?", application.ID).Select("*").Updates(ctx, application)
	if err != nil {
		return fmt.Errorf("failed to update application: %v", err)
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *PostgresService) GetApplicationDependency(ctx context.Context, consumerID, providerID int) (*obj.ApplicationDependency, error) {
	dependency, err := gorm.G[*obj.ApplicationDependency](s.db).
		Preload("Consumer", nil).
		Preload("Provider", nil).
		Where("consumer_id = ? AND provider_id = ?", consumerID, providerID).
		First(ctx)

	if err != nil {
		if errorUtils.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get dependency relationship: %v", err)
	}

	return dependency, nil
}

func (s *PostgresService) GetApplicationDependenciesWithApplicationInvolved(ctx context.Context, applicationName string) ([]*obj.ApplicationDependency, error) {
	application, err := gorm.G[*obj.Application](s.db).Where("name = ?", applicationName).First(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get application: %v", err)
	}

	dependencies, err := gorm.G[*obj.ApplicationDependency](s.db).
		Preload("Consumer", nil).
		Preload("Consumer.Team", nil).
		Preload("Provider", nil).
		Preload("Provider.Team", nil).
		Where("consumer_id = ? OR provider_id = ?", application.ID, application.ID).
		Find(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get application dependencies for application %s: %v", applicationName, err)
	}

	return dependencies, nil
}

func (s *PostgresService) GetApplicationDependenciesWithFilter(ctx context.Context, filters model.ApplicationDependencyFilter) ([]*obj.ApplicationDependency, error) {
	teams := filters.Teams

	if len(teams) == 0 {
		return s.GetAllApplicationDependencies(ctx)
	}

	var applications []*obj.Application
	for _, team := range teams {
		teamApplications, err := s.GetApplicationsByTeam(ctx, team)
		if err != nil {
			return nil, fmt.Errorf("failed to get applications for team %s: %v", team, err)
		}
		applications = append(applications, teamApplications...)
	}

	if len(applications) == 0 {
		return []*obj.ApplicationDependency{}, nil
	}

	appIDs := make([]int, len(applications))
	for i, app := range applications {
		appIDs[i] = int(app.ID)
	}

	query := gorm.G[*obj.ApplicationDependency](s.db).
		Preload("Consumer", nil).
		Preload("Consumer.Team", nil).
		Preload("Provider", nil).
		Preload("Provider.Team", nil)

	if filters.IncludeNeighbors {
		query = query.Where("consumer_id IN ? OR provider_id IN ?", appIDs, appIDs)
	} else {
		query = query.Where("consumer_id IN ? AND provider_id IN ?", appIDs, appIDs)
	}

	dependencies, err := query.
		Find(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get application dependencies with filter: %v", err)
	}

	return dependencies, nil
}

func (s *PostgresService) GetAllApplicationDependencies(ctx context.Context) ([]*obj.ApplicationDependency, error) {
	dependencies, err := gorm.G[*obj.ApplicationDependency](s.db).
		Preload("Consumer", nil).
		Preload("Consumer.Team", nil).
		Preload("Provider", nil).
		Preload("Provider.Team", nil).
		Find(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all application dependencies: %v", err)
	}

	return dependencies, nil
}

func (s *PostgresService) GetApplicationDependenciesByConsumer(ctx context.Context, consumerName string) ([]*obj.ApplicationDependency, error) {
	consumer, err := gorm.G[*obj.Application](s.db).Where("name = ?", consumerName).First(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get consumer application: %v", err)
	}

	dependencies, err := gorm.G[*obj.ApplicationDependency](s.db).
		Preload("Consumer", nil).
		Preload("Consumer.Team", nil).
		Preload("Provider", nil).
		Preload("Provider.Team", nil).
		Where("consumer_id = ?", consumer.ID).
		Find(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get application dependencies for consumer %s: %v", consumerName, err)
	}

	return dependencies, nil
}

func (s *PostgresService) UpdateApplicationDependencies(ctx context.Context, consumerName string, dependenciesToUpsert map[string]*obj.ApplicationDependency, pendingDependencies map[string]*obj.PendingApplicationDependency, dependenciesToDelete []*obj.ApplicationDependency, applicationDependenciesSHA string) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		consumer, err := gorm.G[*obj.Application](tx).Where("LOWER(name) = LOWER(?)", consumerName).First(ctx)
		if err != nil {
			if errorUtils.Is(err, gorm.ErrRecordNotFound) {
				return ErrNotFound
			}
			return fmt.Errorf("failed to get application: %v", err)
		}

		for providerName, dependency := range dependenciesToUpsert {
			if err := s.upsertApplicationDependencyTx(ctx, tx, consumer, providerName, dependency); err != nil {
				return fmt.Errorf("failed to upsert dependency to provider %s: %v", providerName, err)
			}
		}

		if err := s.replacePendingDependenciesTx(ctx, tx, int(consumer.ID), pendingDependencies); err != nil {
			return fmt.Errorf("failed to replace pending dependencies: %v", err)
		}

		if err := s.DeleteApplicationDependenciesTx(ctx, tx, dependenciesToDelete); err != nil {
			return fmt.Errorf("failed to delete dependencies: %v", err)
		}

		rowsAffected, err := gorm.G[*obj.Application](tx).Where("id = ?", consumer.ID).Update(ctx, "dependencies_sha", applicationDependenciesSHA)
		if err != nil {
			return fmt.Errorf("failed to update ApplicationDependenciesSha: %v", err)
		}
		if rowsAffected == 0 {
			return ErrNotFound
		}

		return nil
	})
}

func (s *PostgresService) upsertApplicationDependencyTx(ctx context.Context, tx *gorm.DB, consumer *obj.Application, providerName string, dependency *obj.ApplicationDependency) error {
	provider, err := gorm.G[*obj.Application](tx).Where("name = ?", providerName).First(ctx)
	if err != nil {
		return fmt.Errorf("failed to get provider application %s: %v", providerName, err)
	}

	dependency.ConsumerID = int(consumer.ID)
	dependency.ProviderID = int(provider.ID)

	existing, err := s.GetApplicationDependency(ctx, int(consumer.ID), int(provider.ID))
	if err != nil {
		if errorUtils.Is(err, ErrNotFound) {
			return s.insertApplicationDependencyTx(ctx, tx, dependency)
		}
		return fmt.Errorf("failed to check existing dependency: %v", err)
	}

	dependency.ID = existing.ID
	dependency.CreatedAt = existing.CreatedAt
	rowsAffected, err := gorm.G[*obj.ApplicationDependency](tx).Where("id = ?", existing.ID).Updates(ctx, dependency)
	if err != nil {
		return fmt.Errorf("failed to update application dependency: %v", err)
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *PostgresService) replacePendingDependenciesTx(ctx context.Context, tx *gorm.DB, consumerID int, pendingDependencies map[string]*obj.PendingApplicationDependency) error {
	_, err := gorm.G[obj.PendingApplicationDependency](tx).Where("consumer_id = ?", consumerID).Delete(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete existing pending dependencies: %v", err)
	}

	for _, pendingDependency := range pendingDependencies {
		pendingDependency.ConsumerID = consumerID
		err := gorm.G[obj.PendingApplicationDependency](tx).Create(ctx, pendingDependency)
		if err != nil {
			if errorUtils.Is(err, gorm.ErrDuplicatedKey) {
				return ErrAlreadyExists
			}
			return fmt.Errorf("failed to insert pending application dependency: %v", err)
		}
	}

	return nil
}

func (s *PostgresService) insertApplicationDependencyTx(ctx context.Context, tx *gorm.DB, dependency *obj.ApplicationDependency) error {
	err := gorm.G[obj.ApplicationDependency](tx).Create(ctx, dependency)
	if err != nil {
		if errorUtils.Is(err, gorm.ErrDuplicatedKey) {
			return ErrAlreadyExists
		}
		return fmt.Errorf("failed to insert application dependency: %v", err)
	}

	return nil
}

func (s *PostgresService) DeleteApplicationDependenciesTx(ctx context.Context, tx *gorm.DB, dependenciesToDelete []*obj.ApplicationDependency) error {
	deletedDependenciesIDs := make([]int, 0, len(dependenciesToDelete))
	for _, dependency := range dependenciesToDelete {
		deletedDependenciesIDs = append(deletedDependenciesIDs, int(dependency.ID))
	}

	if len(deletedDependenciesIDs) == 0 {
		return nil
	}

	rowsAffected, err := gorm.G[obj.ApplicationDependency](tx).Where("id IN ?", deletedDependenciesIDs).Delete(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete application dependencies: %v", err)
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *PostgresService) UpsertOpenAPISpecification(ctx context.Context, applicationName string, openAPISpec *obj.ApplicationOpenAPI, applicationOpenApiSHA string) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		application, err := gorm.G[*obj.Application](tx).Preload("Team", nil).Where("LOWER(name) = LOWER(?)", applicationName).First(ctx)
		if err != nil {
			return fmt.Errorf("failed to get application: %v", err)
		}

		openAPISpec.ApplicationID = int(application.ID)
		existing, err := gorm.G[*obj.ApplicationOpenAPI](tx).Where("application_id = ?", application.ID).First(ctx)
		if err != nil {
			if errorUtils.Is(err, gorm.ErrRecordNotFound) {
				if err := gorm.G[obj.ApplicationOpenAPI](tx).Create(ctx, openAPISpec); err != nil {
					return fmt.Errorf("failed to insert OpenAPI specification: %v", err)
				}
			} else {
				return fmt.Errorf("failed to check existing OpenAPI spec: %v", err)
			}
		} else {
			openAPISpec.ID = existing.ID
			openAPISpec.CreatedAt = existing.CreatedAt
			rowsAffected, err := gorm.G[*obj.ApplicationOpenAPI](tx).Where("id = ?", existing.ID).Updates(ctx, openAPISpec)
			if err != nil {
				return fmt.Errorf("failed to update OpenAPI specification: %v", err)
			}
			if rowsAffected == 0 {
				return ErrNotFound
			}
		}

		// Update OpenAPISha field atomically
		rowsAffected, err := gorm.G[*obj.Application](tx).Where("id = ?", application.ID).Update(ctx, "open_api_sha", applicationOpenApiSHA)
		if err != nil {
			return fmt.Errorf("failed to update OpenAPISha: %v", err)
		}
		if rowsAffected == 0 {
			return ErrNotFound
		}

		return nil
	})
}

func (s *PostgresService) CheckPendingDependenciesForApplication(ctx context.Context, applicationName string) error {
	pendingDependencies, err := gorm.G[*obj.PendingApplicationDependency](s.db).Preload("Consumer", nil).Where("provider_name = ?", applicationName).Find(ctx)
	if err != nil {
		return fmt.Errorf("failed to get pending dependencies for application %s: %v", applicationName, err)
	}

	if len(pendingDependencies) == 0 {
		return nil
	}

	application, err := gorm.G[*obj.Application](s.db).Where("name = ?", applicationName).First(ctx)
	if err != nil {
		if errorUtils.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return fmt.Errorf("failed to get application: %v", err)
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, pendingDependency := range pendingDependencies {
			dependency := &obj.ApplicationDependency{
				ConsumerID: pendingDependency.ConsumerID,
				ProviderID: int(application.ID),
				Reasons:    pendingDependency.Reasons,
				Endpoints:  pendingDependency.Endpoints,
			}

			err := s.upsertApplicationDependencyTx(ctx, tx, pendingDependency.Consumer, application.Name, dependency)
			if err != nil {
				return fmt.Errorf("failed to upsert dependency from consumer %s to provider %s: %v", pendingDependency.Consumer.Name, applicationName, err)
			}
		}

		_, err := gorm.G[obj.PendingApplicationDependency](tx).Where("provider_name = ?", applicationName).Delete(ctx)
		if err != nil {
			return fmt.Errorf("failed to delete processed pending dependencies: %v", err)
		}

		return nil
	})
}

func (s *PostgresService) GetOpenAPISpecificationByApplicationName(ctx context.Context, applicationName string) (*obj.ApplicationOpenAPI, error) {
	application, err := gorm.G[*obj.Application](s.db).Where("LOWER(name) = LOWER(?)", applicationName).First(ctx)
	if err != nil {
		if errorUtils.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get application: %v", err)
	}

	openAPISpec, err := gorm.G[*obj.ApplicationOpenAPI](s.db).Preload("Application", nil).Where("application_id = ?", application.ID).First(ctx)
	if err != nil {
		if errorUtils.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get OpenAPI specification for application %s: %v", applicationName, err)
	}

	return openAPISpec, nil
}

func (s *PostgresService) GetSentinelSetting(ctx context.Context, name string) (*obj.SentinelSetting, error) {
	setting, err := gorm.G[*obj.SentinelSetting](s.db).Where("name = ?", name).First(ctx)
	if err != nil {
		if errorUtils.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get sentinel setting: %v", err)
	}

	return setting, nil
}

func (s *PostgresService) InsertSentinelSetting(ctx context.Context, setting *obj.SentinelSetting) error {
	err := gorm.G[obj.SentinelSetting](s.db).Create(ctx, setting)
	if err != nil {
		if errorUtils.Is(err, gorm.ErrDuplicatedKey) {
			return ErrAlreadyExists
		}
		return fmt.Errorf("failed to insert sentinel setting: %v", err)
	}

	return nil
}

func (s *PostgresService) UpdateSentinelSetting(ctx context.Context, setting *obj.SentinelSetting) error {
	rowsAffected, err := gorm.G[*obj.SentinelSetting](s.db).Where("id = ?", setting.ID).Select("*").Updates(ctx, setting)
	if err != nil {
		return fmt.Errorf("failed to update sentinel setting: %v", err)
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *PostgresService) GetApplicationsToMonitor(ctx context.Context) ([]*obj.Application, error) {
	applications, err := gorm.G[*obj.Application](s.db).Where("has_open_api = ? OR has_open_client = ?", true, true).Find(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get applications to monitor: %v", err)
	}

	return applications, nil
}

func (s *PostgresService) InsertToken(ctx context.Context, token *obj.Token) error {
	err := gorm.G[obj.Token](s.db).Create(ctx, token)
	if err != nil {
		if errorUtils.Is(err, gorm.ErrDuplicatedKey) {
			return ErrAlreadyExists
		}
		return fmt.Errorf("failed to insert token: %v", err)
	}

	return nil
}

func (s *PostgresService) GetTokensFromTeam(ctx context.Context, teamName string) ([]*obj.Token, error) {
	team, err := s.GetTeamWithName(ctx, teamName)
	if err != nil {
		return nil, err
	}

	tokens, err := gorm.G[*obj.Token](s.db).Preload("Team", nil).Where("team_id = ?", team.ID).Find(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get tokens for team %s: %v", teamName, err)
	}

	return tokens, nil
}

func (s *PostgresService) GetTokenWithNameAndTeamID(ctx context.Context, name string, teamID int) (*obj.Token, error) {
	token, err := gorm.G[*obj.Token](s.db).Preload("Team", nil).Where("name = ? AND team_id = ?", name, teamID).First(ctx)
	if err != nil {
		if errorUtils.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get token: %v", err)
	}
	return token, nil
}

func (s *PostgresService) DeleteToken(ctx context.Context, name string, team string) error {
	teamObj, err := s.GetTeamWithName(ctx, team)
	if err != nil {
		if errorUtils.Is(err, ErrNotFound) {
			return ErrNotFound
		}
		return err
	}

	rowsAffected, err := gorm.G[obj.Token](s.db).Where("name = ? AND team_id = ?", name, teamObj.ID).Delete(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete token %s for team %s: %v", name, team, err)
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}
