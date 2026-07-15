.PHONY: vet test build generate generate-check

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
	git diff --exit-code
