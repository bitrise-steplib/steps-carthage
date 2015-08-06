#!/bin/bash

set -e

#
# Required parameters
if [ -z "${project_dir}" ] ; then
	echo "[!] Missing required input: project_dir"
	exit 1
fi

#
# Make sure brew has Carthage installed
brew update
brew install carthage

pushd "${project_dir}"

#
# Bootstrap
carthage bootstrap --platform iOS

popd
