package basic

import (
	"github.com/apache/incubator-answer/plugin"
	"github.com/go-resty/resty/v2"
	"github.com/segmentfault/pacman/log"
)

const (
	commentCheckURL = "https://rest.akismet.com/1.1/comment-check"
)

func (r *Reviewer) RequestAkismetToCheck(content *plugin.ReviewContent) (isSpam bool, err error) {
	req := make(map[string]string)
	req["blog"] = plugin.SiteURL()
	req["user_ip"] = content.IP
	req["user_agent"] = content.UserAgent
	req["comment_content"] = content.Title + "\n" + content.Content
	req["comment_type"] = "comment"
	req["is_test"] = "false"
	// This is for test if the akismet is available.
	if content.Title == "akismet-guaranteed-spam" {
		req["comment_content"] = content.Title
	}

	log.Debugf("request akismet: %+v", req)

	req["api_key"] = r.Config.APIKey

	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetFormData(req).
		Post(commentCheckURL)

	if err != nil {
		log.Errorf("request akismet failed: %v", err)
		return false, err
	}

	if resp.StatusCode() != 200 {
		log.Errorf("request akismet failed: %v", resp.String())
		return false, nil
	}

	log.Debugf("akismet response: %v, content title is %s", resp.String(), content.Title)

	if resp.String() == "true" {
		return true, nil
	}
	return false, nil
}
