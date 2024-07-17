/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package dingtalk

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"github.com/apache/incubator-answer-plugins/util"
	"net/http"

	"github.com/apache/incubator-answer-plugins/connector-dingtalk/i18n"
	"github.com/apache/incubator-answer/plugin"
	"github.com/segmentfault/pacman/log"
)

//go:embed  info.yaml
var Info embed.FS

const (
	LogoSVG      = "PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHdpZHRoPSIzMiIgaGVpZ2h0PSIzMiIgdmlld0JveD0iMCAwIDEwMjQgMTAyNCI+PHBhdGggZmlsbD0iIzAwODlmZiIgZD0iTTUxMiA2NEMyNjQuNiA2NCA2NCAyNjQuNiA2NCA1MTJzMjAwLjYgNDQ4IDQ0OCA0NDhzNDQ4LTIwMC42IDQ0OC00NDhTNzU5LjQgNjQgNTEyIDY0bTIyNyAzODUuM2MtMSA0LjItMy41IDEwLjQtNyAxNy44aC4xbC0uNC43Yy0yMC4zIDQzLjEtNzMuMSAxMjcuNy03My4xIDEyNy43cy0uMS0uMi0uMy0uNWwtMTUuNSAyNi44aDc0LjVMNTc1LjEgODEwbDMyLjMtMTI4aC01OC42bDIwLjQtODQuN2MtMTYuNSAzLjktMzUuOSA5LjQtNTkgMTYuOGMwIDAtMzEuMiAxOC4yLTg5LjktMzVjMCAwLTM5LjYtMzQuNy0xNi42LTQzLjRjOS44LTMuNyA0Ny40LTguNCA3Ny0xMi4zYzQwLTUuNCA2NC42LTguMiA2NC42LTguMlM0MjIgNTE3IDM5Mi43IDUxMi41Yy0yOS4zLTQuNi02Ni40LTUzLjEtNzQuMy05NS44YzAgMC0xMi4yLTIzLjQgMjYuMy0xMi4zYzM4LjUgMTEuMSAxOTcuOSA0My4yIDE5Ny45IDQzLjJzLTIwNy40LTYzLjMtMjIxLjItNzguN2MtMTMuOC0xNS40LTQwLjYtODQuMi0zNy4xLTEyNi41YzAgMCAxLjUtMTAuNSAxMi40LTcuN2MwIDAgMTUzLjMgNjkuNyAyNTguMSAxMDcuOWMxMDQuOCAzNy45IDE5NS45IDU3LjMgMTg0LjIgMTA2LjciLz48L3N2Zz4="
	AuthorizeURL = "https://login.dingtalk.com/oauth2/auth"
	TokenURL     = "https://api.dingtalk.com/v1.0/oauth2/userAccessToken"
	UserJsonURL  = "https://api.dingtalk.com/v1.0/contact/users/me"
)

type Connector struct {
	Config *ConnectorConfig
}

type ConnectorConfig struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type TokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int    `json:"expiresIn"`
	CorpID       string `json:"corpId"`
}

type UserInfoResponse struct {
	Nick      string `json:"nick"`
	AvatarUrl string `json:"avatarUrl"`
	Mobile    string `json:"mobile"`
	OpenID    string `json:"openId"`
	UnionId   string `json:"unionId"`
	Email     string `json:"email"`
	StateCode string `json:"stateCode"`
}

func init() {
	plugin.Register(&Connector{
		Config: &ConnectorConfig{},
	})
}

func (g *Connector) Info() plugin.Info {
	info := &util.Info{}
	info.GetInfo(Info)

	return plugin.Info{
		Name:        plugin.MakeTranslator(i18n.InfoName),
		SlugName:    info.SlugName,
		Description: plugin.MakeTranslator(i18n.InfoDescription),
		Author:      info.Author,
		Version:     info.Version,
		Link:        info.Link,
	}
}

func (g *Connector) ConnectorLogoSVG() string {
	return LogoSVG
}

func (g *Connector) ConnectorName() plugin.Translator {
	return plugin.MakeTranslator(i18n.ConnectorName)
}

func (g *Connector) ConnectorSlugName() string {
	return "dingtalk"
}

func (g *Connector) ConnectorSender(ctx *plugin.GinContext, receiverURL string) (redirectURL string) {
	return fmt.Sprintf("%s?redirect_uri=%s&response_type=code&client_id=%s&scope=Contact.User.Read&state=state&prompt=consent",
		AuthorizeURL, receiverURL, g.Config.ClientID)
}

func (g *Connector) ConnectorReceiver(ctx *plugin.GinContext, receiverURL string) (userInfo plugin.ExternalLoginUserInfo, err error) {

	// 1. get code
	code := ctx.Query("code")
	log.Debugf("code: %s", code)

	// 2. get token
	tokenReq := map[string]string{
		"clientId":     g.Config.ClientID,
		"clientSecret": g.Config.ClientSecret,
		"code":         code,
		"grantType":    "authorization_code",
	}
	token, err := getToken(TokenURL, tokenReq)
	if err != nil {
		log.Errorf("fail to get token : %s", err)
		return plugin.ExternalLoginUserInfo{}, err
	}

	// 3. get user info
	user, err := getUserInfo(UserJsonURL, token)
	if err != nil {
		log.Errorf("fail to get user info : %s", err)
		return plugin.ExternalLoginUserInfo{}, err
	}
	return user, nil
}

func getToken(url string, body map[string]string) (token string, err error) {
	jsonBody, _ := json.Marshal(body)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("get token failed, status code: %d", response.StatusCode)
	}

	var resp TokenResponse
	err = json.NewDecoder(response.Body).Decode(&resp)
	if err != nil {
		return "", err
	}

	return resp.AccessToken, nil
}

func getUserInfo(url string, token string) (userInfo plugin.ExternalLoginUserInfo, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return plugin.ExternalLoginUserInfo{}, err
	}

	req.Header.Set("x-acs-dingtalk-access-token", token)
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return plugin.ExternalLoginUserInfo{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return plugin.ExternalLoginUserInfo{}, fmt.Errorf("get user info failed, status code: %d", response.StatusCode)
	}

	var resp UserInfoResponse
	err = json.NewDecoder(response.Body).Decode(&resp)
	if err != nil {
		return plugin.ExternalLoginUserInfo{}, err
	}

	userInfo = plugin.ExternalLoginUserInfo{
		ExternalID:  resp.OpenID,
		DisplayName: resp.Nick,
		Username:    resp.Nick,
		Email:       resp.Email,
		Avatar:      resp.AvatarUrl,
		MetaInfo:    "",
	}

	return userInfo, nil
}

func (g *Connector) ConfigFields() []plugin.ConfigField {
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

func (g *Connector) ConfigReceiver(config []byte) error {
	c := &ConnectorConfig{}
	_ = json.Unmarshal(config, c)
	g.Config = c
	return nil
}
