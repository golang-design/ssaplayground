# Copyright 2020 The golang.design Initiative authors.
# All rights reserved. Use of this source code is governed
# by a GPLv3 license that can be found in the LICENSE file.

FROM golang:1.14-alpine
WORKDIR /app
ADD . /app
# required for runtime/cgo
RUN apk add g++
RUN go get -u golang.org/x/tools/cmd/goimports && go mod tidy
RUN go build -o ssaplayground.service -mod=mod
EXPOSE 6789
ENTRYPOINT [ "/app/ssaplayground.service", "-conf=/app/configs/docker.yaml"]
