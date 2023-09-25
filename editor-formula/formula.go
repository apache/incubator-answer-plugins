package formula

import "github.com/answerdev/answer/plugin"

type FormulaPlugin struct {
}

func init() {
	plugin.Register(&FormulaPlugin{})
}

func (d FormulaPlugin) Info() plugin.Info {
	return plugin.Info{
		Name:        plugin.MakeTranslator("i18n.formula_editor.name"),
		SlugName:    "formula_editor",
		Description: plugin.MakeTranslator("i18n.formula_editor.description"),
		Author:      "answerdev",
		Version:     "0.0.1",
	}
}
