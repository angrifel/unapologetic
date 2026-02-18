.PHONY: git-rev-list-last-release-commit git-diff-since-last-release git-log-since-last-release check-go check-godoc check-docker check-act install-godoc view-godoc-locally install-git-hooks docker-golangci-lint lint-nocache lint test run-github-action prepare-release semver-prerelease

GOLANGCI_LINT_VERSION=v2.3.1

LAST_RELEASE_TAG=git describe --tags --match 'v*' --abbrev=0 2>/dev/null
LAST_RELEASE_REV=$(LAST_RELEASE_TAG) || git rev-list --max-parents=0 HEAD | tail -1

SEMVER_PRERELEASE_CORE={ $(LAST_RELEASE_TAG) || echo "v0.0.0"; } | sed 's/^\(v[0-9]*\.[0-9]*\.[0-9]*\).*/\1/'
SEMVER_PRERELEASE_PRERELEASE=git rev-parse --abbrev-ref HEAD | sed 's/[^a-zA-Z0-9]/--/g'
SEMVER_PRERELEASE_BUILD=git rev-parse --short HEAD

semver-prerelease:
	@echo "$(shell $(SEMVER_PRERELEASE_CORE))-$(shell $(SEMVER_PRERELEASE_PRERELEASE))+$(shell $(SEMVER_PRERELEASE_BUILD))"

git-rev-list-last-release-commit:
	@$(LAST_RELEASE_REV)

git-diff-since-last-release:
	@git diff $(shell $(LAST_RELEASE_REV))...HEAD

git-log-since-last-release:
	@git log --oneline $(shell $(LAST_RELEASE_REV))..HEAD

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

prepare-release:
ifndef VERSION
	$(error VERSION is required. Usage: make prepare-release VERSION=v1.0.0)
endif
	@echo "$(VERSION)" | grep -qE '^v[0-9]+\.[0-9]+\.[0-9]+(-[0-9A-Za-z.-]+)?([+][0-9A-Za-z.-]+)?$$' || \
		(echo "Error: VERSION '$(VERSION)' is not a valid semver string (expected: vMAJOR.MINOR.PATCH[-PRERELEASE][+BUILD])" && exit 1)
	@DATE=$$(date +%Y-%m-%d); \
	sed -i "s/## \[Unreleased\]/## [$(VERSION)] - $$DATE/" CHANGELOG.md; \
	sed -i "s/^\[Unreleased\]: \(.*\/compare\/[a-f0-9]*\)\.\.\..*/[$(VERSION)]: \1...$(VERSION)/" CHANGELOG.md
	@git add CHANGELOG.md
	@git commit -m "release version $(VERSION)"
	@git tag $(VERSION)
