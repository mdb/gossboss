SOURCE=./...
VERSION="0.0.0"

.PHONY: vet \
	fmt \
	test-fmt \
	test \
	goreleaser \
	build \
	tag \
	delete-tag \
	clean

.DEFAULT_GOAL := build

vet:
	go vet $(SOURCE)

fmt:
	go fmt $(SOURCE)

test-fmt:
	test -z $(shell go fmt $(SOURCE))

test: vet test-fmt
	go test -cover $(SOURCE) -count=1

tools:
	echo "Installing tools from tools.go"
	cat tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %

build: tools
	goreleaser release \
		--snapshot \
		--skip-publish \
		--rm-dist

release: tools
	goreleaser release \
		--rm-dist

tag:
	if git rev-parse $(VERSION) >/dev/null 2>&1; then \
		echo "found existing $(VERSION) git tag"; \
		exit 1; \
	else \
		echo "creating git tag $(VERSION)"; \
		git tag $(VERSION); \
		git push origin $(VERSION); \
	fi

delete-tag:
	git tag -d $(VERSION)
	git push --delete origin $(VERSION)

clean:
	rm -rf dist || exit 0
	rm -rf data || exit 0
