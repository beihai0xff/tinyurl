# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=build/tinyurl
BINARY_LINUX=$(BINARY_NAME)_linux


all: clean_win build
build:
		SET CGO_ENABLED=0
		SET GOOS=linux
		SET GOARCH=amd64
		$(GOBUILD)  -o $(BINARY_LINUX) -v -tags=jsoniter .

test:
		$(GOTEST) -v ./...
bench:
		$(GOTEST) -bench=. -benchtime=3s -run=none
clean:
		$(GOCLEAN)
		rm -f $(BINARY_NAME)
		rm -f $(BINARY_UNIX)
clean_win:
		$(GOCLEAN)
		rd /s build
run:
		$(GOBUILD) -o $(BINARY_NAME) -v ./...
		./$(BINARY_NAME)
