#!/bin/bash

set -e

NO_USE_BINARIES=""

#
# Make sure brew has Carthage installed
brew update && brew install carthage

if [ "$no_use_binaries" == "true" ]; then
	"$NO_USE_BINARIES"="--no-use-binaries"
fi

#
# Bootstrap
carthage "$carthage_command" --platform iOS "$NO_USE_BINARIES"
