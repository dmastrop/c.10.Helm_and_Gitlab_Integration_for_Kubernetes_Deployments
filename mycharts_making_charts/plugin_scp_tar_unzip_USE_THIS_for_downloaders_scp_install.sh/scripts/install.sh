#!/usr/bin/env bash
set -euo pipefail

version="$(grep "version" plugin.yaml | cut -d '"' -f 2)"
echo "Downloading and installing helmscp v${version} ..."

binary_url=""
if [ "$(uname)" == "Darwin" ]; then
    binary_url="https://github.com/abohmeed/helmscpplugin/releases/download/${version}/helmscpplugin-${version}-darwin-amd64.tar.gz"
elif [ "$(uname)" == "Linux" ] ; then
    binary_url="https://github.com/abohmeed/helmscpplugin/releases/download/${version}/helmscpplugin-${version}-linux-amd64.tar.gz"
fi

if [ -z "${binary_url}" ]; then
    echo "Unsupported OS type"
    exit 1
fi
mkdir -p "bin"
mkdir -p "releases/v${version}"
binary_filename="releases/v${version}.tar.gz"

(
    if [ -x "$(which curl 2>/dev/null)" ]; then
        curl -sSL "${binary_url}" -o "${binary_filename}"
    elif [ -x "$(which wget 2>/dev/null)" ]; then
        wget -q "${binary_url}" -O "${binary_filename}"
    else
      echo "ERROR: no curl or wget found to download files." > /dev/stderr
    fi
)

# Unpack the binary.
tar xzf "${binary_filename}" -C "releases/v${version}"
mv "releases/v${version}/helmscpplugin" "bin/helmscp"
exit 0