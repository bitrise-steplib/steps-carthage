#!/bin/bash

set -e

THIS_SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
SWIFT_MAJOR_VERSION="$( swift -version | head -n 1 | sed -e "s/^Apple Swift version \([1-9]*\).*$/\1/" )"

if [[ "${SWIFT_MAJOR_VERSION}" == "3" ]]; then
  swift ${THIS_SCRIPT_DIR}/step-swift3.swift
else
  swift ${THIS_SCRIPT_DIR}/step.swift
fi

exit $?
