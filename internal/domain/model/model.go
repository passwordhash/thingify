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
	UpdatedAt time.Ticker
	Labels    []struct {
		Name  string
		Color string
	}
	Repository struct {
		Name     string
		FullName string
		HTMLURL  string
	}
	Assignees []struct {
		Login   string
		HTMLURL string
	}
}
