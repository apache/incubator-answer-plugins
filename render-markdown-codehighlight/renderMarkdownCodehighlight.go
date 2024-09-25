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

package render_markdown_codehighlight

import (
	"embed"
	"encoding/json"
	"github.com/apache/incubator-answer-plugins/render-markdown-codehighlight/i18n"
	"github.com/apache/incubator-answer-plugins/util"
	"github.com/apache/incubator-answer/plugin"
	"github.com/gin-gonic/gin"
	"log"
	"strings"
)

//go:embed info.yaml
var Info embed.FS

type Render struct {
	Config *RenderConfig
}

type RenderConfig struct {
	SelectTheme string `json:"select_theme"`
}

func init() {
	plugin.Register(&Render{
		Config: &RenderConfig{},
	})
}

func (r *Render) Info() plugin.Info {
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

func (r *Render) ConfigFields() []plugin.ConfigField {
	themeOptions := make([]plugin.ConfigFieldOption, len(ThemeList))

	for i, theme := range ThemeList {
		// Split the theme string by the hyphen and take the first part
		themeParts := strings.Split(theme, "-")
		themeValue := themeParts[0]

		themeOptions[i] = plugin.ConfigFieldOption{
			Value: themeValue, // Use the first part as the Value
			Label: plugin.MakeTranslator(theme),
		}
	}

	return []plugin.ConfigField{
		{
			Name:      "select_theme",
			Type:      plugin.ConfigTypeSelect,
			Title:     plugin.MakeTranslator(i18n.ConfigCssFilteringTitle),
			Required:  false,
			UIOptions: plugin.ConfigFieldUIOptions{},
			Value:     r.Config.SelectTheme,
			Options:   themeOptions,
		},
	}
}

func (r *Render) ConfigReceiver(config []byte) error {
	c := &RenderConfig{}
	_ = json.Unmarshal(config, c)
	r.Config = c
	log.Println("Received theme:", r.Config.SelectTheme)
	return nil
}

func (r *Render) GetRenderConfig(ctx *gin.Context) (renderConfig *plugin.RenderConfig) {
	log.Println("Current theme:", r.Config.SelectTheme)
	renderConfig = &plugin.RenderConfig{
		SelectTheme: r.Config.SelectTheme,
	}
	return
}
