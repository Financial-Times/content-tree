#!/bin/bash

for schema in body-tree content-tree transit-tree; do
	for datafile in tests/schema/"$schema"/valid/*; do
		ajv test --valid -s schemas/"$schema".schema.json -d "$datafile"
	done
	for datafile in tests/schema/"$schema"/invalid/*; do
		ajv test --invalid -s schemas/"$schema".schema.json -d "$datafile"
	done
done
