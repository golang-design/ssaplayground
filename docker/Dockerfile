FROM golang:1.14-alpine
WORKDIR /app
ADD . /app
# required for runtime/cgo
RUN apk add g++
RUN go get -u golang.org/x/tools/cmd/goimports && go mod tidy
RUN go build -o ssaplayground.service -mod=vendor
EXPOSE 6789
ENTRYPOINT [ "/app/ssaplayground.service", "-conf=/app/configs/docker.yaml"]
