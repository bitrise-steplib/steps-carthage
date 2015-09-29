#!/bin/bash

set -e

NO_USE_BINARIES=""

if [[ "${no_use_binaries}" == "true" ]] ; then
	NO_USE_BINARIES='--no-use-binaries'
fi

#
# Bootstrap
carthage "${carthage_command}" --platform iOS ${NO_USE_BINARIES}
