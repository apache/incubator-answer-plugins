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

package embed_basic

import (
	"embed"
	"encoding/json"
	"github.com/apache/incubator-answer-plugins/util"
	"github.com/gin-gonic/gin"

	"github.com/apache/incubator-answer-plugins/embed-basic/i18n"
	"github.com/apache/incubator-answer/plugin"
)

//go:embed  info.yaml
var Info embed.FS

//go:embed components
var Build embed.FS

type Embed struct {
	Config *EmbedConfig
}

type EmbedConfig struct {
	Codepen    bool `json:"codepen"`
	Dropbox    bool `json:"dropbox"`
	Excalidraw bool `json:"excalidraw"`
	Figma      bool `json:"figma"`
	Githubgist bool `json:"githubgist"`
	Jsfiddle   bool `json:"jsfiddle"`
	Loom       bool `json:"loom"`
	Twitter    bool `json:"twitter"`
	Youtube    bool `json:"youtube"`
}

func init() {
	plugin.Register(&Embed{
		Config: &EmbedConfig{},
	})
}

func (e *Embed) Info() plugin.Info {
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

func (e *Embed) ConfigFields() []plugin.ConfigField {
	return []plugin.ConfigField{
		{
			Name:  "CodePen",
			Type:  plugin.ConfigTypeSwitch,
			Title: plugin.MakeTranslator(i18n.ConfigEmbedsTitle),
			UIOptions: plugin.ConfigFieldUIOptions{
				Label:          plugin.MakeTranslator(i18n.ConfigOptionCodepen),
				FieldClassName: "mb-0",
			},
			Value: e.Config.Codepen,
		},
		{
			Name: "Dropbox",
			Type: plugin.ConfigTypeSwitch,
			UIOptions: plugin.ConfigFieldUIOptions{
				Label:          plugin.MakeTranslator(i18n.ConfigOptionDropbox),
				FieldClassName: "mb-0",
			},
			Value: e.Config.Dropbox,
		},
		{
			Name: "Excalidraw",
			Type: plugin.ConfigTypeSwitch,
			UIOptions: plugin.ConfigFieldUIOptions{
				Label:          plugin.MakeTranslator(i18n.ConfigOptionExcalidraw),
				FieldClassName: "mb-0",
			},
			Value: e.Config.Excalidraw,
		},
		{
			Name: "Figma",
			Type: plugin.ConfigTypeSwitch,
			UIOptions: plugin.ConfigFieldUIOptions{
				Label:          plugin.MakeTranslator(i18n.ConfigOptionFigma),
				FieldClassName: "mb-0",
			},
			Value: e.Config.Figma,
		},
		{
			Name: "GithubGist",
			Type: plugin.ConfigTypeSwitch,
			UIOptions: plugin.ConfigFieldUIOptions{
				Label:          plugin.MakeTranslator(i18n.ConfigOptionGithubgist),
				FieldClassName: "mb-0",
			},
			Value: e.Config.Githubgist,
		},
		{
			Name: "JSFiddle",
			Type: plugin.ConfigTypeSwitch,
			UIOptions: plugin.ConfigFieldUIOptions{
				Label:          plugin.MakeTranslator(i18n.ConfigOptionJsfiddle),
				FieldClassName: "mb-0",
			},
			Value: e.Config.Jsfiddle,
		},
		{
			Name: "Loom",
			Type: plugin.ConfigTypeSwitch,
			UIOptions: plugin.ConfigFieldUIOptions{
				Label:          plugin.MakeTranslator(i18n.ConfigOptionLoom),
				FieldClassName: "mb-0",
			},
			Value: e.Config.Loom,
		},
		{
			Name: "Twitter",
			Type: plugin.ConfigTypeSwitch,
			UIOptions: plugin.ConfigFieldUIOptions{
				Label:          plugin.MakeTranslator(i18n.ConfigOptionTwitter),
				FieldClassName: "mb-0",
			},
			Value: e.Config.Twitter,
		},
		{
			Name: "YouTube",
			Type: plugin.ConfigTypeSwitch,
			UIOptions: plugin.ConfigFieldUIOptions{
				Label: plugin.MakeTranslator(i18n.ConfigOptionYoutube),
			},
			Value:       e.Config.Youtube,
			Description: plugin.MakeTranslator(i18n.ConfigEmbedsDescription),
		},
	}
}

func (e *Embed) ConfigReceiver(config []byte) error {
	c := &EmbedConfig{}
	_ = json.Unmarshal(config, c)
	e.Config = c
	return nil
}

// GetEmbedConfigs get embed configs
func (e *Embed) GetEmbedConfigs(ctx *gin.Context) (embedConfigs []*plugin.EmbedConfig, err error) {
	embedConfigs = make([]*plugin.EmbedConfig, 0)
	for _, field := range e.ConfigFields() {
		embedConfigs = append(embedConfigs, &plugin.EmbedConfig{
			Platform: field.Name,
			Enable:   field.Value.(bool),
		})
	}
	return
}
