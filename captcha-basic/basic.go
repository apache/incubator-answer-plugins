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
	"github.com/apache/incubator-answer-plugins/util"
	"image/color"

	"github.com/apache/incubator-answer-plugins/captcha-basic/i18n"
	"github.com/apache/incubator-answer/plugin"
	"github.com/mojocn/base64Captcha"
)

//go:embed  info.yaml
var Info embed.FS

type Captcha struct {
}

func init() {
	plugin.Register(&Captcha{})
}

func (c *Captcha) Info() plugin.Info {
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

func (c *Captcha) GetConfig() (config string) {
	return ""
}

func (c *Captcha) Create() (captcha, code string) {
	driverString := base64Captcha.DriverString{
		Height:          60,
		Width:           200,
		NoiseCount:      0,
		ShowLineOptions: 2 | 4,
		Length:          4,
		Source:          "1234567890qwertyuioplkjhgfdsazxcvbnm",
		BgColor:         &color.RGBA{R: 211, G: 211, B: 211, A: 0},
		Fonts:           []string{"wqy-microhei.ttc"},
	}
	driver := driverString.ConvertFonts()
	_, content, answer := driverString.ConvertFonts().GenerateIdQuestionAnswer()
	item, _ := driver.DrawCaptcha(content)
	return item.EncodeB64string(), answer
}

func (c *Captcha) Verify(captcha, userInput string) (pass bool) {
	if len(captcha) == 0 || len(userInput) == 0 {
		return false
	}
	return captcha == userInput
}
