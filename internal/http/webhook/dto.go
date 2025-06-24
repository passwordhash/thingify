package webhook

import (
	"errors"
	"time"

	"thingify/internal/domain/model"
)

const timeFormat = time.RFC3339

type issueWebhookReq struct {
	Action     string        `json:"action"`
	Changes    any           `json:"changes"`
	Issue      issueDTO      `json:"issue"`
	Repository repositoryDTO `json:"repository"`
	Sender     userDTO       `json:"sender"`
}

type issueDTO struct {
	ID        int64      `json:"id"`
	Number    int        `json:"number"`
	Title     string     `json:"title"`
	Body      *string    `json:"body"`
	State     string     `json:"state"`
	Locked    bool       `json:"locked"`
	HTMLURL   string     `json:"html_url"`
	URL       string     `json:"url"`
	CreatedAt string     `json:"created_at"`
	UpdatedAt string     `json:"updated_at"`
	ClosedAt  string     `json:"closed_at"`
	User      userDTO    `json:"user"`
	Assignees []userDTO  `json:"assignees"`
	Labels    []labelDTO `json:"labels"`
}

type labelDTO struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Color       string `json:"color"`
	Description string `json:"description"`
}

type repositoryDTO struct {
	ID       int64   `json:"id"`
	NodeID   string  `json:"node_id"`
	Name     string  `json:"name"`
	FullName string  `json:"full_name"`
	Private  bool    `json:"private"`
	HTMLURL  string  `json:"html_url"`
	Owner    userDTO `json:"owner"`
}

type userDTO struct {
	ID        int64  `json:"id"`
	Login     string `json:"login"`
	AvatarURL string `json:"avatar_url"`
	HTMLURL   string `json:"html_url"`
	Type      string `json:"type"`
}

func (d issueWebhookReq) ToDomain() (model.IssueAction, error) {
	return model.IssueAction{
		Issue:      d.Issue.toDomain(),
		Repository: d.Repository.toDomain(),
		Sender:     d.Sender.toDomain(),
	}, nil
}

func (d issueDTO) toDomain() model.IssueInfo {
	assignees := make([]model.GHUser, len(d.Assignees))
	for i, u := range d.Assignees {
		assignees[i] = u.toDomain()
	}
	labels := make([]model.Label, len(d.Labels))
	for i, l := range d.Labels {
		labels[i] = l.toDomain()
	}

	var convertErr error

	createdAt, err := time.Parse(timeFormat, d.CreatedAt)
	if err != nil {
		errors.Join(convertErr, err)
	}
	updatedAt, err := time.Parse(timeFormat, d.UpdatedAt)
	if err != nil {
		errors.Join(convertErr, err)
	}
	var closedAt time.Time
	if d.ClosedAt != "" {
		closedAt, err = time.Parse(timeFormat, d.ClosedAt)
		if err != nil {
			errors.Join(convertErr, err)
		}
	}

	return model.IssueInfo{
		ID:        d.ID,
		Number:    d.Number,
		Title:     d.Title,
		Body:      d.Body,
		State:     d.State,
		HTMLURL:   d.HTMLURL,
		URL:       d.URL,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		ClosedAt:  &closedAt,
		User:      d.User.toDomain(),
		Assignees: assignees,
		Labels:    labels,
	}
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

func (d labelDTO) toDomain() model.Label {
	return model.Label{
		ID:          d.ID,
		Name:        d.Name,
		Color:       d.Color,
		Description: d.Description,
	}
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
