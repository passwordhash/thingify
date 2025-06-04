package github

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"thingify/internal/converter"
	"thingify/internal/domain/model"
	ghmodel "thingify/internal/github/model"

	jsoniter "github.com/json-iterator/go"
)

const userIssuesQueryFile = "user_issues.graphql"

// TODO: нужен ли здесь логгер?
type GHClient struct {
	log *slog.Logger

	BaseURL       string
	GHQueriesPath string
}

func Register(
	log *slog.Logger,
	baseURL string,
	ghQueriesPath string,
) *GHClient {
	return &GHClient{
		log:           log,
		BaseURL:       baseURL,
		GHQueriesPath: ghQueriesPath,
	}
}

func (c *GHClient) UserIssues(_ context.Context, userToken string) ([]model.Issue, error) {
	const op = "github.UserIssues"

	const num = 1 // TMP: количество записей

	log := c.log.With("op", op)

	queryBytes, err := os.ReadFile(fmt.Sprintf("%s/%s", c.GHQueriesPath, userIssuesQueryFile))
	if err != nil {
		log.Error("error reading queries file", "err", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	query := string(queryBytes)

	payload := map[string]any{
		"query": query,
		"variables": map[string]any{
			"num": num,
		},
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Error("error marshalling payload", "err", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	req, err := http.NewRequest("POST", c.BaseURL+"/graphql", bytes.NewReader(payloadBytes))
	if err != nil {
		log.Error("error creating request", "err", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	req.Header.Set("Authorization", "Bearer "+userToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error("error perfoming request", "err", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("error reading response", "err", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer resp.Body.Close()

	var issuesResp ghmodel.UserIssuesResponse
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	if err := json.Unmarshal(body, &issuesResp); err != nil {
		log.Error("error unmarshalling response", "err", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	issues, err := converter.GithubToDomain(issuesResp.Data.Viewer.Issues.Nodes)
	if err != nil {
		log.Error("error converting GitHub issues to domain model", "err", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return issues, nil
}

func (c *GHClient) Issues(_ context.Context, userToken string) ([]model.Issue, error) {
	const op = "github.Issues"

	url := fmt.Sprintf("%s/issues", c.BaseURL)

	log := c.log.With("op", op, "url", url)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Error("error creating request", "err", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	req.Header.Set("Authorization", "Bearer "+userToken)

	resp, err := client.Do(req)
	if err != nil {
		log.Error("error performing request", "err", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("error reading response", "err", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer resp.Body.Close()

	var issues []model.Issue
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	if err := json.Unmarshal(body, &issues); err != nil {
		log.Error("error unmarshalling response", "err", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	//marshalToFile("out.json", issues[0])

	return issues, nil
}

func marshalToFile(filename string, v interface{}) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return jsoniter.NewEncoder(f).Encode(v)
}
