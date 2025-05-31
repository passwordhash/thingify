package github

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"thingify/server/internal/model"

	jsoniter "github.com/json-iterator/go"
)

type GHClient struct {
	log *slog.Logger

	BaseURL string
	Token   string
}

func Register(log *slog.Logger, baseURL, token string) *GHClient {
	return &GHClient{
		log:     log,
		BaseURL: baseURL,
		Token:   token,
	}
}

func (c *GHClient) Issues() ([]model.Issue, error) {
	const op = "github.Issues"

	url := fmt.Sprintf("%s/issues", c.BaseURL)

	log := c.log.With("op", op, "url", url)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Error("error creating request", "err", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	req.Header.Set("Authorization", "Bearer "+c.Token)

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
