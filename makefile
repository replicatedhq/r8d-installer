

.PHONY: r8d
r8d: fmt vet
	go build -tags r8d -o bin/r8d/cmd/r8d

.PHONY: fmt
fmt:
	go fmt ./pkg/... ./cmd/...

.PHONY: vet
vet:
	go vet ./pkg/... ./cmd/...

.PHONY: deps
deps:
	go run -tags deps ./cmd/deps build --config ./cmd/r8d/manifest.toml

.PHONY: update
update:
	go run -tags deps ./cmd/deps update ./cmd/r8d/manifest.toml

.PHONY: clean
clean:
	rm -rf bin pkg/component/*/assets
