.PHONY: git-rev-list-last-release-commit git-log-since-last-release check-go check-godoc check-docker check-act install-godoc view-godoc-locally install-git-hooks docker-golangci-lint lint-nocache lint test run-github-action

GOLANGCI_LINT_VERSION=v2.3.1

git-rev-list-last-release-commit:
	./release-management.sh last-release-rev

git-log-since-last-release:
	@git log --oneline $(shell ./release-management.sh last-release-rev)..HEAD

install-godoc:
	@go install golang.org/x/tools/cmd/godoc@latest

view-godoc-locally: check-godoc
	@xdg-open http://localhost:6060/
	@godoc -http=:6060

install-git-hooks:
	@mkdir -p ./.git/hooks/
	@cp ./.git-hooks/* ./.git/hooks/

check-go:
	@command -v go 2>/dev/null 1>&2 || (echo "Error: 'go' not found" && exit 1)

check-godoc:
	@command -v godoc 2>/dev/null 1>&2 || (echo "Error: 'godoc' not found" && exit 1)

check-docker:
	@command -v docker 2>/dev/null 1>&2 || (echo "Error: 'docker' not found" && exit 1)

check-act:
	@command -v act 2>/dev/null 1>&2 || (echo "Error: 'act' not found" && exit 1)

run-github-action: check-docker check-act
	@act

tidy: check-go
	@go mod tidy

test: check-go
	@go test -v -race -cover ./...

docker-golangci-lint: check-docker
	@docker pull golangci/golangci-lint:$(GOLANGCI_LINT_VERSION)

lint-nocache: tidy docker-golangci-lint
	@docker run -t --rm -v $(CURDIR):/app -w /app golangci/golangci-lint:$(GOLANGCI_LINT_VERSION) golangci-lint run

lint: check-go tidy docker-golangci-lint
	@mkdir -p $(HOME)/.cache/golangci-lint # creates the cache directory, so it is created with the users permissions and not docker container root user.
	@docker run -t --rm \
		-v $(CURDIR):/app \
		-v $(shell go env GOMODCACHE):/go/pkg/mod:ro \
		-v $(HOME)/.cache/golangci-lint:/root/.cache/golangci-lint \
		-w /app \
		golangci/golangci-lint:$(GOLANGCI_LINT_VERSION) golangci-lint run
