# Copyright 2020 The golang.design Initiative authors.
# All rights reserved. Use of this source code is governed
# by a GPLv3 license that can be found in the LICENSE file.

FROM golang:1.19.4
WORKDIR /app
ADD . /app
# required for runtime/cgo
RUN apt install g++ && go build -o gossa
CMD [ "/app/gossa", "-conf=/app/configs/docker.yaml"]
