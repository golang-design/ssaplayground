NAME=gossa
VERSION = $(shell git describe --always --tags)
all: clean
	go build -o $(NAME) -mod vendor
	./$(NAME) -conf configs/config.yaml
build:
	docker build -t $(NAME):$(VERSION) -t $(NAME):latest -f docker/Dockerfile .
up: down
	docker-compose -f docker/deploy.yml up -d
down:
	docker-compose -f docker/deploy.yml down
clean: down
	rm -rf $(NAME)
	docker rmi -f $(shell docker images -f "dangling=true" -q) 2> /dev/null; true
	docker rmi -f $(NAME):latest $(NAME):$(VERSION) 2> /dev/null; true
