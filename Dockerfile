##### Build the Go binary first from Go build image
FROM golang:alpine AS builder
COPY ./certs/* ./certs/
RUN apk --no-cache --no-progress add git upx ca-certificates && \
  cat ./certs/*.pem >> /etc/ssl/certs/ca-certificates.crt


COPY . /go/src/inspectmx/
WORKDIR /go/src/inspectmx/

RUN echo Building && \
  cd ./src && \
  mkdir -p tmp && \
  env CGO_ENABLED=0 GO111MODULE=on go build \
  -trimpath \
  -ldflags="all=-s -w -X main.Version=$(cat ../VERSION)-$(git rev-parse --short HEAD) -X main.BuildTime=$(date +%FT%T%z)" \
  -o tmp/inspectmx \
  cmd/inspectmx/main.go && \
  echo Compressing && \
  upx tmp/inspectmx > /dev/null


##### Build final inspectmx image
FROM alpine:3.15.4 AS final

COPY --from=builder /go/src/inspectmx/src/.config.yml /.config.yml
COPY --from=builder /go/src/inspectmx/src/tmp/inspectmx /app/inspectmx

RUN addgroup -S imx \
      && adduser -S -u 10000 -g imx imx

USER imx

EXPOSE 3000 8443

HEALTHCHECK --interval=15m --timeout=60s --retries=10 \
  CMD wget --spider --no-verbose http://localhost:8080/ping || exit 1

CMD ["/app/inspectmx"]