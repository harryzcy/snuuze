FROM golang:1.21.5-alpine3.19 as builder

WORKDIR /app

COPY . ./

RUN set -ex && \
  go mod download && \
  go build \
  -ldflags="-w -s" \
  -o /bin/snuuze

FROM alpine:3.19.0

COPY --from=builder /bin/snuuze /bin/snuuze

EXPOSE 1323

CMD ["/bin/snuuze"]
