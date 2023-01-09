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

### `SourceSet`

```ts
interface ImageSource {
	url: string
	width: number
	dpr: number
}
```

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
interface Text extends Node {
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
	type: "thematic-break"
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

Paragraph represents a unit of text.

### `Heading`

```ts
interface Heading extends Parent {
	type: "heading"
	children: Text[]
	level: "chapter" | "subheading" | "label"
}
```

**Heading** represents a unit of text that marks the beginning of an article section.

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
	type: "list-item"
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

### `Pullquote`

```ts
interface Pullquote extends Parent {
	type: "pullquote"
	text: string
	source: string
}
```

**Pullquote** represents a brief quotation taken from the main text of an article.

_non normative note:_ the reason this is string properties and not children is that it is more confusing if a pullquote falls back to text than if it doesn't. The text is taken from elsewhere in the article.

### `Recommended`

```ts
interface Recommended extends Parent {
	type: "recommended"
	children: [/*TODO*/]
}
```

- Recommended represents a list of recommended links.
- TODO: this has a list of things and the list items are ...?

### `ImageSetReference`

```ts
interface ImageSetReference extends Reference {
	referencedType: "image-set"
	layoutWidth: "inset-left" | "full-text" | "full-grid" | "full-bleed"
}
```

ImageSetReference represents a reference to an external tweet. The `id` is a URL.

### `ImageSet`

```ts
interface ImageSet extends Node {
	type: "image-set"
	id: string
	imageType: "graphic" | "image"
	layoutWidth: "inset-left" | "full-text" | "full-grid" | "full-bleed"
	alt: string
	caption: string
	credit: string
	images: Image[]
	// fallbackImage: ???
}
```

### `Image`

// TODO why "originalWidth" "originalHeight" and "binaryUrl" rather than "width", "height" and "url"?

```ts
interface Image extends Node {
	type: "image"
	id: string
	originalWidth: number
	originalHeight: number
	format:
		| "desktop"
		| "mobile"
		| "square"
		| "standard"
		| "wide"
      | "standard-inline"
	binaryUrl: "string"
	sourceSet: ImageSource[]
}
```

### `TweetReference`

```ts
interface TweetReference extends Reference {
	referencedType: "tweet"
}
```

**TweetReference** represents a reference to an external tweet. The `id` is a URL.

### `Tweet`

```ts
interface Tweet extends Node {
	type: "tweet"
	id: string
	html: string
}
```

**Tweet** represents a tweet.

### `FlourishReference`

```ts
interface FlourishReference extends Reference {
	referencedType: "flourish"
	flourishType: string
	layoutWidth: "full-text" | "full-grid"
	description: string
	timestamp: string
}
```

**FlourishReference** represents a reference to an external **Flourish**.

### `Flourish`

```ts
interface Flourish extends Node {
	type: "flourish"
	id: string
	layoutWidth: "full-text" | "full-grid"
	flourishType: string
	description: string
	timestamp: string
	fallbackImage: Image
}
```

**Flourish** represents a flourish chart.

### `BigNumber`

```ts
interface BigNumber extends Parent {
	type: "big-number"
	children: [BigNumberNumber, BigNumberDescription]
}

// TODO: consider making these children two paragraphs
```

**BigNumber** provides a description for a big number.

### `BigNumberNumber`

```ts
interface BigNumberNumber extends Parent {
	type: "big-number-number"
	children: Phrasing[]
}
```

**BigNumberNumber** represents the number itself.

### `BigNumberDescription`

```ts
interface BigNumberDescription extends Parent {
	type: "big-number-description"
	children: Phrasing[]
}
```

**BigNumberDescription** represents the description of the big number.

### `ScrollyBlock`

```ts
interface ScrollyBlock extends Parent {
	type: "scrolly-block"
	theme: "sans" | "serif"
	children: ScrollySection[]
}
```

**ScrollyBlock** represents a block for telling stories through scroll position.

### `ScrollySection`

```ts
interface ScrollySection extends Parent {
	type: "scrolly-section"
	theme: "dark-text" | "light-text" | "dark-text-no-box" | "light-text-no-box"
	position: "left" | "center" | "right"
	transition?: "delay-before" | "delay-after"
	children: [ImageSet, ...ScrollyCopy[]]
}
```

**ScrollySection** represents a section of a [ScrollyBlock](#scrollyblock)

- TODO: could `transition` have a `"none"` value so it isn't optional?

### `ScrollyCopy`

```ts
interface ScrollyCopy extends Parent {
	type: "scrolly-copy"
	children: ScrollyText[]
}
```

TODO is this badly named?

**ScrollyCopy** represents a collection of **ScrollyText** nodes.

```ts
interface ScrollyText extends Parent {
	type: "scrolly-text"
	level: string
}

interface ScrollyHeading extends ScrollyText {
	type: "scrolly-text"
	level: "chapter" | "heading" | "subheading"
	children: Text[]
}

interface ScrollyParagraph extends ScrollyText {
	type: "scrolly-text"
	level: "text"
	children: Phrasing[]
}
```

**ScrollyText** represents an individual unit of copy for a [ScrollyBlock](#scrollableblock)

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
	type: "table-head"
}
interface TableBody {
	type: "table-body"
}
```

**Table** represents 2d data.

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
