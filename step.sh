#!/bin/bash

set -e

NO_USE_BINARIES=""
VERBOSE_MODE=""

if [[ "${carthage_version}" == "0.9.4" ]] ; then
	curl -OlL "https://github.com/Carthage/Carthage/releases/download/0.9.4/Carthage.pkg"
	sudo installer -pkg "Carthage.pkg" -target /
	rm "Carthage.pkg"
else
	brew update && brew install carthage
fi

if [[ "${no_use_binaries}" == "true" ]] ; then
	NO_USE_BINARIES='--no-use-binaries'
fi

if [[ "${verbose_output}" == "true" ]] ; then
	VERBOSE_MODE='--verbose'
fi

#
# Bootstrap
carthage "${carthage_command}" --platform iOS ${NO_USE_BINARIES} ${VERBOSE_MODE}
