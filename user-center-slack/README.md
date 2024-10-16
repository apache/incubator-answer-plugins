# Slack User Center

## Feature

- User login via slack Account

## Config

To use this plugin, you need to create [a Slack App](https://api.slack.com/quickstart) first, set the Scope and Redirect URL correctly, and copy the `Client ID`, `Client Secrect`, `Signing Secret` and `Webhook URL`. To activate the Slash Command function, you also need to set the `slash command` in your app. Here are default settings you can try:

> Scope: chat:write, commands, groups:write, im:write, incoming-webhook, mpim:write, users:read, users:read.email
>
> RedirectURL: https://Your_Site_URL/answer/api/v1/user-center/login/callback
>
> Slash command: 
>
> * Command: /ask
> * Requesti URL: https://Your_Site_URL/answer/api/v1/slack/slash
> * Usage Hint: [Title][Content\][Tag1,Tag2...\]



- `Client ID`:  Slack App Client ID

- `Client Secret`: Slack App Secret

- `Signing Secret`: Slack App Signing Secret

- `Webhook URL`: find in the `Incoming Webhooks` feature, such as `https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX`


Note: A Redirect URL must also use HTTPS. You can configure a Redirect URL and scope in the **App Management** page under **OAuth & Permissions**. 

## Document
- https://api.slack.com/quickstart
- https://api.slack.com/authentication/oauth-v2
- https://api.slack.com/messaging/webhooks
- https://api.slack.com/interactivity/slash-commands