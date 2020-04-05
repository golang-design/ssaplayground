BUILD=ssaplayground.out
all:
	go build -o $(BUILD) -mod vendor
start:
	./$(BUILD) -conf config/config.yaml
docker:
	docker build -t ssaplayground:v0.1 -f docker/Dockerfile .
update:
	docker run -itd -v ~/dev/ssaplayground/data:/app/public/buildbox -p 6789:6789 ssaplayground:v0.1
clean:
	rm -rf $(BUILD)
.PHONY: all start docker update clean