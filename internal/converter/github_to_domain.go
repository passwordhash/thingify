package converter

import (
	"errors"
	"fmt"
	"thingify/internal/domain/model"
	ghmodel "thingify/internal/github/model"
	"time"
)

// var convertErr = errors.New("convert error")
var convertTimeErr = errors.New("time conversion error")

func GithubToDomain(nodes []ghmodel.IssueNode) ([]model.Issue, error) {
	if len(nodes) == 0 {
		return nil, nil
	}

	var issues []model.Issue
	for _, node := range nodes {
		_ = "2025-06-02T16:39:26Z"
		createdAt, err := time.Parse(time.RFC3339, node.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", convertTimeErr, err)
		}

		// updatedAt, err := time.Parse(time.RFC3339, node.UpdatedAt)
		// if err != nil {
		// 	return nil, fmt.Errorf("%w: %v", convertTimeErr, err)
		// }

		issue := model.Issue{
			ID:        node.ID,
			Title:     node.Title,
			Body:      node.Body,
			State:     node.State,
			URL:       node.URL,
			CreatedAt: createdAt,
			// UpdatedAt: updatedAt, // TODO: uncomment when needed
			Number: node.Number,
			// Labels:
			// Repository:
			// Assignees:
		}
		issues = append(issues, issue)
	}

	return issues, nil
}
