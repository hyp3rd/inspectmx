#!/bin/bash
# Shell script to initialize the dev certs to run local testing on 'NIX like systems'.
#

set -euo pipefail

#######################################
# Globals:
#   TLS_SECRETS_FOLDER
#   TEST_DOMAIN
#######################################
TLS_SECRETS_FOLDER="$(pwd)/.secrets/tls"
TEST_DOMAIN="hyperd.thekingfisher.io"

GREEN='\033[1;32m'
ORANGE='\033[1;33m'
NC='\033[0m' # No Color

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
  echo -e "${ORANGE}I cowardly aborted:${NC} $*" && exit 1
}

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

#######################################
# Generate the tls certificates leveraging openssl
# Globals:
#   TLS_SECRETS_FOLDER
#   TEST_DOMAIN
# Arguments:
#   None
# Returns:
#   None
#######################################
generate_certs () {
  if ! __command_exists openssl; then
    abort "To run this program you are required to install openssl"
  fi

  # create a private key and a self signed certificate
  pushd "$TLS_SECRETS_FOLDER" || return
    if [[ ! -f $TLS_SECRETS_FOLDER/$TEST_DOMAIN.key ]] || [[ ! -f $TLS_SECRETS_FOLDER/$TEST_DOMAIN.crt ]]; then
      info "${GREEN}[INFO]${NC} generating self-signed certs for $TEST_DOMAIN"
      echo
      openssl req -new -newkey rsa:4096 -days 365 -nodes -x509 \
      -subj "/C=NL/ST=Amsterdam/L=Amsterdam/O=ING/CN=$TEST_DOMAIN" \
      -keyout $TEST_DOMAIN.key -out $TEST_DOMAIN.crt
    fi

    # create a strong Diffie-Hellman group, used in negotiating Perfect Forward Secrecy with clients
    if [[ ! -f $TLS_SECRETS_FOLDER/dhparam.pem ]]; then
      info "${GREEN}[INFO]${NC} generating dhparam"
      echo
      openssl dhparam -out dhparam.pem 2048
    fi
  popd || return

  info "Successfully generated the self-signed certs for $TEST_DOMAIN"
}

pushd "$(git rev-parse --show-toplevel)"  || return
  generate_certs
popd || return
