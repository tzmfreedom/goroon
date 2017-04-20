NAME := goroon
SRCS := $(shell find . -type d -name vendor -prune -o -type f -name "*.go" -print)
VERSION := 0.1.2
REVISION := $(shell git rev-parse --short HEAD)
LDFLAGS := -ldflags="-s -w -X \"main.Version=$(VERSION)\" -X \"main.Revision=$(REVISION)\"" 
DIST_DIRS := find * -type d -exec
SHA256_386 = $(shell cat dist/goroon-$(VERSION)-darwin-386.tar.gz | openssl dgst -sha256 | sed 's/^.* //')
SHA256_AMD64 = $(shell cat dist/goroon-$(VERSION)-darwin-amd64.tar.gz | openssl dgst -sha256 | sed 's/^.* //')

.DEFAULT_GOAL := test

.PHONY: test
test: glide
	@go test -cover -v `glide novendor`

.PHONY: install
install: build
	@go install ./cmd

.PHONY: uninstall
uninstall:

.PHONY: clean
clean:
	@rm -rf bin/*
	@rm -rf vendor/*
	@rm -rf dist/*

.PHONY: dist-clean
dist-clean: clean
	@rm -f $(NAME).tar.gz

.PHONY: build
build:
	-@goimports -w .
	@gofmt -w .
	@go build $(LDFLAGS)

.PHONY: cross-build
cross-build: deps
	@for os in darwin linux windows; do \
	    for arch in amd64 386; do \
	        GOOS=$$os GOARCH=$$arch CGO_ENABLED=0 go build -a -tags netgo \
	        -installsuffix netgo $(LDFLAGS) -o dist/$$os-$$arch/$(NAME) ./cmd; \
	    done; \
	done

.PHONY: glide
glide:
ifeq ($(shell command -v glide 2> /dev/null),)
	curl https://glide.sh/get | sh
endif

.PHONY: deps
deps: glide
	glide install

.PHONY: bin/$(NAME) 
bin/$(NAME): $(SRCS)
	go build -a -tags netgo -installsuffix netgo $(LDFLAGS) -o bin/$(NAME) ./cmd

.PHONY: dist
dist:
	@cd dist && \
	$(DIST_DIRS) cp ../LICENSE {} \; && \
	$(DIST_DIRS) cp ../README.md {} \; && \
	$(DIST_DIRS) cp ../completions/zsh/_goroon {} \; && \
	$(DIST_DIRS) tar zcf $(NAME)-$(VERSION)-{}.tar.gz {} \;

.PHONY: release
release:
	@cat formula/goroon.rb.tmpl | \
	sed -e 's/{VERSION}/$(VERSION)/g' -e 's/{SHA256_AMD64}/$(SHA256_AMD64)/g' \
	  -e 's/{SHA256_386}/$(SHA256_386)/g' > formula.rb
	@wget https://raw.githubusercontent.com/tzmfreedom/homebrew-$(NAME)/master/$(NAME).rb
	@curl -i -X PUT \
	  -H "Content-Type:application/json" \
	  -H "Authorization:token ${GH_TOKEN}" \
	  -d \
	  "{\"path\":\"$(NAME).rb\",\"sha\":\"$$(git hash-object $(NAME).rb)\",\"content\":\"$$(cat formula.rb | openssl enc -e -base64 | tr -d '\n ')\",\"message\":\"Update version $(VERSION)\"}" \
	  'https://api.github.com/repos/tzmfreedom/homebrew-$(NAME)/contents/$(NAME).rb'
	@rm formula.rb $(NAME).rb
