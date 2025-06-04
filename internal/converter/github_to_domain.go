package converter

import (
	"thingify/internal/domain/model"
	ghmodel "thingify/internal/github/model"
)

func GithubToDomain(nodes []ghmodel.IssueNode) ([]model.Issue, error) {
	if len(nodes) == 0 {
		return nil, nil
	}

	var issues []model.Issue
	for _, node := range nodes {
		issue := model.Issue{
			ID:    node.ID,
			Title: node.Title,
			Body:  node.Body,
			State: node.State,
			URL:   node.URL,
			// CreatedAt: node.CreatedAt,
			// UpdatedAt: node.UpdatedAt,
			Number: node.Number,
			// Labels:
			// Repository:
			// Assignees:
		}
		issues = append(issues, issue)
	}

	return issues, nil
}
