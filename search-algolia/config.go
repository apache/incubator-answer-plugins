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

package algolia

import (
	_ "embed"
	"encoding/json"
	"github.com/apache/incubator-answer-plugins/search-algolia/i18n"

	"github.com/apache/incubator-answer/plugin"
)

var (
	NewestIndex = "newest"
	ActiveIndex = "active"
	ScoreIndex  = "score"
)

type AlgoliaSearchConfig struct {
	APPID        string `json:"app_id"`
	PublicAPIKey string `json:"public_api_key"`
	APIKey       string `json:"api_key"`
	Index        string `json:"index"`
	ShowLogo     bool   `json:"show_logo"`
}

// ConfigFields return config fields
func (s *SearchAlgolia) ConfigFields() []plugin.ConfigField {
	return []plugin.ConfigField{
		{
			Name:        "app_id",
			Type:        plugin.ConfigTypeInput,
			Title:       plugin.MakeTranslator(i18n.ConfigAPPIDTitle),
			Description: plugin.MakeTranslator(i18n.ConfigAPPIDDescription),
			Required:    true,
			UIOptions: plugin.ConfigFieldUIOptions{
				InputType: plugin.InputTypeText,
			},
			Value: s.Config.APPID,
		},
		{
			Name:        "public_api_key",
			Type:        plugin.ConfigTypeInput,
			Title:       plugin.MakeTranslator(i18n.ConfigPublicAPIKeyTitle),
			Description: plugin.MakeTranslator(i18n.ConfigPublicAPIKeyDescription),
			Required:    true,
			UIOptions: plugin.ConfigFieldUIOptions{
				InputType: plugin.InputTypePassword,
			},
			Value: s.Config.PublicAPIKey,
		},
		{
			Name:        "api_key",
			Type:        plugin.ConfigTypeInput,
			Title:       plugin.MakeTranslator(i18n.ConfigAPIKeyTitle),
			Description: plugin.MakeTranslator(i18n.ConfigAPIKeyDescription),
			Required:    true,
			UIOptions: plugin.ConfigFieldUIOptions{
				InputType: plugin.InputTypePassword,
			},
			Value: s.Config.APIKey,
		},
		{
			Name:        "index",
			Type:        plugin.ConfigTypeInput,
			Title:       plugin.MakeTranslator(i18n.ConfigIndexTitle),
			Description: plugin.MakeTranslator(i18n.ConfigIndexDescription),
			Required:    true,
			UIOptions: plugin.ConfigFieldUIOptions{
				InputType: plugin.InputTypeText,
			},
			Value: s.Config.Index,
		},
		{
			Name:        "show_logo",
			Type:        plugin.ConfigTypeSwitch,
			Title:       plugin.MakeTranslator(i18n.ConfigShowLogoTitle),
			Description: plugin.MakeTranslator(i18n.ConfigShowLogoDescription),
			UIOptions: plugin.ConfigFieldUIOptions{
				Label: plugin.MakeTranslator(i18n.ConfigShowLogoLabel),
			},
			Value: s.Config.ShowLogo,
		},
	}
}

// ConfigReceiver receive config from admin
func (s *SearchAlgolia) ConfigReceiver(config []byte) error {
	c := &AlgoliaSearchConfig{}
	_ = json.Unmarshal(config, c)
	s.Config = c
	err := s.connect()
	if err != nil {
		return err
	}
	// if config update, re-init settings
	err = s.initSettings()
	return err
}
