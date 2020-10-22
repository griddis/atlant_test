FROM golang:1.15-alpine AS build

ENV APP=./cmd
ENV BIN=/bin/atlant_test
ENV PATH_ROJECT=${GOPATH}/src/github.com/griddis/atlant_test
ENV GO111MODULE=on
ARG VERSION
ENV VERSION ${VERSION:-0.1.0}
ARG BUILD_TIME
ENV BUILD_TIME ${BUILD_TIME:-unknown}
ARG COMMIT
ENV COMMIT ${COMMIT:-unknown}

WORKDIR ${PATH_ROJECT}
COPY . ${PATH_ROJECT}

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w \
        -X github.com/griddis/atlant_test/pkg/health.Version=${VERSION} \
        -X github.com/griddis/atlant_test/pkg/health.Commit=${COMMIT} \
        -X github.com/griddis/atlant_test/pkg/health.BuildTime=${BUILD_TIME}" \
    -a -o ${BIN} ${APP}

FROM alpine:3.12 as production
COPY --from=build /bin/atlant_test /bin/atlant_test
ENTRYPOINT ["/bin/atlant_test"]
