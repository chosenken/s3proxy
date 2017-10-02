IMAGE_NAME=chosenken/s3proxy

build:
	go build -o s3proxy main.go


image:
	docker build -t $(IMAGE_NAME) .