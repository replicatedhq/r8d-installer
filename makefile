

.PHONY: r8d-installer
r8d-installer: fmt vet
	go build -o bin/r8d-installer ./cmd/installer

.PHONY: fmt
fmt:
	go fmt ./pkg/... ./cmd/...

.PHONY: vet
vet:
	go vet ./pkg/... ./cmd/...

.PHONY: deps
deps:
	go run -tags deps ./cmd/deps build --config ./cmd/installer/manifest.toml

.PHONY: update
update:
	go run -tags deps ./cmd/deps update ./cmd/installer/manifest.toml

.PHONY: clean
clean:
	rm -rf bin pkg/component/*/assets
