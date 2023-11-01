# Google connector
> Google connector is a OAuth plug-in designed to support Google OAuth login. 

## How to use

### Build
```bash
./answer build --with github.com/apache/incubator-answer-plugins/connector-google
```

### Configuration
- `ClientID` - Google OAuth client ID
- `ClientSecret` - Google OAuth client secret

You need to configure the **redirect URI** such as:
https://example.com/answer/api/v1/connector/redirect/google