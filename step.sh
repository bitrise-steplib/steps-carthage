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

brew update && brew install carthage

if [ ! -z "${working_dir}" ] ; then
	echo "==> Switching to working directory: ${working_dir}"
	cd "${working_dir}"
	if [ $? -ne 0 ] ; then
		echo " [!] Failed to switch to working directory: ${working_dir}"
		exit 1
	fi
fi

#
# Run carthage command
carthage "${carthage_command}" --"${platform}" ${NO_USE_BINARIES} ${VERBOSE_MODE}
