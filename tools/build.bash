#!/bin/bash
# Shell script to wrap the build of the golang code.
#

set -euo pipefail

RED='\033[0;31m'
ORANGE='\033[1;33m'
GREEN='\033[1;32m'
NC='\033[0m' # No Color

#######################################
# Display an error message and redirects output to stderr.
# Globals:
#   None
# Arguments:
#   Error message
#######################################
err() {
  echo -e "${RED}[ERR]${NC} [$(date +'%Y-%m-%dT%H:%M:%S%z')]: $*" >&2
}

#######################################
# Displays a warning and redirects the output to stderr.
# Globals:
#   None
# Arguments:
#   Warn message
#######################################
warn() {
  echo -e "${ORANGE}[WARN]${NC} [$(date +'%Y-%m-%dT%H:%M:%S%z')]: $*" >&2
}

#######################################
# Displays an info message.
# Globals:
#   None
# Arguments:
#   Info message
#######################################
info() {
  echo -e "${GREEN}[INFO]${NC} [$(date +'%Y-%m-%dT%H:%M:%S%z')]: $*" >&1
}

#######################################
# Displays log message.
# Globals:
#   None
# Arguments:
#   Log message
#######################################
log() {
  echo -e "${GREEN}[+]${NC} $*" >&1
}

#######################################
# Displays exit message, and exits the program.
# Globals:
#   None
# Arguments:
#   None
# Returns:
#   non-zero.
#######################################
abort() {
  echo -e "${ORANGE}I cowardly aborted${NC}: $*" && exit 1
}

#######################################
# Globals:
#   PROJECT_NAME
#######################################
PROJECT_NAME="inspectmx"

### Internal Helpers

#######################################
# Determine if a command exists on the system
# Globals:
#   None
# Arguments:
#   A directory or a file
# Returns:
#   0 if the command exists, non-zero on error.
#######################################
__command_exists() {
  command -v "$@" >/dev/null 2>&1
}

# shellcheck disable=SC1091
[[ -f "$(git rev-parse --show-toplevel)/env.bash" ]] && source "$(git rev-parse --show-toplevel)/env.bash"

#######################################
# Cross-compile leveraging docker that generates alpine linux compatible binaries
# Globals:
#   PROJECT_NAME
# Arguments:
#   None
# Returns:
#   0 if the command exists, non-zero on error.
#######################################
build () {
    if ! __command_exists docker; then
        abort "docker doesn't seem to be present on your system"
    fi
    pushd "$(git rev-parse --show-toplevel)/src" || return
    rm -rf "$(pwd)/releases/*"
    ls -al "$(pwd)/releases/"
    info "the releases and data folders are empty"
    info "build process started"

    docker run --rm -it -v "$(PWD)":/usr/local/go/src/app -w /usr/local/go/src/app golang:1.18.1-bullseye bash -c '
    GO111MODULE=off go get golang.org/x/net/http2
    mkdir -p /usr/local/go/src/vendor/golang.org/x/net
    cp -r /go/src/golang.org/x/net/http2 /usr/local/go/src/vendor/golang.org/x/net

    for GOOS in darwin linux; do
        for GOARCH in amd64; do
         export GOOS GOARCH
         GOPATH="$(pwd)/vendor:$(pwd)" CGO_ENABLED=0 GO111MODULE=on go build -ldflags="-w -s -X main.minversion=$(date -u +.%Y%m%d.%H%M%S)" \
          -a -installsuffix "static" -o releases/inspectmx-$GOOS-$GOARCH cmd/inspectmx/main.go
        done
    done
    '
  popd || return
  info "successfully cross-compiled the code"
}

#######################################
# Build the docker images and pushes the changes to ACR
# Globals:
#   PROJECT_NAME
# Arguments:
#   None
# Returns:
#   0 if the command exists, non-zero on error.
#######################################
build_docker_image() {
  pushd "$(git rev-parse --show-toplevel)"
     info "building the docker images"

      # build the api image with our modifications (see Dockerfile) and tag for private ACR
      # docker build --no-cache --file "$(git rev-parse --show-toplevel)/Dockerfile" -t "$REGISTRY/$PURPOSE_ID/$PROJECT_NAME:latest" .
      docker build --file "$(git rev-parse --show-toplevel)/Dockerfile" -t "$PROJECT_NAME:latest" .

      # info "pushing the docker image to $REGISTRY/$PURPOSE_ID"

      # docker push "$REGISTRY/$PROJECT_ID/$PROJECT_NAME:latest"
   popd
}

# build
build_docker_image