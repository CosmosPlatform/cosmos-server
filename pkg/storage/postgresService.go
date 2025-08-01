package storage

import (
	"context"
	"cosmos-server/pkg/config"
	"cosmos-server/pkg/log"
	"cosmos-server/pkg/storage/obj"
	errorUtils "errors"
	"fmt"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresService struct {
	db     *gorm.DB
	logger log.Logger
}

func NewPostgresService(config config.StorageConfig, logger log.Logger) (*PostgresService, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.DatabaseName)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
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
	user, err := gorm.G[obj.User](s.db).Where("email = ?", email).First(ctx)

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
	users, err := gorm.G[*obj.User](s.db).Where("username ILIKE ? OR email ILIKE ?", "%"+filter+"%", "%"+filter+"%").Order("username").Find(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get users with filter '%s': %v", filter, err)
	}

	return users, nil
}

func (s *PostgresService) InsertTeam(ctx context.Context, team *obj.Team) error {
	err := gorm.G[obj.Team](s.db).Create(ctx, team)
	if err != nil {
		if errorUtils.Is(err, gorm.ErrDuplicatedKey) {
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
