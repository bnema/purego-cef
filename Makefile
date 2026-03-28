.PHONY: vet test build generate

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
