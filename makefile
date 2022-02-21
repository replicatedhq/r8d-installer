

.PHONY: r8d-installer
r8d-installer: fmt vet
	go build -o bin/r8d-installer ./cmd/r8d-installer

.PHONY: fmt
fmt:
	go fmt ./pkg/... ./cmd/...

.PHONY: vet
vet:
	go vet ./pkg/... ./cmd/...

.PHONY: deps
deps:
	go run -tags deps ./cmd/r8d-deps build --config ./cmd/r8d-installer/manifest.toml

.PHONY: update
update:
	go run -tags deps ./cmd/r8d-deps update ./cmd/r8d-installer/manifest.toml
