package model

import "time"

type Issue struct {
	ID        string
	Number    int
	Title     string
	Body      string
	State     string
	URL       string
	CreatedAt time.Time
	UpdatedAt time.Time
	Labels    []Label
	Asignees  []Assignee
	Repository
}

type Label struct {
	Name  string
	Color string
}

type Repository struct {
	Name     string
	FullName string
	HTMLURL  string
}

type Assignee struct {
	Login   string
	HTMLURL string
}
