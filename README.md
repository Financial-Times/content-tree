![content-tree][logo]
===============

A tree for Financial Times article content.

***

**content-tree** is a specification for representing Financial Times article content as an abstract tree.
It implements the **[unist][unist]** spec.


## Contents

* [Introduction](#introduction)
* [Types](#types)
* [Mixins](#mixins)
* [Nodes](#nodes)
* [TODO](#todo)
* [License](#license)


## Introduction

This document defines a format for representing Financial Times article content as a tree.
This specification is written in a [Web IDL][webidl]-like grammar.


### What is `content-tree`?

`content-tree` extends [unist][unist], a format for syntax trees, to benefit from its [ecosystem of utilities][unist-utilities].

`content-tree` relates to [JavaScript][js] in that it has an [ecosystem of utilities][unist-utilities] for working with trees in JavaScript.
However, `content-tree` is not limited to JavaScript and can be used in other programming languages.


## Types

These abstract helper types define special types a [Parent](#parent) can use as [children][term-child].

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

## Nodes

### `Parent`

```idl
interface Parent <: UnistParent {
  children: [Node]
}
```

**Parent** (**[UnistParent][term-parent]**) represents a node in content-tree containing other nodes (said to be *[children][term-child]*).

Its content is limited to only other content-tree content.


### `Literal`

```idl
interface Literal <: UnistLiteral {
  value: string
}
```

**Literal** (**[UnistLiteral][term-literal]**) represents a node in content-tree containing a value.


### `Root`

```idl
interface Root <: Parent {
  type: "root"
}
```

**Root** (**[Parent][term-parent]**) represents a document.

**Root** can be used as the *[root][term-root]* of a *[tree][term-tree]*. Always a *[root][term-root]*, never a *[child][term-child]*.


### `Text`

```idl
interface Text <: Literal {
  type: "text"
}
```

**Text** (**[Literal][term-literal]**) represents text.


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

### `Chapter`

```idl
interface Chapter <: Parent {
  type: "chapter",
  children: [Text]
}
```

A **Chapter** represents a chapter-level heading.

### `Heading`

```idl
interface Heading <: Parent {
  type: "heading",
  children: [Text]
}
```

A **Heading** represents a heading-level heading.

### `Subheading`

```idl
interface Subheading <: Parent {
  type: "subheading",
  children: [Text]
}
```

A **Subheading** represents a subheading-level heading.

### `Label`

```idl
interface Label <: Parent {
  type: "label",
  children: [Text]
}
```

A **Label** represents a label-level heading.

- TODO: is this name ok?

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
  url: string,
  title: string,
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

}
```

- A **Recommended** node represents a list of recommended links.
- TODO: this has a list of things and the list items are 


### `ImageSet`

```idl
interface ImageSet <: Node {
  type: "imageSet",
  content?: ImageSetContent
}

ImageSet includes Content
```


### `ImageSetContent`

```idl
interface ImageSetContent <: Node {
  type: "imageSetContent",
  alt: string,
  caption?: string,
  imageType: "Image" | "Graphic",
  images: [Image]
}
```

- TODO: should we be using the full url as the `image`/`graphic` (like 'http://www.ft.com/ontology/content/Image')? might be better

### `Image`

- TODO: we want this to look like this [https://raw.githubusercontent.com/Financial-Times/cp-content-pipeline/main/packages/schema/src/picture.ts](https://github.com/Financial-Times/cp-content-pipeline/blob/main/packages/schema/src/picture.ts#L12-L99)
- TODO: should i call this `Picture`???? maybe.

### `TweetReference`

```idl
interface TweetReference <: Node {
  type: "tweetReference",
  id: string
}
```

A **TweetReference** node represents a reference to an external tweet. The `id` is a URL.

### `Tweet`

```idl
interface Tweet <: Node {
  type: "tweet",
  id: string,
  children: [Phrasing]
}
```

A **Tweet** node represents a tweet.

TODO: what are the valid children here? 


### `FlourishReference`

```idl
interface FlourishReference <: Node {
  type: "flourishReference",
  id: string,
  flourishType: string
}
```

A **FlourishReference** node represents a reference to an external **Flourish**.

### `Flourish`

```idl
interface Flourish <: Node {
  type: "flourish",
  id: string,
  fullGrid: boolean,
  flourishType: ,
  description: string,
  fallbackImage: TODO
}
```

A **FlourishReference** node represents a Flourish.

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
- TODO: rather than this "style" property on ScrollableText, what if we made these the same Paragraph, Chapter, Heading and Subheading nodes as above?

## TODO

- define all heading types as straight-up Nodes (like, Chapter y SubHeading y et cetera)
- do we need an `HTML` node that has a raw html string to __dangerously insert like markdown for some embed types?
- promo-box??? podcast promo? concept? ~content??????~ do we allow inline img, b, u? (spark doesn't. maybe no. what does this mean for embeds?)

### TODO: `LayoutContainer`

TODO: what is this container for? why does the data need a container in addition to the Layout?

### TODO: `Layout`### TODO: `LayoutSlot`### TODO: `LayoutImage`

TODO: okay so we're going to not do this ! we'll be defining ImagePair, Timeline, etc 

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
