package group

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
	CreateGroup(ctx context.Context, name, description string, members []string) error
	GetAllGroups(ctx context.Context) ([]*model.Group, error)
	GetGroupByName(ctx context.Context, name string) (*model.Group, error)
	DeleteGroup(ctx context.Context, name string) error
	UpdateGroup(ctx context.Context, groupName string, updateData *model.GroupUpdate) error
}

type groupService struct {
	storageService storage.Service
	translator     Translator
	logger         log.Logger
}

func NewGroupService(storageService storage.Service, translator Translator, logger log.Logger) Service {
	return &groupService{
		storageService: storageService,
		translator:     translator,
		logger:         logger,
	}
}

func (s *groupService) CreateGroup(ctx context.Context, name, description string, members []string) error {
	applications := make([]*obj.Application, 0)

	for _, member := range members {
		applicationObj, err := s.storageService.GetApplicationWithName(ctx, member)
		if err != nil {
			if errorUtils.Is(err, storage.ErrNotFound) {
				return errors.NewNotFoundError(fmt.Sprintf("application %s not found", member))
			}
			s.logger.Errorf("Failed to retrieve application %s: %v", member, err)
			return fmt.Errorf("failed to retrieve application %s: %v", member, err)
		}
		applications = append(applications, applicationObj)
	}

	err := s.storageService.CreateGroup(ctx, name, description, applications)
	if err != nil {
		s.logger.Errorf("Failed to create group: %v", err)
		return err
	}
	return nil
}

func (s *groupService) GetAllGroups(ctx context.Context) ([]*model.Group, error) {
	groupsObj, err := s.storageService.GetGroups(ctx)
	if err != nil {
		s.logger.Errorf("Failed to retrieve groups: %v", err)
		return nil, fmt.Errorf("failed to retrieve groups: %v", err)
	}

	return s.translator.ToGroupModels(groupsObj), nil
}

func (s *groupService) GetGroupByName(ctx context.Context, name string) (*model.Group, error) {
	groupObj, err := s.storageService.GetGroupByName(ctx, name)
	if err != nil {
		if errorUtils.Is(err, storage.ErrNotFound) {
			return nil, errors.NewNotFoundError(fmt.Sprintf("group %s not found", name))
		}
		s.logger.Errorf("Failed to retrieve group %s: %v", name, err)
		return nil, fmt.Errorf("failed to retrieve group %s: %v", name, err)
	}

	return s.translator.ToGroupModel(groupObj), nil
}

func (s *groupService) DeleteGroup(ctx context.Context, name string) error {
	err := s.storageService.DeleteGroupByName(ctx, name)
	if err != nil {
		if errorUtils.Is(err, storage.ErrNotFound) {
			return errors.NewNotFoundError(fmt.Sprintf("group %s not found", name))
		}
		s.logger.Errorf("Failed to delete group %s: %v", name, err)
		return fmt.Errorf("failed to delete group %s: %v", name, err)
	}
	return nil
}

func (s *groupService) UpdateGroup(ctx context.Context, groupName string, updateData *model.GroupUpdate) error {
	if updateData == nil {
		return nil
	}

	existingGroup, err := s.storageService.GetGroupByName(ctx, groupName)
	if err != nil {
		if errorUtils.Is(err, storage.ErrNotFound) {
			return errors.NewNotFoundError(fmt.Sprintf("group %s not found", groupName))
		}
		s.logger.Errorf("Failed to retrieve group %s: %v", groupName, err)
		return fmt.Errorf("failed to retrieve group %s: %v", groupName, err)
	}

	if updateData.Name != nil {
		existingGroup.Name = *updateData.Name
	}

	if updateData.Description != nil {
		existingGroup.Description = *updateData.Description
	}

	if updateData.Members != nil {
		applications := make([]*obj.Application, 0)
		for _, member := range updateData.Members {
			applicationObj, err := s.storageService.GetApplicationWithName(ctx, member)
			if err != nil {
				if errorUtils.Is(err, storage.ErrNotFound) {
					return errors.NewNotFoundError(fmt.Sprintf("application %s not found", member))
				}
				s.logger.Errorf("Failed to retrieve application %s: %v", member, err)
				return fmt.Errorf("failed to retrieve application %s: %v", member, err)
			}
			applications = append(applications, applicationObj)
		}
		existingGroup.Applications = applications
	}

	err = s.storageService.UpdateGroup(ctx, existingGroup)
	if err != nil {
		s.logger.Errorf("Failed to update group %s: %v", groupName, err)
		return fmt.Errorf("failed to update group %s: %v", groupName, err)
	}

	return nil
}
