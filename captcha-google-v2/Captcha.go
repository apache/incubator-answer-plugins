package chart

import (
	"github.com/apache/incubator-answer-plugins/captcha_google_v2/i18n"
	"github.com/apache/incubator-answer/plugin"
)

type ChartPlugin struct {
}

func init() {
	plugin.Register(&ChartPlugin{})
}

func (d ChartPlugin) Info() plugin.Info {
	return plugin.Info{
		Name:        plugin.MakeTranslator(i18n.InfoName),
		SlugName:    "captcha_google_v2",
		Description: plugin.MakeTranslator(i18n.InfoDescription),
		Author:      "answerdev",
		Version:     "0.0.1",
	}
}
