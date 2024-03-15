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
	"encoding/json"

	"github.com/apache/incubator-answer-plugins/reviewer-basic/i18n"
	"github.com/apache/incubator-answer/plugin"
	myI18n "github.com/segmentfault/pacman/i18n"
)

type Reviewer struct {
	Config *ReviewerConfig
}

type ReviewerConfig struct {
	PostNeedReview bool `json:"review_post"`
}

func init() {
	plugin.Register(&Reviewer{
		Config: &ReviewerConfig{},
	})
}

func (r *Reviewer) Info() plugin.Info {
	return plugin.Info{
		Name:        plugin.MakeTranslator(i18n.InfoName),
		SlugName:    "basic_reviewer",
		Description: plugin.MakeTranslator(i18n.InfoDescription),
		Author:      "answerdev",
		Version:     "1.0.0",
		Link:        "https://github.com/apache/incubator-answer-plugins/tree/main/reviewer-basic",
	}
}

func (r *Reviewer) Review(content *plugin.ReviewContent) (result *plugin.ReviewResult) {
	result = &plugin.ReviewResult{Approved: true}
	// If the author is admin, no need to review
	if content.Author.Role > 1 {
		return result
	}
	if content.Author.ApprovedQuestionAmount+content.Author.ApprovedAnswerAmount > 1 {
		return result
	}
	return &plugin.ReviewResult{
		Approved: false,
		Reason:   plugin.TranslateWithData(myI18n.Language(content.Language), i18n.CommentNeedReview, nil),
	}
}

func (r *Reviewer) ConfigFields() []plugin.ConfigField {
	return []plugin.ConfigField{
		{
			Name:        "review_post",
			Type:        plugin.ConfigTypeSwitch,
			Title:       plugin.MakeTranslator(i18n.ConfigReviewPostTitle),
			Description: plugin.MakeTranslator(i18n.ConfigReviewPostDescription),
			UIOptions: plugin.ConfigFieldUIOptions{
				Label: plugin.MakeTranslator(i18n.ConfigReviewPostLabel),
			},
			Value: r.Config.PostNeedReview,
		},
	}
}

func (r *Reviewer) ConfigReceiver(config []byte) error {
	c := &ReviewerConfig{}
	_ = json.Unmarshal(config, c)
	r.Config = c
	return nil
}
