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

package slack_notification

import (
	"encoding/json"

	"github.com/apache/incubator-answer-plugins/notification-slack/i18n"
	"github.com/apache/incubator-answer/plugin"
)

type NotificationConfig struct {
	Notification bool `json:"notification"`
}

func (n *Notification) ConfigFields() []plugin.ConfigField {
	return []plugin.ConfigField{
		{
			Name:        "notification",
			Type:        plugin.ConfigTypeSwitch,
			Title:       plugin.MakeTranslator(i18n.ConfigNotificationTitle),
			Description: plugin.MakeTranslator(i18n.ConfigNotificationDescription),
			UIOptions: plugin.ConfigFieldUIOptions{
				Label: plugin.MakeTranslator(i18n.ConfigNotificationLabel),
			},
			Value: n.Config.Notification,
		},
	}
}

func (n *Notification) ConfigReceiver(config []byte) error {
	c := &NotificationConfig{}
	_ = json.Unmarshal(config, c)
	n.Config = c
	return nil
}
