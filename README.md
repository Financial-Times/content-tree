# ![content-tree][logo]

A tree for Financial Times article content.

---

**content-tree** is a specification for representing Financial Times article content as an abstract tree.
It implements the **[unist][unist]** spec.

## Contents

- [Introduction](#introduction)
- [Types](#types)
- [Mixins](#mixins)
- [Nodes](#nodes)
- [TODO](#todo)
- [License](#license)

## Introduction

This document defines a format for representing Financial Times article content as a tree.
This specification is written in a [typescript][typescript] grammar.

### What is `content-tree`?

`content-tree` extends [unist][unist], a format for syntax trees, to benefit from its [ecosystem of utilities][unist-utilities].

`content-tree` relates to [JavaScript][js] in that it has an [ecosystem of utilities][unist-utilities] for working with trees in JavaScript.
However, `content-tree` is not limited to JavaScript and can be used in other programming languages.

## Types

These abstract helper types define special types a [Parent](#parent) can use as [children][term-child].

### `Block`

```ts
type Block = Node // TODO
```

### `Phrasing`

```ts
type Phrasing = Text | Break | Strong | Emphasis | Strikethrough | Link
```

Phrasing nodes cannot have an ancestor of their same type.

TODO: clarify that i mean Strong cannot have an ancestor of Strong etc

## Nodes

### `Node`

```ts
interface Node {
	type: string
}
```

The abstract node.

### `Parent`

```ts
interface Parent extends Node {
	children: Node[]
}
```

**Parent** (**[UnistParent][term-parent]**) represents a node in content-tree containing other nodes (said to be _[children][term-child]_).

Its content is limited to only other content-tree content.

### `Literal`

```ts
interface Literal extends Node {
	value: any
}
```

**Literal** (**[UnistLiteral][term-literal]**) represents a node in content-tree containing a value.

### `Reference`

```ts
interface Reference extends Node {
	type: "reference"
	referencedType: string
	id: string
	alt?: string
}
```

**Reference** nodes represent a reference to a piece of external content. The `alt` field is an optional string to be used if the external resource was not available. The `kind` field is the `type` of the node that the reference dereferences to.

### `Root`

```ts
interface Root extends Parent {
	type: "root"
	children: [Body]
}
```

**Root** (**[Parent][term-parent]**) represents the root of a content-tree.

**Root** can be used as the _[root][term-root]_ of a _[tree][term-tree]_.

### `Body`

```ts
interface Body extends Parent {
	type: "body"
	version: number
	children: Block[]
}
```

**Body** (**[Parent][term-parent]**) represents the body of an article.

(note: `bodyTree` is just this part)

### `Text`

```ts
interface Text extends Literal {
	type: "text"
	value: string
}
```

**Text** (**[Literal][term-literal]**) represents text.

### `Break`

```ts
interface Break extends Node {
	type: "break"
}
```

**Break** Node represents a break in the text, such as in a poem.

_Non-normative note: this would be represented by a `<br>` in the html._

### `ThematicBreak`

```ts
interface ThematicBreak extends Node {
	type: "thematicBreak"
}
```

**ThematicBreak** Node represents a break in the text, such as in a shift of topic within a section.

_Non-normative note: this would be represented by an `<hr>` in the html._

### `Paragraph`

```ts
interface Paragraph extends Parent {
	type: "paragraph"
	children: Phrasing[]
}
```

A **Paragraph** represents a unit of text.

### `Chapter`

```ts
interface Chapter extends Parent {
	type: "chapter"
	children: Text[]
}
```

A **Chapter** represents a chapter-level heading.

### `Heading`

```ts
interface Heading extends Parent {
	type: "heading"
	children: Text[]
}
```

A **Heading** represents a heading-level heading.

### `Subheading`

```ts
interface Subheading extends Parent {
	type: "subheading"
	children: Text[]
}
```

**Subheading** represents a subheading-level heading.

### `Label`

```ts
interface Label extends Parent {
	type: "label"
	children: Text[]
}
```

**Label** represents a label-level heading.

- TODO: is this name ok?

### `Strong`

```ts
interface Strong extends Parent {
	type: "strong"
	children: Phrasing[]
}
```

**Strong** represents contents with strong importance, seriousness or urgency.

### `Emphasis`

```ts
interface Emphasis extends Parent {
	type: "emphasis"
	children: Phrasing[]
}
```

**Emphasis** represents stressed emphasis of its contents.

### `Strikethrough`

```ts
interface Strikethrough extends Parent {
	type: "strikethrough"
	children: Phrasing[]
}
```

**Strikethrough** represents a piece of text that has been stricken.

### `Link`

```ts
interface Link extends Parent {
	type: "link"
	url: string
	title: string
	children: Phrasing[]
}
```

**Link** represents a hyperlink.

### `List`

```ts
interface List extends Parent {
	type: "list"
	ordered: boolean
	children: ListItem[]
}
```

**List** represents a list of items.

### `ListItem`

```ts
interface ListItem extends Parent {
	type: "listItem"
	children: Phrasing[]
}
```

### `Blockquote`

```ts
interface Blockquote extends Parent {
	type: "blockquote"
	children: Phrasing[]
}
```

**BlockQuote** represents a quotation.

### `PullQuote`

```ts
interface PullQuote extends Parent {
	type: "pullQuote"
	children: [PullQuoteText, PullQuoteSource]
}
```

**PullQuote** represents a brief quotation taken from the main text of an article.

### `PullQuoteText`

```ts
interface PullQuoteText extends Parent {
	type: "pullQuoteText"
	children: Text[]
}
```

**PullQuoteText** represents the text of a pullquote.

### `PullQuoteSource`

```ts
interface PullQuoteSource extends Parent {
	type: "pullQuoteSource"
	children: Text[]
}
```

**PullQuoteText** represents the source of a pullquote.

### `Recommended`

```ts
interface Recommended extends Parent {
	type: "recommended"
	children: [/*TODO*/]
}
```

- A **Recommended** node represents a list of recommended links.
- TODO: this has a list of things and the list items are ...?

### `ImageSetReference`

```ts
interface ImageSetReference extends Reference {
	kind: "imageSet"
	imageType: "Image" | "Graphic"
}
```

A **ImageSetReference** node represents a reference to an external tweet. The `id` is a URL.

### `ImageSet`

```ts
interface ImageSet extends Node {
	type: "imageSet"
	alt: string
	caption?: string
	imageType: "Image" | "Graphic"
	images: Image[]
}
```

- TODO: should we be using the full url as the `image`/`graphic` (like 'http://www.ft.com/ontology/content/Image')? might be better

### `Image`

```ts
interface Image extends Node {
	type: "image"
}
```

- TODO: we want this to look like this [https://raw.githubusercontent.com/Financial-Times/cp-content-pipeline/main/packages/schema/src/picture.ts](https://github.com/Financial-Times/cp-content-pipeline/blob/main/packages/schema/src/picture.ts#L12-L99)
- TODO: should i call this `Picture`???? maybe.

### `TweetReference`

```ts
interface TweetReference extends Reference {
	kind: "tweet"
}
```

A **TweetReference** node represents a reference to an external tweet. The `id` is a URL.

### `Tweet`

```ts
interface Tweet extends Node {
	type: "tweet"
	id: string
	children: Phrasing[]
}
```

A **Tweet** node represents a tweet.

TODO: what are the valid children here? Should we allow a tweet to contain a hast document root as its child?

### `FlourishReference`

```ts
interface FlourishReference extends Reference {
	kind: "flourish"
	flourishType: string
}
```

A **FlourishReference** node represents a reference to an external **Flourish**.

### `Flourish`

```ts
interface Flourish extends Node {
	type: "flourish"
	id: string
	layoutWidth: "" | "full-grid"
	flourishType: string
	description: string
	fallbackImage: Image
}
```

A **Flourish** node represents a flourish chart.

### `BigNumber`

```ts
interface BigNumber extends Parent {
	type: "bigNumber"
	children: [BigNumberNumber, BigNumberDescription]
}
```

**BigNumber** provides a description for a big number.

### `BigNumberNumber`

```ts
interface BigNumberNumber extends Parent {
	type: "bigNumberNumber"
	children: Phrasing[]
}
```

**BigNumberNumber** represents the number itself.

### `BigNumberDescription`

```ts
interface BigNumberDescription extends Parent {
	type: "bigNumberDescription"
	children: Phrasing[]
}
```

**BigNumberNumber** represents the description of the big number.

### `ScrollableBlock`

```ts
interface ScrollableBlock extends Parent {
	type: "scrollableBlock"
	theme: "sans" | "serif"
	children: ScrollableSection[]
}
```

A **ScrollableBlock** node represents a block for telling stories through scroll position.

### `ScrollableSection`

```ts
interface ScrollableSection extends Parent {
	type: "scrollableSection"
	display: "dark" | "light"
	position: "left" | "centre" | "right"
	transition?: "delay-before" | "delay-after"
	noBox?: boolean
	children: Array<ImageSet | ScrollableText>
}
```

- TODO: define these children properly

A **ScrollableBlock** node represents a section of a [ScrollableBlock](#scrollableblock)

- TODO: why is noBox not a display option? like "dark" | "light" | "transparent"?
- TODO: does this need to be more specific about its children?
- TODO: should each section have 1 `imageSet` field and then children of any number of ScrollableText?
- TODO: could `transition` have a `"none"` value so it isn't optional?

### `ScrollableText`

```ts
interface ScrollableText extends Parent {
	type: "scrollableText"
	style: "text"
	children: Phrasing[]
}

interface ScrollableHeading extends Parent {
	type: "scrollableText"
	style: "chapter" | "heading" | "subheading"
	children: Text[]
}
```

A **ScrollableBlock** node represents a piece of copy for a [ScrollableBlock](#scrollableblock)

- TODO: heading doesn't
- TODO: i'm a little confused by this part of the spec, i need to look at some scrollable-text blocks
  https://github.com/Financial-Times/body-validation-service/blob/fddc5609b15729a0b60e06054d1b7749cc70c62b/src/main/resources/xsd/ft-types.xsd#L224-L263
- TODO: rather than this "style" property on ScrollableText, what if we made these the same Paragraph, Chapter, Heading and Subheading nodes as above?

## TODO

- define all heading types as straight-up Nodes (like, Chapter y SubHeading y et cetera)
- do we need an `HTML` node that has a raw html string to \_\_dangerously insert like markdown for some embed types? <-- YES
- promo-box??? podcast promo? concept? ~content??????~ do we allow inline img, b, u? (spark doesn't. maybe no. what does this mean for embeds?)

### TODO: `LayoutContainer`

TODO: what is this container for? why does the data need a container in addition to the Layout?

### TODO: `Layout`### TODO: `LayoutSlot`### TODO: `LayoutImage`

TODO: okay so we're going to not do this ! we'll be defining ImagePair, Timeline, etc

### TODO: `Table`

```ts
interface Table extends Parent {
	type: "table"
	children: [Caption | TableHead | TableBody]
}

interface Caption {
	type: "caption"
}
interface TableHead {
	type: "tableHead"
}
interface TableBody {
	type: "tableBody"
}
```

A **Table** represents 2d data.

look here https://github.com/Financial-Times/body-validation-service/blob/master/src/main/resources/xsd/ft-html-types.xsd#L214

maybe we can be more strict than this? i don't know. we might not be able to because we don't know what old articles have done. however, we could find out what old articles have done... we could validate all old articles by trying to convert their bodyxml to this format, validating them etc,... and then make changes. maybe we want to restrict old articles from being able to do anything Spark can't do? who knows. we need more eyes on this whole document.

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
