#!/bin/bash

# This is a wrapper for easily running the generation scripts for each of our
# protobuf stub packages in a single command.

set -euo pipefail

SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ]; do SOURCE="$(readlink "$SOURCE")"; done
DIR="$(cd -P "$(dirname "$SOURCE")" && pwd)"

for pkgname in tfplugin5 tfplugin6; do
    echo "Generating protobuf stubs for ${pkgname}..."
    "${DIR}/${pkgname}/generate.sh"
done
echo "All done! Remember to run 'go generate ./...' to update the client mocks too."
