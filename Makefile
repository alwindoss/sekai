.PHONY: build build-ui generate run dev build-backend

build:
	go tool task build
build-backend:
	go tool task build-backend
build-ui:
	go tool task build-ui

generate:
	go tool task generate

run: build-ui generate build-backend
	go tool task run

dev:
	go tool task backend-watch