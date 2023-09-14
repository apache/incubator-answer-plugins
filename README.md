# Answer Plugins

> Answer Plugins are built to enhance the feature of [Answer](https://github.com/answerdev/answer).

[![LICENSE](https://img.shields.io/github/license/answerdev/answer)](https://github.com/answerdev/answer/blob/main/LICENSE)
[![Language](https://img.shields.io/badge/language-go-blue.svg)](https://golang.org/)
[![Language](https://img.shields.io/badge/language-react-blue.svg)](https://reactjs.org/)
[![Go Report Card](https://goreportcard.com/badge/github.com/answerdev/answer)](https://goreportcard.com/report/github.com/answerdev/answer)
[![Discord](https://img.shields.io/badge/discord-chat-5865f2?logo=discord&logoColor=f5f5f5)](https://discord.gg/Jm7Y4cbUej)

## Contributing your own plugin

If you have a plugin you would like included in the repo create a PR to add the repository to plugins.txt

## How to build the Answer with your need plugins?

Learn more about the plugin, please visit our docs [answer.dev](https://answer.dev/docs/development/extending/plugin_config).

## Want to try the plugin early?

If you want to try it out earlier, you can use the all-in-one docker image. Note that this image will contain **the latest version of answer** and all official plugins, **which may not have been released yet**.

```bash
$ docker run -d -p 9080:80 -v answer-data:/data --name answer answerdev/answer:all-in-one
```
