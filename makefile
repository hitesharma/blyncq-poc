VERSION ?= latest
REPO ?= hitesharma
IMG = $(REPO)/blyncq-poc:$(VERSION)

run:
	go run main.go

build: go-format go-vet
	docker build -t ${IMG} -f build/dockerfile .

push:
	docker push ${IMG}

build-push: build push

go-format:
	go fmt ./...

go-vet:
	go vet ./...