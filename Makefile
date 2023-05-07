default: help
.PHONY: help
help: # Show help for each of the Makefile recipes.
	@grep -E '^[a-zA-Z0-9_-]+:.*#'  Makefile | sort | while read -r l; do printf "\033[1;32m$$(echo $$l | cut -f 1 -d':')\033[00m:$$(echo $$l | cut -f 2- -d'#')\n"; done

.PHONY: vet
vet: # Run go vet
	go vet -vettool=$(shell where fieldalignment) ./...

.PHONY: install-tools
install-tools: # Install tools to lint/test the code in this repository
	go install golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment@latest
