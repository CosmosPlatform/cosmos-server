package token

import (
	"context"
	"cosmos-server/pkg/errors"
	"cosmos-server/pkg/log"
	"cosmos-server/pkg/model"
	"cosmos-server/pkg/storage"
	"cosmos-server/pkg/storage/obj"
	errorUtils "errors"
	"fmt"
)

type Service interface {
	CreateToken(ctx context.Context, teamName, name, secret string) error
	GetTokensFromTeam(ctx context.Context, teamName string) ([]*model.Token, error)
	DeleteToken(ctx context.Context, teamName, name string) error
	UpdateToken(ctx context.Context, teamName string, tokenName string, updateTokenModel *model.TokenUpdate) error
}

type tokenService struct {
	storageService storage.Service
	encriptor      Encryptor
	translator     Translator
	logger         log.Logger
}

func NewTokenService(encryptor Encryptor, storageService storage.Service, translator Translator, logger log.Logger) Service {
	return &tokenService{
		storageService: storageService,
		encriptor:      encryptor,
		translator:     translator,
		logger:         logger,
	}
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

func (s *tokenService) GetTokensFromTeam(ctx context.Context, teamName string) ([]*model.Token, error) {
	tokens, err := s.storageService.GetTokensFromTeam(ctx, teamName)
	if err != nil {
		return nil, errors.NewInternalServerError("failed to retrieve tokens: " + err.Error())
	}

	return s.translator.ToModelTokens(tokens), nil
}

func (s *tokenService) DeleteToken(ctx context.Context, teamName, name string) error {
	err := s.storageService.DeleteToken(ctx, name, teamName)
	if err != nil {
		if errorUtils.Is(err, storage.ErrNotFound) {
			return errors.NewNotFoundError("token not found")
		}
		return errors.NewInternalServerError("failed to delete token: " + err.Error())
	}

	return nil
}

func (s *tokenService) UpdateToken(ctx context.Context, teamName string, tokenName string, updateTokenModel *model.TokenUpdate) error {
	teamObj, err := s.storageService.GetTeamWithName(ctx, teamName)
	if err != nil {
		if errorUtils.Is(err, storage.ErrNotFound) {
			return errors.NewNotFoundError(fmt.Sprintf("Team %s not found", teamName))
		}
		return errors.NewInternalServerError("failed to retrieve team: " + err.Error())
	}

	tokenToUpdate, err := s.storageService.GetTokenWithNameAndTeamID(ctx, tokenName, int(teamObj.ID))
	if err != nil {
		if errorUtils.Is(err, storage.ErrNotFound) {
			return errors.NewNotFoundError("token not found")
		}
		return errors.NewInternalServerError("failed to retrieve token: " + err.Error())
	}

	if updateTokenModel.Name != nil {
		tokenToUpdate.Name = *updateTokenModel.Name
	}

	if updateTokenModel.Value != nil {
		encryptedSecret, err := s.encriptor.Encrypt(*updateTokenModel.Value)
		if err != nil {
			return errors.NewInternalServerError("failed to encrypt token secret")
		}
		tokenToUpdate.EncryptedValue = encryptedSecret
	}

	err = s.storageService.UpdateToken(ctx, tokenToUpdate)
	if err != nil {
		return errors.NewInternalServerError("failed to update token: " + err.Error())
	}

	return nil
}
