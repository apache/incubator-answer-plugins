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

package google

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"github.com/apache/incubator-answer-plugins/util"
	"io"
	"time"

	"github.com/apache/incubator-answer-plugins/connector-google/i18n"
	"github.com/apache/incubator-answer/plugin"
	"golang.org/x/oauth2"
	oauth2Google "golang.org/x/oauth2/google"
)

//go:embed  info.yaml
var Info embed.FS

type Connector struct {
	Config *ConnectorConfig
}

type ConnectorConfig struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type AuthUserInfo struct {
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Profile       string `json:"profile"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Gender        string `json:"gender"`
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
	return `PHN2ZyBpZD0iQ2FwYV8xIiBzdHlsZT0iZW5hYmxlLWJhY2tncm91bmQ6bmV3IDAgMCAxNTAgMTUwOyIgdmVyc2lvbj0iMS4xIiB2aWV3Qm94PSIwIDAgMTUwIDE1MCIgeG1sOnNwYWNlPSJwcmVzZXJ2ZSIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIiB4bWxuczp4bGluaz0iaHR0cDovL3d3dy53My5vcmcvMTk5OS94bGluayI+PHN0eWxlIHR5cGU9InRleHQvY3NzIj4KCS5zdDB7ZmlsbDojMUE3M0U4O30KCS5zdDF7ZmlsbDojRUE0MzM1O30KCS5zdDJ7ZmlsbDojNDI4NUY0O30KCS5zdDN7ZmlsbDojRkJCQzA0O30KCS5zdDR7ZmlsbDojMzRBODUzO30KCS5zdDV7ZmlsbDojNENBRjUwO30KCS5zdDZ7ZmlsbDojMUU4OEU1O30KCS5zdDd7ZmlsbDojRTUzOTM1O30KCS5zdDh7ZmlsbDojQzYyODI4O30KCS5zdDl7ZmlsbDojRkJDMDJEO30KCS5zdDEwe2ZpbGw6IzE1NjVDMDt9Cgkuc3QxMXtmaWxsOiMyRTdEMzI7fQoJLnN0MTJ7ZmlsbDojRjZCNzA0O30KCS5zdDEze2ZpbGw6I0U1NDMzNTt9Cgkuc3QxNHtmaWxsOiM0MjgwRUY7fQoJLnN0MTV7ZmlsbDojMzRBMzUzO30KCS5zdDE2e2NsaXAtcGF0aDp1cmwoI1NWR0lEXzJfKTt9Cgkuc3QxN3tmaWxsOiMxODgwMzg7fQoJLnN0MTh7b3BhY2l0eTowLjI7ZmlsbDojRkZGRkZGO2VuYWJsZS1iYWNrZ3JvdW5kOm5ldyAgICA7fQoJLnN0MTl7b3BhY2l0eTowLjM7ZmlsbDojMEQ2NTJEO2VuYWJsZS1iYWNrZ3JvdW5kOm5ldyAgICA7fQoJLnN0MjB7Y2xpcC1wYXRoOnVybCgjU1ZHSURfNF8pO30KCS5zdDIxe29wYWNpdHk6MC4zO2ZpbGw6dXJsKCNfNDVfc2hhZG93XzFfKTtlbmFibGUtYmFja2dyb3VuZDpuZXcgICAgO30KCS5zdDIye2NsaXAtcGF0aDp1cmwoI1NWR0lEXzZfKTt9Cgkuc3QyM3tmaWxsOiNGQTdCMTc7fQoJLnN0MjR7b3BhY2l0eTowLjM7ZmlsbDojMTc0RUE2O2VuYWJsZS1iYWNrZ3JvdW5kOm5ldyAgICA7fQoJLnN0MjV7b3BhY2l0eTowLjM7ZmlsbDojQTUwRTBFO2VuYWJsZS1iYWNrZ3JvdW5kOm5ldyAgICA7fQoJLnN0MjZ7b3BhY2l0eTowLjM7ZmlsbDojRTM3NDAwO2VuYWJsZS1iYWNrZ3JvdW5kOm5ldyAgICA7fQoJLnN0Mjd7ZmlsbDp1cmwoI0ZpbmlzaF9tYXNrXzFfKTt9Cgkuc3QyOHtmaWxsOiNGRkZGRkY7fQoJLnN0Mjl7ZmlsbDojMEM5RDU4O30KCS5zdDMwe29wYWNpdHk6MC4yO2ZpbGw6IzAwNEQ0MDtlbmFibGUtYmFja2dyb3VuZDpuZXcgICAgO30KCS5zdDMxe29wYWNpdHk6MC4yO2ZpbGw6IzNFMjcyMztlbmFibGUtYmFja2dyb3VuZDpuZXcgICAgO30KCS5zdDMye2ZpbGw6I0ZGQzEwNzt9Cgkuc3QzM3tvcGFjaXR5OjAuMjtmaWxsOiMxQTIzN0U7ZW5hYmxlLWJhY2tncm91bmQ6bmV3ICAgIDt9Cgkuc3QzNHtvcGFjaXR5OjAuMjt9Cgkuc3QzNXtmaWxsOiMxQTIzN0U7fQoJLnN0MzZ7ZmlsbDp1cmwoI1NWR0lEXzdfKTt9Cgkuc3QzN3tmaWxsOiNGQkJDMDU7fQoJLnN0Mzh7Y2xpcC1wYXRoOnVybCgjU1ZHSURfOV8pO2ZpbGw6I0U1MzkzNTt9Cgkuc3QzOXtjbGlwLXBhdGg6dXJsKCNTVkdJRF8xMV8pO2ZpbGw6I0ZCQzAyRDt9Cgkuc3Q0MHtjbGlwLXBhdGg6dXJsKCNTVkdJRF8xM18pO2ZpbGw6I0U1MzkzNTt9Cgkuc3Q0MXtjbGlwLXBhdGg6dXJsKCNTVkdJRF8xNV8pO2ZpbGw6I0ZCQzAyRDt9Cjwvc3R5bGU+PGc+PHBhdGggY2xhc3M9InN0MTQiIGQ9Ik0xMjAsNzYuMWMwLTMuMS0wLjMtNi4zLTAuOC05LjNINzUuOXYxNy43aDI0LjhjLTEsNS43LTQuMywxMC43LTkuMiwxMy45bDE0LjgsMTEuNSAgIEMxMTUsMTAxLjgsMTIwLDkwLDEyMCw3Ni4xTDEyMCw3Ni4xeiIvPjxwYXRoIGNsYXNzPSJzdDE1IiBkPSJNNzUuOSwxMjAuOWMxMi40LDAsMjIuOC00LjEsMzAuNC0xMS4xTDkxLjUsOTguNGMtNC4xLDIuOC05LjQsNC40LTE1LjYsNC40Yy0xMiwwLTIyLjEtOC4xLTI1LjgtMTguOSAgIEwzNC45LDk1LjZDNDIuNywxMTEuMSw1OC41LDEyMC45LDc1LjksMTIwLjl6Ii8+PHBhdGggY2xhc3M9InN0MTIiIGQ9Ik01MC4xLDgzLjhjLTEuOS01LjctMS45LTExLjksMC0xNy42TDM0LjksNTQuNGMtNi41LDEzLTYuNSwyOC4zLDAsNDEuMkw1MC4xLDgzLjh6Ii8+PHBhdGggY2xhc3M9InN0MTMiIGQ9Ik03NS45LDQ3LjNjNi41LTAuMSwxMi45LDIuNCwxNy42LDYuOUwxMDYuNiw0MUM5OC4zLDMzLjIsODcuMywyOSw3NS45LDI5LjFjLTE3LjQsMC0zMy4yLDkuOC00MSwyNS4zICAgbDE1LjIsMTEuOEM1My44LDU1LjMsNjMuOSw0Ny4zLDc1LjksNDcuM3oiLz48L2c+PC9zdmc+`
}

func (g *Connector) ConnectorName() plugin.Translator {
	return plugin.MakeTranslator(i18n.ConnectorName)
}

func (g *Connector) ConnectorSlugName() string {
	return "google"
}

func (g *Connector) ConnectorSender(ctx *plugin.GinContext, receiverURL string) (redirectURL string) {
	oauth2Config := &oauth2.Config{
		ClientID:     g.Config.ClientID,
		ClientSecret: g.Config.ClientSecret,
		Endpoint:     oauth2Google.Endpoint,
		RedirectURL:  receiverURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
			"openid",
		},
	}
	return oauth2Config.AuthCodeURL("state")
}

func (g *Connector) ConnectorReceiver(ctx *plugin.GinContext, receiverURL string) (userInfo plugin.ExternalLoginUserInfo, err error) {
	code := ctx.Query("code")
	oauth2Config := &oauth2.Config{
		ClientID:     g.Config.ClientID,
		ClientSecret: g.Config.ClientSecret,
		Endpoint:     oauth2Google.Endpoint,
		RedirectURL:  receiverURL,
	}

	token, err := oauth2Config.Exchange(context.Background(), code)
	if err != nil {
		return userInfo, err
	}

	client := oauth2Config.Client(context.TODO(), token)
	client.Timeout = 15 * time.Second
	response, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return userInfo, err
	}
	defer response.Body.Close()
	data, _ := io.ReadAll(response.Body)

	respGoogleAuthUserInfo := &AuthUserInfo{}
	if err = json.Unmarshal(data, respGoogleAuthUserInfo); err != nil {
		return userInfo, fmt.Errorf("parse google oauth user info response failed: %v", err)
	}

	userInfo = plugin.ExternalLoginUserInfo{
		ExternalID:  respGoogleAuthUserInfo.Sub,
		DisplayName: respGoogleAuthUserInfo.Name,
		Username:    respGoogleAuthUserInfo.Name,
		Email:       respGoogleAuthUserInfo.Email,
		Avatar:      respGoogleAuthUserInfo.Picture,
		MetaInfo:    string(data),
	}
	// If email is not verified, set it to empty
	if !respGoogleAuthUserInfo.EmailVerified {
		userInfo.Email = ""
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
