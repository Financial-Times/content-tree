# ![content-tree][logo]

**content-tree** is a specification that describes the shape of an FT article as an abstract tree. It implements the [unist](https://github.com/syntax-tree/unist) spec. It is intended to be a shared data contract across our CMS ([Spark](https://github.com/Financial-Times/spark)), Content APIs, and Front End systems.

To read the spec, go to [SPEC.md](./SPEC.md). You should read the spec if you are implementing an article renderer, or adding or amending an article component.

---

## Contents

- [Overview](#overview)
- [Concepts](#concepts)
- [Using `content-tree`](#using-content-tree)
	- [Typescript](#typescript)
	- [Go](#go-libraries)
	- [JSON Schema](#json-schemas)
- [Contributing](#contributing)
- [Releasing](#releasing)
- [License](#license)

## Overview

The main content-tree specification is defined in a markdown file

- [**`SPEC.md`**](./SPEC.md): a markdown document defining the specification for the tree. Written in a typescript-like grammar, augmented with a custom `[external](#concepts)` property modifier.

From this spec, we automatically generate the following as part of the build process:

- **`content-tree.d.ts`** - Typescript types, which are automatically generated from the markdown spec
- **`/schemas`** - JSON schemas, which are automatically generated from the markdown spec

The [Content & Metadata](https://biz-ops.in.ft.com/Team/content) team maintain some [Go](https://go.dev/) code to support working with content-tree inside Golang:
- **`content_tree.go`** - Go structs containing type definitions for the content tree spec to use in Golang applications. Manually updated by Content & Metadata. 
- **`/libraries`** - Go libraries used to transform the content-tree data structure between other formats

Supporting code:
- **`/tests`** - tests to validate the schema and transformers
- **`/tools`** - utilities for building and generating the schemas and types


## Concepts

### (Abstract) Syntax Tree

An Abstract Syntax Tree, or AST, is a structured representation of content where the meaning and structure of content are represented as a tree, independent of any particular rendering instruction or representation.

Content Tree is an AST that represents the _semantic_ structure of an article - for example paragraphs, headings, images, other editorial components - without enforcing a particular markup language or visual representation. This makes it easier to use the same content across different products and platforms and contexts.

- **`Node`** - an element in the tree with a `type` that represents a unit of content
- **`Parent`** - a `Node` which has `children` nodes
- **`Root`** - the single, top-level `Node` in a tree 

### In Transit vs At Rest

An FT article may contain supplementary assets that are published independently of the content (e.g. images, video clips), or editorial components that pull in data from external systems (e.g. flourish charts, social media posts). When an article is being transmitted over the network (e.g. being fetched via the FT Content API), these external resources are typically referenced by `id`, and it is up to the consuming application to fetch these resources. 

To support both of these use cases, content-tree can exist in different states: a **full** tree (with resolved external data) and a **transit** tree (where external resources are referenced by ID). The content and editorial intent remain the same across both states.

### `external`

To distinguish between transit and full node properties, we use an additional modifier on our typescript properties, called `external`, which indicates that the property is omitted from the transit representation and must be supplied by the consuming application in order to construct a full tree.

**Example:** a tweet/X post will be published with an ID, but the actual content of the post should be fetched at render time from the X API.

<table>
	<thead>
		<tr>	
			<td><strong>spec</strong></td>
			<td><strong>transit</strong></td>
			<td><strong>full</strong></td>
			<td><strong>loose</strong></td>
		</tr>
	<thead>
		<tbody>
		<tr>
<td>

```ts
interface Tweet extends Node {
	id: string
	type: "tweet"
	external html: string
}
```
</td>

<td>

```ts
interface Tweet extends Node {
	id: string
	type: "tweet"
}
```
</td>
<td>

```ts
interface Tweet extends Node {
	id: string
	type: "tweet"
	html: string
}
```
</td>
<td>

```ts
interface Tweet extends Node {
	id: string
	type: "tweet"
	html?: string
}
```
</td>
</tr>
	<thead>
</table>

## Using `content-tree`

### Typescript

This package provides typescript types in `content-tree.d.ts` that can be used to validate the shape of data in a JS/TS application. These types are automatically generated from [SPEC.md](./SPEC.md).

 There are three different namespaces exposed, for the different states a tree can be in (see [In Transit vs At Rest](#in-transit-vs-at-rest))

- `ContentTree.transit` - contains only the fields that are published and available from the Content API's response
- `ContentTree.full` - contains the full representation of the content, including any data required from external resources (also exposed on the top level `ContentTree` namespace)
- `ContentTree.loose` - contains the full representation of the content, including any data required from external resources as optional 


1. Install this repository as a dependency:

```sh notangle
npm install https://github.com/Financial-Times/content-tree
```

2. Use it in your Typescript / JSDoc code:

```ts notangle
import type { ContentTree } from "@financial-times/content-tree"

function makeBigNumber(): ContentTree.BigNumber {
	return {
		type: "big-number", //will autocomplete in code editor, and be valid
		number: "1.2m",
		description: "People affected worldwide"
	}
}

function makeImageSetNoFixins(): ContentTree.transit.ImageSet {
	return {
		type: "image-set",
		id: "79acd774-6ca7-487d-a257-cbf64d2498d9",
		// if you try to add a `picture` here it will get mad
	}
}
```

### JSON Schema

There are also a few JSON schemas generated from the spec. 

- `content-tree.schema.json` - JSON schema for the full content-tree shape
- `transit-tree.schema.json` - JSON schema for content-tree excluding any `external` properties
- `body-tree.schema.json` - JSON schema for the transit tree _without_ the `Root` node, containing just the [Body](./SPEC.md#body) definition. This schema defines the data returned in the `bodyTree` field of the FT `/content-tree` API.


### Go Libraries

These libraries are designed for internal use by the Content & Metadata teams, and are used to ensure that all of our Content API representations are aligned with changes in the content-tree spec.

- `from-bodyxml` - converts the legacy `bodyXML` field from C&M's internal representation, to a valid `bodyTree` JSON.
- `to-external-bodyxml` - converts content-tree to a stable XML representation, used as a the `bodyXML` field in the FT's `/enrichedcontent` API
- `to-string` - converts content-tree to a plain text string

(TODO: how to install / use)


## Contributing

To make a change to the content tree spec:

- Clone this repo and run `npm install`
- Update [SPEC.md](./SPEC.md) with your changes:
	- To add a new formatting node, add the definition under the [`Formatting Blocks`](./SPEC.md#formatting-blocks). If it is formatting that can be applied to text in a paragraph, ensure it is added to the [`Phrasing`](./SPEC.md#phrasing) type
	- To add a new storyblock, add the definition under the [`Storyblocks`](./SPEC.md#storyblocks). If the block can appear at the top level of the article body, ensure it is also added to the [`BodyBlock`](./SPEC.md#bodyblock) type definition
- Run `npm run build` to update `content-tree.d.ts` and (if required) the `schemas` files

Once the PR is created, liaise with the [Content & Metadata](https://biz-ops.in.ft.com/Team/content) team to ensure the relevant changes are made in the Go libraries and transformers.

For major or non-standard changes, consider creating an issue first, or discussing in the `#content-pipeline` Slack channel.

## Releasing

(TODO: no good release process yet)


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
