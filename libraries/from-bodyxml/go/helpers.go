package tocontenttree

import (
	"fmt"
	"strings"

	"github.com/beevik/etree"
)

type layoutwidth string
type scrollytheme string
type scrollydisplay string
type scrollyposition string
type scrollytransition string
type scrollyheadinglevel string

func toValidLayoutWidth(w string) layoutwidth {
	switch w {
	case "auto", "in-line", "inset-left", "inset-right",
		"full-bleed", "full-grid", "mid-grid", "full-width":
		return layoutwidth(w)
	default:
		return "full-width"
	}
}

func toValidFlourishLayoutWidth(w string) layoutwidth {
	switch w {
	case "auto", "in-line", "inset-left", "inset-right",
		"full-bleed", "full-grid", "mid-grid", "full-width":
		return layoutwidth(w)
	default:
		return "in-line"
	}
}

func toValidClipLayoutWidth(w string) layoutwidth {
	switch w {
	case "in-line", "full-grid", "mid-grid":
		return layoutwidth(w)
	default:
		return "in-line"
	}

}

func findChild(el *etree.Element, tag string) *etree.Element {
	for _, ch := range el.ChildElements() {
		if ch.Tag == tag {
			return ch
		}
		if found := findChild(ch, tag); found != nil {
			return found
		}
	}
	return nil
}

func textContent(el *etree.Element) string {
	var b strings.Builder
	for _, tok := range el.Child {
		switch t := tok.(type) {
		case *etree.CharData:
			b.WriteString(t.Data)
		case *etree.Element:
			b.WriteString(textContent(t))
		}
	}
	return b.String()
}

func flattenedChildren(el *etree.Element) []etree.Token {
	out := make([]etree.Token, 0, len(el.Child))
	for _, tok := range el.Child {
		if d, ok := tok.(*etree.Element); ok && d.Tag == "div" {
			out = append(out, d.Child...)
		} else {
			out = append(out, tok)
		}
	}
	return out
}

func valueOr(v, fallback string) string {
	if v != "" {
		return v
	}
	return fallback
}

func attr(el *etree.Element, name string) string {
	return el.SelectAttrValue(name, "")
}

func hasParentTag(el *etree.Element, tag string) bool {
	parent := el.Parent()
	if parent == nil {
		return false
	}
	return parent.Tag == tag
}

func optScrollBlockNoBoxBool(v string) (*bool, error) {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "0", "":
		return nil, nil
	case "1":
		b := true
		return &b, nil
	default:
		return nil, fmt.Errorf("unsupported scrolly no-box value %q", v)
	}
}

func toScrollyTheme(v string) (scrollytheme, error) {
	switch strings.TrimSpace(v) {
	case "1":
		return scrollytheme("sans"), nil
	case "2":
		return scrollytheme("serif"), nil
	default:
		return "", fmt.Errorf("unsupported scrolly theme value %q", v)
	}
}

func toScrollyDisplay(v string) (scrollydisplay, error) {
	switch strings.TrimSpace(v) {
	case "1":
		return scrollydisplay("dark-background"), nil
	case "2":
		return scrollydisplay("light-background"), nil
	default:
		return "", fmt.Errorf("unsupported scrolly display value %q", v)
	}
}

func toScrollyPosition(v string) (scrollyposition, error) {
	switch strings.TrimSpace(v) {
	case "1":
		return scrollyposition("left"), nil
	case "2":
		return scrollyposition("center"), nil
	case "3":
		return scrollyposition("right"), nil
	default:
		return "", fmt.Errorf("unsupported scrolly position value %q", v)
	}
}

func toScrollyTransition(v string) (scrollytransition, error) {
	switch strings.TrimSpace(v) {
	case "":
		return "", nil
	case "1":
		return scrollytransition("delay-before"), nil
	case "2":
		return scrollytransition("delay-after"), nil
	default:
		return "", fmt.Errorf("unsupported scrolly transition value %q", v)
	}
}

func toScrollyHeadingLevel(v string) (scrollyheadinglevel, error) {
	switch strings.TrimSpace(v) {
	case "1":
		return scrollyheadinglevel("chapter"), nil
	case "2":
		return scrollyheadinglevel("heading"), nil
	case "3":
		return scrollyheadinglevel("subheading"), nil
	default:
		return "", fmt.Errorf("unsupported scrolly heading level value %q", v)
	}
}
