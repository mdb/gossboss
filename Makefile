SOURCE=./...

.PHONY: vet \
	fmt \
	test-fmt \
	test \
	goreleaser \
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

goreleaser:
	if ! which goreleaser &> /dev/null; then \
		go get github.com/goreleaser/goreleaser; \
	fi

build: goreleaser
	goreleaser release \
		--snapshot \
		--skip-publish \
		--rm-dist

clean:
	rm -rf dist || exit 0
	rm -rf data || exit 0
