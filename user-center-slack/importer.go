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

package slack_user_center

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/apache/incubator-answer/plugin"
	"github.com/gin-gonic/gin"
	"github.com/segmentfault/pacman/log"
)

func (uc *UserCenter) parseText(text string) (string, string, []string, error) {
	re := regexp.MustCompile(`\[(.*?)\]`)
	matches := re.FindAllStringSubmatch(text, -1)

	if len(matches) != 3 {
		return "", "", nil, fmt.Errorf("text field does not conform to the required format")
	}

	part1 := matches[0][1]
	part2 := matches[1][1]
	rawTags := strings.Split(matches[2][1], ",")

	var tags []string
	for _, tag := range rawTags {
		if tag != "" {
			tags = append(tags, tag)
		}
	}

	// if part1 or part2 or tags in empty return error
	if part1 == "" || part2 == "" || len(tags) == 0 {
		return "", "", nil, fmt.Errorf("text field does not be empty")
	}
	return part1, part2, tags, nil
}
func getSlackUserEmail(userID, token string) (string, error) {
	url := fmt.Sprintf("https://slack.com/api/users.info?user=%s", userID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var userResponse SlackUserResponse
	if err := json.Unmarshal(body, &userResponse); err != nil {
		return "", err
	}
	if !userResponse.Ok {
		return "", fmt.Errorf("failed to get user info from Slack")
	}

	return userResponse.User.Profile.Email, nil
}

func (uc *UserCenter) verifySlackRequest(ctx *gin.Context) error {
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		return fmt.Errorf("could not read request body: %v", err)
	}
	timestamp := ctx.GetHeader("X-Slack-Request-Timestamp")
	slackSignature := ctx.GetHeader("X-Slack-Signature")

	// check the timestamp validity in 5 minutes
	ts, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid timestamp: %v", err)
	}
	if time.Now().Unix()-ts > 60*5 {
		return fmt.Errorf("timestamp is too old")
	}
	// Reset the request body for further processing
	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	sigBaseString := fmt.Sprintf("v0:%s:%s", timestamp, string(body))

	h := hmac.New(sha256.New, []byte(uc.Config.SigningSecret))
	h.Write([]byte(sigBaseString))
	computedSignature := "v0=" + hex.EncodeToString(h.Sum(nil))

	if !hmac.Equal([]byte(computedSignature), []byte(slackSignature)) {
		return fmt.Errorf("invalid signature")
	}

	return nil
}
func (uc *UserCenter) GetQuestion(ctx *gin.Context) (questionInfo plugin.QuestionImporterInfo, err error) {
	questionInfo = plugin.QuestionImporterInfo{}

	err = uc.verifySlackRequest(ctx)
	if err != nil {
		return questionInfo, err
	}

	text := ctx.PostForm("text")
	part1, part2, tags, err := uc.parseText(text)
	if err != nil {
		return questionInfo, err
	}

	questionInfo.Title = part1
	questionInfo.Content = part2
	questionInfo.Tags = tags
	userID := ctx.PostForm("user_id")

	token := uc.SlackClient.AccessToken
	email, err := getSlackUserEmail(userID, token)
	if err != nil {
		return questionInfo, err
	}

	questionInfo.UserEmail = email
	return questionInfo, nil
}

func (uc *UserCenter) SlashCommand(ctx *gin.Context) {
	body, _ := io.ReadAll(ctx.Request.Body)
	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(body))
	cmd := ctx.PostForm("command")
	if cmd != "/ask" {
		log.Errorf("error: Invalid command")
		ctx.JSON(http.StatusBadRequest, gin.H{"text": "Invalid command"})
		return
	}
	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(body))
	err := uc.verifySlackRequest(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"text": "Slack request verification faild"})
		log.Errorf("error: %v", err)
		return
	}
	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(body))
	questionInfo, err := uc.GetQuestion(ctx)
	if err != nil {
		log.Errorf("error: %v", err)
		ctx.JSON(200, gin.H{"text": err.Error()})
		return
	}
	if uc.importerFunc == nil {
		log.Errorf("error: importerFunc is not initialized")
		return
	}
	err = uc.importerFunc.AddQuestion(ctx, questionInfo)
	if err != nil {
		log.Errorf("error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"text": "Failed to add question"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"text": "Question has been added successfully"})
}

func (uc *UserCenter) RegisterImporterFunc(ctx context.Context, importerFunc plugin.ImporterFunc) {
	uc.importerFunc = importerFunc
}
