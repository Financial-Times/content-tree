package tostring

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	contenttree "github.com/Financial-Times/content-tree"
)

type Schema interface {
	fmt.Stringer
}

type schema string

func (s schema) String() string { return string(s) }

var (
	TransitTree Schema = schema("transit-tree")
	BodyTree    Schema = schema("body-tree")
)

var ErrUnknownKind = errors.New("unknown tree kind")

// Transform extracts and returns plain text from a content tree represented as unmarshalled JSON(json.RawMessage).
func Transform(tree json.RawMessage, s Schema) (string, error) {
	switch s {
	case TransitTree:
		n := contenttree.Root{}
		return unmarshalAndTransform(tree, &n)
	case BodyTree:
		n := contenttree.Body{}
		return unmarshalAndTransform(tree, &n)
	default:
		return "", fmt.Errorf("%w: %q (expected %q or %q)", ErrUnknownKind, s, TransitTree, BodyTree)
	}
}

func unmarshalAndTransform(tree json.RawMessage, n contenttree.Node) (string, error) {
	if err := json.Unmarshal(tree, n); err != nil {
		return "", fmt.Errorf("failed to instantiate content tree: %w", err)
	}
	text, err := transformNode(n)
	if err != nil {
		return "", fmt.Errorf("failed to transform tree to string: %w", err)
	}

	// Deduplicate whitespace.
	text = strings.Join(strings.Fields(text), " ")

	// Trim leading and trailing whitespace.
	text = strings.TrimSpace(text)

	return text, nil
}

func transformNode(n contenttree.Node) (string, error) {
	if n == nil {
		return "", errors.New("nil node")
	}

	result := n.GetString()
	if currentElementChildren := n.GetChildren(); currentElementChildren != nil && len(currentElementChildren) > 0 {
		childrenString, err := transformChildren(currentElementChildren)
		if err != nil {
			return "", err
		}
		result = fmt.Sprintf("%s %s ", result, childrenString)
	}
	return result, nil
}

func transformChildren(childrenNodes []contenttree.Node) (string, error) {
	if len(childrenNodes) == 0 {
		return "", nil
	}

	childrenStrs := []string{}
	for _, child := range childrenNodes {
		s, err := transformNode(child)
		if err != nil {
			return "", fmt.Errorf("failed to transform child node to string: %w", err)
		}

		childrenStrs = append(childrenStrs, s)
	}

	return strings.Join(childrenStrs, " "), nil
}
