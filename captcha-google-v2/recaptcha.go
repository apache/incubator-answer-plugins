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

package recaptcha

import (
	"embed"
	"github.com/apache/incubator-answer-plugins/util"
	"io"
	"net/http"
	"time"

	"github.com/apache/incubator-answer-plugins/captcha-google-v2/i18n"
	"github.com/apache/incubator-answer/plugin"

	"encoding/json"

	"github.com/segmentfault/pacman/log"
)

//go:embed  info.yaml
var Info embed.FS

type Captcha struct {
	Config *CaptchaConfig
}

type CaptchaConfig struct {
	SiteKey            string `json:"site_key"`
	SecretKey          string `json:"secret_key"`
	SiteVerifyEndpoint string `json:"site_verify_endpoint"`
}

type GoogleCaptchaResponse struct {
	Success    bool     `json:"success"`
	ErrorCodes []string `json:"error-codes"`
}

func init() {
	plugin.Register(&Captcha{
		Config: &CaptchaConfig{},
	})
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
	data, _ := json.Marshal(map[string]interface{}{
		"key": c.Config.SiteKey,
	})
	return string(data)
}

func (c *Captcha) Create() (captcha, code string) {
	return "", ""
}

func (c *Captcha) Verify(captcha, userInput string) (pass bool) {
	if len(userInput) == 0 {
		return false
	}
	cli := &http.Client{}
	cli.Timeout = 10 * time.Second
	siteVerifyEndpoint := c.Config.SiteVerifyEndpoint
	if siteVerifyEndpoint == "" {
		siteVerifyEndpoint = "https://www.google.com/recaptcha/api/siteverify"
	}
	resp, err := cli.PostForm(siteVerifyEndpoint, map[string][]string{
		"secret":   {c.Config.SecretKey},
		"response": {userInput},
	})
	if err != nil {
		log.Errorf("verify captcha error %s", err.Error())
		return false
	}
	defer resp.Body.Close()
	all, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("verify captcha, read body error %s", err.Error())
		return false
	}
	r := &GoogleCaptchaResponse{}
	_ = json.Unmarshal(all, r)
	if r.Success {
		return true
	}
	log.Debugf("user input is wrong %s, user input is %s", string(all), userInput)
	return false
}

func (c *Captcha) ConfigFields() []plugin.ConfigField {
	return []plugin.ConfigField{
		{
			Name:        "site_key",
			Type:        plugin.ConfigTypeInput,
			Title:       plugin.MakeTranslator(i18n.ConfigSiteKeyTitle),
			Description: plugin.MakeTranslator(i18n.ConfigSiteKeyDescription),
			Required:    true,
			UIOptions: plugin.ConfigFieldUIOptions{
				InputType: plugin.InputTypeText,
			},
			Value: c.Config.SiteKey,
		},
		{
			Name:        "secret_key",
			Type:        plugin.ConfigTypeInput,
			Title:       plugin.MakeTranslator(i18n.ConfigSecretKeyTitle),
			Description: plugin.MakeTranslator(i18n.ConfigSecretKeyDescription),
			Required:    true,
			UIOptions: plugin.ConfigFieldUIOptions{
				InputType: plugin.InputTypeText,
			},
			Value: c.Config.SecretKey,
		},
		{
			Name:        "site_verify_endpoint",
			Type:        plugin.ConfigTypeInput,
			Title:       plugin.MakeTranslator(i18n.ConfigSiteVerifyEndpointTitle),
			Description: plugin.MakeTranslator(i18n.ConfigSiteVerifyEndpointDescription),
			Required:    false,
			UIOptions: plugin.ConfigFieldUIOptions{
				InputType: plugin.InputTypeText,
			},
			Value: c.Config.SiteVerifyEndpoint,
		},
	}
}

func (c *Captcha) ConfigReceiver(config []byte) error {
	conf := &CaptchaConfig{}
	_ = json.Unmarshal(config, conf)
	c.Config = conf
	return nil
}
