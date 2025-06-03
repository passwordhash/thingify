package model

type UserIssuesResponse struct {
	Data struct {
		Viewer struct {
			Issues struct {
				Nodes []Issue `json:"nodes"`
			} `json:"issues"`
		} `json:"viewer"`
	} `json:"data"`
}

type Issue struct {
	ID         string     `json:"id"`
	Number     int        `json:"number"`
	Title      string     `json:"title"`
	Body       string     `json:"body"`
	State      string     `json:"state"`
	URL        string     `json:"url"`
	CreatedAt  string     `json:"createdAt"`
	UpdatedAt  string     `json:"updatedAt"`
	Labels     Labels     `json:"labels"`
	Repository Repository `json:"repository"`
	Assignees  Assignees  `json:"assignees"`
}

type Labels struct {
	Nodes []Label `json:"nodes"`
}

type Label struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

type Assignees struct {
	Login     string `json:"login"`
	AvatarURL string `json:"avatar_url"`
	HTMLURL   string `json:"html_url"`
}

type Repository struct {
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	HTMLURL  string `json:"html_url"`
}
