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
	"embed"
	"encoding/json"
	"github.com/apache/incubator-answer-plugins/util"

	"github.com/apache/incubator-answer-plugins/reviewer-baidu/i18n"
	"github.com/apache/incubator-answer/plugin"
	"github.com/lufei/baidu-golang-sdk/aip/censor"
	myI18n "github.com/segmentfault/pacman/i18n"
	"github.com/segmentfault/pacman/log"
)

//go:embed  info.yaml
var Info embed.FS

type Reviewer struct {
	Config *ReviewerConfig
}

type ReviewerConfig struct {
	APIKey        string `json:"api_key"`
	SecretKey     string `json:"secret_key"`
	SpamFiltering string `json:"span_filtering"`
}

func init() {
	plugin.Register(&Reviewer{
		Config: &ReviewerConfig{},
	})
}

func (r *Reviewer) Info() plugin.Info {
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

func (r *Reviewer) Review(content *plugin.ReviewContent) (result *plugin.ReviewResult) {
	result = &plugin.ReviewResult{Approved: true}
	if len(r.Config.APIKey) == 0 {
		return result
	}
	// If the author is admin, no need to review
	if content.Author.Role > 1 {
		return result
	}

	client := censor.NewClient(r.Config.APIKey, r.Config.SecretKey)
	textCensorResult, err := client.TextCensor(content.Title+"\n"+content.Content, content.IP)
	if err != nil {
		log.Errorf("Request baidu to check failed: %v", err)
		return handleReviewError(content, plugin.ReviewStatusNeedReview)
	}

	var jsonMap map[string]interface{}
	err = json.Unmarshal([]byte(textCensorResult), &jsonMap)
	if err != nil {
		return handleReviewError(content, plugin.ReviewStatusNeedReview)
	}

	if conclusionType, ok := jsonMap["conclusionType"].(float64); ok {
		if conclusionType == 1.0 {
			return result
		}
	}

	if r.Config.SpamFiltering == "delete" {
		return handleReviewError(content, plugin.ReviewStatusDeleteDirectly)
	}

	return handleReviewError(content, plugin.ReviewStatusNeedReview)
}

func (r *Reviewer) ConfigFields() []plugin.ConfigField {
	return []plugin.ConfigField{
		{
			Name:        "api_key",
			Type:        plugin.ConfigTypeInput,
			Title:       plugin.MakeTranslator(i18n.ConfigAPIKeyLabel),
			Description: plugin.MakeTranslator(i18n.ConfigAPIKeyDescription),
			Required:    true,
			UIOptions: plugin.ConfigFieldUIOptions{
				InputType: plugin.InputTypeText,
				Label:     plugin.MakeTranslator(i18n.ConfigAPIKeyLabel),
			},
			Value: r.Config.APIKey,
		},
		{
			Name:        "secret_key",
			Type:        plugin.ConfigTypeInput,
			Title:       plugin.MakeTranslator(i18n.ConfigSecretKeyTitle),
			Description: plugin.MakeTranslator(i18n.ConfigSecretKeyDescription),
			Required:    true,
			UIOptions: plugin.ConfigFieldUIOptions{
				InputType: plugin.InputTypeText,
				Label:     plugin.MakeTranslator(i18n.ConfigSecretKeyLabel),
			},
			Value: r.Config.SecretKey,
		},
		{
			Name:      "span_filtering",
			Type:      plugin.ConfigTypeSelect,
			Title:     plugin.MakeTranslator(i18n.ConfigSpanFilteringTitle),
			Required:  false,
			UIOptions: plugin.ConfigFieldUIOptions{},
			Value:     r.Config.SpamFiltering,
			Options: []plugin.ConfigFieldOption{
				{
					Value: "review",
					Label: plugin.MakeTranslator(i18n.ConfigSpanFilteringReview),
				},
				{
					Value: "delete",
					Label: plugin.MakeTranslator(i18n.ConfigSpanFilteringDelete),
				},
			},
		},
	}
}

func handleReviewError(content *plugin.ReviewContent, ReviewStatus plugin.ReviewStatus) *plugin.ReviewResult {
	return &plugin.ReviewResult{
		Approved:     false,
		ReviewStatus: ReviewStatus,
		Reason:       plugin.TranslateWithData(myI18n.Language(content.Language), i18n.CommentNeedReview, nil),
	}
}

func (r *Reviewer) ConfigReceiver(config []byte) error {
	c := &ReviewerConfig{}
	_ = json.Unmarshal(config, c)
	r.Config = c
	return nil
}
