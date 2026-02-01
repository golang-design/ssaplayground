# Copyright 2020 The golang.design Initiative authors.
# All rights reserved. Use of this source code is governed
# by a GPLv3 license that can be found in the LICENSE file.

# Build stage
FROM golang:1.22-alpine AS builder
RUN apk add --no-cache gcc musl-dev
WORKDIR /app
ADD . /app
RUN CGO_ENABLED=1 go build -o gossa

# Runtime stage
FROM alpine:latest
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /app/gossa /app/gossa
COPY --from=builder /app/configs /app/configs
COPY --from=builder /app/public /app/public
CMD ["/app/gossa", "-conf=/app/configs/docker.yaml"]
