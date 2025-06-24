package model

import (
	"time"
)

type IssueAction struct {
	Issue      IssueInfo
	Repository GHRepository
	Sender     GHUser
}

type IssueInfo struct {
	ID        int64
	Number    int
	Title     string
	Body      *string
	State     string
	HTMLURL   string
	URL       string
	CreatedAt time.Time
	UpdatedAt time.Time
	ClosedAt  *time.Time
	User      GHUser
	Assignees []GHUser
	Labels    []Label
}

type GHRepository struct {
	ID       int64
	Name     string
	FullName string
	Private  bool
	HTMLURL  string
	Owner    GHUser
}

type GHUser struct {
	ID        int64
	Login     string
	AvatarURL string
	HTMLURL   string
	Type      string
}

type Label struct {
	ID          int64
	Name        string
	Color       string
	Description string
}
