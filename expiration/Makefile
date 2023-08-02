.PHONY: dev clean build

GOCMD=go
GODIR=$(shell pwd)/cmd/expiration
GOOUT=$(shell pwd)/bin
GOBUILD=CGO_ENABLED=0 GOOS=linux $(GOCMD) build -a -installsuffix cgo -o $(GOOUT) $(GODIR)
GODEV=$(GOCMD) build -o $(GOOUT) $(GODIR)
GOCLEAN=$(GOCMD) clean
GOLOCAL=$(GOCMD) run $(GODIR)/main.go
GOTEST=$(GOCMD) test -count=1 -v ./tests

clean: 
	$(GOCLEAN)
	rm -rf $(GOOUT)/*

local:
	$(GOLOCAL)

test:
	$(GOTEST)

testone:
	@echo "Running test $(test)"
	$(GOTEST) -run $(test)

dev:
	$(GODEV)

build:
	$(GOBUILD)
