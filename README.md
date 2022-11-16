![ftcast][logo]
===============

**F**inancial Times **C**ontent **A**bstract **S**yntax **T**ree format

***

**ftcast** is a specification for representing Financial Times article content as an abstract [syntax tree][syntax-tree].
It implements the **[unist][]** spec.


## Contents

*   [Introduction](#introduction)
*   [Types](#types)
*   [Mixins](#mixins)
*   [Nodes](#nodes)
*   [TODO](#todo)
*   [License](#license)


## Introduction

This document defines a format for representing Financial Times article content as an [abstract syntax tree][syntax-tree].
This specification is written in a [Web IDL][webidl]-like grammar.


### What is ftcast?

ftcast extends [unist][], a format for syntax trees, to benefit from its [ecosystem of utilities][utilities].

ftcast relates to [JavaScript][js] in that it has an [ecosystem of utilities][list-of-utilities] for working with compliant syntax trees in JavaScript.
However, ftcast is not limited to JavaScript and can be used in other programming languages.


## Types

These abstract helper types define special types a [Parent][#parent] can use as [children][#term-child].

### `Node`

```idl
type Node = UnistNode
```

The abstract node.


### `Phrasing`

```idl
type Phrasing = Text | Strong | Emphasis | Strikethrough | Break
```


### `Transparent`

**Transparent** children can contain whatever their parent can contain.

This is used to prohibit nested links.


## Mixins

### Content


```idl
interface mixin Content {
	id: string,
	url?: string,
	content?: Node
}
```

The **Content** mixin is used by nodes that represent a piece of external content. Examples include [Tweet](#tweet) and [ImageSet](#ImageSet).

- TODO: should `id` and the optional `url` both be here? or just `id`? is a tweet's url its id or does it have another? 
what about flourish charts and imagesets etc? perhaps only one required `id` field is needed.


###

## Nodes

### `Parent`

```idl
interface Parent <: UnistParent {
  children: [Node]
}
```

**Parent** (**[UnistParent][dfn-unist-parent]**) represents a node in ftcast containing other nodes (said to be *[children][term-child]*).

Its content is limited to only other ftcast content.


### `Literal`

```idl
interface Literal <: UnistLiteral {
  value: string
}
```

**Literal** (**[UnistLiteral][dfn-unist-literal]**) represents a node in ftcast containing a value.


### `Root`

```idl
interface Root <: Parent {
  type: "root"
}
```

**Root** (**[Parent][dfn-parent]**) represents a document.

**Root** can be used as the *[root][term-root]* of a *[tree][term-tree]*. Always a *[root][term-root]*, never a *[child][term-child]*.


### `Text`

```idl
interface Text <: Literal {
  type: "text"
}
```

**Text** (**[Literal][dfn-literal]**) represents text.


### `Break`

```idl
interface Break <: Node {
  type: "break"
}
```

**Break** Node represents a break in the text, such as in a poem.


### `ThematicBreak`

```idl
interface ThematicBreak <: Node {
  type: "thematicBreak"
}
```

**ThematicBreak** Node represents a break in the text, such as in a shift of topic within a section.

_Non-normative note: this would be represented by an `<hr>` in the html._


### `Paragraph`

```idl
interface Paragraph <: Parent {
  type: "paragraph",
  children: [Phrasing | Link]
}
```

A **Paragraph** represents a unit of text.


### `Strong`

```idl
interface Strong <: Parent {
  type: "strong",
  children: [Transparent] 
}
```

A **Strong** node represents contents with strong importance, seriousness or urgency.


### `Emphasis`

```idl
interface Emphasis <: Parent {
  type: "emphasis"
  children: [Transparent]
}
```

An **Emphasis** node represents stressed emphasis of its contents.


### `Link`

```idl
interface Link <: Parent {
  type: "link",
  children: [Phrasing]
}
```

A **Link** represents a hyperlink.

### `List`

```idl
interface List <: Parent {
  type: "list",
  ordered: boolean,
  children: [ListItem]
}
```

An **List** node represents a list of items.

### `ListItem`

```idl
interface ListItem <: Parent {
  type: "listItem",
  children: [Phrasing | Link]
}
```

### `Blockquote`

```idl
interface PullQuote <: Parent {
  type: "pullQuote",
  citation?: string,
  children: [Phrasing]
}
```

A **BlockQuote** represents a quotation and optional citation.


### `PullQuote`

```idl
interface PullQuote <: Parent {
  type: "pullQuote",
  citation: string,
  children: [PullQuoteImage | PullQuoteText]
}
```

A **PullQuote** node represents a brief quotation taken from the main text of an article.

- TODO: make sure all the casing of these is consistent with C&M's casing.
- TODO: Spark doesn't seem to have a concept of PullQuoteImage, and the text appears to only be a string. maybe PullQuote should only contain Paragraph nodes rather than a PullQuoteText containing a Paragraph node.


### `PullQuoteImage`

```idl
interface PullQuoteImage <: Node {
  type: "pullQuoteImage",
  source: string
}
```

- TODO: what's all this then?


### `PullQuoteText`

```idl
interface PullQuote <: Parent {
  type: "pullQuoteText",
  citation: string,
  children: [Paragraph]
}
```

- TODO: see [pullquote](#pullquote)


### `Recommended`

```idl
interface Recommended <: Parent {
  type: "recommended",
  title?: "string",
  children: [List]
}
```

- A **Recommended** node represents a list of recommended links.


### `ImageSet`

```idl
interface ImageSet <: Node {
  type: "imageSet",
  content?: ImageSetContent
}

ImageSet includes Content
```


### `ImageSetContent`

- TODO: get the de-referenced imageset shape from cp and define


### `Tweet`

```idl
interface Tweet <: Node {
  type: "tweet",
  content?: TweetContent
}

Tweet includes Content
```

- TODO: get the de-referenced tweet shape from cp and define


### `Flourish`

```idl
interface Flourish <: Node {
  type: "flourish",
  flourishType: string
  content?: FlourishContent
}

Flourish includes Content
```

- TODO: is flourish-type a thing here?
- TODO: get the de-referenced flourish shape from cp and define

### `BigNumber`

```idl
interface BigNumber <: Node {
  type: "bigNumber",
  number: Paragraph,
  description: Paragraph
}
```

A **BigNumber** node is used to provide a description for a big number.

### `ScrollableBlock`

```idl
interface ScrollableBlock <: Parent {
  type: "scrollableBlock",
  theme: "sans" | "serif",
  children: [ScrollableSection]
}
```

A **ScrollableBlock** node represents a block for telling stories through scroll position.

### `ScrollableSection`

```idl
interface ScrollableSection <: Parent {
  type: "scrollableSection",
  display: "dark" | "light"
  position: "left" | "centre" | "right"
  transition?: "delay-before" | "delay-after"
  noBox?: boolean
  children: [ImageSet | ScrollableText]
}
```

A **ScrollableBlock** node represents a section of a [ScrollableBlock](#scrollableblock)

- TODO: why is noBox not a display option? like "dark" | "light" | "transparent"?
- TODO: does this need to be more specific about its children?
- TODO: should each section have 1 `imageSet` field and then children of any number of ScrollableText?
- TODO: could `transition` have a `"none"` value so it isn't optional?

### `ScrollableText`

```idl
interface ScrollableHeading <: Parent {
  type: "scrollableHeading",
  style: "chapter" | "heading" | "subheading" | "text"
  children: [Paragraph]
}
```

A **ScrollableBlock** node represents a piece of copy for a [ScrollableBlock](#scrollableblock)

- TODO: heading doesn't 
- TODO: i'm a little confused by this part of the spec, i need to look at some scrollable-text blocks
https://github.com/Financial-Times/body-validation-service/blob/fddc5609b15729a0b60e06054d1b7749cc70c62b/src/main/resources/xsd/ft-types.xsd#L224-L263

## TODO

- define all heading types as straight-up Nodes (like, Chapter y SubHeading y et cetera)
- do we need an `HTML` node that has a raw html string to __dangerously insert like markdown for some embed types?
- promo-box??? podcast promo? concept? ~content??????~ do we allow inline img, b, u? (spark doesn't. maybe no. what does this mean for embeds?)
- define all the Experimental things like ImagePair

### TODO: `LayoutContainer`

TODO: what is this container for? why does the data need a container in addition to the Layout?

### TODO: `Layout`

TODO: Don't know anything about layouts or how the data is shaped.

### TODO: `LayoutSlot`

### TODO: `LayoutImage`

### TODO: `Table`

```idl
interface Table <: Parent {
  type: "table",
  children: [Caption | TableHead | TableBody]
}
```

A **Table** represents 2d data.

look here https://github.com/Financial-Times/body-validation-service/blob/master/src/main/resources/xsd/ft-html-types.xsd#L214

maybe we can be more strict than this? i don't know. we might not be able to because we don't know what old articles have done. however, we could find out what old articles have done... we could validate all old articles by trying to convert their bodyxml to this format, validating them etc,... and then make changes. maybe we want to restrict old articles from being able to do anything Spark can't do? who knows. we need more eyes on this whole document.


## License

This software is published by the Financial Times under the [MIT licence](mit).

Derived from [unist][unist] © [Titus Wormer][titus]

[mit]: http://opensource.org/licenses/MIT
[ideas]: https://github.com/syntax-tree/ideas
[titus]: https://wooorm.com
[logo]: ./logo.png
[unist]: https://github.com/syntax-tree/unist
[syntax-tree]: https://github.com/syntax-tree/unist#syntax-tree
[js]: https://www.ecma-international.org/ecma-262/9.0/index.html
[webidl]: https://heycam.github.io/webidl/
[term-tree]: https://github.com/syntax-tree/unist#tree
[term-child]: https://github.com/syntax-tree/unist#child
[term-root]: https://github.com/syntax-tree/unist#root
[term-leaf]: https://github.com/syntax-tree/unist#lea
