ifndef GOPATH
$(error GOPATH is not set)
endif

PACKAGE_NAME = github.com/puppetlabs/transparent-containers/cli

LDFLAGS += -X "$(PACKAGE_NAME)/version.BuildTime=$(shell date -u '+%Y-%m-%d %I:%M:%S %Z')"
LDFLAGS += -X "$(PACKAGE_NAME)/version.BuildVersion=development"
LDFLAGS += -X "$(PACKAGE_NAME)/version.BuildSHA=$(shell git rev-parse HEAD)"
# Strip debug information
LDFLAGS += -s

TESTLDFLAGS += -X "$(PACKAGE_NAME)/version.BuildTime=testdatestring"
TESTLDFLAGS += -X "$(PACKAGE_NAME)/version.BuildVersion=testversionstring"
TESTLDFLAGS += -X "$(PACKAGE_NAME)/version.BuildSHA=testbuildsha"
# Strip debug information - see https://github.com/golang/go/issues/19734
TESTLDFLAGS += -s

GOARCH ?= amd64
GOOS ?= linux

clean:
	rm -rf bin/*;
	go clean -i ./...

dependencies: bootstrap
	glide install

test: lint vet
	go test -v -cover `glide novendor` -ldflags '$(TESTLDFLAGS)'

watch: bootstrap
	goconvey

build: bootstrap
	mkdir -p bin
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=0 go build -a -ldflags '$(LDFLAGS)' -o bin/lumogon lumogon.go

image: bootstrap
	docker build -t local/lumogon -f ./Dockerfile.build .

todo:
	grep -rnw "TODO" .

lint: bootstrap $(GOPATH)/src/github.com/golang/lint/golint
	golint `glide novendor`

vet: bootstrap
	go vet `glide novendor`

licenses: $(GOPATH)/bin/licenses
	@licenses  $(PACKAGE_NAME) | grep $(PACKAGE_NAME)/vendor

all: clean dependencies test build image

buildimage: clean build image

$(GOPATH)/bin/glide:
	go get -u github.com/Masterminds/glide

$(GOPATH)/src/github.com/golang/lint/golint:
	go get -u github.com/golang/lint/golint

$(GOPATH)/bin/licenses:
	go get -u github.com/pmezard/licenses

$(GOPATH)/bin/goconvey:
	go get -u github.com/smartystreets/goconvey

bootstrap: $(GOPATH)/bin/glide $(GOPATH)/src/github.com/golang/lint/golint $(GOPATH)/bin/licenses $(GOPATH)/bin/goconvey

.PHONY: build image buildimage test todo clean dependencies bootstrap licenses watch
