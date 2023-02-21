package github

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/answerdev/answer/plugin"
	"github.com/answerdev/plugins/connector/github/i18n"
	"github.com/google/go-github/v48/github"
	"golang.org/x/oauth2"
	oauth2GitHub "golang.org/x/oauth2/github"
)

type GitHubConnector struct {
	Config *GitHubConnectorConfig
}

type GitHubConnectorConfig struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

func init() {
	plugin.Register(&GitHubConnector{
		Config: &GitHubConnectorConfig{},
	})
}

func (g *GitHubConnector) Info() plugin.Info {
	return plugin.Info{
		Name:        plugin.MakeTranslator(i18n.InfoName),
		SlugName:    "github_connector",
		Description: plugin.MakeTranslator(i18n.InfoDescription),
		Author:      "answerdev",
		Version:     "0.0.1",
	}
}

func (g *GitHubConnector) ConnectorLogoSVG() string {
	return `PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHdpZHRoPSIyNCIgaGVpZ2h0PSIyNCIgdmlld0JveD0iMCAwIDI0IDI0Ij48cGF0aCBkPSJNMTIgMGMtNi42MjYgMC0xMiA1LjM3My0xMiAxMiAwIDUuMzAyIDMuNDM4IDkuOCA4LjIwNyAxMS4zODcuNTk5LjExMS43OTMtLjI2MS43OTMtLjU3N3YtMi4yMzRjLTMuMzM4LjcyNi00LjAzMy0xLjQxNi00LjAzMy0xLjQxNi0uNTQ2LTEuMzg3LTEuMzMzLTEuNzU2LTEuMzMzLTEuNzU2LTEuMDg5LS43NDUuMDgzLS43MjkuMDgzLS43MjkgMS4yMDUuMDg0IDEuODM5IDEuMjM3IDEuODM5IDEuMjM3IDEuMDcgMS44MzQgMi44MDcgMS4zMDQgMy40OTIuOTk3LjEwNy0uNzc1LjQxOC0xLjMwNS43NjItMS42MDQtMi42NjUtLjMwNS01LjQ2Ny0xLjMzNC01LjQ2Ny01LjkzMSAwLTEuMzExLjQ2OS0yLjM4MSAxLjIzNi0zLjIyMS0uMTI0LS4zMDMtLjUzNS0xLjUyNC4xMTctMy4xNzYgMCAwIDEuMDA4LS4zMjIgMy4zMDEgMS4yMy45NTctLjI2NiAxLjk4My0uMzk5IDMuMDAzLS40MDQgMS4wMi4wMDUgMi4wNDcuMTM4IDMuMDA2LjQwNCAyLjI5MS0xLjU1MiAzLjI5Ny0xLjIzIDMuMjk3LTEuMjMuNjUzIDEuNjUzLjI0MiAyLjg3NC4xMTggMy4xNzYuNzcuODQgMS4yMzUgMS45MTEgMS4yMzUgMy4yMjEgMCA0LjYwOS0yLjgwNyA1LjYyNC01LjQ3OSA1LjkyMS40My4zNzIuODIzIDEuMTAyLjgyMyAyLjIyMnYzLjI5M2MwIC4zMTkuMTkyLjY5NC44MDEuNTc2IDQuNzY1LTEuNTg5IDguMTk5LTYuMDg2IDguMTk5LTExLjM4NiAwLTYuNjI3LTUuMzczLTEyLTEyLTEyeiIvPjwvc3ZnPg==`
}

func (g *GitHubConnector) ConnectorName() plugin.Translator {
	return plugin.MakeTranslator(i18n.ConnectorName)
}

func (g *GitHubConnector) ConnectorSlugName() string {
	return "github"
}

func (g *GitHubConnector) ConnectorSender(ctx *plugin.GinContext, receiverURL string) (redirectURL string) {
	oauth2Config := &oauth2.Config{
		ClientID:     g.Config.ClientID,
		ClientSecret: g.Config.ClientSecret,
		Endpoint:     oauth2GitHub.Endpoint,
		RedirectURL:  receiverURL,
		Scopes:       nil,
	}
	return oauth2Config.AuthCodeURL("state")
}

func (g *GitHubConnector) ConnectorReceiver(ctx *plugin.GinContext) (userInfo plugin.ExternalLoginUserInfo, err error) {
	code := ctx.Query("code")
	// Exchange code for token
	oauth2Config := &oauth2.Config{
		ClientID:     g.Config.ClientID,
		ClientSecret: g.Config.ClientSecret,
		Endpoint:     oauth2GitHub.Endpoint,
	}
	token, err := oauth2Config.Exchange(context.Background(), code)
	if err != nil {
		return userInfo, err
	}

	// Exchange token for user info
	cli := github.NewClient(oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token.AccessToken},
	)))
	resp, _, err := cli.Users.Get(context.Background(), "")
	if err != nil {
		return userInfo, err
	}

	metaInfo, _ := json.Marshal(resp)
	userInfo = plugin.ExternalLoginUserInfo{
		ExternalID: fmt.Sprintf("%d", resp.GetID()),
		Name:       resp.GetName(),
		Email:      resp.GetEmail(),
		MetaInfo:   string(metaInfo),
	}
	return userInfo, nil
}

func (g *GitHubConnector) ConfigFields() []plugin.ConfigField {
	return []plugin.ConfigField{
		{
			Name:        "client_id",
			Type:        plugin.ConfigTypeInput,
			Title:       plugin.MakeTranslator(i18n.ConfigClientIDTitle),
			Description: plugin.MakeTranslator(i18n.ConfigClientIDDescription),
			Required:    true,
			UIOptions: plugin.ConfigFieldUIOptions{
				InputType: plugin.InputTypeText,
			},
			Value: g.Config.ClientID,
		},
		{
			Name:        "client_secret",
			Type:        plugin.ConfigTypeInput,
			Title:       plugin.MakeTranslator(i18n.ConfigClientSecretTitle),
			Description: plugin.MakeTranslator(i18n.ConfigClientSecretDescription),
			Required:    true,
			UIOptions: plugin.ConfigFieldUIOptions{
				InputType: plugin.InputTypeText,
			},
			Value: g.Config.ClientSecret,
		},
	}
}

func (g *GitHubConnector) ConfigReceiver(config []byte) error {
	c := &GitHubConnectorConfig{}
	_ = json.Unmarshal(config, c)
	g.Config = c
	return nil
}
