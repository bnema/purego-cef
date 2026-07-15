.PHONY: vet test build generate generate-check generate-check-test

# purego converts uintptr callback args to unsafe.Pointer — this is
# required by the calling convention and triggers false positives.
vet:
	go vet -unsafeptr=false ./...

test:
	go test ./...

build:
	go build ./...

generate:
	go generate ./...

generate-check: generate
	./scripts/check-generation-drift.sh

generate-check-test:
	./scripts/test-generate-check.sh
