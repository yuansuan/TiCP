all: darwin-bin

ifeq ($(type),lite)
    SRC_FILE := ysctl_lite.go
else
    SRC_FILE := ysctl.go
endif

linux-bin:
	GOOS=linux go build -o ysctl $(SRC_FILE)
darwin-bin:
	GOOS=darwin go build -o ysctl $(SRC_FILE)
windows-bin:
	GOOS=windows go build -o ysctl.exe $(SRC_FILE)
