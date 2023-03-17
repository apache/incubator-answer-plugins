# OAuth2 Basic Connector
> OAuth2 Basic Connector is a generic OAuth plug-in designed to support any of the OAuth login functions.  
> For example: Google, GitHub, Facebook, Twitter, etc.

## How to use
```bash
./answer build --with github.com/answerdev/plugins/connector/basic
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
https://example.com/answer/api/v1/connector/login/basic

## GitHub OAuth Configuration Example
> The following list is not mentioned can be configured according to your actual situation, not required.

- Name: GitHub
- Client ID: 8cb9dxxxxxc24de9
- Client Secret: 9a3e055xxxxxxxxxxxxxxxxxxxxxxxxxxb78978bc
- Authorize URL: https://github.com/login/oauth/authorize
- Token URL: https://github.com/login/oauth/access_token
- User Json Url: https://api.github.com/user
- User ID Json Path: id
- User Username Json Path: name
- User Display Name Json Path: login
- User Email Json Path: email
- User Avatar Json Path: avatar_url

In the [https://github.com/settings/applications/new](https://github.com/settings/applications/new) page, 
config the `Authorization callback URL` as `https://example.com/answer/api/v1/connector/login/basic`