# make traffic
.PHONY: default image

default: image

image: bin
	docker buildx build --platform linux/arm/v7 -t matlock08/kubeedge-relay:0.0.1 . --push	

.PHONY: bin
bin:
	GOOS=linux GOARCH=arm go build -o relay

.PHONY: run
run: bin
	./relay
