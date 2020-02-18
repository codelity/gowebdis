GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=gowebdis
SHELL=/bin/bash
PYTHONCMD=python3
PIPCMD=$(PYTHONCMD) -m pip
ROBOTCMD=$(PYTHONCMD) -m robot

ROOT := $$(git rev-parse --show-toplevel)

all: test build
build: 
	$(GOBUILD) -o $(ROOT)/bin/$(BINARY_NAME) -v
test: 
	$(GOTEST) -v ./...
e2e:
	$(PIPCMD) install -r requirements.txt
	$(ROBOTCMD) $(ROOT)/test/e2e/iam.robot
clean: 
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME)
