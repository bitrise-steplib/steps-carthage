#!/bin/bash

set -e

#
# Make sure brew has Carthage installed
brew update
brew install carthage

#
# Bootstrap
carthage bootstrap --platform iOS
