package tocontenttree

import (
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

func optScrollBlockNoBoxBool(v string) *bool {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "1":
		b := true
		return &b
	default:
		return nil
	}
}

func toScrollyTheme(v string) scrollytheme {
	switch strings.TrimSpace(v) {
	case "1":
		return scrollytheme("sans")
	case "2":
		return scrollytheme("serif")
	default:
		return scrollytheme(v)
	}
}

func toScrollyDisplay(v string) scrollydisplay {
	switch strings.TrimSpace(v) {
	case "1":
		return scrollydisplay("dark-background")
	case "2":
		return scrollydisplay("light-background")
	default:
		return scrollydisplay(v)
	}
}

func toScrollyPosition(v string) scrollyposition {
	switch strings.TrimSpace(v) {
	case "1":
		return scrollyposition("left")
	case "2":
		return scrollyposition("center")
	case "3":
		return scrollyposition("right")
	default:
		return scrollyposition(v)
	}
}

func toScrollyTransition(v string) scrollytransition {
	switch strings.TrimSpace(v) {
	case "1":
		return scrollytransition("delay-before")
	case "2":
		return scrollytransition("delay-after")
	default:
		return scrollytransition(v)
	}
}

func toScrollyHeadingLevel(v string) scrollyheadinglevel {
	switch strings.TrimSpace(v) {
	case "1":
		return scrollyheadinglevel("chapter")
	case "2":
		return scrollyheadinglevel("heading")
	case "3":
		return scrollyheadinglevel("subheading")
	default:
		return scrollyheadinglevel("")
	}
}
