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

package slack_user_center

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/apache/incubator-answer/plugin"
	"github.com/go-resty/resty/v2"
	"github.com/segmentfault/pacman/log"
)

type SlackClient struct {
	AccessToken  string
	ClientID     string
	ClientSecret string
	RedirectURI  string
	AuthedUserID string

	UserInfoMapping map[string]*UserInfo
	ChannelMapping  string
}

func NewSlackClient(clientID, clientSecret string) *SlackClient {
	return &SlackClient{
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}
}

// OAuthV2ResponseTeam
type OAuthV2ResponseTeam struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// OAuthResponseIncomingWebhook
type OAuthResponseIncomingWebhook struct {
	URL              string `json:"url"`
	Channel          string `json:"channel"`
	ChannelID        string `json:"channel_id,omitempty"`
	ConfigurationURL string `json:"configuration_url"`
}

// OAuthV2ResponseEnterprise
type OAuthV2ResponseEnterprise struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// OAuthV2ResponseAuthedUser
type OAuthV2ResponseAuthedUser struct {
	ID           string `json:"id"`
	Scope        string `json:"scope"`
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
}

type TokenResponse struct {
	AccessToken         string                       `json:"access_token"`
	TokenType           string                       `json:"token_type"`
	Scope               string                       `json:"scope"`
	BotUserID           string                       `json:"bot_user_id"`
	AppID               string                       `json:"app_id"`
	Team                OAuthV2ResponseTeam          `json:"team"`
	IncomingWebhook     OAuthResponseIncomingWebhook `json:"incoming_webhook"`
	Enterprise          OAuthV2ResponseEnterprise    `json:"enterprise"`
	IsEnterpriseInstall bool                         `json:"is_enterprise_install"`
	AuthedUser          OAuthV2ResponseAuthedUser    `json:"authed_user"`
	RefreshToken        string                       `json:"refresh_token"`
	ExpiresIn           int                          `json:"expires_in"`
	Error               string                       `json:"error,omitempty"`
}

type Member struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	TeamID string `json:"team_id"`
}

// ExchangeCodeForUser through OAuthToken
func (sc *SlackClient) AuthUser(code string) (info *UserInfo, err error) {
	clientID := sc.ClientID
	clientSecret := sc.ClientSecret
	redirectURI := fmt.Sprintf("%s/answer/api/v1/user-center/login/callback", plugin.SiteURL())

	data := url.Values{}
	data.Set("code", code)
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("redirect_uri", redirectURI)

	resp, err := http.PostForm("https://slack.com/api/oauth.v2.access", data)
	if err != nil {
		log.Errorf("Failed to exchange code for token: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read response body: %v", err)
	}

	var tokenResp TokenResponse
	err = json.Unmarshal([]byte(body), &tokenResp)
	if err != nil {
		fmt.Println("Error parsing response:", err)
		return nil, err
	}

	if tokenResp.Error != "" {
		return nil, fmt.Errorf("Slack API error in AuthUser: %s", tokenResp.Error)
	}

	sc.AccessToken = tokenResp.AccessToken
	sc.AuthedUserID = tokenResp.AuthedUser.ID

	return sc.GetUserDetailInfo(sc.AuthedUserID)
}

func (sc *SlackClient) GetUserDetailInfo(userid string) (info *UserInfo, err error) {
	getUserInfoResp, err := resty.New().R().
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", sc.AccessToken)).
		SetHeader("Accept", "application/json").
		Get("https://slack.com/api/users.info?user=" + userid)
	if err != nil {
		log.Errorf("Failed to get user info: %v", err)
		return nil, err
	}

	var authUserResp *AuthUserResp
	err = json.Unmarshal([]byte(getUserInfoResp.String()), &authUserResp)
	if err != nil {
		log.Errorf("Error unmarshaling user info: %v", err)
		return nil, err
	}
	if !authUserResp.Ok {
		log.Errorf("Failed to get valid user info, Slack API error: %s", authUserResp.Errmsg)
		return nil, fmt.Errorf("Get user info failed: %s", authUserResp.Errmsg)
	}
	log.Debugf("Get user info for UserID: %s", userid)

	if authUserResp.User == nil {
		log.Errorf("No user data available in the response")
		return nil, fmt.Errorf("No user data available in the response")
	}

	authUserResp.User.IsAvailable = true
	authUserResp.User.Status = 1

	// Directly returning the user data parsed from the response
	return authUserResp.User, nil
}

func (sc *SlackClient) UpdateUserInfo() (err error) {
	log.Debug("Try to update slack client")

	userInfo, err := sc.GetUserDetailInfo(sc.AuthedUserID)
	if err != nil {
		log.Errorf("Failed to update user info: %v", err)
		return err
	}

	if sc.UserInfoMapping == nil {
		sc.UserInfoMapping = make(map[string]*UserInfo)
	}
	sc.UserInfoMapping[sc.AuthedUserID] = userInfo
	log.Infof("Updated user info for UserID: %s", sc.AuthedUserID)

	return nil
}
