package tocontenttree

import (
	"errors"
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
	transformErr := convertToContentTree(root, m)

	out := &contenttree.Root{
		Type: contenttree.RootType,
		Body: m,
	}

	return out, transformErr
}

func convertToContentTree(elem etree.Token, m contenttree.Node) error {
	switch t := elem.(type) {
	case *etree.Element:
		if t.Tag == "body" {
			var errs []error
			for _, child := range t.Child {
				err := convertToContentTree(child, m)
				if err != nil {
					errs = append(errs, err)
				}
			}
			return errors.Join(errs...)
		}

		tag := t.Tag
		if t.Tag == "content" {
			for _, attr := range t.Attr {
				if attr.Key == "type" {
					tag = attr.Value
					break
				}
			}
		}

		transformer, ok := defaultTransformers[tag]
		if !ok {
			return fmt.Errorf("skipped unsupported element <%s>", tag)
		}
		switch transformed := transformer(t).(type) {
		case *unknownNode:
			{
				return fmt.Errorf("skipped unsupported element <%s>", t.Tag)
			}
		case *liftChildrenNode:
			{
				var errs []error
				for _, child := range t.Child {
					err := convertToContentTree(child, m)
					if err != nil {
						errs = append(errs, err)
					}
				}
				return errors.Join(errs...)
			}
		default:
			{
				err := m.AppendChild(transformed)
				if err != nil {
					return fmt.Errorf(
						"skipped invalid child node %q under %q: %w",
						transformed.GetType(),
						m.GetType(),
						err,
					)
				}
				if transformed.GetChildren() != nil {
					var errs []error
					for _, child := range t.Child {
						err := convertToContentTree(child, transformed)
						if err != nil {
							errs = append(errs, err)
						}
					}
					return errors.Join(errs...)
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
			return fmt.Errorf("skipped invalid text node under %q: %w", m.GetType(), err)
		}
	}
	return nil
}
