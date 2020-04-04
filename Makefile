BUILD=ssaplayground
all:
	go build -o $(BUILD)
start:
	./ssaplayground -conf config.yaml
docker:
	echo "NOT IMPLEMENTED"
clean:
	rm -rf $(BUILD)