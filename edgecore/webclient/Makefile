# make traffic
.PHONY: default image

default: image

image: bin
	docker build -t matlock08/kubeedge-client:0.0.1 . 
	docker image push matlock08/kubeedge-client:0.0.1

.PHONY: bin
bin:
	GOOS=linux go build -o client

.PHONY: run
run: bin
	./client
