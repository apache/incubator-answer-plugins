# OAuth2 Basic Connector
> OAuth2 Basic Connector is a generic OAuth plug-in designed to support any of the OAuth login functions.  
> For example: Google, GitHub, Facebook, Twitter, etc.

## How to use
```bash
./answer build --with github.com/apache/incubator-answer-plugins/connector-basic
```

## How to config
> The following configuration items are in the plugin tab of the admin pag.

- Name: Name of your connector which will be shown in the login page
- ClientID: Client ID of your application 
- ClientSecret: Client secret of your application
- Authorize URL: Authorize URL of your application
- Token URL: Token URL of your application
- User JSON URL: Get user info from this URL
- User ID JSON Path: Path in the OAuth2 User JSON to the user id. eg: user.id
- User Display Name JSON Path: Path in the OAuth2 User JSON to the user display name. eg: user.name
- User Username JSON Path: Path in the OAuth2 User JSON to the user username. eg: user.login
- User Email JSON Path: Path in the OAuth2 User JSON to the user email. eg: user.email
- User Avatar JSON Path: Path in the OAuth2 User JSON to the user avatar. eg: user.avatar_url
- Check Email Verified: If set to true, the email will be verified by email_verified_json_path. If not, the email is always believed to have been verified.
- Email Verified JSON Path: Path in the OAuth2 User JSON to the email verified. eg: user.email_verified
- Scope: OAuth Scope of your application. Multiple scopes separated by `,` e.g. user.email,user.age
- Logo SVG: SVG of your application logo which format is base64

You need to configure the **redirect URI** in a third-party platform, such as google oauth, such as:
https://example.com/answer/api/v1/connector/redirect/basic

## GitHub OAuth Configuration Example
> The following list is not mentioned can be configured according to your actual situation, not required.

- Name: `GitHub`
- Client ID: `8cb9dxxxxxc24de9`
- Client Secret: `9a3e055xxxxxxxxxxxxxxxxxxxxxxxxxxb78978bc`
- Authorize URL: `https://github.com/login/oauth/authorize`
- Token URL: `https://github.com/login/oauth/access_token`
- User Json Url: `https://api.github.com/user`
- User ID Json Path: `id`
- User Display Name Json Path: `login`
- User Username Json Path: `name`
- User Email Json Path: `email`
- User Avatar Json Path: `avatar_url`

In the [https://github.com/settings/applications/new](https://github.com/settings/applications/new) page, 
config the `Authorization callback URL` as `https://example.com/answer/api/v1/connector/redirect/basic`

## Google OAuth Configuration Example

- Name: `Google`
- Client ID: `xxx.apps.googleusercontent.com`
- Client Secret: `GOCSPX-xxx-xxxx`
- Authorize URL: `https://accounts.google.com/o/oauth2/auth`
- Token URL: `https://oauth2.googleapis.com/token`
- User Json Url: `https://www.googleapis.com/oauth2/v3/userinfo`
- User ID Json Path: `sub`
- User Display Name Json Path: `name`
- User Username Json Path: `name`
- User Email Json Path: `email`
- User Avatar Json Path: `picture`
- Email Verified Json Path: `email_verified`
- Scope: `https://www.googleapis.com/auth/userinfo.email,https://www.googleapis.com/auth/userinfo.profile,openid`

In the [https://console.developers.google.com/apis/credentials](https://console.developers.google.com/apis/credentials) page, config the `Authorized redirect URIs` as `https://example.com/answer/api/v1/connector/redirect/basic`

## Discord OAuth Configuration Example

- Name: `Discord`
- Client ID: `1126xxx`
- Client Secret: `NfmIMMcxxx`
- Authorize URL: `https://discord.com/oauth2/authorize`
- Token URL: `https://discord.com/api/oauth2/token`
- User Json Url: `https://discord.com/api/users/@me`
- User ID Json Path: `id`
- User Display Name Json Path: `username`
- User Username Json Path: `username`
- User Email Json Path: `email`
- User Avatar Json Path: `avatar`
- Scope: `email,identify`

In the [https://discord.com/developers/applications](https://discord.com/developers/applications) page, config the `Redirects` as `https://example.com/answer/api/v1/connector/redirect/basic`

## Okta Workforce Identity Cloud (WIC) OAuth Configuration Example

- Name: `Okta`
- Client ID: `0oa666666`
- Client Secret: `UGqYGya5GJ4E`
- Authorize URL: `https://example.okta.com/oauth2/v1/authorize`
- Token URL: `https://example.okta.com/oauth2/v1/token`
- User Json Url: `https://example.okta.com/oauth2/v1/userinfo`
- User ID Json Path: `sub`
- User Display Name Json Path: `name`
- User Username Json Path: `email`
- User Email Json Path: `email`
- Email Verified JSON Path: `email_verified`
- Scope: `openid,email,groups`

In the Okta Application setup; config the `Sign-in redirect URIs` as `https://example.com/answer/api/v1/connector/redirect/basic` and the `Initiate login URI` as `https://example.com/answer/api/v1/connector/login/basic`
In the `Admin \ General` in `Answers` ensure that the `Site URL` matches the page adddress as above (`https://example.com/answer`) or `Okta` will return a `4xx` error.
