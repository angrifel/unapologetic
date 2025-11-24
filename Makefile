.PHONY: docker-golangci-lint lint test

test:
	go test -v -race -cover ./...

docker-golangci-lint:
	docker pull golangci/golangci-lint:v2.3.1

lint: docker-golangci-lint
	docker run -t --rm -v $(CURDIR):/app -w /app golangci/golangci-lint:v2.3.1 golangci-lint run

