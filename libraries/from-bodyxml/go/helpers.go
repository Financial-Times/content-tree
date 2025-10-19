package tocontenttree

import (
	"strings"

	"github.com/beevik/etree"
)

type layoutwidth string

func toValidLayoutWidth(w string) layoutwidth {
	switch w {
	case "auto", "in-line", "inset-left", "inset-right",
		"full-bleed", "full-grid", "mid-grid", "full-width":
		return layoutwidth(w)
	default:
		return "full-width"
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
