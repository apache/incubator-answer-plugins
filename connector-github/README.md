# GitHub connector
> GitHub connector is a OAuth plug-in designed to support GitHub OAuth login.

## How to use

### Build
```bash
./answer build --with github.com/apache/incubator-answer-plugins/connector-github
```

### Configuration
- `ClientID` - GitHub OAuth client ID
- `ClientSecret` - GitHub OAuth client secret

In the https://github.com/settings/applications/new page, config the Authorization callback URL as https://example.com/answer/api/v1/connector/redirect/github