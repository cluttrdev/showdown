GIT_DIR := `git rev-parse --show-toplevel`

BIN_NAME := "showdown"
BIN_DIR := "bin"
DIST_DIR := "dist"

# list available recipes
default:
    @just --list

# format code
fmt:
    go fmt ./...

generate:
    go generate ./...

# lint code
lint: generate
    golangci-lint run ./...

# vet code
vet: generate
    go vet ./...

# build application
build *args="": generate
    {{GIT_DIR}}/scripts/build.sh {{args}}

# create binary distribution
dist *args="":
    {{GIT_DIR}}/scripts/dist.sh {{args}}

# create a new release
release *args="":
    {{GIT_DIR}}/scripts/release.sh {{args}}

clean:
    @# generated files
    @echo "rm -r pkg/server/assets"
    @rm -r pkg/server/assets 2>/dev/null || true

    @# build artifacts
    @echo "rm {{BIN_DIR}}/{{BIN_NAME}}"
    @-[ -f {{BIN_DIR}}/{{BIN_NAME}} ] && rm {{BIN_DIR}}/{{BIN_NAME}}
    @-[ -d {{BIN_DIR}} ] && rmdir {{BIN_DIR}}

    @# distribution archives
    @echo "rm {{DIST_DIR}}/{{BIN_NAME}}_*.tar.gz"
    @rm {{DIST_DIR}}/{{BIN_NAME}}_*.tar.gz 2>/dev/null || true
    @echo "rm {{DIST_DIR}}/{{BIN_NAME}}_*.zip"
    @rm {{DIST_DIR}}/{{BIN_NAME}}_*.zip 2>/dev/null || true
    @-[ -d {{DIST_DIR}} ] && rmdir {{DIST_DIR}}

# ---

_unreleased:
    @git log --oneline $(git describe --tags --abbrev=0)...HEAD

_system-info:
    @echo "{{os()}}_{{arch()}}"
