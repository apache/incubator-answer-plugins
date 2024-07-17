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

package wecom

import (
	"embed"
	"fmt"
	"github.com/apache/incubator-answer-plugins/util"
	"net/http"
	"sync"
	"time"

	"github.com/apache/incubator-answer-plugins/user-center-wecom/i18n"
	"github.com/apache/incubator-answer/plugin"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"github.com/segmentfault/pacman/log"
)

//go:embed  info.yaml
var Info embed.FS

type UserCenter struct {
	Config          *UserCenterConfig
	Company         *Company
	UserConfigCache *UserConfigCache
	Cache           *cache.Cache
	syncLock        sync.Mutex
	syncing         bool
	syncSuccess     bool
	syncTime        time.Time
}

func (uc *UserCenter) RegisterUnAuthRouter(r *gin.RouterGroup) {
	r.GET("/wecom/login/url", uc.GetRedirectURL)
	r.GET("/wecom/login/check", uc.CheckUserLogin)
}

func (uc *UserCenter) RegisterAuthUserRouter(r *gin.RouterGroup) {
}

func (uc *UserCenter) RegisterAuthAdminRouter(r *gin.RouterGroup) {
	r.GET("/wecom/sync", uc.Sync)
	r.GET("/wecom/data", uc.Data)
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
	userDetailInfo := uc.Company.UserDetailInfoMapping[externalID]
	if userDetailInfo == nil {
		userDetailInfo, err = uc.Company.GetUserDetailInfo(externalID)
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
		Company:         NewCompany("", "", ""),
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
	redirectURL := "/user-center/auth"
	desc := plugin.UserCenterDesc{
		Name:                      "WeCom",
		DisplayName:               plugin.MakeTranslator(i18n.InfoName),
		Icon:                      "PHN2ZyB3aWR0aD0iMTYiIGhlaWdodD0iMTQiIHZpZXdCb3g9IjAgMCAxNiAxNCIgZmlsbD0iY3VycmVudENvbG9yIiB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciPgo8cGF0aCBmaWxsLXJ1bGU9ImV2ZW5vZGQiIGNsaXAtcnVsZT0iZXZlbm9kZCIgZD0iTTkuMDY4NTkgMTAuOTkwMkM4Ljc1NDMxIDEwLjgzMjcgOC41OTQ2MiAxMC41NzQ0IDguNTczODYgMTAuMjI3NUM4LjU2NjE5IDEwLjA5OTggOC41MDkwMiAxMC4wODU0IDguMzk0NjggMTAuMTE4OUM3Ljg0NTY2IDEwLjI3OTIgNy4yODQ4MiAxMC4zNTUzIDYuNzExNTIgMTAuMzY0OEM2LjEyMTI5IDEwLjM3NDcgNS41NDA2NSAxMC4zMTE4IDQuOTc3MjUgMTAuMTQ0MUM0LjY1MTQ4IDEwLjA0NzEgNC4zODg5NSAxMC4xMjg1IDQuMTIwMDIgMTAuMjk5QzMuODc4NTcgMTAuNDUyMyAzLjYzMDExIDEwLjU5NDggMy4zODE3OSAxMC43MzczTDMuMzgxNzcgMTAuNzM3M0MzLjMxMTgxIDEwLjc3NzUgMy4yNDE4NyAxMC44MTc2IDMuMTcyMDkgMTAuODU4QzMuMTYzNjcgMTAuODYyOSAzLjE1NDk4IDEwLjg2ODcgMy4xNDYwOCAxMC44NzQ4TDMuMTQ2MDUgMTAuODc0OEMzLjEwOTA1IDEwLjg5OTggMy4wNjgzMyAxMC45Mjc0IDMuMDI3NDEgMTAuODk3M0MyLjk3OTU0IDEwLjg2MiAyLjk5NTE1IDEwLjgwOCAzLjAwOTQ1IDEwLjc1ODVDMy4wMTE4NCAxMC43NTAyIDMuMDE0MjEgMTAuNzQyIDMuMDE2MjMgMTAuNzM0QzMuMDM2MTEgMTAuNjU3NiAzLjA1NTEyIDEwLjU4MDggMy4wNzQxMiAxMC41MDRDMy4xMjI4OSAxMC4zMDY5IDMuMTcxNjcgMTAuMTA5OCAzLjIzNTMzIDkuOTE3N0MzLjM0MTY4IDkuNTk2NCAzLjI0OTA2IDkuMzYxNjUgMi45ODUyNSA5LjE2Nzc4QzIuNTc0ODQgOC44NjU5NiAyLjIyOTI2IDguNDk4MzUgMS45NDM0MSA4LjA3NzRDMS4yMTg3MiA3LjAwOTA2IDEuMDIwMzggNS44NDkwNSAxLjM5MDIzIDQuNjEwNDdDMS42NTA1MyAzLjczOTE4IDIuMTY3NjIgMy4wMjg4NyAyLjg3NjY1IDIuNDY2MTFDNC4wNzkxNCAxLjUxMTc5IDUuNDU5ODUgMS4xMjU5NyA2Ljk4Mzk2IDEuMjAwMzlDNy45ODcxNSAxLjI0OTI1IDguOTIwNzEgMS41MjA0MSA5Ljc4NDMzIDIuMDMwNDdDMTAuNDIxMiAyLjQwNjcxIDEwLjk1MjMgMi44OTg1NiAxMS4zNjE4IDMuNTE1OTNDMTEuODIwNCA0LjIwNzA4IDEyLjA3MDIgNC45NjQwMyAxMi4wNTYxIDUuNzk5NTRDMTIuMDU0MiA1LjkxOTMxIDEyLjA4MSA1Ljk2MDUxIDEyLjIxMzkgNS45MTgwM0MxMi41NDQxIDUuODEyOTYgMTIuODY1MSA1LjgyMDk0IDEzLjE0NDYgNi4wNjE3NkMxMy4yNDM5IDYuMTQ3MDMgMTMuMjY5MSA2LjA5ODE3IDEzLjI3MzkgNS45OTgyQzEzLjI4NTEgNS43NjQwOSAxMy4yNzMgNS41MzEyNiAxMy4yNTM4IDUuMjk3NDdDMTMuMTk0NCA0LjU2OTU5IDEyLjk4NDYgMy44ODY3NCAxMi42MjM0IDMuMjU3MjNDMTEuNzc3IDEuNzgyNjMgMTAuNDkzNCAwLjg1MjU3NiA4Ljg5NDg0IDAuMzQyMTk4QzguMTUzNTUgMC4xMDUyMTQgNy4zODc2NiAtMC4wMjM4MTgxIDYuNzU3NTEgMC4wMDM2NDkwN0M1LjYxNTcxIDAuMDAxNzMyNzYgNC42NjY4MSAwLjE4NjAxOCAzLjc1OTQ0IDAuNTc0MzkxQzIuMzAzNjggMS4xOTgxNSAxLjE1ODM2IDIuMTY2ODUgMC40ODU3MzMgMy42MjIyOUMtMC4xNTkxMDcgNS4wMTc2OSAtMC4xNTc4MjkgNi40MzM4NCAwLjQ2NTI5MiA3Ljg0MTA2QzAuODA4NjMyIDguNjE2NTIgMS4zMzIxMSA5LjI2NDg4IDEuOTQzNDEgOS44NDQyNEMyLjA1NTE5IDkuOTQ5OTYgMi4xMDQwNiAxMC4wNTU3IDIuMDgxNyAxMC4yMDhDMi4wNTIwOCAxMC40MDggMi4wMjUwMyAxMC42MDg0IDEuOTk3OTcgMTAuODA4OEMxLjk2MTcyIDExLjA3NzMgMS45MjU0NyAxMS4zNDU4IDEuODgzMDUgMTEuNjEzM0MxLjg0NjMyIDExLjg0NDkgMS45MTQwMyAxMi4wMjIxIDIuMDk2NzEgMTIuMTU4NUMyLjI3Mzk3IDEyLjI5MDcgMi40NjI3MyAxMi4yODM3IDIuNjUzMDggMTIuMTg2OUMzLjExNzU1IDExLjk1MTEgMy41ODIwMSAxMS43MTQ2IDQuMDQ2MjggMTEuNDc4Mkw0LjI3Mzk3IDExLjM2MjNDNC40NTE1NSAxMS4yNzE5IDQuNjIxNzggMTEuMjAyNiA0LjgzMzIxIDExLjI4NjZDNC45NTgzMiAxMS4zMzYyIDUuMDk1NjkgMTEuMzU2OCA1LjIzMjEgMTEuMzc3M0g1LjIzMjExTDUuMjMyMTMgMTEuMzc3M0w1LjIzMjE1IDExLjM3NzNMNS4yMzIxNyAxMS4zNzczQzUuMjcwNTUgMTEuMzgzIDUuMzA4ODUgMTEuMzg4OCA1LjM0Njc4IDExLjM5NTJDNi4wNzQ5OCAxMS41MTY5IDYuODA3NjUgMTEuNTMyOCA3LjU0MDMzIDExLjQ1MjdDOC4wNDg3OSAxMS4zOTcxIDguNTQ1NzUgMTEuMjc4OSA5LjAzMTU0IDExLjExODZDOS4wMzgwNSAxMS4xMTY0IDkuMDQ0ODggMTEuMTE0NSA5LjA1MTc4IDExLjExMjVDOS4wODY5IDExLjEwMjUgOS4xMjQwOSAxMS4wOTE4IDkuMTMxODMgMTEuMDUxNUM5LjEzNzcgMTEuMDIxMSA5LjEwOTQ1IDExLjAwODUgOS4wODM2MSAxMC45OTdDOS4wNzg1MiAxMC45OTQ3IDkuMDczNTMgMTAuOTkyNSA5LjA2ODkxIDEwLjk5MDJIOS4wNjg1OVpNMTEuNjMxIDguNDAzNUMxMS43NjEzIDguMjc3MDIgMTEuNzczNSA4LjE0ODMxIDExLjY3MTkgOC4wMzk0QzExLjU3MzIgNy45MzM2OCAxMS40NCA3Ljk0Mjk1IDExLjMxIDguMDY0NjNDMTEuMjkzOCA4LjA3OTc5IDExLjI3ODIgOC4wOTU0IDExLjI2MjYgOC4xMTEwN0wxMS4yNjI2IDguMTExMDhMMTEuMjYyNiA4LjExMTA5TDExLjI1MzIgOC4xMjA1MkMxMC44ODg0IDguNDg2MjIgMTAuNDYxNCA4LjczOTQ5IDkuOTU1MiA4Ljg2MDg2QzkuOTExMDUgOC44NzE0MyA5Ljg2NjU4IDguODgwNSA5LjgyMjEyIDguODg5NTZDOS42NjY1MSA4LjkyMTI3IDkuNTEwODkgOC45NTI5OCA5LjM2ODggOS4wNDk2MkM5LjAyNzM4IDkuMjgxODEgOC44NTkwNiA5LjY4MiA4Ljk2OTg5IDEwLjA3MzJDOS4wOTEyNiAxMC41MDE5IDkuNDIzNzQgMTAuNzYwNiA5Ljg1MzYzIDEwLjc2MTVDMTAuMzM5MSAxMC43NjIyIDEwLjcwOTYgMTAuNDU4MSAxMC43OTcxIDkuOTgzODJDMTAuOTEwNSA5LjM2OTk2IDExLjE4MDcgOC44Mzk0NiAxMS42MzEgOC40MDMxOFY4LjQwMzVaTTE0LjE1NzcgOS42OTA2MkMxNC4yOTA5IDkuMDYwNzkgMTQuODc2OSA4Ljc1MzU1IDE1LjQ0MjMgOC45ODI4NlY4Ljk4MzE4QzE1Ljc5NTIgOS4xMjY1OSAxNi4wMDUzIDkuNDY2NDEgMTUuOTk5OSA5Ljg4NjA5QzE1Ljk5NDEgMTAuMzM3NCAxNS42NjEgMTAuNzA2MyAxNS4xODMyIDEwLjc4NDhDMTQuNjQ4MyAxMC44NzI3IDE0LjE4NDggMTEuMTAyNiAxMy43ODg1IDExLjQ3MDlDMTMuNzU3MSAxMS40OTk5IDEzLjcyNjUgMTEuNTMgMTMuNjk2IDExLjU2MDFMMTMuNjYyNiAxMS41OTI5QzEzLjUxOCAxMS43MzMxIDEzLjM3NzEgMTEuNzUzMiAxMy4yNzIgMTEuNjQ4MUMxMy4xNzA1IDExLjU0NjMgMTMuMTc5MSAxMS40MDY3IDEzLjMxOTkgMTEuMjY5N0MxMy43NzAzIDEwLjgzMjQgMTQuMDI4IDEwLjMwMjkgMTQuMTU3NyA5LjY5MDYyWk0xMi42MzAxIDExLjUxMjFDMTIuNTk2OSAxMS41MDUzIDEyLjU2MzYgMTEuNDk5MiAxMi41MzAzIDExLjQ5MzJDMTIuNDUzNyAxMS40NzkyIDEyLjM3NzIgMTEuNDY1MyAxMi4zMDMzIDExLjQ0MjhDMTEuODE2NiAxMS4yOTQ5IDExLjM4NTEgMTEuMDU2NiAxMS4wMzk4IDEwLjY3MjdDMTAuOTMxNiAxMC41NTIzIDEwLjc2NDUgMTAuNTUyIDEwLjY2ODQgMTAuNjQ4NUMxMC41Nzg2IDEwLjczODggMTAuNTg1MyAxMC44OTE4IDEwLjcwNjEgMTEuMDAzOUMxMS4xODcxIDExLjQ0ODggMTEuNDQyOSAxMi4wMDY4IDExLjU2NjggMTIuNjQwOEMxMS42NTQzIDEzLjA4OTIgMTIuMDI1NSAxMy4zNzYzIDEyLjQ5NTYgMTMuMzc3M0MxMi42MDYxIDEzLjM4NCAxMi43MjQzIDEzLjM1NTYgMTIuODM3MyAxMy4zMDQ4QzEzLjIyNjcgMTMuMTI5OCAxMy40NDkgMTIuNzUyNiAxMy40MDU1IDEyLjM0MTVDMTMuMzYyNyAxMS45MzY5IDEzLjA0NDkgMTEuNTk3NCAxMi42Mjk3IDExLjUxMjRMMTIuNjMwMSAxMS41MTIxWk0xMS41NDM4IDcuMTM0NThDMTEuNTk0MyA2LjY3MDE5IDEyLjAxMDggNi4zMDA5OCAxMi40ODQxIDYuMzAxNjJIMTIuNDg0NEMxMi45MjggNi4zMDIyNiAxMy4zMjkyIDYuNjM2NjYgMTMuNDAyNiA3LjA4NDc2QzEzLjUwNTUgNy43MTIwMyAxMy43OTggOC4yMzQ1NCAxNC4yNDcxIDguNjc2MjVDMTQuMzczMyA4LjgwMDgxIDE0LjM3OSA4LjkwMjA2IDE0LjI2OTEgOS4wMjI3OUMxNC4xNzAxIDkuMTMxNyAxNC4wNjE1IDkuMTM4NDEgMTMuOTM2MyA5LjAxNjcyQzEzLjY0IDguNzI4MzEgMTMuMzE1OCA4LjQ4NDMgMTIuOTMgOC4zMjcxN0MxMi43NDE2IDguMjUwMTkgMTIuNTQ0OCA4LjIwNTMyIDEyLjM0ODEgOC4xNjA0OEMxMi4zMjIgOC4xNTQ1MyAxMi4yOTU5IDguMTQ4NTggMTIuMjY5OCA4LjE0MjU2QzExLjgxNCA4LjAzNzQ4IDExLjQ5NCA3LjU5NDUgMTEuNTQzOCA3LjEzNDU4WiIgZmlsbD0iY3VycmVudENvbG9yIi8+Cjwvc3ZnPg==",
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
	code := ctx.Query("code")
	if len(code) == 0 {
		return nil, fmt.Errorf("code is empty")
	}
	state := ctx.Query("state")
	if len(state) == 0 {
		return nil, fmt.Errorf("state is empty")
	}
	log.Debugf("request code: %s, state: %s", code, state)

	info, err := uc.Company.AuthUser(code)
	if err != nil {
		return nil, fmt.Errorf("auth user failed: %w", err)
	}
	if !info.IsAvailable {
		return nil, fmt.Errorf("user is not available")
	}
	if len(info.GetEmail()) == 0 {
		ctx.Redirect(http.StatusFound, "/user-center/auth-failed")
		ctx.Abort()
		return nil, fmt.Errorf("user email is empty")
	}

	userInfo = &plugin.UserCenterBasicUserInfo{}
	userInfo.ExternalID = info.Userid
	userInfo.Username = info.Userid
	userInfo.DisplayName = info.Name
	userInfo.Email = info.GetEmail()
	userInfo.Rank = 0
	userInfo.Avatar = info.Avatar
	userInfo.Mobile = info.Mobile

	uc.Cache.Set(state, userInfo.ExternalID, time.Minute*5)
	return userInfo, nil
}

func (uc *UserCenter) SignUpCallback(ctx *plugin.GinContext) (userInfo *plugin.UserCenterBasicUserInfo, err error) {
	return uc.LoginCallback(ctx)
}

func (uc *UserCenter) UserInfo(externalID string) (userInfo *plugin.UserCenterBasicUserInfo, err error) {
	userDetailInfo := uc.Company.UserDetailInfoMapping[externalID]
	if userDetailInfo == nil {
		userDetailInfo, err = uc.Company.GetUserDetailInfo(externalID)
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
		Username:    userDetailInfo.Userid,
		DisplayName: userDetailInfo.Name,
		Bio:         uc.Company.formatDepartmentAndPosition(userDetailInfo.Department, userDetailInfo.Position),
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

func (uc *UserCenter) asyncCompany() {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Errorf("sync data panic: %s", err)
			}
		}()
		uc.syncCompany()
	}()
}

func (uc *UserCenter) syncCompany() {
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

	log.Info("start sync company data")
	uc.syncing = true
	uc.syncSuccess = true

	if err := uc.Company.ListDepartmentAll(); err != nil {
		log.Errorf("list department error: %s", err)
		uc.syncSuccess = false
		return
	}
	if err := uc.Company.ListUser(); err != nil {
		log.Errorf("list user error: %s", err)
		uc.syncSuccess = false
		return
	}
	log.Info("end sync company data")
}
