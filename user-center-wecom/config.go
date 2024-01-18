package wecom

import (
	"encoding/json"
	"time"

	"github.com/apache/incubator-answer-plugins/user-center-wecom/i18n"
	"github.com/apache/incubator-answer/plugin"
)

type UserCenterConfig struct {
	CorpID       string `json:"corp_id"`
	CorpSecret   string `json:"corp_secret"`
	AgentID      string `json:"agent_id"`
	AutoSync     bool   `json:"auto_sync"`
	Notification bool   `json:"notification"`
}

func (uc *UserCenter) ConfigFields() []plugin.ConfigField {
	syncState := plugin.LoadingActionStateNone
	lastSuccessfulSyncAt := "None"
	if !uc.syncTime.IsZero() {
		syncState = plugin.LoadingActionStateComplete
		lastSuccessfulSyncAt = uc.syncTime.In(time.FixedZone("GMT", 8*3600)).Format("2006-01-02 15:04:05")
	}
	t := func(ctx *plugin.GinContext) string {
		return plugin.Translate(ctx, i18n.ConfigSyncNowDescription) + ": " + lastSuccessfulSyncAt
	}
	syncNowDesc := plugin.Translator{Fn: t}

	syncNowLabel := plugin.MakeTranslator(i18n.ConfigSyncNowLabel)

	if uc.syncing {
		syncNowLabel = plugin.MakeTranslator(i18n.ConfigSyncNowLabelForDoing)
		syncState = plugin.LoadingActionStatePending
	}

	return []plugin.ConfigField{
		{
			Name:        "auto_sync",
			Type:        plugin.ConfigTypeSwitch,
			Title:       plugin.MakeTranslator(i18n.ConfigAutoSyncTitle),
			Description: plugin.MakeTranslator(i18n.ConfigAutoSyncDescription),
			Required:    false,
			UIOptions: plugin.ConfigFieldUIOptions{
				Label: plugin.MakeTranslator(i18n.ConfigAutoSyncLabel),
			},
			Value: uc.Config.AutoSync,
		},
		{
			Name:        "sync_now",
			Type:        plugin.ConfigTypeButton,
			Title:       plugin.MakeTranslator(i18n.ConfigSyncNowTitle),
			Description: syncNowDesc,
			UIOptions: plugin.ConfigFieldUIOptions{
				Text: syncNowLabel,
				Action: &plugin.UIOptionAction{
					Url:    "/answer/admin/api/wecom/sync",
					Method: "get",
					Loading: &plugin.LoadingAction{
						Text:  plugin.MakeTranslator(i18n.ConfigSyncNowLabelForDoing),
						State: syncState,
					},
					OnComplete: &plugin.OnCompleteAction{
						ToastReturnMessage: true,
						RefreshFormConfig:  true,
					},
				},
				Variant: "outline-secondary",
			},
		},
		{
			Name:     "corp_id",
			Type:     plugin.ConfigTypeInput,
			Title:    plugin.MakeTranslator(i18n.ConfigCorpIdTitle),
			Required: true,
			UIOptions: plugin.ConfigFieldUIOptions{
				InputType: plugin.InputTypeText,
			},
			Value: uc.Config.CorpID,
		},
		{
			Name:     "corp_secret",
			Type:     plugin.ConfigTypeInput,
			Title:    plugin.MakeTranslator(i18n.ConfigCorpSecretTitle),
			Required: true,
			UIOptions: plugin.ConfigFieldUIOptions{
				InputType: plugin.InputTypePassword,
			},
			Value: uc.Config.CorpSecret,
		},
		{
			Name:     "agent_id",
			Type:     plugin.ConfigTypeInput,
			Title:    plugin.MakeTranslator(i18n.ConfigAgentIDTitle),
			Required: true,
			UIOptions: plugin.ConfigFieldUIOptions{
				InputType: plugin.InputTypeText,
			},
			Value: uc.Config.AgentID,
		},
		{
			Name:        "notification",
			Type:        plugin.ConfigTypeSwitch,
			Title:       plugin.MakeTranslator(i18n.ConfigNotificationTitle),
			Description: plugin.MakeTranslator(i18n.ConfigNotificationDescription),
			UIOptions: plugin.ConfigFieldUIOptions{
				Label: plugin.MakeTranslator(i18n.ConfigNotificationLabel),
			},
			Value: uc.Config.Notification,
		},
	}
}

func (uc *UserCenter) ConfigReceiver(config []byte) error {
	c := &UserCenterConfig{}
	_ = json.Unmarshal(config, c)
	uc.Config = c
	uc.Company = NewCompany(c.CorpID, c.CorpSecret, c.AgentID)
	uc.asyncCompany()
	return nil
}
