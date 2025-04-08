windows-build:
	GOOS=windows GOARCH=amd64 go build -o bin/agent.exe -v cmd/main.go

linux-build:
	GOOS=linux GOARCH=amd64 go build -o bin/ys-agent -v cmd/main.go
