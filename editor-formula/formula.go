
package formula

import "github.com/answerdev/answer/plugin"

type FormulaPlugin struct {
}

func init() {
	plugin.Register(&FormulaPlugin{})
}

func (d FormulaPlugin) Info() plugin.Info {
	return plugin.Info{
		Name:        plugin.MakeTranslator("i18n.formula.name"),
		SlugName:    "formula",
		Description: plugin.MakeTranslator("i18n.formula.description"),
		Author:      "Answer.dev",
		Version:     "0.0.1",
	}
}
