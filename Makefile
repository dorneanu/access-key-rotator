build-cli:
	go build -o build/access-key-rotator cmd/access-key-rotator/main.go

all: build-cli
