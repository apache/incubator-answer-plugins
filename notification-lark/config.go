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

package lark

import (
	"context"
	"encoding/json"

	"github.com/apache/incubator-answer-plugins/notification-lark/i18n"
	"github.com/apache/incubator-answer/plugin"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkCore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkWebSocket "github.com/larksuite/oapi-sdk-go/v3/ws"
	"github.com/segmentfault/pacman/log"
)

type NotificationConfig struct {
	Version           string `json:"version"`
	AppID             string `json:"app_id"`
	AppSecret         string `json:"app_secret"`
	VerificationToken string `json:"verification_token"`
	EventEncryptKey   string `json:"event_encrypt_key"`
}

func (n *NotificationConfig) GetVersion() string {
	if n == nil {
		return ""
	}

	return n.Version
}

func (n *NotificationConfig) GetAppID() string {
	if n == nil {
		return ""
	}

	return n.AppID
}

func (n *NotificationConfig) GetAppSecret() string {
	if n == nil {
		return ""
	}

	return n.AppSecret
}

func (n *NotificationConfig) GetVerificationToken() string {
	if n == nil {
		return ""
	}

	return n.VerificationToken
}

func (n *NotificationConfig) GetEventEncryptKey() string {
	if n == nil {
		return ""
	}

	return n.EventEncryptKey
}

func (n *Notification) ConfigFields() []plugin.ConfigField {
	return []plugin.ConfigField{
		{
			Name:        "version",
			Type:        plugin.ConfigTypeSelect,
			Title:       plugin.MakeTranslator(i18n.ConfigVersionTitle),
			Description: plugin.MakeTranslator(i18n.ConfigVersionDescription),
			Required:    true,
			Value:       n.config.GetVersion(),
			Options: []plugin.ConfigFieldOption{
				{
					Label: plugin.MakeTranslator(i18n.ConfigVersionOptionsFeishu),
					Value: i18n.ConfigVersionOptionsFeishu,
				},
				{
					Label: plugin.MakeTranslator(i18n.ConfigVersionOptionsLark),
					Value: i18n.ConfigVersionOptionsLark,
				},
			},
		},
		{
			Name:        "app_id",
			Type:        plugin.ConfigTypeInput,
			Title:       plugin.MakeTranslator(i18n.ConfigAppIdTitle),
			Description: plugin.MakeTranslator(i18n.ConfigAppIdDescription),
			Required:    true,
			Value:       n.config.GetAppID(),
		},
		{
			Name:        "app_secret",
			Type:        plugin.ConfigTypeInput,
			Title:       plugin.MakeTranslator(i18n.ConfigAppSecretTitle),
			Description: plugin.MakeTranslator(i18n.ConfigAppSecretDescription),
			Required:    true,
			Value:       n.config.GetAppSecret(),
			UIOptions: plugin.ConfigFieldUIOptions{
				InputType: plugin.InputTypePassword,
			},
		},
		{
			Name:        "event_encrypt_key",
			Type:        plugin.ConfigTypeInput,
			Title:       plugin.MakeTranslator(i18n.ConfigEventEncryptKeyTitle),
			Description: plugin.MakeTranslator(i18n.ConfigEventEncryptKeyDescription),
			Required:    false,
			Value:       n.config.GetEventEncryptKey(),
			UIOptions: plugin.ConfigFieldUIOptions{
				InputType: plugin.InputTypePassword,
			},
		},
		{
			Name:        "verification_token",
			Type:        plugin.ConfigTypeInput,
			Title:       plugin.MakeTranslator(i18n.ConfigVerificationTokenTitle),
			Description: plugin.MakeTranslator(i18n.ConfigVerificationTokenDescription),
			Required:    false,
			Value:       n.config.GetVerificationToken(),
			UIOptions: plugin.ConfigFieldUIOptions{
				InputType: plugin.InputTypePassword,
			},
		},
	}
}

type LarkLogger struct{}

func (l LarkLogger) Debug(ctx context.Context, args ...interface{}) {
	log.Debug(args...)
}

func (l LarkLogger) Info(ctx context.Context, args ...interface{}) {
	log.Info(args...)
}

func (l LarkLogger) Warn(ctx context.Context, args ...interface{}) {
	log.Warn(args...)
}

func (l LarkLogger) Error(ctx context.Context, args ...interface{}) {
	log.Error(args...)
}

func (n *Notification) ConfigReceiver(config []byte) error {
	c := &NotificationConfig{}
	if err := json.Unmarshal(config, c); err != nil {
		return err
	}

	n.config = c

	larkDomain := lark.FeishuBaseUrl
	if c.Version == i18n.ConfigVersionOptionsLark {
		larkDomain = lark.LarkBaseUrl
	}

	n.client = &LarkClient{
		ws: larkWebSocket.NewClient(
			n.config.GetAppID(),
			n.config.GetAppSecret(),
			larkWebSocket.WithDomain(larkDomain),
			larkWebSocket.WithLogger(LarkLogger{}),
			larkWebSocket.WithLogLevel(larkCore.LogLevelDebug),
			larkWebSocket.WithEventHandler(n.LarkWsEventHub()),
		),
		http: lark.NewClient(n.config.GetAppID(), n.config.GetAppSecret()),
	}
	go n.client.Start()
	return nil
}
