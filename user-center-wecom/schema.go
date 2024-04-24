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

type Department struct {
	Id               int      `json:"id"`
	Name             string   `json:"name"`
	ParentID         int      `json:"parentid"`
	Order            int      `json:"order"`
	DepartmentLeader []string `json:"department_leader"`
}

type Employee struct {
	Name       string `json:"name"`
	Department []int  `json:"department"`
	Userid     string `json:"userid"`
}

type AuthUserInfoResp struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
	Userid  string `json:"userid"`
	Mobile  string `json:"mobile"`
	Gender  string `json:"gender"`
	Email   string `json:"email"`
	Avatar  string `json:"avatar"`
	QrCode  string `json:"qr_code"`
	Address string `json:"address"`
}

type UserInfo struct {
	Userid        string `json:"userid"`
	Mobile        string `json:"mobile"`
	Gender        string `json:"gender"`
	Email         string `json:"email"`
	BizEmail      string `json:"biz_mail"`
	Avatar        string `json:"avatar"`
	QrCode        string `json:"qr_code"`
	Address       string `json:"address"`
	Name          string `json:"name"`
	DepartmentIDs []int  `json:"department"`
	Position      string `json:"position"`
	IsAvailable   bool   `json:"is_available"`
}

func (u *UserInfo) GetEmail() string {
	if len(u.BizEmail) > 0 {
		return u.BizEmail
	}
	return u.Email
}

type UserDetailInfo struct {
	Errcode        int    `json:"errcode"`
	Errmsg         string `json:"errmsg"`
	Userid         string `json:"userid"`
	Name           string `json:"name"`
	Department     []int  `json:"department"`
	Position       string `json:"position"`
	Status         int    `json:"status"`
	Isleader       int    `json:"isleader"`
	EnglishName    string `json:"english_name"`
	Telephone      string `json:"telephone"`
	Enable         int    `json:"enable"`
	HideMobile     int    `json:"hide_mobile"`
	Order          []int  `json:"order"`
	MainDepartment int    `json:"main_department"`
	Alias          string `json:"alias"`
}
