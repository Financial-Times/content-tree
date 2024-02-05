#!/bin/bash
for schema in body-tree content-tree transit-tree; do
	for datafile in tests/"$schema"/valid/*; do
		ajv test --valid -s "$schema".schema.json -d "$datafile"
	done
	for datafile in tests/"$schema"/invalid/*; do
		ajv test --invalid -s "$schema".schema.json -d "$datafile"
	done
done
