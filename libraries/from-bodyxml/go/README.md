# XML to Content Tree Transformer

## Overview
The Transformer converts external XHTML-formatted document into content tree.  
It supports format stored in the **internalComponent** collection as well as the one returned by the **Internal Content API**.
The latter is produced by the content-public-read service after applying certain transformations to the bodyXML it retrieves from the internalComponents collection.
These transformations include renaming the content, related, and concept tags to ft-content, ft-related, and ft-concept, respectively, and replacing the id attribute with url, with a few caveats.

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