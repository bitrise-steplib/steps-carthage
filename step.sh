#!/bin/bash

set -e

NO_USE_BINARIES=""
VERBOSE_MODE=""

if [[ "${no_use_binaries}" == "true" ]] ; then
	NO_USE_BINARIES='--no-use-binaries'
fi

if [[ "${verbose_output}" == "true" ]] ; then
	VERBOSE_MODE='--verbose'
fi

#
# Bootstrap
carthage "${carthage_command}" --platform iOS ${NO_USE_BINARIES} ${VERBOSE_MODE}
