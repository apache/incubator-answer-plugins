package slack_user_center

import (
	"bytes"
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
func (uc *UserCenter) GetQuestion(ctx *gin.Context) (questionInfo *plugin.QuestionImporterInfo, err error) {
	questionInfo = &plugin.QuestionImporterInfo{}

	err = uc.verifySlackRequest(ctx)
	if err != nil {
		return nil, err
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
