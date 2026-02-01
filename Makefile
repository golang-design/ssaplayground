NAME=gossa
VERSION = $(shell git describe --always --tags)
all: clean
	go build -o $(NAME) -mod vendor
	./$(NAME) -conf configs/config.yaml
build:
	docker build -t $(NAME):latest -f docker/Dockerfile .
up:
	docker compose -f docker/deploy.yml up -d
down:
	docker compose -f docker/deploy.yml down
clean: down
	rm -rf $(NAME)
	docker image prune -f 2> /dev/null; true
	docker image rm -f $(NAME):latest $(NAME):$(VERSION) 2> /dev/null; true
