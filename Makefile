# Configuration
GO111MODULES=on
MAIN_BINARY=access-key-rotator
MAIN_LAMBDA=access-key-rotator.lambda
LAMBDA_NAME=AccessKeyRotator
SRC_FILES=cmd/access-key-rotator/main.go
SRC_FILES_LAMBDA=cmd/lambda/access-key-rotator.go
DIST_DIR=./build

zip: build-lambda # Make zip for AWS lambda deployment
	zip ${DIST_DIR}/${LAMBDA_NAME}.zip ${DIST_DIR}/${MAIN_LAMBDA}

build-cli: # Build CLI binary
	go build -o ${DIST_DIR}/${MAIN_BINARY} ${SRC_FILES}

build-lambda: # Build Lambda
	go build -o ${DIST_DIR}/${MAIN_LAMBDA} ${SRC_FILES_LAMBDA}

all: build-cli
