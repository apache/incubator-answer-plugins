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
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/apache/incubator-answer/plugin"
	"github.com/gin-gonic/gin"
	"github.com/segmentfault/pacman/log"
)

// RespBody response body.
type RespBody struct {
	// http code
	Code int `json:"code"`
	// reason key
	Reason string `json:"reason"`
	// response message
	Message string `json:"msg"`
	// response data
	Data interface{} `json:"data"`
}

// NewRespBodyData new response body with data
func NewRespBodyData(code int, reason string, data interface{}) *RespBody {
	return &RespBody{
		Code:   code,
		Reason: reason,
		Data:   data,
	}
}

func (uc *UserCenter) BuildSlackBaseRedirectURL() string {
	clientID := uc.Config.ClientID
	log.Debug("Get client ID:", clientID)
	scope := "chat:write,commands,groups:write,im:write,incoming-webhook,mpim:write,users:read,users:read.email"
	response_type := "code"
	redirect_uri := fmt.Sprintf("%s/answer/api/v1/user-center/login/callback", plugin.SiteURL())

	base_redirectURL := fmt.Sprintf(
		"https://slack.com/oauth/v2/authorize?client_id=%s&scope=%s&response_type=%s&redirect_uri=%s",
		clientID, scope, response_type, redirect_uri,
	)

	state := genState()
	nonce := genNonce()
	uc.Cache.Set("oauth_state_"+state, state, time.Minute*5)

	redirectURL := fmt.Sprintf("%s&state=%s&nonce=%s", base_redirectURL, state, nonce)
	log.Debug("RedirectURL from BuildSlackBaseRedirectURL:", redirectURL)

	return redirectURL
}

func (uc *UserCenter) GetSlackRedirectURL(ctx *gin.Context) {
	redirectURL := uc.BuildSlackBaseRedirectURL()
	log.Debug("Processing GetSlackRedirectURL")

	ctx.Writer.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(ctx.Writer)
	encoder.SetEscapeHTML(false)

	respData := NewRespBodyData(http.StatusOK, "success", map[string]string{
		"redirect_url": redirectURL,
	})
	ctx.Writer.WriteHeader(http.StatusOK)
	if err := encoder.Encode(respData); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encode response"})
		return
	}
}

func genNonce() string {
	bytes := make([]byte, 10)
	_, _ = rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func genState() string {
	bytes := make([]byte, 32)
	_, _ = rand.Read(bytes)
	return base64.URLEncoding.EncodeToString(bytes)
}

func (uc *UserCenter) Sync(ctx *gin.Context) {
	uc.syncSlackClient()

	if uc.syncSuccess {
		ctx.JSON(http.StatusOK, NewRespBodyData(http.StatusOK, "success", map[string]any{
			"message": "User data synced successfully",
		}))
		return
	}

	errRespBodyData := NewRespBodyData(http.StatusBadRequest, "error", map[string]any{
		"err_type": "toast",
	})
	errRespBodyData.Message = "Failed to sync user data"
	ctx.JSON(http.StatusBadRequest, errRespBodyData)
}

func (uc *UserCenter) syncSlackClient() {
	if !uc.syncLock.TryLock() {
		log.Infof("sync data is running")
		return
	}
	defer func() {
		uc.syncing = false
		if uc.syncSuccess {
			uc.syncTime = time.Now()
		}
		uc.syncLock.Unlock()
	}()

	log.Info("start sync slack data")
	uc.syncing = true
	uc.syncSuccess = true

	if err := uc.SlackClient.UpdateUserInfo(); err != nil {
		log.Errorf("list user error: %s", err)
		uc.syncSuccess = false
		return
	}
	log.Info("end sync slack data")
}
