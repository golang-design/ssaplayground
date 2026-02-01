# Copyright 2020 The golang.design Initiative authors.
# All rights reserved. Use of this source code is governed
# by a GPLv3 license that can be found in the LICENSE file.

FROM golang:1.25-alpine
RUN apk add --no-cache gcc musl-dev
WORKDIR /app
ADD . /app
RUN CGO_ENABLED=1 go build -o gossa
CMD ["/app/gossa", "-conf=/app/configs/docker.yaml"]
