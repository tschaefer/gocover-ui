.PHONY: all
all: fmt lint test

.PHONY: fmt
fmt:
	test -z $(shell gofmt -l .) || (echo "[WARN] Fix format issues" && exit 1)

.PHONY: lint
lint:
	test -z $(shell golangci-lint run >/dev/null || echo 1) || (echo "[WARN] Fix lint issues" && exit 1)

.PHONY: test
test:
	test -z $(shell go test -v ./... 2>&1 >/dev/null || echo 1) || (echo "[WARN] Fix test issues" && exit 1)

.PHONY: coverage
coverage:
	test -z $(shell go test -coverprofile=coverage.out ./... 2>&1 >/dev/null || echo 1) || (echo "[WARN] Fix test issues" && exit 1)

.PHONY: build
build:
	goreleaser build --clean --snapshot --single-target

.PHONY: clean
clean:
	rm -rf coverage.out coverage/ dist/
