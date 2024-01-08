package stringifier

import (
	"encoding/json"
	"fmt"
	"strings"

	"golang.org/x/exp/slices"
)

var separate = []string{"heading", "paragraph"}

// node is overly simplified and generic representation of content-tree node. This representation is applicable for the
// root node as well. node has only fields related to the task of transforming content-tree to text.
type node struct {
	Type     string      `json:"type"`
	Children []*node     `json:"children"`
	Value    interface{} `json:"value"`
	Body     *node       `json:"body"`
}

func (n *node) isText() bool {
	return n.Type == "text"
}

func (n *node) isParent() bool {
	return len(n.Children) > 0
}

// Stringify receives a content-tree as a unmarshalled JSON([]byte) and returns the plain text from the content-tree.
func Stringify(body json.RawMessage) (string, error) {
	ct := node{}
	err := json.Unmarshal(body, &ct)
	if err != nil {
		return "", fmt.Errorf("failed to parse string to content tree: %w", err)
	}

	s, err := stringifyNode(&ct)
	return strings.TrimSpace(dedupWhitespace(s)), err
}

func stringifyNode(n *node) (string, error) {
	if n.isText() {
		valStr, ok := n.Value.(string)
		if !ok {
			return "", fmt.Errorf("failed to parse content-tree's text node value to string")
		}

		return valStr, nil
	}

	if n.Body != nil {
		return stringifyNode(n.Body)
	}

	if n.isParent() {
		children := []string{}
		for _, child := range n.Children {
			s, err := stringifyNode(child)
			if err != nil {
				return "", fmt.Errorf("failed to stringify content-tree child node: %w", err)
			}

			children = append(children, s)
		}

		childrenStr := strings.Join(children, "")

		if slices.Contains(separate, n.Type) && len(childrenStr) > 0 {
			childrenStr += " "
		}

		return childrenStr, nil
	}

	return "", nil
}

func dedupWhitespace(s string) string {
	return strings.Join(strings.Fields(s), " ")
}
