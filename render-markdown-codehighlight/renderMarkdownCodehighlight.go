package answer
import (
  "embed"
  "github.com/apache/incubator-answer/plugin"

  "github.com/apache/incubator-answer-plugins/util"
)

//go:embed  info.yaml
var Info embed.FS

type RenderMarkdownCodehighlight struct {
}
func init() {
	plugin.Register(&RenderMarkdownCodehighlight{})
}
func (RenderMarkdownCodehighlight) Info() plugin.Info {
  info := &util.Info{}
	info.GetInfo(Info)

  return plugin.Info{
    Name:        plugin.MakeTranslator("RenderMarkdownCodehighlight"),
    SlugName:    info.SlugName,
    Description: plugin.MakeTranslator(""),
    Author:      info.Author,
    Version:     info.Version,
    Link:        info.Link,
  }
}
