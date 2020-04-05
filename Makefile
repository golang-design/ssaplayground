NAME=ssaplayground
BUILD=$(NAME).out
IMAGE=$(NAME)
VERSION = $(shell git describe --always --tags)
all: clean
	go build -o $(BUILD) -mod vendor
	./$(BUILD) -conf configs/config.yaml
docker:
	docker build -t $(IMAGE):$(VERSION) -t $(IMAGE):latest -f docker/Dockerfile .
run:
	docker ps -a | grep $(NAME) && ([ $$? -eq 0 ] && (docker stop $(NAME) && docker rm -f $(NAME))) || echo "no running container."
	docker run -itd -v $(shell pwd)/data:/app/public/buildbox -p 6789:6789 --name $(NAME) ssaplayground:latest
tidy-docker:
	docker ps -a | grep $(NAME)
	[ $$? -eq 0 ] && docker stop $(NAME) && docker rm -f $(NAME)
	docker images -f "dangling=true" -q | xargs docker rmi -f
	docker image prune -f
clean:
	rm -rf $(BUILD)
.PHONY: all start docker update clean