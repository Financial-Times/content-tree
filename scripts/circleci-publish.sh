#!/usr/bin/env bash
# script to publish to npm *IF* the tag matches the package json version
set -euo pipefail

tag="${CIRCLE_TAG:-}"
if [[ -z "$tag" ]]; then
  echo "CIRCLE_TAG is not set (script should run on tag builds in CircleCI)."
  exit 1
fi

pkg_version="$(node -p "require('./package.json').version")"
tag_version="${tag#v}"

if [[ "$tag_version" != "$pkg_version" ]]; then
  echo "Tag ${tag} does not match package.json version ${pkg_version}"
  exit 1
fi

npm publish "$@"
