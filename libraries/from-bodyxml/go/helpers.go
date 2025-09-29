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

var contentTypeTemplates = map[string]string{
	"http://www.ft.com/ontology/content/Article":             "/content/{{id}}",
	"http://www.ft.com/ontology/content/ImageSet":            "/content/{{id}}",
	"http://www.ft.com/ontology/content/ClipSet":             "/content/{{id}}",
	"http://www.ft.com/ontology/content/CustomCodeComponent": "/content/{{id}}",
	"http://www.ft.com/ontology/content/MediaResource":       "/content/{{id}}",
	"http://www.ft.com/ontology/content/Video":               "/content/{{id}}",
	"http://www.ft.com/ontology/company/PublicCompany":       "/organisations/{{id}}",
	"http://www.ft.com/ontology/content/ContentPackage":      "/content/{{id}}",
	"http://www.ft.com/ontology/content/Content":             "/content/{{id}}",
	"http://www.ft.com/ontology/content/Image":               "/content/{{id}}",
	"http://www.ft.com/ontology/content/DynamicContent":      "/content/{{id}}",
	"http://www.ft.com/ontology/content/Graphic":             "/content/{{id}}",
	"http://www.ft.com/ontology/content/Audio":               "/content/{{id}}",
	"http://www.ft.com/ontology/company/Organisation":        "/organisations/{{id}}",
}

func generateUrl(t, id string) string {
	const host = "http://api.ft.com"
	template, ok := contentTypeTemplates[t]
	if !ok {
		return ""
	}
	path := strings.Replace(template, "{{id}}", id, 1)
	return host + path
}
