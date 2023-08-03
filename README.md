# Answer Official Plugins
> Answer Official Plugins for enhance the feature of [Answer](https://github.com/answerdev/answer).

[![LICENSE](https://img.shields.io/github/license/answerdev/answer)](https://github.com/answerdev/answer/blob/main/LICENSE)
[![Language](https://img.shields.io/badge/language-go-blue.svg)](https://golang.org/)
[![Language](https://img.shields.io/badge/language-react-blue.svg)](https://reactjs.org/)
[![Go Report Card](https://goreportcard.com/badge/github.com/answerdev/answer)](https://goreportcard.com/report/github.com/answerdev/answer)
[![Discord](https://img.shields.io/badge/discord-chat-5865f2?logo=discord&logoColor=f5f5f5)](https://discord.gg/Jm7Y4cbUej)

## Types of plugin
> Our plugin is under development and the interface definition of the plugin can be viewed at this [link](https://github.com/answerdev/answer/tree/main/plugin).

### Connector 
> The Connector plugin helps us to implement third-party login functionality.   
> For example: Google or GitHub OAuth login.

- [x] [OAuth2 Basic](https://github.com/answerdev/plugins/tree/main/connector/basic)
- [x] [GitHub](https://github.com/answerdev/plugins/tree/main/connector/github)
- [x] [Google](https://github.com/answerdev/plugins/tree/main/connector/google)

### Storage (preview)
> The Storage plugin helps us to upload files to third-party storage.  
> For example: Aliyun OSS or AWS S3.

- [ ] [Aliyun](https://github.com/answerdev/plugins/tree/main/storage/aliyunoss)
- [ ] [S3](https://github.com/answerdev/plugins/tree/main/storage/s3)

### Cache (preview)
> Using the Cache plugin allows you to store cached data in a different location.  
> For example: Redis or Memcached.

- [ ] [Redis](https://github.com/answerdev/plugins/tree/main/cache/redis)

### Finder (preview)
- [ ] [Elasticsearch](https://github.com/answerdev/plugins/tree/main/search/es)

### Filter (coming soon)

### Render (coming soon)

### Exporter (coming soon)

### Importer (coming soon)

## How to build the Answer with your need plugins?
To learn more about the plugin, visit [answer.dev](https://answer.dev).

## Want to try the plugin early?
If you want to try it out early, you can use the all-in-one docker image. Note that this image will contain **the latest version of answer** and all official plugins, **which may not have been released yet**.

```bash
$ docker run -d -p 9080:80 -v answer-data:/data --name answer answerdev/answer:all-in-one
```