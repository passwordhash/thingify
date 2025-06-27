package dto

import (
	"errors"
	"time"

	"thingify/internal/domain/model"
)

const timeFormat = time.RFC3339

type IssueWebhookReq struct {
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

func (d IssueWebhookReq) ToDomain() (model.IssueAction, error) {
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
