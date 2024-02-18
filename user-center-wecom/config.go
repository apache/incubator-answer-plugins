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
	"time"

	"github.com/apache/incubator-answer-plugins/user-center-wecom/i18n"
	"github.com/apache/incubator-answer/plugin"
)

type UserCenterConfig struct {
	CorpID       string `json:"corp_id"`
	CorpSecret   string `json:"corp_secret"`
	AgentID      string `json:"agent_id"`
	AutoSync     bool   `json:"auto_sync"`
	Notification bool   `json:"notification"`
}

func (uc *UserCenter) ConfigFields() []plugin.ConfigField {
	syncState := plugin.LoadingActionStateNone
	lastSuccessfulSyncAt := "None"
	if !uc.syncTime.IsZero() {
		syncState = plugin.LoadingActionStateComplete
		lastSuccessfulSyncAt = uc.syncTime.In(time.FixedZone("GMT", 8*3600)).Format("2006-01-02 15:04:05")
	}
	t := func(ctx *plugin.GinContext) string {
		return plugin.Translate(ctx, i18n.ConfigSyncNowDescription) + ": " + lastSuccessfulSyncAt
	}
	syncNowDesc := plugin.Translator{Fn: t}

	syncNowLabel := plugin.MakeTranslator(i18n.ConfigSyncNowLabel)

	if uc.syncing {
		syncNowLabel = plugin.MakeTranslator(i18n.ConfigSyncNowLabelForDoing)
		syncState = plugin.LoadingActionStatePending
	}

	return []plugin.ConfigField{
		{
			Name:        "auto_sync",
			Type:        plugin.ConfigTypeSwitch,
			Title:       plugin.MakeTranslator(i18n.ConfigAutoSyncTitle),
			Description: plugin.MakeTranslator(i18n.ConfigAutoSyncDescription),
			Required:    false,
			UIOptions: plugin.ConfigFieldUIOptions{
				Label: plugin.MakeTranslator(i18n.ConfigAutoSyncLabel),
			},
			Value: uc.Config.AutoSync,
		},
		{
			Name:        "sync_now",
			Type:        plugin.ConfigTypeButton,
			Title:       plugin.MakeTranslator(i18n.ConfigSyncNowTitle),
			Description: syncNowDesc,
			UIOptions: plugin.ConfigFieldUIOptions{
				Text: syncNowLabel,
				Action: &plugin.UIOptionAction{
					Url:    "/answer/admin/api/wecom/sync",
					Method: "get",
					Loading: &plugin.LoadingAction{
						Text:  plugin.MakeTranslator(i18n.ConfigSyncNowLabelForDoing),
						State: syncState,
					},
					OnComplete: &plugin.OnCompleteAction{
						ToastReturnMessage: true,
						RefreshFormConfig:  true,
					},
				},
				Variant: "outline-secondary",
			},
		},
		{
			Name:     "corp_id",
			Type:     plugin.ConfigTypeInput,
			Title:    plugin.MakeTranslator(i18n.ConfigCorpIdTitle),
			Required: true,
			UIOptions: plugin.ConfigFieldUIOptions{
				InputType: plugin.InputTypeText,
			},
			Value: uc.Config.CorpID,
		},
		{
			Name:     "corp_secret",
			Type:     plugin.ConfigTypeInput,
			Title:    plugin.MakeTranslator(i18n.ConfigCorpSecretTitle),
			Required: true,
			UIOptions: plugin.ConfigFieldUIOptions{
				InputType: plugin.InputTypePassword,
			},
			Value: uc.Config.CorpSecret,
		},
		{
			Name:     "agent_id",
			Type:     plugin.ConfigTypeInput,
			Title:    plugin.MakeTranslator(i18n.ConfigAgentIDTitle),
			Required: true,
			UIOptions: plugin.ConfigFieldUIOptions{
				InputType: plugin.InputTypeText,
			},
			Value: uc.Config.AgentID,
		},
		{
			Name:        "notification",
			Type:        plugin.ConfigTypeSwitch,
			Title:       plugin.MakeTranslator(i18n.ConfigNotificationTitle),
			Description: plugin.MakeTranslator(i18n.ConfigNotificationDescription),
			UIOptions: plugin.ConfigFieldUIOptions{
				Label: plugin.MakeTranslator(i18n.ConfigNotificationLabel),
			},
			Value: uc.Config.Notification,
		},
	}
}

func (uc *UserCenter) ConfigReceiver(config []byte) error {
	c := &UserCenterConfig{}
	_ = json.Unmarshal(config, c)
	uc.Config = c
	uc.Company = NewCompany(c.CorpID, c.CorpSecret, c.AgentID)
	uc.asyncCompany()
	return nil
}
