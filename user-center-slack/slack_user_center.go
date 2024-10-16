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
	"embed"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/apache/incubator-answer-plugins/util"

	"github.com/apache/incubator-answer-plugins/user-center-slack/i18n"
	"github.com/apache/incubator-answer/plugin"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"github.com/segmentfault/pacman/log"
)

//go:embed  info.yaml
var Info embed.FS

type Importer struct{}

type UserCenter struct {
	Config          *UserCenterConfig
	SlackClient     *SlackClient
	UserConfigCache *UserConfigCache
	Cache           *cache.Cache
	syncLock        sync.Mutex
	syncing         bool
	syncSuccess     bool
	syncTime        time.Time
	importerFunc    plugin.ImporterFunc
}

func (uc *UserCenter) RegisterUnAuthRouter(r *gin.RouterGroup) {
	r.GET("/slack/login/url", uc.GetSlackRedirectURL)
	r.POST("/slack/slash", uc.SlashCommand)
}

func (uc *UserCenter) RegisterAuthUserRouter(r *gin.RouterGroup) {
}

func (uc *UserCenter) RegisterAuthAdminRouter(r *gin.RouterGroup) {
	r.GET("/slack/sync", uc.Sync)
}

func (uc *UserCenter) AfterLogin(externalID, accessToken string) {
	log.Debugf("user %s is login", externalID)
	uc.Cache.Set(externalID, accessToken, time.Minute*5)
}

func (uc *UserCenter) UserStatus(externalID string) (userStatus plugin.UserStatus) {
	if len(externalID) == 0 {
		return plugin.UserStatusAvailable
	}

	var err error
	userDetailInfo := uc.SlackClient.UserInfoMapping[externalID]
	if userDetailInfo == nil {
		userDetailInfo, err = uc.SlackClient.GetUserDetailInfo(externalID)
		if err != nil {
			log.Errorf("get user detail info failed: %v", err)
		}
	}
	if userDetailInfo == nil {
		return plugin.UserStatusDeleted
	}
	switch userDetailInfo.Status {
	case 1:
		return plugin.UserStatusAvailable
	case 2:
		return plugin.UserStatusSuspended
	default:
		return plugin.UserStatusDeleted
	}
}

func init() {
	uc := &UserCenter{
		Config:          &UserCenterConfig{},
		UserConfigCache: NewUserConfigCache(),
		SlackClient:     NewSlackClient("", ""),
		Cache:           cache.New(5*time.Minute, 10*time.Minute),
		syncLock:        sync.Mutex{},
	}

	plugin.Register(uc)
	uc.CronSyncData()
}

func (uc *UserCenter) Info() plugin.Info {
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

func (uc *UserCenter) Description() plugin.UserCenterDesc {
	redirectURL := uc.BuildSlackBaseRedirectURL()
	desc := plugin.UserCenterDesc{
		Name:                      "Slack",
		DisplayName:               plugin.MakeTranslator(i18n.InfoName),
		Icon:                      "",
		Url:                       "",
		LoginRedirectURL:          redirectURL,
		SignUpRedirectURL:         redirectURL,
		RankAgentEnabled:          false,
		UserStatusAgentEnabled:    false,
		UserRoleAgentEnabled:      false,
		MustAuthEmailEnabled:      true,
		EnabledOriginalUserSystem: true,
	}
	return desc
}

func (uc *UserCenter) ControlCenterItems() []plugin.ControlCenter {
	var controlCenterItems []plugin.ControlCenter
	return controlCenterItems
}

func (uc *UserCenter) LoginCallback(ctx *plugin.GinContext) (userInfo *plugin.UserCenterBasicUserInfo, err error) {
	log.Debugf("Processing LoginCallback")
	CallbackURL := ctx.Request.URL.String()
	log.Debugf("callbackURL in SlackLoginCallback:", CallbackURL)
	code := ctx.Query("code")
	if len(code) == 0 {
		return nil, fmt.Errorf("code is empty")
	}

	state := ctx.Query("state")
	if len(state) == 0 {
		return nil, fmt.Errorf("state is empty")
	}
	log.Debugf("request code: %s, state: %s", code, state)

	expectedState, exist := uc.Cache.Get("oauth_state_" + state)
	if !exist {
		fmt.Println("State not found in cache or expired")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired state"})
		return
	}
	if state != expectedState {
		fmt.Println("State mismatch")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid state"})
		return
	}
	log.Debugf("State validated successfully")

	info, err := uc.SlackClient.AuthUser(code)
	if err != nil {
		return nil, fmt.Errorf("auth user failed: %w", err)
	}
	if !info.IsAvailable {
		return nil, fmt.Errorf("user is not available")
	}
	//Get Email
	if len(info.Profile.Email) == 0 {
		ctx.Redirect(http.StatusFound, "/user-center/auth-failed")
		return nil, fmt.Errorf("user email is empty")
	}

	userInfo = &plugin.UserCenterBasicUserInfo{}
	userInfo.ExternalID = info.ID
	userInfo.Username = info.ID
	userInfo.DisplayName = info.Name
	userInfo.Email = info.Profile.Email
	userInfo.Rank = 0
	userInfo.Mobile = ""
	userInfo.Avatar = info.Profile.Image192

	uc.Cache.Set(state, userInfo.ExternalID, time.Minute*5)
	return userInfo, nil
}

func (uc *UserCenter) SignUpCallback(ctx *plugin.GinContext) (userInfo *plugin.UserCenterBasicUserInfo, err error) {
	return uc.LoginCallback(ctx)
}

func (uc *UserCenter) UserInfo(externalID string) (userInfo *plugin.UserCenterBasicUserInfo, err error) {
	userDetailInfo := uc.SlackClient.UserInfoMapping[externalID]
	if userDetailInfo == nil {
		userDetailInfo, err = uc.SlackClient.GetUserDetailInfo(externalID)
		if err != nil {
			log.Errorf("get user detail info failed: %v", err)
			userInfo = &plugin.UserCenterBasicUserInfo{
				ExternalID: externalID,
				Status:     plugin.UserStatusDeleted,
			}
			return userInfo, nil
		}
	}

	userInfo = &plugin.UserCenterBasicUserInfo{
		ExternalID:  externalID,
		Username:    userDetailInfo.ID,
		DisplayName: userDetailInfo.Name,
		Bio:         "",
	}
	switch userDetailInfo.Status {
	case 1:
		userInfo.Status = plugin.UserStatusAvailable
	case 2:
		userInfo.Status = plugin.UserStatusSuspended
	default:
		userInfo.Status = plugin.UserStatusDeleted
	}
	return userInfo, nil
}

func (uc *UserCenter) UserList(externalIDs []string) (userList []*plugin.UserCenterBasicUserInfo, err error) {
	userList = make([]*plugin.UserCenterBasicUserInfo, 0)
	return userList, nil
}

func (uc *UserCenter) UserSettings(externalID string) (userSettings *plugin.SettingInfo, err error) {
	return &plugin.SettingInfo{
		ProfileSettingRedirectURL: "",
		AccountSettingRedirectURL: "",
	}, nil
}

func (uc *UserCenter) PersonalBranding(externalID string) (branding []*plugin.PersonalBranding) {
	return branding
}
