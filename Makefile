.PHONY: build
build:
	docker build -t global-room-chat .
run:
	make build
	docker run --rm -p 8764:8764 global-room-chat