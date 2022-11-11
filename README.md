![ftcast][logo]
===============

**F**inancial Times **C**ontent **A**bstract **S**yntax **T**ree format

***

**ftcast** is a specification for representing Financial Times article content as an abstract [syntax tree][syntax-tree].
It implements the **[unist][]** spec.


## Contents

*   [Introduction](#introduction)
*   [Types](#types)
*   [Nodes](#nodes)
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

### List

```idl
interface List <: Parent {
  type: "list",
  ordered: boolean,
  children: [ListItem]
}
```

### ListItem

```idl
interface ListItem <: Parent {
  type: "listItem",
  children: [Phrasing | Link]
}
```

### Blockquote

### PullQuote

TODO: PullQuoteSource is a plain string, so make it a property of recommended

### PullQuoteImage

### PullQuoteText

TODO: can contain Paragraph only

### LayoutContainer

### Layout

### LayoutSlot

### LayoutImage

### Recommended

TODO: RecommendedTitle is a plain string, so make it a property of recommended

### ImageSet

### Tweet

### Flourish

### BigNumber

```idl
interface BigNumber <: Parent {
  type: "bigNumber",
  children: [Phrasing | Link]
}
```

### ScrollableBlock

TODO: this has so many rules and children. can we make it simpler as part of this?????

### TODO: define all heading types as straight-up Nodes (like, Chapter y SubHeading y et cetera)

### TODO: promo-box??? podcast promo? concept? content?????? do we allow inline img, b, u?

### TODO: `Table`

```idl
interface Table <: Parent {
  type: "table",
  children: [Caption | TableHead | TableBody]
}
```

A **Table** represents 2d data.

wip. look here https://github.com/Financial-Times/body-validation-service/blob/master/src/main/resources/xsd/ft-html-types.xsd#L214


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
