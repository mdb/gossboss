SOURCE=./...
VERSION="0.0.2"

.DEFAULT_GOAL := build

vet:
	go vet $(SOURCE)
.PHONY: vet

fmt:
	go fmt $(SOURCE)
.PHONY: fmt

test-fmt:
	test -z $(shell go fmt $(SOURCE))
.PHONY: test-fmt

test: vet test-fmt
	go test -race -cover $(SOURCE) -count=1
.PHONY: test

benchmark:
	go test -bench=.
.PHONY: benchmark

tools:
	echo "Installing tools from tools.go"
	cat tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %
.PHONY: tools

build: tools
	goreleaser release \
		--snapshot \
		--skip-publish \
		--rm-dist
.PHONY: build

release: tools
	goreleaser release \
		--rm-dist
.PHONY: release

check-tag:
	./scripts/ensure_unique_version.sh "$(VERSION)"
.PHONY: check-tag

tag:
	if git rev-parse $(VERSION) >/dev/null 2>&1; then \
		echo "found existing $(VERSION) git tag"; \
		exit 1; \
	else \
		echo "creating git tag $(VERSION)"; \
		git tag $(VERSION); \
		git push origin $(VERSION); \
	fi
.PHONY: tag

delete-tag:
	git tag -d $(VERSION)
	git push --delete origin $(VERSION)
.PHONY: delete-tag

clean:
	rm -rf dist || exit 0
	rm -rf data || exit 0
.PHONY: clean
