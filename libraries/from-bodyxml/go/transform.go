package tocontenttree

import (
	"fmt"
	"io"
	"strings"

	contenttree "github.com/Financial-Times/content-tree"
	"github.com/beevik/etree"
)

// Transform converts an external XHTML-formatted document into a content tree.
// It returns an error if the input contains unsupported HTML elements
// or does not comply with the content tree schema.
func Transform(bodyXML string) (*contenttree.Root, error) {
	return fromETreeReader(strings.NewReader(bodyXML))
}

func fromETreeReader(r io.Reader) (*contenttree.Root, error) {
	doc := etree.NewDocument()
	_, err := doc.ReadFrom(r)
	if err != nil {
		return nil, err
	}

	root := doc.Root()
	if root == nil {
		return nil, fmt.Errorf("no root element found")
	}

	m := &contenttree.Body{Type: contenttree.BodyType, Version: 1}
	err = convertToContentTree(root, m)
	if err != nil {
		return nil, err
	}

	out := &contenttree.Root{
		Type: contenttree.RootType,
		Body: m,
	}

	return out, nil
}

func convertToContentTree(elem etree.Token, m contenttree.Node) error {
	switch t := elem.(type) {
	case *etree.Element:
		if t.Tag == "body" {
			for _, child := range t.Child {
				err := convertToContentTree(child, m)
				if err != nil {
					return err
				}
			}
			return nil
		}

		if t.Tag == "content" || t.Tag == "related" || t.Tag == "concept" {
			id := attr(t, "id")
			typeAttr := attr(t, "type")
			if id != "" {
				t.CreateAttr("url", generateUrl(typeAttr, id))
				if attr(t, "data-asset-type") != "flourish" {
					t.RemoveAttr("id")
				}
			}
		}

		tag := t.Tag
		if t.Tag == "content" || t.Tag == "ft-content" {
			for _, attr := range t.Attr {
				if attr.Key == "type" {
					tag = attr.Value
					break
				}
			}
		}

		transformer, ok := defaultTransformers[tag]
		if !ok {
			return fmt.Errorf("unknownNode transformer for tag <%s>", t.Tag)
		}

		switch transformed := transformer(t).(type) {
		case *unknownNode:
			{
				return fmt.Errorf("unknownNode div node with class '%s'", transformed.Class)
			}
		case *liftChildrenNode:
			{
				for _, child := range t.Child {
					err := convertToContentTree(child, m)
					if err != nil {
						return err
					}
				}
				return nil
			}
		default:
			{
				err := m.AppendChild(transformed)
				if err != nil {
					return fmt.Errorf(
						"failed to append transformed child of type <%s> for parent <%s>: %w",
						transformed.GetType(),
						m.GetType(),
						err,
					)
				}
				if transformed.GetChildren() != nil {
					for _, child := range t.Child {
						err := convertToContentTree(child, transformed)
						if err != nil {
							return err
						}
					}
				}
				return nil
			}

		}
	case *etree.CharData:
		data := t.Data
		tx := &contenttree.Text{
			Value: data,
			Type:  contenttree.TextType,
		}
		err := m.AppendChild(tx)
		if err != nil {
			return err
		}
	}
	return nil
}
