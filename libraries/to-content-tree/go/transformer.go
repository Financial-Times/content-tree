package xmltotree

import (
	contenttree "github.com/Financial-Times/content-tree"
	"strings"

	"github.com/beevik/etree"
)

// ---------- Domain-y bits you referenced ----------

const DO_NOT_PROCESS_CHILDREN = "DO_NOT_PROCESS_CHILDREN"

type LayoutWidth string

func toValidLayoutWidth(w string) LayoutWidth {
	switch w {
	case "auto", "in-line", "inset-left", "inset-right",
		"full-bleed", "full-grid", "mid-grid", "full-width":
		return LayoutWidth(w)
	default:
		return LayoutWidth("full-width")
	}
}

var ContentType = struct {
	ImageSet            string
	Video               string
	Content             string
	Article             string
	CustomCodeComponent string
}{
	ImageSet:            "http://www.ft.com/ontology/content/ImageSet",
	Video:               "http://www.ft.com/ontology/content/Video",
	Content:             "http://www.ft.com/ontology/content/Content",
	Article:             "http://www.ft.com/ontology/content/Article",
	CustomCodeComponent: "http://www.ft.com/ontology/content/CustomCodeComponent",
}

// We return plain maps so the JSON shape mirrors the TS output.
// If you have real Go structs (e.g., transit.Heading), swap map[string]any for those types.
type Transformer func(el *etree.Element) map[string]any

type TransformerWithType func(el *etree.Element) contenttree.Node

// ---------- Helpers (etree-based) ----------

func findChild(el *etree.Element, tag string) *etree.Element {
	// Check direct children first
	for _, ch := range el.ChildElements() {
		if ch.Tag == tag {
			return ch
		}
		// Recursively search in this child's subtree
		if found := findChild(ch, tag); found != nil {
			return found
		}
	}
	return nil
}

// textContent returns concatenated text for an element, recursively including children.
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

// flattens one level of nested <div> children (used in layout-slot hack)
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

// small util
func valueOr(v, fallback string) string {
	if v != "" {
		return v
	}
	return fallback
}

func attr(el *etree.Element, name string) string {
	return el.SelectAttrValue(name, "")
}

// ---------- The translation of defaultTransformers (etree version) ----------

var DefaultTransformers = map[string]Transformer{
	"h1": func(h1 *etree.Element) map[string]any {
		return map[string]any{"type": "heading", "level": "chapter"}
	},
	"h2": func(h2 *etree.Element) map[string]any {
		return map[string]any{"type": "heading", "level": "subheading"}
	},
	"h3": func(h3 *etree.Element) map[string]any {
		return map[string]any{"type": "heading", "level": "subheading"}
	},
	"h4": func(h4 *etree.Element) map[string]any {
		return map[string]any{"type": "heading", "level": "label"}
	},
	"p": func(p *etree.Element) map[string]any {
		return map[string]any{"type": "paragraph"}
	},
	"em": func(em *etree.Element) map[string]any {
		return map[string]any{"type": "emphasis"}
	},
	"strong": func(strong *etree.Element) map[string]any {
		return map[string]any{"type": "strong"}
	},
	"s": func(s *etree.Element) map[string]any {
		return map[string]any{"type": "strikethrough"}
	},
	"br": func(br *etree.Element) map[string]any {
		return map[string]any{"type": "break", "children": DO_NOT_PROCESS_CHILDREN}
	},
	"hr": func(hr *etree.Element) map[string]any {
		return map[string]any{"type": "thematic-break", "children": DO_NOT_PROCESS_CHILDREN}
	},
	"a": func(a *etree.Element) map[string]any {
		if attr(a, "data-asset-type") == "video" {
			url := attr(a, "href")
			if strings.Contains(url, "youtube.com") {
				return map[string]any{
					"type":     "youtube-video",
					"url":      url,
					"children": DO_NOT_PROCESS_CHILDREN,
				}
			}
			// NOTE: Vimeo not yet in spec (same as TODO in TS)
		}
		return map[string]any{
			"type":  "link",
			"title": attr(a, "title"),
			"url":   attr(a, "href"),
		}
	},
	"ol": func(ol *etree.Element) map[string]any {
		return map[string]any{"type": "list", "ordered": true}
	},
	"ul": func(ul *etree.Element) map[string]any {
		return map[string]any{"type": "list", "ordered": false}
	},
	"li": func(li *etree.Element) map[string]any {
		return map[string]any{"type": "list-item"}
	},
	"blockquote": func(bq *etree.Element) map[string]any {
		return map[string]any{"type": "blockquote"}
	},
	"pull-quote": func(pq *etree.Element) map[string]any {
		textEl := findChild(pq, "pull-quote-text")
		sourceEl := findChild(pq, "pull-quote-source")
		return map[string]any{
			"type": "pullquote",
			"text": func() string {
				if textEl != nil {
					return textContent(textEl)
				}
				return ""
			}(),
			"source": func() string {
				if sourceEl != nil {
					return textContent(sourceEl)
				}
				return ""
			}(),
			"children": DO_NOT_PROCESS_CHILDREN,
		}
	},
	"big-number": func(bn *etree.Element) map[string]any {
		numEl := findChild(bn, "big-number-headline")
		descEl := findChild(bn, "big-number-intro")
		return map[string]any{
			"type": "big-number",
			"number": func() string {
				if numEl != nil {
					return textContent(numEl)
				}
				return ""
			}(),
			"description": func() string {
				if descEl != nil {
					return textContent(descEl)
				}
				return ""
			}(),
			"children": DO_NOT_PROCESS_CHILDREN,
		}
	},
	"img": func(img *etree.Element) map[string]any {
		return map[string]any{
			"type":     "layout-image",
			"id":       attr(img, "src"),
			"credit":   attr(img, "data-copyright"),
			"alt":      attr(img, "alt"),
			"caption":  attr(img, "longdesc"),
			"children": DO_NOT_PROCESS_CHILDREN,
		}
	},

	// content types keyed by their ontology URIs
	ContentType.ImageSet: func(content *etree.Element) map[string]any {
		return map[string]any{
			"type":     "image-set",
			"id":       attr(content, "url"),
			"children": DO_NOT_PROCESS_CHILDREN,
		}
	},
	ContentType.Video: func(content *etree.Element) map[string]any {
		return map[string]any{
			"type":     "video",
			"id":       attr(content, "url"),
			"children": DO_NOT_PROCESS_CHILDREN,
		}
	},

	// NOTE: The TS had a TODO here: "what is a 'content' or an 'article'?"
	ContentType.Content: func(content *etree.Element) map[string]any {
		id := attr(content, "url")
		parts := strings.Split(id, "/")
		uuid := ""
		if len(parts) > 0 {
			uuid = parts[len(parts)-1]
		}
		if attr(content, "data-asset-type") == "flourish" {
			return map[string]any{
				"type":         "flourish",
				"id":           uuid,
				"flourishType": attr(content, "data-flourish-type"),
				"layoutWidth":  string(toValidLayoutWidth(attr(content, "data-layout-width"))),
				"description":  attr(content, "alt"),
				"timestamp":    attr(content, "data-time-stamp"),
				"children":     DO_NOT_PROCESS_CHILDREN,
			}
		}
		return map[string]any{
			"type":  "link",
			"url":   "https://www.ft.com/content/" + uuid,
			"title": attr(content, "dataTitle"),
		}
	},
	ContentType.Article: func(content *etree.Element) map[string]any {
		id := attr(content, "url")
		parts := strings.Split(id, "/")
		uuid := ""
		if len(parts) > 0 {
			uuid = parts[len(parts)-1]
		}
		return map[string]any{
			"type":  "link",
			"url":   "https://www.ft.com/content/" + uuid,
			"title": attr(content, "dataTitle"),
		}
	},
	ContentType.CustomCodeComponent: func(content *etree.Element) map[string]any {
		id := attr(content, "url")
		parts := strings.Split(id, "/")
		uuid := ""
		if len(parts) > 0 {
			uuid = parts[len(parts)-1]
		}
		return map[string]any{
			"type":        "custom-code-component",
			"id":          uuid,
			"layoutWidth": string(toValidLayoutWidth(attr(content, "data-layout-width"))),
			"children":    DO_NOT_PROCESS_CHILDREN,
		}
	},

	"recommended": func(rl *etree.Element) map[string]any {
		link := findChild(rl, "ft-content")
		heading := findChild(rl, "recommended-title")
		id := ""
		teaser := ""
		if link != nil {
			id = attr(link, "url")
			teaser = textContent(link)
		}
		return map[string]any{
			"type": "recommended",
			"id":   id,
			"heading": func() string {
				if heading != nil {
					return textContent(heading)
				}
				return ""
			}(),
			"teaserTitleOverride": teaser,
			"children":            DO_NOT_PROCESS_CHILDREN,
		}
	},

	"div": func(div *etree.Element) map[string]any {
		switch attr(div, "class") {
		case "n-content-layout":
			return map[string]any{
				"type":        "layout",
				"layoutName":  valueOr(attr(div, "data-layout-name"), "auto"),
				"layoutWidth": string(toValidLayoutWidth(attr(div, "data-layout-width"))),
			}
		case "n-content-layout__container":
			return map[string]any{"type": "__LIFT_CHILDREN__"}
		case "n-content-layout__slot":
			// Flatten doubled-up layout-slot divs as in the TS code.
			div.Child = flattenedChildren(div)
			return map[string]any{"type": "layout-slot"}
		default:
			//node code likes to set the div as it is but we can't do that since we are using etree and node is using xast
			// { type: "__UNKNOWN__", data: div };
			return map[string]any{"type": "__UNKNOWN__", "data": div, "class": attr(div, "class")}
		}
	},
	"experimental": func(_ *etree.Element) map[string]any {
		return map[string]any{"type": "__LIFT_CHILDREN__"}
	},
}
