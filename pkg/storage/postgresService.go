package storage

import (
	"context"
	"cosmos-server/pkg/config"
	"cosmos-server/pkg/log"
	"cosmos-server/pkg/storage/obj"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type PostgresService struct {
	db     *sqlx.DB
	logger log.Logger
}

func NewPostgresService(config config.StorageConfig, logger log.Logger) (*PostgresService, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.DatabaseName)

	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping PostgreSQL: %v", err)
	}

	return &PostgresService{
		db:     db,
		logger: logger,
	}, nil
}

func (s *PostgresService) InsertUser(ctx context.Context, user *obj.User) error {
	query := `INSERT INTO users (username, email, encrypted_password, role) 
             VALUES ($1, $2, $3, $4)`

	_, err := s.db.ExecContext(ctx, query, user.Username, user.Email, user.EncryptedPassword, user.Role)
	if err != nil {
		if err.Error() == "pq: duplicate key value violates unique constraint" {
			return ErrAlreadyExists
		}
		return fmt.Errorf("failed to insert user: %v", err)
	}

	return nil
}

func (s *PostgresService) GetUserWithEmail(ctx context.Context, email string) (*obj.User, error) {
	var user obj.User
	query := `SELECT username, email, encrypted_password, role FROM users WHERE email = $1`

	err := s.db.GetContext(ctx, &user, query, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get user with email %s: %v", email, err)
	}

	return &user, nil
}

func (s *PostgresService) GetUserWithRole(ctx context.Context, role string) (*obj.User, error) {
	var user obj.User
	query := `SELECT username, email, encrypted_password, role FROM users WHERE role = $1 LIMIT 1`

	err := s.db.GetContext(ctx, &user, query, role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user with role %s: %v", role, err)
	}

	return &user, nil
}
