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

package lark_test

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	answer "github.com/apache/incubator-answer-plugins/notification-lark"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkIM "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

func TestLarkCardMessage(t *testing.T) {
	appId := os.Getenv("LARK_APP_ID")
	appSecret := os.Getenv("LARK_APP_SECRET")
	openId := os.Getenv("LARK_OPEN_ID")
	if appId == "" || appSecret == "" || openId == "" {
		t.Skip("LARK_APP_ID, LARK_APP_SECRET, LARK_OPEN_ID are required")
	}

	larkClient := lark.NewClient(appId, appSecret)

	contentData, _ := json.Marshal(answer.Card{
		Config: &answer.Config{
			WidthMode:            "compact",
			UseCustomTranslation: answer.PtrBool(true),
			EnableForward:        answer.PtrBool(false),
		},
		Header: &answer.Header{
			Title: &answer.Text{
				Tag: "plain_text",
				I18n: &answer.I18n{
					ZhCn: "新通知",
					EnUs: "New Notification",
				},
			},
			UdIcon: &answer.Icon{
				Tag:   "icon",
				Token: "bell_outlined",
				Color: "blue",
			},
			Template: answer.ThemeGreen,
		},
		I18nElements: &answer.I18nElements{
			ZhCn: []answer.ColumnSet{
				{
					Show: &answer.Show{
						Tag:      "column_set",
						FlexMode: "flex_mode",
						Columns: []answer.Column{
							{
								Elements: []answer.Element{
									{
										PlainText: &answer.PlainText{
											Tag: "div",
											Text: &answer.Text{
												Tag:     "lark_md",
												Content: "[@Answer](https://answer.apache.org/) 创建了问题 [如何使用 Answer?](https://answer.apache.org/docs/)",
											},
										},
									},
								},
							},
						},
					},
				},
				{
					Action: &answer.Action{
						Tag: "action",
						Actions: []*answer.Button{
							{
								Width: "fill",
								Text: &answer.Text{
									Tag:     "plain_text",
									Content: "查看详情",
									Icon: &answer.Icon{
										Tag:   "standard_icon",
										Token: "link-copy_outlined",
									},
								},
								Behaviors: []answer.Behavior{
									{
										Type:       "open_url",
										DefaultURL: "https://answer.apache.org/docs/",
									},
								},
							},
						},
					},
				},
			},
			EnUs: []answer.ColumnSet{
				{
					Show: &answer.Show{
						FlexMode: "flex_mode",
						Columns: []answer.Column{
							{
								Elements: []answer.Element{
									{
										PlainText: &answer.PlainText{
											Tag: "div",
											Text: &answer.Text{
												Tag:     "lark_md",
												Content: "[@Answer](https://answer.apache.org/) created a question [How to use Answer?](https://answer.apache.org/docs/)",
											},
										},
									},
								},
							},
						},
					},
				},
				{
					Action: &answer.Action{
						Tag: "action",
						Actions: []*answer.Button{
							{
								Width: "fill",
								Text: &answer.Text{
									Tag:     "plain_text",
									Content: "View details",
									Icon: &answer.Icon{
										Tag:   "standard_icon",
										Token: "link-copy_outlined",
									},
								},
								Behaviors: []answer.Behavior{
									{
										Type:       "open_url",
										DefaultURL: "https://answer.apache.org/docs/",
									},
								},
							},
						},
					},
				},
			},
		},
	})

	t.Logf("Content: %s", string(contentData))
	req := larkIM.NewCreateMessageReqBuilder().
		ReceiveIdType("open_id").
		Body(larkIM.NewCreateMessageReqBodyBuilder().
			ReceiveId(openId).
			MsgType("interactive").
			Content(string(contentData)).
			Build()).
		Build()

	resp, err := larkClient.Im.Message.Create(context.Background(), req)
	if err != nil {
		t.Errorf("Failed to send message: %v", err)
	}
	if !resp.Success() {
		t.Errorf("Failed to send message: %v", resp.Error())
	}
}
