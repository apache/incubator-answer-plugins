package answer

import (
  "embed"
  "github.com/apache/incubator-answer-plugins/connector-wallet-route/i18n"
  "github.com/apache/incubator-answer-plugins/util"
  "github.com/apache/incubator-answer/plugin"
)

//go:embed  info.yaml
var Info embed.FS

type ConnectorWalletRoute struct {
}

func init() {
	plugin.Register(&ConnectorWalletRoute{})
}

func (ConnectorWalletRoute) Info() plugin.Info {
  info := &util.Info{}
	info.GetInfo(Info)

  return plugin.Info{
    Name:        plugin.MakeTranslator(i18n.InfoName),
    SlugName:    info.SlugName,
    Description: plugin.MakeTranslator(i18n.InfoDescription),
    Author:      info.Author,
    Version:     info.Version,
    Link:        info.Link,
  }
}
