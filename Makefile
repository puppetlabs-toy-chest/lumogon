ifndef GOPATH
$(error GOPATH is not set)
endif

PACKAGE_NAME = github.com/puppetlabs/lumogon
CONTAINER_NAME = puppet/lumogon

LDFLAGS += -X "$(PACKAGE_NAME)/version.BuildTime=$(shell date -u '+%Y-%m-%d %I:%M:%S %Z')"
LDFLAGS += -X "$(PACKAGE_NAME)/version.BuildVersion=$(shell date +%Y%m%d%H%M%S)-$(shell git describe --tags)"
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

build: bin/lumogon

test: lint vet
	go test -v -cover ./... -ldflags '$(TESTLDFLAGS)'

watch: bootstrap
	goconvey

bin/lumogon: bootstrap
	mkdir -p bin
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=0 go build -a -ldflags '$(LDFLAGS)' -o bin/lumogon lumogon.go

clean:
	rm -rf bin/*;
	go clean -i ./...

image:
	docker build -t $(CONTAINER_NAME) .

deploy: image
	script/deploy

todo:
	grep -rnw "TODO" .

lint: bootstrap
	@echo "Linting..."
	@for d in $$(go list ./... | grep -v /vendor/); do golint $${d}; done

vet: bootstrap
	@for d in $$(go list ./... | grep -v /vendor/); do go vet $${d}; done

licenses: $(GOPATH)/bin/licenses
	@licenses  $(PACKAGE_NAME) | grep $(PACKAGE_NAME)/vendor

all: clean test build image puppet-module

$(GOPATH)/bin/golint:
	go get -u golang.org/x/lint/golint

$(GOPATH)/bin/licenses:
	go get -u github.com/pmezard/licenses

$(GOPATH)/bin/goconvey:
	go get -u github.com/smartystreets/goconvey

bootstrap: $(GOPATH)/bin/golint $(GOPATH)/bin/licenses $(GOPATH)/bin/goconvey

puppet-module:
	cd contrib/puppetlabs-lumogon; make all

.PHONY: build image test todo clean bootstrap licenses watch deploy puppet-module
