package formula

import (
	"github.com/answerdev/answer/plugin"
	"github.com/answerdev/plugins/editor-formula/i18n"
)

type FormulaPlugin struct {
}

func init() {
	plugin.Register(&FormulaPlugin{})
}

func (d FormulaPlugin) Info() plugin.Info {
	return plugin.Info{
		Name:        plugin.MakeTranslator(i18n.InfoName),
		SlugName:    "formula_editor",
		Description: plugin.MakeTranslator(i18n.InfoDescription),
		Author:      "answerdev",
		Version:     "0.0.1",
	}
}
