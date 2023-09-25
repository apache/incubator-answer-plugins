package chart

import "github.com/answerdev/answer/plugin"

type ChartPlugin struct {
}

func init() {
	plugin.Register(&ChartPlugin{})
}

func (d ChartPlugin) Info() plugin.Info {
	return plugin.Info{
		Name:        plugin.MakeTranslator("i18n.chart_editor.name"),
		SlugName:    "chart_editor",
		Description: plugin.MakeTranslator("i18n.chart_editor.description"),
		Author:      "answerdev",
		Version:     "0.0.1",
	}
}
