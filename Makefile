PROJECT?=atlant_test
PATH_ROJECT?=github.com/griddis/${PROJECT}
APP?=bin/${PROJECT}

VERSION?=0.1.0
COMMIT?=$(shell git rev-parse --short HEAD)
BUILD_TIME?=$(shell date -u '+%Y-%m-%d_%H:%M:%S')

proto:
	sh ./scripts/protoc-gen.sh
clean:
	rm -f ${APP}
build: clean
	cd ./cmd && \
	CGO_ENABLED=0 GOOS=darwin go build -ldflags "-s -w \
        -X github.com/griddis/atlant_test/pkg/health.Version=${VERSION} \
        -X github.com/griddis/atlant_test/pkg/health.Commit=${COMMIT} \
        -X github.com/griddis/atlant_test/pkg/health.BuildTime=${BUILD_TIME}" \
		-a -installsuffix cgo -o ../${APP} .
run: build
	./${APP} --logger.level=debug
test:
	go test -v -race ./...
