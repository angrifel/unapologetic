.PHONY: install-git-hooks docker-golangci-lint lint test run-github-action check-docker check-act

install-git-hooks:
	@mkdir -p ./.git/hooks/
	@cp ./.git-hooks/* ./.git/hooks/

check-docker:
	@command -v docker 2>/dev/null || (echo "Error: 'docker' not found" && exit 1)

check-act:
	@command -v act 2>/dev/null || (echo "Error: 'act' not found" && exit 1)

run-github-action: check-docker check-act
	@act

tidy:
	@go mod tidy

test:
	@go test -v -race -cover ./...

docker-golangci-lint: check-docker
	@docker pull golangci/golangci-lint:v2.3.1

lint: tidy docker-golangci-lint
	@docker run -t --rm -v $(CURDIR):/app -w /app golangci/golangci-lint:v2.3.1 golangci-lint run

