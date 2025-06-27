package dto

import "thingify/internal/domain/model"

const (
	ActionCreated = "created"
	ActionOpened  = "opened"
	ActionDeleted = "deleted"
)

type userDTO struct {
	ID        int64  `json:"id"`
	Login     string `json:"login"`
	AvatarURL string `json:"avatar_url"`
	HTMLURL   string `json:"html_url"`
	Type      string `json:"type"`
}

func (d userDTO) toDomain() model.GHUser {
	return model.GHUser{
		ID:        d.ID,
		Login:     d.Login,
		AvatarURL: d.AvatarURL,
		HTMLURL:   d.HTMLURL,
		Type:      d.Type,
	}
}

type repositoryDTO struct {
	ID       int64   `json:"id"`
	NodeID   string  `json:"node_id,omitempty"`
	Name     string  `json:"name"`
	FullName string  `json:"full_name"`
	Private  bool    `json:"private"`
	HTMLURL  string  `json:"html_url,omitempty"`
	Owner    userDTO `json:"owner,omitempty"`
}

func (d repositoryDTO) toDomain() model.GHRepository {
	return model.GHRepository{
		ID:       d.ID,
		Name:     d.Name,
		FullName: d.FullName,
		Private:  d.Private,
		HTMLURL:  d.HTMLURL,
		Owner:    d.Owner.toDomain(),
	}
}

type labelDTO struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Color       string `json:"color"`
	Description string `json:"description"`
}

func (d labelDTO) toDomain() model.Label {
	return model.Label{
		ID:          d.ID,
		Name:        d.Name,
		Color:       d.Color,
		Description: d.Description,
	}
}
