BUILD=ssaplayground.out
all:
	go build -o $(BUILD) -mod vendor
start:
	./ssaplayground -conf config.yaml
docker:
	echo "NOT IMPLEMENTED"
clean:
	rm -rf $(BUILD)