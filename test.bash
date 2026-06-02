#!/bin/bash

set -euo pipefail

AJV=./node_modules/.bin/ajv
AJV_CUSTOM_KEYWORDS=./tests/schema/ajv-custom-keywords.js

for schema in body-tree content-tree transit-tree; do
	for datafile in tests/schema/"$schema"/valid/*; do
		"$AJV" test --valid --allow-union-types -c "$AJV_CUSTOM_KEYWORDS" -s schemas/"$schema".schema.json -d "$datafile"
	done
	for datafile in tests/schema/"$schema"/invalid/*; do
		"$AJV" test --invalid --allow-union-types -c "$AJV_CUSTOM_KEYWORDS" -s schemas/"$schema".schema.json -d "$datafile"
	done
done
