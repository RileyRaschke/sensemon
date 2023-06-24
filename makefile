.DEFAULT_GOAL := all

GO=$(shell which go)
DISTVER=$(shell git describe --always --dirty --long --tags)
PKG=$(shell head -1 go.mod | sed 's/^module //')

all: dist

test:
	$(GO) test -v ./...

dist: clean
	$(GO) build -ldflags "-X main.Version=$(DISTVER)"

air: 
	mkdir -p var/air_temp
	$(RM) var/air_temp/sensemon
	$(GO) build -o var/air_temp/sensemon -ldflags "-X $(PKG)/main.Version=$(DISTVER)"

race:
	$(GO) run -ldflags "-X main.Version=$(DISTVER)" --race .

upgrade:
	$(GO) get -u && $(GO) mod tidy

clean:
	$(RM) sensemon

goformat:
	gofmt -s -w .

