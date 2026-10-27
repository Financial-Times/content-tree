# XML to Content Tree Transformer

## Overview

The Transformer converts external XHTML-formatted document into content tree. It supports the bodyXML format used in the main content store within the Content & Metadata platform — specifically, in the **internalComponent** collection.

## Usage

```go
package main

import (
    "fmt"
    "log"

    tocontenttree "github.com/Financial-Times/content-tree"
)

func main() {
    xmlInput := `<body><p>Hello World</p></body>`

    out, err := tocontenttree.Transform(xmlInput)
    if err != nil {
        log.Fatalf("Transform (XmlToTree) failed: %v", err)
    }

    fmt.Printf("Transformed content tree: %+v\n", out)
}
```

## Known Limitations and Behavior

The current implementation of the transformer has the following limitations:

- content-tree only components, such as Carousel, have no legacy internal BodyXML representation and are not supported by this transformer. Adding a content-tree node does not automatically add BodyXML import support.
- If the transformer encounters an HTML tag that does not have a corresponding definition in the content tree, that tag is skipped.
- If an HTML element contains child elements that are not allowed, those disallowed children are ignored.
