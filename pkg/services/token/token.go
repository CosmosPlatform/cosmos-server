package token

import (
	"context"
	"cosmos-server/pkg/config"
	"cosmos-server/pkg/errors"
	"cosmos-server/pkg/log"
	"cosmos-server/pkg/storage"
	"cosmos-server/pkg/storage/obj"
	errorUtils "errors"
	"fmt"
)

type Service interface {
	CreateToken(ctx context.Context, teamName, name, secret string) error
}

type tokenService struct {
	storageService storage.Service
	encriptor      Encryptor
	logger         log.Logger
}

func NewTokenService(conf config.TokenConfig, storageService storage.Service, logger log.Logger) (Service, error) {
	encryptor, err := NewAESEncryptor(conf.EncryptionKey)
	if err != nil {
		return nil, err
	}

	return &tokenService{
		storageService: storageService,
		encriptor:      encryptor,
		logger:         logger,
	}, nil
}

func (s *tokenService) CreateToken(ctx context.Context, teamName, name, secret string) error {
	teamObj, err := s.storageService.GetTeamWithName(ctx, teamName)
	if err != nil {
		if errorUtils.Is(err, storage.ErrNotFound) {
			return errors.NewNotFoundError(fmt.Sprintf("Team %s not found", teamName))
		}
		return errors.NewInternalServerError("failed to retrieve team: " + err.Error())
	}

	encryptedSecret, err := s.encriptor.Encrypt(secret)
	if err != nil {
		return errors.NewInternalServerError("failed to encrypt token secret")
	}

	tokenObj := &obj.Token{
		Name:           name,
		EncryptedValue: encryptedSecret,
		TeamID:         int(teamObj.ID),
	}

	err = s.storageService.InsertToken(ctx, tokenObj)
	if err != nil {
		return errors.NewInternalServerError("failed to create token: " + err.Error())
	}

	return nil
}
