#!/usr/bin/env bash

set -euo pipefail

act_version=v0.2.92

cd "$(dirname "$0")"

tools_dir=./.tools
bin_dir="${tools_dir}/bin"
src_dir="${tools_dir}/act"

mkdir -p "${bin_dir}"

if [ ! -d "${src_dir}" ]; then
  git clone https://github.com/nektos/act.git "${src_dir}"
fi

git -C "${src_dir}" fetch --tags
git -C "${src_dir}" checkout --quiet "${act_version}"
make -C "${src_dir}" install "PREFIX=../"
