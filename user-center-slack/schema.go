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

type AuthUserResp struct {
	Ok     bool      `json:"ok"`
	Errmsg string    `json:"error"`
	User   *UserInfo `json:"user"`
}

type UserProfile struct {
	AvatarHash    string `json:"avatar_hash"`
	StatusText    string `json:"status_text"`
	StatusEmoji   string `json:"status_emoji"`
	RealName      string `json:"real_name"`
	DisplayName   string `json:"display_name"`
	Email         string `json:"email"`
	ImageOriginal string `json:"image_original"`
	Image24       string `json:"image_24"`
	Image32       string `json:"image_32"`
	Image48       string `json:"image_48"`
	Image72       string `json:"image_72"`
	Image192      string `json:"image_192"`
	Image512      string `json:"image_512"`
}

type UserInfo struct {
	ID                string      `json:"id"`
	TeamID            string      `json:"team_id"`
	Name              string      `json:"name"`
	RealName          string      `json:"real_name"`
	Deleted           bool        `json:"deleted"`
	TimeZone          string      `json:"tz"`
	TimeZoneLabel     string      `json:"tz_label"`
	TimeZoneOffset    int         `json:"tz_offset"`
	Profile           UserProfile `json:"profile"`
	IsAdmin           bool        `json:"is_admin"`
	IsOwner           bool        `json:"is_owner"`
	IsPrimaryOwner    bool        `json:"is_primary_owner"`
	IsRestricted      bool        `json:"is_restricted"`
	IsUltraRestricted bool        `json:"is_ultra_restricted"`
	IsBot             bool        `json:"is_bot"`
	Updated           int64       `json:"updated"`
	IsAppUser         bool        `json:"is_app_user"`
	Has2FA            bool        `json:"has_2fa"`

	LastLogin   int64 `json:"last_login,omitempty"`
	IsAvailable bool  `json:"is_available"`
	Enable      bool  `json:"true"`
	Status      int   `json:"status"`
}

type WebhookReq struct {
	Blocks []struct {
		Type string `json:"type"`
		Text struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"text"`
	} `json:"blocks"`
}

func NewWebhookReq(content string) *WebhookReq {
	return &WebhookReq{
		Blocks: []struct {
			Type string `json:"type"`
			Text struct {
				Type string `json:"type"`
				Text string `json:"text"`
			} `json:"text"`
		}{
			{
				Type: "section",
				Text: struct {
					Type string `json:"type"`
					Text string `json:"text"`
				}{
					Type: "mrkdwn",
					Text: content,
				},
			},
		},
	}
}

type SlackUserResponse struct {
	Ok   bool `json:"ok"`
	User struct {
		Profile struct {
			Email string `json:"email"`
		} `json:"profile"`
	} `json:"user"`
}
