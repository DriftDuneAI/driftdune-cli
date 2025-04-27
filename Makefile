# Makefile for DriftDune CLI

NAME = driftdune
VERSION = $(shell git describe --tags --always --dirty)
PLATFORMS = darwin-amd64 linux-amd64

.PHONY: all build build-all package clean

all: build

build:
	@echo "Building $(NAME)"
	go build -o $(NAME) .

build-all: $(addprefix build-,$(PLATFORMS))

build-%:
	@echo "Building $(NAME) for $*"
	GOOS=$(word 1,$(subst -, ,$*)) \
	GOARCH=$(word 2,$(subst -, ,$*)) \
	go build -o $(NAME)_$* .

package: build-all
	@echo "Packaging version $(VERSION)"
	@for platform in $(PLATFORMS); do \
		tar czf $(NAME)-$(VERSION)-$$platform.tar.gz $(NAME)_$$platform; \
		zip -j $(NAME)-$(VERSION)-$$platform.zip $(NAME)_$$platform; \
	done

clean:
	@echo "Cleaning up"
	@rm -f $(NAME) \
		$(addsuffix .tar.gz,$(addprefix $(NAME)-$(VERSION)-,$(PLATFORMS))) \
		$(addsuffix .zip,$(addprefix $(NAME)-$(VERSION)-,$(PLATFORMS))) \
		$(addprefix $(NAME)_,$(PLATFORMS))