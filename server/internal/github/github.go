package github

type GHClient struct {
	BaseURL string
	Token   string
}

func Register(baseURL, token string) *GHClient {
	return &GHClient{
		BaseURL: baseURL,
		Token:   token,
	}
}

func (c *GHClient) Issues() {

}
