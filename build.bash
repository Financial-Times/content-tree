#!/bin/bash
node tools/maketypes <README.md> content-tree.ts
tsc -d content-tree.ts
typescript-json-schema --validationKeywords sparkGenerateStoryblock --noExtraProps --required content-tree.ts ContentTree.full.Root > schemas/content-tree.schema.json
typescript-json-schema --validationKeywords sparkGenerateStoryblock --noExtraProps --required  content-tree.ts ContentTree.transit.Root > schemas/transit-tree.schema.json
typescript-json-schema --validationKeywords sparkGenerateStoryblock --noExtraProps --required content-tree.ts ContentTree.transit.Body > schemas/body-tree.schema.json
rm content-tree.ts
rm content-tree.js
