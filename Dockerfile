FROM golang:1.24.1-alpine3.21 as builder

WORKDIR /app

COPY . ./

RUN set -ex && \
  go mod download && \
  go build \
  -ldflags="-w -s" \
  -o /bin/snuuze

FROM alpine:3.21.3

RUN addgroup -S snuuze && adduser -S snuuze -G snuuze
USER snuuze

COPY --from=builder --chown=snuuze:snuuze /bin/snuuze /bin/snuuze

EXPOSE 1323
HEALTHCHECK --interval=60s --timeout=30s --start-period=5s --retries=3 \
  CMD [ "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:1323/ping", "||", "exit", "1" ]

CMD ["/bin/snuuze"]
