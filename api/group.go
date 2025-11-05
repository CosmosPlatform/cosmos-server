package api

import validation "github.com/go-ozzo/ozzo-validation/v4"

type CreateGroupRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Members     []string `json:"members"`
}

func (r *CreateGroupRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Name, validation.Required, validation.Length(1, 100)),
		validation.Field(&r.Description, validation.Length(0, 500)),
	)
}

type GetGroupsResponse struct {
	Groups []*GroupReduced `json:"groups"`
}

type GroupReduced struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Group struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Members     []*Application `json:"members"`
}

type GetGroupResponse struct {
	Group *Group `json:"group"`
}

type UpdateGroupRequest struct {
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	Members     []string `json:"members"`
}
