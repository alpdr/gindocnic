
.DEFAULT_GOAL := help

.PHONY: help lint test build index


##@
##@ Utilities
##@
test: build lint ##@ Runs the tests
	go test -v ./...

build: ##@ Builds the project
	go build ./...

lint: ##@ Examines source code and reports suspicious constructs
	go vet ./...

GIT_TAG ?= v0.0.0

index: ##@ GIT_TAG=<release version> Prompt Go to update its index of modules with information about the module you’re publishing.
	GOPROXY=proxy.golang.org go list -m github.com/alpdr/gindocnic@$(GIT_TAG)

##@
##@ Misc commands
##@
# Reference: https://gist.github.com/prwhite/8168133?permalink_comment_id=4744119#gistcomment-4744119
help: ##@ (Default) Display this message
	@printf "\nUsage: make <command>\n"
	@grep -F -h "##@" $(MAKEFILE_LIST) | grep -F -v grep -F | sed -e 's/\\$$//' | awk 'BEGIN {FS = ":*[[:space:]]*##@[[:space:]]*"}; \
	{ \
		if($$2 == "") \
			printf ""; \
		else if($$0 ~ /^#/) \
			printf "\n%s\n", $$2; \
		else if($$1 == "") \
			printf "     %-20s%s\n", "", $$2; \
		else \
			printf "\n    \033[34m%-20s\033[0m %s\n", $$1, $$2; \
	}'

