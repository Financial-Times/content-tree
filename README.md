![ftcast][logo]
===============

**F**inancial Times **C**ontent **A**bstract **S**yntax **T**ree format

***

**ftcast** is a specification for representing Financial Times article content as an abstract [syntax tree][syntax-tree].
It implements the **[unist][]** spec.

## Contents

*   [Introduction](#introduction)
*   [Nodes](#nodes)
*   [License](#license)

## Introduction

This document defines a format for representing Financial Times article content as an [abstract syntax tree][syntax-tree].
This specification is written in a [Web IDL][webidl]-like grammar.

### What is ftcast?

ftcast extends [unist][], a format for syntax trees, to benefit from its [ecosystem of utilities][utilities].

ftcast relates to [JavaScript][js] in that it has an [ecosystem of utilities][list-of-utilities] for working with compliant syntax trees in JavaScript.
However, ftcast is not limited to JavaScript and can be used in other programming languages.

## Nodes

### `Parent`

```idl
interface Parent <: UnistParent {
  children: [Parent]
}
```

**Parent** (**[UnistParent][dfn-unist-parent]**) represents a node in ftcast containing other nodes (said to be *[children][term-child]*).

Its content is limited to only other ftcast content.

### `Parent`

```idl
interface Parent <: UnistParent {
  children: [Parent]
}
```


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
[term-leaf]: https://github.com/syntax-tree/unist#leaf
