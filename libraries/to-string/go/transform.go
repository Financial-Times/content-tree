package tostring

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"
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

var toSeparate = []string{contenttree.HeadingType, contenttree.ParagraphType}

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

	switch n.GetType() {
	case contenttree.BodyBlockType:
		return transformNode(n.GetEmbedded())
	case contenttree.BlockquoteChildType:
		return transformNode(n.GetEmbedded())
	case contenttree.LayoutChildType:
		return transformNode(n.GetEmbedded())
	case contenttree.LayoutSlotChildType:
		return transformNode(n.GetEmbedded())
	case contenttree.ListItemChildType:
		return transformNode(n.GetEmbedded())
	case contenttree.PhrasingType:
		return transformNode(n.GetEmbedded())
	case contenttree.ScrollyCopyChildType:
		return transformNode(n.GetEmbedded())
	case contenttree.ScrollySectionChildType:
		return transformNode(n.GetEmbedded())
	case contenttree.TableChildType:
		return transformNode(n.GetEmbedded())
	case contenttree.RecommendedType:
		return "", nil
	case contenttree.RecommendedListType:
		return "", nil
	case contenttree.TextType:
		text, ok := n.(*contenttree.Text)
		if !ok {
			return "", errors.New("failed to parse node to text")
		}

		return text.Value, nil
	case contenttree.RootType:
		root, ok := n.(*contenttree.Root)
		if !ok {
			return "", errors.New("failed to parse node to root")
		}

		return transformNode(root.Body)
	case contenttree.BigNumberType:
		bigNumber, ok := n.(*contenttree.BigNumber)
		if !ok {
			return "", errors.New("failed to parse node to bigNumber")
		}

		return fmt.Sprintf("%s %s", bigNumber.Number, bigNumber.Description), nil
	case contenttree.TimelineEventType:
		te, ok := n.(*contenttree.TimelineEvent)
		if !ok {
			return "", errors.New("failed to parse node to text")
		}

		resultChildred, err := transformChildren(te.GetChildren())
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("%s %s", te.Title, resultChildred), nil
	default:

		result, err := transformChildren(n.GetChildren())
		if err != nil {
			return "", err
		}
		if slices.Contains(toSeparate, n.GetType()) && len(result) > 0 {
			result += " "
		}

		return result, nil
	}
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

	return strings.Join(childrenStrs, ""), nil
}
