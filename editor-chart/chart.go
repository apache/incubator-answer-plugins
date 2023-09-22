
package chart

import "github.com/answerdev/answer/plugin"

type ChartPlugin struct {
}

func init() {
	plugin.Register(&ChartPlugin{})
}

func (d ChartPlugin) Info() plugin.Info {
	return plugin.Info{
		Name:        plugin.MakeTranslator("i18n.chart.name"),
		SlugName:    "chart",
		Description: plugin.MakeTranslator("i18n.chart.description"),
		Author:      "Answer.dev",
		Version:     "0.0.1",
	}
}
