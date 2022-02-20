# make traffic
.PHONY: default image

default: image

image: bin
	docker buildx build --platform linux/arm/v7 -t matlock08/kubeedge-client:0.0.1 . --push	

.PHONY: bin
bin:
	GOOS=linux GOARCH=arm go build -o client

.PHONY: run
run: bin
	./client
