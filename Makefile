# Definitions
ROOT                    := $(PWD)
GO_HTML_COV             := ./coverage.html
GO_TEST_OUTFILE         := ./c.out
GOLANG_DOCKER_IMAGE     := golang:1.13
GOLANG_DOCKER_CONTAINER := goairtable-container

#   Format according to gofmt: https://github.com/cytopia/docker-gofmt
#   Usage:
#       make fmt
#       make fmt path=src/elastic/index_setup.go
fmt:
ifdef path
	docker run --rm -v ${ROOT}:/data cytopia/gofmt -w ${path}
else
	docker run --rm -v ${ROOT}:/data cytopia/gofmt -w .
endif

#   Deletes container if exists
#   Usage:
#       make clean
clean:
	docker rm -f ${GOLANG_DOCKER_CONTAINER} || true

build:
	docker run -w /app -v ${ROOT}:/app ${GOLANG_DOCKER_IMAGE} go mod download
	docker run -w /app -v ${ROOT}:/app ${GOLANG_DOCKER_IMAGE} go build cmd/cli/airtbl.go

#   Deletes container if exists
#   Usage:
#       make run arguments="-help"
run:
	docker run -w /app -v ${ROOT}:/app ${GOLANG_DOCKER_IMAGE} ./airtbl ${arguments}

#   Usage:
#       make test
test:
	docker run -w /app -v ${ROOT}:/app ${GOLANG_DOCKER_IMAGE} go test ./... -coverprofile=${GO_TEST_OUTFILE}
	docker run -w /app -v ${ROOT}:/app ${GOLANG_DOCKER_IMAGE} go tool cover -html=${GO_TEST_OUTFILE} -o ${GO_HTML_COV}

#   Usage:
#       make test-dev
#		make test-dev arguments="-v"
#   Note: initial invocation will be slow as it compiles for
#   	  the first time.
stop:
	docker stop ${GOLANG_DOCKER_CONTAINER} || true

BUILDDOCKER = "docker run --rm -d --name ${GOLANG_DOCKER_CONTAINER} -w /app -v ${ROOT}:/app ${GOLANG_DOCKER_IMAGE} tail -f /dev/null"

test-dev:
	docker top ${GOLANG_DOCKER_CONTAINER} > /dev/null 2>&1 || eval ${BUILDDOCKER}
	docker exec ${GOLANG_DOCKER_CONTAINER} go test ./... -coverprofile=${GO_TEST_OUTFILE} ${arguments}
	docker exec ${GOLANG_DOCKER_CONTAINER} go tool cover -html=${GO_TEST_OUTFILE} -o ${GO_HTML_COV}

#   Usage:
#       make lint
lint:
	docker run --rm -v ${ROOT}:/app -w /app golangci/golangci-lint:v1.27.0 golangci-lint run -v