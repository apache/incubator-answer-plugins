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

package wallet

import (
	"embed"
	"fmt"
	"strconv"
	"time"

	"github.com/apache/incubator-answer-plugins/connector-wallet/i18n"
	"github.com/apache/incubator-answer-plugins/util"
	"github.com/apache/incubator-answer/plugin"
	"github.com/i-lucifer/crypto"
	"golang.org/x/exp/rand"
)

//go:embed  info.yaml
var Info embed.FS

type Connector struct {
}

func init() {
	plugin.Register(&Connector{})
}

func (g *Connector) Info() plugin.Info {
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

func (g *Connector) ConnectorLogoSVG() string {
	return `PHN2ZyB0PSIxNzE3ODM1NzkwNTM1IiBjbGFzcz0iaWNvbiIgdmlld0JveD0iMCAwIDEwMjQgMTAyNCIgdmVyc2lvbj0iMS4xIiB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHAtaWQ9IjMyNjU3IiB3aWR0aD0iMjAwIiBoZWlnaHQ9IjIwMCI+PHBhdGggZD0iTTgxMC42NjY2NjcgMjk4LjY2NjY2N2gtNDIuNjY2NjY3VjI1NmExMjggMTI4IDAgMCAwLTEyOC0xMjhIMjEzLjMzMzMzM2ExMjggMTI4IDAgMCAwLTEyOCAxMjh2NTEyYTEyOCAxMjggMCAwIDAgMTI4IDEyOGg1OTcuMzMzMzM0YTEyOCAxMjggMCAwIDAgMTI4LTEyOHYtMzQxLjMzMzMzM2ExMjggMTI4IDAgMCAwLTEyOC0xMjh6TTIxMy4zMzMzMzMgMjEzLjMzMzMzM2g0MjYuNjY2NjY3YTQyLjY2NjY2NyA0Mi42NjY2NjcgMCAwIDEgNDIuNjY2NjY3IDQyLjY2NjY2N3Y0Mi42NjY2NjdIMjEzLjMzMzMzM2E0Mi42NjY2NjcgNDIuNjY2NjY3IDAgMCAxIDAtODUuMzMzMzM0eiBtNjQwIDQyNi42NjY2NjdoLTQyLjY2NjY2NmE0Mi42NjY2NjcgNDIuNjY2NjY3IDAgMCAxIDAtODUuMzMzMzMzaDQyLjY2NjY2NnogbTAtMTcwLjY2NjY2N2gtNDIuNjY2NjY2YTEyOCAxMjggMCAwIDAgMCAyNTZoNDIuNjY2NjY2djQyLjY2NjY2N2E0Mi42NjY2NjcgNDIuNjY2NjY3IDAgMCAxLTQyLjY2NjY2NiA0Mi42NjY2NjdIMjEzLjMzMzMzM2E0Mi42NjY2NjcgNDIuNjY2NjY3IDAgMCAxLTQyLjY2NjY2Ni00Mi42NjY2NjdWMzc2Ljc0NjY2N0ExMjggMTI4IDAgMCAwIDIxMy4zMzMzMzMgMzg0aDU5Ny4zMzMzMzRhNDIuNjY2NjY3IDQyLjY2NjY2NyAwIDAgMSA0Mi42NjY2NjYgNDIuNjY2NjY3eiIgcC1pZD0iMzI2NTgiPjwvcGF0aD48L3N2Zz4=`
}

func (g *Connector) ConnectorName() plugin.Translator {
	return plugin.MakeTranslator(i18n.ConnectorName)
}

func (g *Connector) ConnectorSlugName() string {
	return "wallet"
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func generateRandomString(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func (g *Connector) ConnectorSender(ctx *plugin.GinContext, receiverURL string) (redirectURL string) {
	randomString := fmt.Sprintf("%d", time.Now().Unix()) + generateRandomString(8)
	redirectURL = "/connector-wallet-auth" + "?nonce=" + randomString
	return redirectURL
}

func (g *Connector) ConnectorReceiver(ctx *plugin.GinContext, receiverURL string) (userInfo plugin.ExternalLoginUserInfo, err error) {
	message := ctx.Query("message")
	signature := ctx.Query("signature")
	address := ctx.Query("address")

	if !verifySignature(message, signature, address) {
		return userInfo, fmt.Errorf("Signature verification failed")
	}
	userInfo.ExternalID = address
	return userInfo, nil
}

func (g *Connector) ConfigFields() []plugin.ConfigField {
	return []plugin.ConfigField{}
}

func (g *Connector) ConfigReceiver(config []byte) error {
	return nil
}

func (g *Connector) guaranteeEmail(email string, accessToken string) string {
	return email
}

func verifySignature(message, signature, address string) bool {
	defer func() {
		recover()
	}()
	if len(message) != 18 {
		return false
	}

	timestamp, err := strconv.ParseInt(message[0:10], 10, 64)
	if err != nil {
		return false
	}
	if timestamp == 0 {
		return false
	}
	nowTime := time.Now().Unix()
	diffTime := nowTime - timestamp
	if diffTime < 0 || diffTime > 300 {
		return false
	}

	valid := crypto.ValidateSignature(message, signature, address)
	return valid
}
