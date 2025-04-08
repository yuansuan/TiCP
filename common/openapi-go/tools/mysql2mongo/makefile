# Project name
PROJECT_NAME := mysql2mongo

# Default make command
all: go-linux-bin image

go-linux-bin:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o $(PROJECT_NAME) main.go

go-bin:
	CGO_ENABLED=0 go build -o $(PROJECT_NAME) main.go

# Docker build command
image:
	docker build --build-arg PROJECT_NAME=$(PROJECT_NAME) -t $(PROJECT_NAME):latest .

local:
	go run ./cmd

.DEFAULT_GOAL = all
.PHONY: all go-linux-bin go-bin image local
