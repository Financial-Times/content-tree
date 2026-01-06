# ![content-tree][logo]

A tree for Financial Times article content.

---

**content-tree** is a specification for representing Financial Times article
content as an abstract tree. It implements the **[unist][unist]** spec.

## Contents

- [Introduction](#introduction)
- [Types](#types)
- [Mixins](#mixins)
- [Nodes](#nodes)
- [License](#license)

## Introduction

This document defines a format for representing Financial Times article content
as a tree. This specification is written in a
[typescript](https://www.typescriptlang.org/)-like grammar, augmented by the
addition of the `external` property modifier.

The `external` property modifier indicates that the specified field is absent
when the `content-tree` is in
[**transit**](#what-does-it-mean-to-be-in-transit), and required when the
**content-tree** is at rest.

### What is `content-tree`?

`content-tree` extends [unist][unist], a format for syntax trees, to benefit
from its [ecosystem of utilities][unist-utilities].

`content-tree` relates to [JavaScript][js] in that it has an [ecosystem of
utilities][unist-utilities] for working with trees in JavaScript.  However,
`content-tree` is not limited to JavaScript and can be used in other programming
languages.

### How to use the types

We provide two namespaces in `content-tree.d.ts`, which is automatically
generated from this README. `ContentTree` and `ContentTree.transit`.

Install this repository as a dependency:

```sh notangle
npm install https://github.com/Financial-Times/content-tree
```

Use it in your code:

```ts notangle
import type {ContentTree} from "@financial-times/content-tree"

function makeBigNumber(): ContentTree.BigNumber {
	return {
		type: "|<tab>"
		        +--------------+
		        | "big-number" |
		        +--------------+
	}
}

function makeImageSetNoFixins(): ContentTree.transit.ImageSet {
	return {
		type: "image-set",
		id: string,
		// if you try to add a `picture` here it will get mad
	}
}
```

### Contributing

To add or update an item in the content tree:
- clone this repo and run `npm install`
- add your definitions to the content below under the `Nodes` section, using the same pattern as other definitions (subheading, definitions, notes). If they can appear at the top level of the body block, add them to the definition under Abstract Types
- run `npm run build` to update `content-tree.d.ts` and (if required) the `schemas` files

### What does it mean to be in transit?

When a `content-tree` is being rendered visually, external resources have been
fetched and added to the tree. When the `content-tree` is being transmitted
across the network, these external resources are referenced only by their `id`.

It is the state of the tree in the network that we call "in transit".


## Schema

The spec for content-tree can be found at [SPEC.md](./SPEC.md)
## License

This software is published by the Financial Times under the [MIT licence](mit).

Derived from [unist][unist] © [Titus Wormer][titus]

[mit]: http://opensource.org/licenses/MIT
[titus]: https://wooorm.com
[logo]: ./logo.png
[unist]: https://github.com/syntax-tree/unist
[js]: https://www.ecma-international.org/ecma-262/9.0/index.html
[webidl]: https://heycam.github.io/webidl/
[term-tree]: https://github.com/syntax-tree/unist#tree
[term-literal]: https://github.com/syntax-tree/unist#tree
[term-parent]: https://github.com/syntax-tree/unist#parent
[term-child]: https://github.com/syntax-tree/unist#child
[term-root]: https://github.com/syntax-tree/unist#root
[term-leaf]: https://github.com/syntax-tree/unist#leaf
[unist-utilities]: https://github.com/syntax-tree/unist#utilities
