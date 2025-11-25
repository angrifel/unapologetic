.PHONY: check-godoc install-godoc view-godoc-locally install-git-hooks docker-golangci-lint lint test run-github-action check-docker check-act

install-godoc:
	@go install golang.org/x/tools/cmd/godoc@latest

view-godoc-locally: check-godoc
	@xdg-open http://localhost:6060/
	@godoc -http=:6060

install-git-hooks:
	@mkdir -p ./.git/hooks/
	@cp ./.git-hooks/* ./.git/hooks/

check-godoc:
	@command -v godoc 2>/dev/null 1>&2 || (echo "Error: 'godoc' not found" && exit 1)

check-docker:
	@command -v docker 2>/dev/null 1>&2 || (echo "Error: 'docker' not found" && exit 1)

check-act:
	@command -v act 2>/dev/null 1>&2 || (echo "Error: 'act' not found" && exit 1)

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

