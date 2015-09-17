#!/bin/bash

set -e

#
# Make sure brew has Carthage installed
brew update
brew install carthage

pushd "${project_dir}"

#
# Bootstrap
carthage bootstrap --platform iOS

popd
