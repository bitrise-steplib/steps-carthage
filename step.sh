#!/bin/bash

set -e

THIS_SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

swift ${THIS_SCRIPT_DIR}/step.swift
exit $?
