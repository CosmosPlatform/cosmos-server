package token

import (
	"context"
	"cosmos-server/api"
	"cosmos-server/pkg/errors"
	"cosmos-server/pkg/log"
	"cosmos-server/pkg/services/auth"
	"cosmos-server/pkg/services/token"
	"cosmos-server/pkg/services/user"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type handler struct {
	tokenService token.Service
	userService  user.Service
	logger       log.Logger
}

func AddAuthenticatedTokenHandler(e *gin.RouterGroup, tokenService token.Service, userService user.Service, logger log.Logger) {
	h := &handler{
		tokenService: tokenService,
		userService:  userService,
		logger:       logger,
	}

	e.POST("/tokens/:team", h.handlePostToken)
}

func (h *handler) handlePostToken(c *gin.Context) {
	var req api.CreateTokenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorf("Failed to bind JSON for create token request: %v", err)
		_ = c.Error(errors.NewBadRequestError(fmt.Sprintf("Invalid request format: %v", err)))
		return
	}

	err := req.Validate()
	if err != nil {
		_ = c.Error(errors.NewBadRequestError(err.Error()))
		return
	}

	teamName := c.Param("team")
	if teamName == "" {
		_ = c.Error(errors.NewBadRequestError("Team name is required"))
		return
	}

	role, email, err := getRoleAndEmailFromContext(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	if role != user.AdminUserRole {
		if isFromTeam, err := h.userIsFromTeam(c, email, teamName); err != nil {
			_ = c.Error(err)
			return
		} else if !isFromTeam {
			_ = c.Error(errors.NewForbiddenError("regular users cannot create token for other teams"))
			return
		}
	}

	err = h.tokenService.CreateToken(c, teamName, req.Name, req.Value)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.Status(http.StatusCreated)
}

func getRoleAndEmailFromContext(c *gin.Context) (string, string, error) {
	role, exists := c.Get(auth.UserRoleContextKey)
	if !exists {
		return "", "", errors.NewUnauthorizedError("role not found in token")
	}

	email, exists := c.Get(auth.UserEmailContextKey)
	if !exists {
		return "", "", errors.NewUnauthorizedError("email not found in token")
	}

	return role.(string), email.(string), nil
}

func (h *handler) userIsFromTeam(context context.Context, userEmail, team string) (bool, error) {
	user, err := h.userService.GetUserWithEmail(context, userEmail)
	if err != nil {
		h.logger.Errorf("failed to get user by email: %v", err)
		return false, errors.NewInternalServerError(fmt.Sprintf("failed to retrieve user: %v", err))
	}

	if user.Team == nil {
		return false, nil
	}

	return user.Team.Name == team, nil
}
