# Dingtalk connector
> Dingtalk connector is a OAuth plug-in designed to support Dingtalk OAuth login.

## How to use

### Build
```bash
./answer build --with github.com/apache/incubator-answer-plugins/connector-dingtalk
```

### Configuration
- `ClientID` - Dingtalk OAuth client ID
- `ClientSecret` - Dingtalk OAuth client secret

Authorization callback URL as https://example.com/answer/api/v1/connector/redirect/dingtalk

Dingtalk OAuth API documentation: https://open.dingtalk.com/document/orgapp-server/use-dingtalk-account-to-log-on-to-third-party-websites-1

### Build docker image with plugin from answer base image

```Dockerfile
FROM apache/answer as answer-builder

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories

RUN apk --no-cache add \
    build-base git bash nodejs npm go && \
    npm install -g pnpm

RUN go env -w GOPROXY=https://goproxy.cn,direct

RUN answer build \
    --with github.com/apache/incubator-answer-plugins/connector-dingtalk \
    --output /usr/bin/new_answer

FROM alpine

ARG TIMEZONE
ENV TIMEZONE=${TIMEZONE:-"Asia/Shanghai"}

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories

RUN apk update \
    && apk --no-cache add \
        bash \
        ca-certificates \
        curl \
        dumb-init \
        gettext \
        openssh \
        sqlite \
        gnupg \
        tzdata \
    && ln -sf /usr/share/zoneinfo/${TIMEZONE} /etc/localtime \
    && echo "${TIMEZONE}" > /etc/timezone

COPY --from=answer-builder /usr/bin/new_answer /usr/bin/answer
COPY --from=answer-builder /data /data
COPY --from=answer-builder /entrypoint.sh /entrypoint.sh
RUN chmod 755 /entrypoint.sh

VOLUME /data
EXPOSE 80
ENTRYPOINT ["/entrypoint.sh"]
```

You can update the --with parameter to add more plugins that you need.

```bash
docker build -t answer-with-plugin .
docker run -d -p 9080:80 -v answer-data:/data --name answer answer-with-plugin
```
