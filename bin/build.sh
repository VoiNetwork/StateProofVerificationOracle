#!/usr/bin/env bash

source ./bin/set_vars.sh

# Public: Inject the version and uilds the binary to build/weaver.
#
# $1 - [optional] a version to inject, otherwise the version from the VERSION file is read.
#
# Examples
#
#   ./bin/build.sh # reads the version in the VERSION file
#   ./bin/build.sh "1.2.3"
#
# Returns exit code 0.
function main() {
  local version

  set_vars

  version=$(<VERSION)

  # If the version argument exists, use it.
  if [ -n "$1" ]; then
    version="$1"
  fi

  printf "%b compiling binary...\n" "${INFO_PREFIX}"
  go build -o build/weaver -ldflags "-X main.Version=$version" cmd/weaver/main.go


  printf "%b done!\n" "${INFO_PREFIX}"
}

# And so, it begins...
main "$1"
