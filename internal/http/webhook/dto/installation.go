package dto



type InstallationWebhookReq struct {
	Action       string          `json:"action"`
	Installation installationDTO `json:"installation"`
	Sender       userDTO         `json:"sender"`
	Repositories []repositoryDTO `json:"repositories,omitempty"`
}

type installationDTO struct {
	ID         int64   `json:"id"`
	AppID      int64   `json:"app_id"`
	TargetID   int64   `json:"target_id"`
	TargetType string  `json:"target_type"` // "User" или "Organization"
	Account    userDTO `json:"account"`
	CreatedAt  string  `json:"created_at,omitempty"`
	UpdatedAt  string  `json:"updated_at,omitempty"`
}
