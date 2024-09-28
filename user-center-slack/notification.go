package slack_user_center

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	slackI18n "github.com/Anan1225/incubator-answer-plugins/user-center-slack/i18n"
	"github.com/apache/incubator-answer/plugin"
	"github.com/segmentfault/pacman/i18n"
	"github.com/segmentfault/pacman/log"
)

// GetNewQuestionSubscribers returns the subscribers of the new question notification
func (uc *UserCenter) GetNewQuestionSubscribers() (userIDs []string) {
	for userID, conf := range uc.UserConfigCache.userConfigMapping {
		if conf.AllNewQuestions {
			userIDs = append(userIDs, userID)
		}
	}
	return userIDs
}

// Notify sends a notification to the user using Slack
func (uc *UserCenter) Notify(msg plugin.NotificationMessage) {
	log.Debugf("try to send notification %+v", msg)

	if !uc.Config.Notification {
		return
	}

	// get user config
	userConfig, err := uc.getUserConfig(msg.ReceiverUserID)
	if err != nil {
		log.Errorf("get user config failed: %v", err)
		return
	}
	if userConfig == nil {
		log.Debugf("user %s has no config", msg.ReceiverUserID)
		return
	}

	// check if the notification is enabled
	switch msg.Type {
	case plugin.NotificationNewQuestion:
		if !userConfig.AllNewQuestions {
			log.Debugf("user %s not config the new question", msg.ReceiverUserID)
			return
		}
	case plugin.NotificationNewQuestionFollowedTag:
		if !userConfig.NewQuestionsForFollowingTags {
			log.Debugf("user %s not config the new question followed tag", msg.ReceiverUserID)
			return
		}
	default:
		if !userConfig.InboxNotifications {
			log.Debugf("user %s not config the inbox notification", msg.ReceiverUserID)
			return
		}
	}

	log.Debugf("user %s config the notification", msg.ReceiverUserID)
}

// renderNotification generates the notification message based on type
func renderNotification(msg plugin.NotificationMessage) string {
	lang := i18n.Language(msg.ReceiverLang)
	switch msg.Type {
	case plugin.NotificationUpdateQuestion:
		return plugin.TranslateWithData(lang, slackI18n.TplUpdateQuestion, msg)
	case plugin.NotificationAnswerTheQuestion:
		return plugin.TranslateWithData(lang, slackI18n.TplAnswerTheQuestion, msg)
	case plugin.NotificationUpdateAnswer:
		return plugin.TranslateWithData(lang, slackI18n.TplUpdateAnswer, msg)
	case plugin.NotificationAcceptAnswer:
		return plugin.TranslateWithData(lang, slackI18n.TplAcceptAnswer, msg)
	case plugin.NotificationCommentQuestion:
		return plugin.TranslateWithData(lang, slackI18n.TplCommentQuestion, msg)
	case plugin.NotificationCommentAnswer:
		return plugin.TranslateWithData(lang, slackI18n.TplCommentAnswer, msg)
	case plugin.NotificationReplyToYou:
		return plugin.TranslateWithData(lang, slackI18n.TplReplyToYou, msg)
	case plugin.NotificationMentionYou:
		return plugin.TranslateWithData(lang, slackI18n.TplMentionYou, msg)
	case plugin.NotificationInvitedYouToAnswer:
		return plugin.TranslateWithData(lang, slackI18n.TplInvitedYouToAnswer, msg)
	case plugin.NotificationNewQuestion, plugin.NotificationNewQuestionFollowedTag:
		msg.QuestionTags = strings.Join(strings.Split(msg.QuestionTags, ","), ", ")
		return plugin.TranslateWithData(lang, slackI18n.TplNewQuestion, msg)
	}
	return ""
}

// SendMessage sends a message to a Slack user using Slack API
func (sc *SlackClient) SendMessage(userID string, message string) error {
	data := url.Values{}
	data.Set("channel", userID) // Slack user ID or channel ID
	data.Set("text", message)

	req, err := http.NewRequest("POST", "https://slack.com/api/chat.postMessage", strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Bearer "+sc.AccessToken)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var result struct {
		OK    bool   `json:"ok"`
		Error string `json:"error"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return err
	}

	if !result.OK {
		return fmt.Errorf("Slack API error: %s", result.Error)
	}

	return nil
}
