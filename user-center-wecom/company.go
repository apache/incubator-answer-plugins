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
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/segmentfault/pacman/log"
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/work"
	workConfig "github.com/silenceper/wechat/v2/work/config"
	"github.com/tidwall/gjson"
)

type Company struct {
	CorpID      string
	CorpSecret  string
	AgentID     string
	CallbackURL string

	Work                  *work.Work
	DepartmentMapping     map[int]*Department
	EmployeeMapping       map[string]*Employee
	UserDetailInfoMapping map[string]*UserDetailInfo
}

func NewCompany(corpID, corpSecret, agentID string) *Company {
	memory := cache.NewMemory()
	cfg := &workConfig.Config{
		CorpID:     corpID,
		CorpSecret: corpSecret,
		AgentID:    agentID,
		Cache:      memory,
	}
	newWork := work.NewWork(cfg)
	return &Company{
		CorpID:                corpID,
		CorpSecret:            corpSecret,
		AgentID:               agentID,
		Work:                  newWork,
		DepartmentMapping:     make(map[int]*Department),
		EmployeeMapping:       make(map[string]*Employee),
		UserDetailInfoMapping: make(map[string]*UserDetailInfo),
	}
}

func (c *Company) ListDepartmentAll() (err error) {
	log.Debugf("try to list department all")
	token, err := c.Work.GetOauth().GetAccessToken()
	if err != nil {
		return fmt.Errorf("get access token failed: %w", err)
	}
	log.Debugf("get access token success")

	resp, err := resty.New().R().Get("https://qyapi.weixin.qq.com/cgi-bin/department/list?access_token=" + token)
	if err != nil {
		return fmt.Errorf("get department list failed: %w", err)
	}
	if resp.StatusCode() != 200 {
		return fmt.Errorf("get department list failed: %s", resp.String())
	}

	log.Debugf("get department success: %s", resp.String())

	department := gjson.Get(resp.String(), "department").String()
	var departments []*Department
	err = json.Unmarshal([]byte(department), &departments)
	if err != nil {
		return fmt.Errorf("unmarshal department failed: %w", err)
	}

	departmentMapping := make(map[int]*Department)
	for _, d := range departments {
		departmentMapping[d.Id] = d
	}
	c.DepartmentMapping = departmentMapping
	log.Debugf("get department list: %d", len(departments))
	return nil
}

func (c *Company) ListUser() (err error) {
	token, err := c.Work.GetOauth().GetAccessToken()
	if err != nil {
		return fmt.Errorf("get access token failed: %w", err)
	}
	log.Debugf("get access token success")

	for _, department := range c.DepartmentMapping {
		log.Debugf("try to get department user list: %d %s", department.Id, department.Name)
		resp, err := resty.New().R().Get(fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/user/simplelist?department_id=%d&access_token=%s",
			department.Id, token))
		if err != nil {
			log.Errorf("get department user list failed: %v", err)
			continue
		}
		if gjson.Get(resp.String(), "errcode").Int() != 0 {
			log.Errorf("get department user list failed: %v", resp.String())
			continue
		}

		userList := gjson.Get(resp.String(), "userlist").String()
		var employees []*Employee
		err = json.Unmarshal([]byte(userList), &employees)
		if err != nil {
			log.Errorf("unmarshal userList failed: %w", err)
			continue
		}
		log.Debugf("get department user list: %d", len(employees))
		for _, employee := range employees {
			c.EmployeeMapping[employee.Userid] = employee
			log.Debugf(employee.Userid)
		}
	}

	log.Debugf("all user amount: %d", len(c.EmployeeMapping))
	return nil
}

func (c *Company) AuthUser(code string) (info *UserInfo, err error) {
	token, err := c.Work.GetOauth().GetAccessToken()
	if err != nil {
		return nil, fmt.Errorf("get access token failed: %w", err)
	}

	getUserInfoResp, err := resty.New().R().Get(fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/auth/getuserinfo?access_token=" + token + "&code=" + code))
	if err != nil {
		log.Errorf("get user info failed: %v", err)
		return nil, err
	}
	log.Debugf("get user info: %s", getUserInfoResp.String())

	userTicket := gjson.Get(getUserInfoResp.String(), "user_ticket").String()

	getUserDetailResp, err := resty.New().R().
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]string{"user_ticket": userTicket}).
		Post(fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/auth/getuserdetail?access_token=" + token))
	if err != nil {
		log.Errorf("get user info failed: %v", err)
		return nil, err
	}

	var userInfoResp *AuthUserInfoResp
	err = json.Unmarshal([]byte(getUserDetailResp.String()), &userInfoResp)
	if err != nil {
		log.Errorf("unmarshal user info failed: %s", err)
		return nil, err
	}
	if userInfoResp.Errcode != 0 {
		log.Errorf("get user info failed: %v", getUserDetailResp.String())
		return nil, fmt.Errorf("get user info failed")
	}
	log.Debugf("get user info: %s", getUserDetailResp.String())

	employee := c.EmployeeMapping[userInfoResp.Userid]
	if employee == nil {
		return nil, fmt.Errorf("user %s not found in employee list", userInfoResp.Userid)
	}

	userDetailInfo, err := c.GetUserDetailInfo(userInfoResp.Userid)
	if err != nil {
		return nil, err
	}

	userInfo := &UserInfo{
		Userid:        userInfoResp.Userid,
		Mobile:        userInfoResp.Mobile,
		Gender:        userInfoResp.Gender,
		Email:         userInfoResp.Email,
		Avatar:        userInfoResp.Avatar,
		QrCode:        userInfoResp.QrCode,
		Address:       userInfoResp.Address,
		Name:          employee.Name,
		Position:      userDetailInfo.Position,
		IsAvailable:   userDetailInfo.Status == 1,
		DepartmentIDs: employee.Department,
	}
	return userInfo, nil
}

func (c *Company) GetRedirectURL(callbackURl string) (redirectURL string) {
	return c.Work.GetOauth().GetTargetPrivateURL(callbackURl, c.AgentID)
}

func (c *Company) GetUserDetailInfo(userid string) (info *UserDetailInfo, err error) {
	token, err := c.Work.GetOauth().GetAccessToken()
	if err != nil {
		return nil, fmt.Errorf("get access token failed: %w", err)
	}

	userDetailInfoResp, err := resty.New().R().Get(fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/user/get?access_token=" + token + "&userid=" + userid))
	if err != nil {
		log.Errorf("get user info failed: %v", err)
		return nil, err
	}
	var userDetailInfo *UserDetailInfo
	_ = json.Unmarshal([]byte(userDetailInfoResp.String()), &userDetailInfo)
	if userDetailInfo.Errcode != 0 {
		log.Errorf("get user info failed: %v", userDetailInfoResp.String())
		return nil, fmt.Errorf("get user info failed")
	}
	log.Debugf("get user detail info: %s", userDetailInfoResp.String())
	c.UserDetailInfoMapping[userid] = userDetailInfo
	return userDetailInfo, nil
}

func (c *Company) formatDepartmentAndPosition(departmentIDs []int, position string) string {
	var departmentName []string
	for _, t := range departmentIDs {
		name := c.fullDepartmentName(t)
		if len(name) == 0 {
			continue
		}
		departmentName = append(departmentName, name)
	}
	var desc []string
	if dep := strings.Join(departmentName, "，"); len(dep) > 0 {
		desc = append(desc, fmt.Sprintf("部门：%s", dep))
	}
	if len(position) > 0 {
		desc = append(desc, fmt.Sprintf("职位：%s", position))
	}
	return strings.Join(desc, "\n")
}

func (c *Company) fullDepartmentName(departmentID int) string {
	departmentNames := make([]string, 0)
	for {
		department := c.DepartmentMapping[departmentID]
		if department == nil {
			break
		}
		departmentNames = append([]string{department.Name}, departmentNames...)
		if department.ParentID == 0 {
			break
		}
		departmentID = department.ParentID
	}
	return strings.Join(departmentNames, "/")
}
