#!/bin/bash

set -x

NO_USE_BINARIES=""

#
# Make sure brew has Carthage installed
brew update && brew install carthage

if [ $no_use_binaries == 1 ]; then
	$"NO_USE_BINARIES"="--no-use-binaries"
fi

#
# Bootstrap
carthage $"carthage_command" --platform iOS $"NO_USE_BINARIES"
