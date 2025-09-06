package storage

import (
	"context"
	"cosmos-server/pkg/config"
	"cosmos-server/pkg/log"
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
	application, err := gorm.G[*obj.Application](s.db).Preload("Team", nil).Where("LOWER(name) = LOWER(?)", name).First(ctx)
	if err != nil {
		if errorUtils.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get application with name %s: %v", name, err)
	}

	return application, nil
}

func (s *PostgresService) GetApplicationsWithFilter(ctx context.Context, filter string) ([]*obj.Application, error) {
	applications, err := gorm.G[*obj.Application](s.db).Preload("Team", nil).Where("name ILIKE ?", "%"+filter+"%").Order("name").Find(ctx)
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

	applications, err := gorm.G[*obj.Application](s.db).Preload("Team", nil).Where("team_id = ?", teamObj.ID).Order("name").Find(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get applications for team %s: %v", team, err)
	}

	return applications, nil
}

func (s *PostgresService) UpdateApplication(ctx context.Context, application *obj.Application) error {
	rowsAffected, err := gorm.G[*obj.Application](s.db).Where("id = ?", application.ID).Updates(ctx, application)
	if err != nil {
		return fmt.Errorf("failed to update application: %v", err)
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *PostgresService) InsertApplicationDependency(ctx context.Context, dependency *obj.ApplicationDependency) error {
	err := gorm.G[obj.ApplicationDependency](s.db).Create(ctx, dependency)
	if err != nil {
		if errorUtils.Is(err, gorm.ErrDuplicatedKey) {
			return ErrAlreadyExists
		}
		return fmt.Errorf("failed to insert application dependency: %v", err)
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

func (s *PostgresService) UpsertApplicationDependency(ctx context.Context, consumerName, providerName string, dependency *obj.ApplicationDependency) error {
	consumer, err := gorm.G[*obj.Application](s.db).Where("name = ?", consumerName).First(ctx)
	if err != nil {
		return fmt.Errorf("failed to get consumer application: %v", err)
	}

	provider, err := gorm.G[*obj.Application](s.db).Where("name = ?", providerName).First(ctx)
	if err != nil {
		return fmt.Errorf("failed to get provider application: %v", err)
	}

	dependency.ConsumerID = int(consumer.ID)
	dependency.ProviderID = int(provider.ID)

	existing, err := s.GetApplicationDependency(ctx, int(consumer.ID), int(provider.ID))
	if err != nil && !errorUtils.Is(err, ErrNotFound) {
		return fmt.Errorf("failed to check existing dependency: %v", err)
	}

	if existing != nil {
		dependency.ID = existing.ID
		dependency.CreatedAt = existing.CreatedAt
		rowsAffected, err := gorm.G[*obj.ApplicationDependency](s.db).Where("id = ?", existing.ID).Updates(ctx, dependency)
		if err != nil {
			return fmt.Errorf("failed to update application dependency: %v", err)
		}
		if rowsAffected == 0 {
			return ErrNotFound
		}
		return nil
	}

	return s.InsertApplicationDependency(ctx, dependency)
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

func (s *PostgresService) GetApplicationDependenciesByProvider(ctx context.Context, providerName string) ([]*obj.ApplicationDependency, error) {
	provider, err := gorm.G[*obj.Application](s.db).Where("name = ?", providerName).First(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get provider application: %v", err)
	}

	dependencies, err := gorm.G[*obj.ApplicationDependency](s.db).
		Preload("Consumer", nil).
		Preload("Consumer.Team", nil).
		Preload("Provider", nil).
		Preload("Provider.Team", nil).
		Where("provider_id = ?", provider.ID).
		Find(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get application dependencies for provider %s: %v", providerName, err)
	}

	return dependencies, nil
}
