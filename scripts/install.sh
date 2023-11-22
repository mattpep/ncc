#!/usr/bin/env bash

set -e

path=$(dirname $0)
install $path/pre-commit ${path}/../.git/hooks/pre-commit

echo All done
