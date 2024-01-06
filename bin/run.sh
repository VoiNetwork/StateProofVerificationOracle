#!/usr/bin/env bash

source ./bin/set_vars.sh

# Public: Injects the version and runs the program.
#
# $1 - [optional] a version to inject, otherwise the version from the VERSION file is read.
#
# Examples
#
#   ./bin/run.sh # reads the version in the VERSION file
#   ./bin/run.sh "1.2.3"
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

  printf "%b running program...\n" "${INFO_PREFIX}"
  go run -ldflags "-X main.Version=$version" cmd/state-proof-relayer/main.go


  printf "%b done!\n" "${INFO_PREFIX}"
}

# And so, it begins...
main "$1"
