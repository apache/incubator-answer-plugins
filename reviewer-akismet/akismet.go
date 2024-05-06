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

package basic

import (
	"github.com/apache/incubator-answer/plugin"
	"github.com/go-resty/resty/v2"
	"github.com/segmentfault/pacman/log"
)

const (
	commentCheckURL = "https://rest.akismet.com/1.1/comment-check"
)

func (r *Reviewer) RequestAkismetToCheck(content *plugin.ReviewContent) (isSpam bool, err error) {
	req := make(map[string]string)
	req["blog"] = plugin.SiteURL()
	req["user_ip"] = content.IP
	req["user_agent"] = content.UserAgent
	req["comment_content"] = content.Title + "\n" + content.Content
	req["comment_type"] = "comment"
	req["is_test"] = "false"
	// This is for test if the akismet is available.
	if content.Title == "akismet-guaranteed-spam" {
		req["comment_content"] = content.Title
	}

	log.Debugf("request akismet: %+v", req)

	req["api_key"] = r.Config.APIKey

	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetFormData(req).
		Post(commentCheckURL)

	if err != nil {
		log.Errorf("request akismet failed: %v", err)
		return false, err
	}

	if resp.StatusCode() != 200 {
		log.Errorf("request akismet failed: %v", resp.String())
		return false, nil
	}

	log.Debugf("akismet response: %v, content title is %s", resp.String(), content.Title)

	if resp.String() == "true" {
		return true, nil
	}
	return false, nil
}
