GOLANGCI_LINT ?= golangci-lint

GOLANGCI_LINT_VERSION := v2.13.0
GOVULNCHECK_VERSION := v1.7.0

.PHONY: tools fmt fmt-check lint security vuln test check

tools:
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)
	go install golang.org/x/vuln/cmd/govulncheck@$(GOVULNCHECK_VERSION)

fmt:
	$(GOLANGCI_LINT) fmt

fmt-check:
	@diff="$$( $(GOLANGCI_LINT) fmt --diff )"; \
	if [ -n "$$diff" ]; then \
		printf '%s\n' "$$diff"; \
		exit 1; \
	fi

lint:
	$(GOLANGCI_LINT) run ./...

security:
	$(GOLANGCI_LINT) run --enable-only=gosec,bidichk,asciicheck ./...

vuln:
	go run golang.org/x/vuln/cmd/govulncheck@$(GOVULNCHECK_VERSION) ./...

test:
	go test -race -covermode=atomic ./...

check: fmt-check lint security vuln test
