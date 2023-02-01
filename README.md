# ![content-tree][logo]

A tree for Financial Times article content.

---

**content-tree** is a specification for representing Financial Times article
content as an abstract tree.  It implements the **[unist][unist]** spec.

## Contents

- [Introduction](#introduction)
- [Types](#types)
- [Mixins](#mixins)
- [Nodes](#nodes)
- [TODO](#todo)
- [License](#license)

## Introduction

This document defines a format for representing Financial Times article content
as a tree.  This specification is written in a [typescript][typescript] grammar.

### What is `content-tree`?

`content-tree` extends [unist][unist], a format for syntax trees, to benefit
from its [ecosystem of utilities][unist-utilities].

`content-tree` relates to [JavaScript][js] in that it has an [ecosystem of
utilities][unist-utilities] for working with trees in JavaScript.  However,
`content-tree` is not limited to JavaScript and can be used in other programming
languages.

## Types

These abstract helper types define special types a [Parent](#parent) can use as
[children][term-child].

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

### `Teaser`

```ts
type TeaserConcept = {
	apiUrl: string
	directType: string
	id: string
	predicate: string
	prefLabel: string
	type: string
	types: string[]
	url: string
}

type TeaserImage = {
	url: string
	width: number
	height: number
}

type Indicators = {
	accessLevel: "premium" | "subscribed" | "registered" | "free"
	isOpinion?: boolean
	isColumn?: boolean
	isPodcast?: boolean
	isEditorsChoice?: boolean
	isExclusive?: boolean
	isScoop?: boolean
}

type Teaser = {
	id: string
	url: string
	type:
		| "article"
		| "video"
		| "podcast"
		| "audio"
		| "package"
		| "liveblog"
		| "promoted-content"
		| "paid-post"
	title: string
	publishedDate: string
	firstPublishedDate: string
	metaLink?: TeaserConcept
	metaAltLink?: TeaserConcept
	metaPrefixText?: string
	metaSuffixText?: string
	indicators: Indicators
	image: Image
}
```

- The above types are adapted from the data structure used by
  [x-teaser](https://github.com/Financial-Times/x-dash/blob/main/components/x-teaser/Props.d.ts),
  limited to the types required for rendering teasers used within a content-tree
  (i.e. recommended links)
- TODO: consider having x-teaser use types from content-tree


### `ImageSet`

```ts
type ImageSetPicture = {
	layoutWidth: string
	imageType: "image" | "graphic"
	alt: string
	caption: string
	credit: string
	images: Image[]
	fallbackImage: Image
}
```

```ts
type Image = {
	id: string
	width: number
	height: number
	format:
		| "desktop"
		| "mobile"
		| "square"
		| "standard"
		| "wide"
		| "standard-inline"
	url: string
	sourceSet?: ImageSource[]
}
```

### `ImageSource`

```ts
type ImageSource = {
	url: string
	width: number
	dpr: number
}
```

**ImageSource** are the shapes of things in an [image](#image)'s sourceSet


## Nodes

### `Node`

```ts
interface Node {
	type: string
	data?: any
}
```

The abstract node. The data field is for internal implementation information and
will never be defined in the content-tree spec.

### `Parent`

```ts
interface Parent extends Node {
	children: Node[]
}
```

**Parent** (**[UnistParent][term-parent]**) represents a node in content-tree
containing other nodes (said to be _[children][term-child]_).

Its content is limited to only other content-tree content.

### `Root`

```ts
interface Root extends Node {
	type: "root"
	body: Body
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

**ThematicBreak** Node represents a break in the text, such as in a shift of
topic within a section.

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

**Heading** represents a unit of text that marks the beginning of an article
section.

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
interface Pullquote extends Node {
	type: "pullquote"
	text: string
	source?: string
}
```

**Pullquote** represents a brief quotation taken from the main text of an article.

_non normative note:_ the reason this is string properties and not children is
that it is more confusing if a pullquote falls back to text than if it
doesn't. The text is taken from elsewhere in the article.

### `Recommended`

```ts
interface Recommended extends Node {
	type: "recommended"
	id: string
	heading?: string
	teaserTitleOverride?: string
	teaser?: Teaser
}
```

- Recommended represents a reference to an FT content that has been recommended
  by editorial.
- The `heading`, when present, is used where the purpose of the link is more
  specific than being "Recommended" (an example might be "In depth")
- The `teaserTitleOverride`, when present, is used in place of the content title
  of the link.

_non normative note:_ historically, recommended links used to be a list of up to
three content items. Testing later showed that having one more prominent link
was more engaging, and Spark (and therefore content-tree)now only supports that use case.


### `ImageSet`

```ts
interface ImageSet extends Node {
	type: "image-set"
	id: string
	picture?: ImageSetPicture
}
```

### `Tweet`

```ts
interface Tweet extends Node {
	id: string
	type: "tweet"
	html?: string
}
```

**Tweet** represents a tweet.

### `Flourish`

```ts
interface Flourish extends Node {
	type: "flourish"
	id: string
	layoutWidth: string
	flourishType: string
	description?: string
	timestamp?: string
	fallbackImage?: Image
}
```

**Flourish** represents a flourish chart.

### `BigNumber`

```ts
interface BigNumber extends Node {
	type: "big-number"
	number: string
	description: string
}
```

**BigNumber** provides a description for a big number.

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

**ScrollyText** represents an individual unit of copy for a
[ScrollyBlock](#scrollableblock)

### `Group`

```ts
interface Group extends Parent {
	type: "group"
	children: Block[]
}
```

**Group** wraps multiple 

### `Grid`

```ts
interface Grid extends Parent {
	type: "grid"
	layoutWidth: string
	rows: number
	columns: number
	children: Block[]
}
```

**Grid** represents content that is presented together in a grid. The children is a 1-dimensional array of blocks.

### `Timeline`

```ts
interface Timeline extends Parent {
	type: "timeline"
	title?: string
	image?: ImageSet
	children: TimelineEvent[]
}
```

### `TimelineEvent`

```ts
interface TimelineEvent extends Parent {
	type: "timeline-event"
	date: string
	children: Block[]
}
```

**Timeline** represents content that takes place over a series of dated events.

### TODO `Table`

**Table** represents 2d data.

look here https://github.com/Financial-Times/body-validation-service/blob/master/src/main/resources/xsd/ft-html-types.xsd#L214

and in Spark, which is more strict.

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
