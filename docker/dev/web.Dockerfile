FROM golang:1.17-alpine

# Required because go requires gcc to build
RUN apk add build-base

RUN apk add inotify-tools

RUN echo $GOPATH

COPY . /clean_web

ARG VERSION="4.13.0"

RUN set -x \
    && apk add --no-cache git \
    && git clone --branch "v${VERSION}" --depth 1 --single-branch https://github.com/golang-migrate/migrate /tmp/go-migrate

WORKDIR /tmp/go-migrate

RUN set -x \
    && CGO_ENABLED=0 go build -tags 'postgres' -ldflags="-s -w" -o ./migrate ./cmd/migrate \
    && ./migrate -version

RUN cp /tmp/go-migrate/migrate /usr/bin/migrate

WORKDIR /clean_web

RUN --mount=type=cache,mode=0755,target=/go/pkg/mod go mod download -x -json
RUN --mount=type=cache,mode=0755,target=/go/pkg/mod go mod vendor


RUN go get github.com/go-delve/delve/cmd/dlv@master

CMD sh /clean_web/docker/dev/web.sh
