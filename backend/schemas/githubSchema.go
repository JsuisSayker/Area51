package schemas

type GithubAction string

const (
	GithubPullRequest GithubAction = "pull_request"
)

type GithubReaction string

const (
	GithubReactionCreateNewRelease GithubReaction = "create_new_release"
)

type GitHubResponseToken struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

type GithubUserInfo struct {
	Login     string `json:"login"`
	Id        uint64 `json:"id"         gorm:"primaryKey"`
	AvatarUrl string `json:"avatar_url"`
	Type      string `json:"type"`
	HtmlUrl   string `json:"html_url"`
	Name      string `json:"name"`
	Email     string `json:"email"`
}
