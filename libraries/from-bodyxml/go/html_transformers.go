package tocontenttree

import (
	"strings"

	contenttree "github.com/Financial-Times/content-tree"
	"github.com/beevik/etree"
)

type unknownNode struct {
	Type  string         `json:"type"`
	Data  *etree.Element `json:"data,omitempty"`
	Class string         `json:"class,omitempty"`
}

func (n *unknownNode) GetType() string                      { return n.Type }
func (n *unknownNode) GetEmbedded() contenttree.Node        { return nil }
func (n *unknownNode) GetChildren() []contenttree.Node      { return nil }
func (n *unknownNode) AppendChild(_ contenttree.Node) error { return contenttree.ErrCannotHaveChildren }

func newUnknownNode(class string, data *etree.Element) *unknownNode {
	return &unknownNode{
		Type:  "__UNKNOWN__",
		Class: class,
		Data:  data,
	}
}

type liftChildrenNode struct {
	Type  string      `json:"type"`
	Data  interface{} `json:"data,omitempty"`
	Class string      `json:"class,omitempty"`
}

func (n *liftChildrenNode) GetType() string                 { return n.Type }
func (n *liftChildrenNode) GetEmbedded() contenttree.Node   { return nil }
func (n *liftChildrenNode) GetChildren() []contenttree.Node { return nil }
func (n *liftChildrenNode) AppendChild(child contenttree.Node) error {
	return contenttree.ErrCannotHaveChildren
}

func newLiftChildrenNode() *liftChildrenNode {
	return &liftChildrenNode{
		Type: "__LIFT_CHILDREN__",
	}
}

var contentType = struct {
	ImageSet            string
	Video               string
	Content             string
	Article             string
	CustomCodeComponent string
	ClipSet             string
}{
	ImageSet:            "http://www.ft.com/ontology/content/ImageSet",
	Video:               "http://www.ft.com/ontology/content/Video",
	Content:             "http://www.ft.com/ontology/content/Content",
	Article:             "http://www.ft.com/ontology/content/Article",
	CustomCodeComponent: "http://www.ft.com/ontology/content/CustomCodeComponent",
	ClipSet:             "http://www.ft.com/ontology/content/ClipSet",
}

type transformer func(el *etree.Element) contenttree.Node

var defaultTransformers = map[string]transformer{
	"h1": func(h1 *etree.Element) contenttree.Node {
		dfrgId := attr(h1, "data-fragment-identifier")
		heading := &contenttree.Heading{
			Type:               contenttree.HeadingType,
			Level:              "chapter",
			Children:           []*contenttree.Text{},
			FragmentIdentifier: dfrgId,
		}
		return heading
	},
	"h2": func(h2 *etree.Element) contenttree.Node {
		dfrgId := attr(h2, "data-fragment-identifier")
		return &contenttree.Heading{
			Type:               contenttree.HeadingType,
			Level:              "subheading",
			Children:           []*contenttree.Text{},
			FragmentIdentifier: dfrgId,
		}
	},
	"h3": func(h3 *etree.Element) contenttree.Node {
		dfrgId := attr(h3, "data-fragment-identifier")
		return &contenttree.Heading{
			Type:               contenttree.HeadingType,
			Level:              "subheading",
			Children:           []*contenttree.Text{},
			FragmentIdentifier: dfrgId,
		}
	},
	"h4": func(h4 *etree.Element) contenttree.Node {
		dfrgId := attr(h4, "data-fragment-identifier")
		return &contenttree.Heading{
			Type:               contenttree.HeadingType,
			Level:              "label",
			Children:           []*contenttree.Text{},
			FragmentIdentifier: dfrgId,
		}
	},
	"p": func(p *etree.Element) contenttree.Node {
		return &contenttree.Paragraph{
			Type:     contenttree.ParagraphType,
			Children: []*contenttree.Phrasing{},
		}
	},
	"em": func(em *etree.Element) contenttree.Node {
		return &contenttree.Emphasis{
			Type:     contenttree.EmphasisType,
			Children: []*contenttree.Phrasing{},
		}
	},
	"strong": func(strong *etree.Element) contenttree.Node {
		return &contenttree.Strong{
			Type:     contenttree.StrongType,
			Children: []*contenttree.Phrasing{},
		}
	},
	"s": func(s *etree.Element) contenttree.Node {
		return &contenttree.Strikethrough{
			Type:     contenttree.StrikethroughType,
			Children: []*contenttree.Phrasing{},
		}
	},
	"br": func(br *etree.Element) contenttree.Node {
		return &contenttree.Break{
			Type: contenttree.BreakType,
		}
	},
	"hr": func(hr *etree.Element) contenttree.Node {
		return &contenttree.ThematicBreak{
			Type: contenttree.ThematicBreakType,
		}
	},
	"a": func(a *etree.Element) contenttree.Node {
		if attr(a, "data-asset-type") == "video" {
			url := attr(a, "href")
			if strings.Contains(url, "youtube.com") {
				return &contenttree.YoutubeVideo{
					Type: contenttree.YoutubeVideoType,
					URL:  url,
				}
			}
			// NOTE: Vimeo not yet in spec
		} else if attr(a, "data-asset-type") == "tweet" {
			url := attr(a, "href")
			return &contenttree.Tweet{
				Type: contenttree.TweetType,
				ID:   url,
			}
		}
		return &contenttree.Link{
			Type:     contenttree.LinkType,
			Title:    attr(a, "title"),
			URL:      attr(a, "href"),
			Children: []*contenttree.Phrasing{},
		}
	},
	"ol": func(ol *etree.Element) contenttree.Node {
		dataType := attr(ol, "data-type")
		if dataType == "timeline_events" {
			return newLiftChildrenNode()
		}
		return &contenttree.List{
			Type:     contenttree.ListType,
			Ordered:  true,
			Children: []*contenttree.ListItem{},
		}
	},
	"ul": func(ul *etree.Element) contenttree.Node {
		return &contenttree.List{
			Type:     contenttree.ListType,
			Ordered:  false,
			Children: []*contenttree.ListItem{},
		}
	},
	"li": func(li *etree.Element) contenttree.Node {
		dataType := attr(li, "data-type")
		if dataType == "timeline_event" {
			timelineEventTitle := ""
			if h4Element := findChild(li, "h4"); h4Element != nil {
				timelineEventTitle = textContent(h4Element)
				//extract title but don't treat like a child element
				li.RemoveChild(h4Element)
			}
			return &contenttree.TimelineEvent{
				Type:     contenttree.TimelineEventType,
				Title:    timelineEventTitle,
				Children: []*contenttree.TimelineEventChild{},
			}
		}
		return &contenttree.ListItem{
			Type:     contenttree.ListItemType,
			Children: []*contenttree.ListItemChild{},
		}
	},
	"blockquote": func(bq *etree.Element) contenttree.Node {
		return &contenttree.Blockquote{
			Type:     contenttree.BlockquoteType,
			Children: []*contenttree.BlockquoteChild{},
		}
	},
	"pull-quote": func(pq *etree.Element) contenttree.Node {
		textEl := findChild(pq, "pull-quote-text")
		sourceEl := findChild(pq, "pull-quote-source")
		return &contenttree.Pullquote{
			Type: contenttree.PullquoteType,
			Text: func() string {
				if textEl != nil {
					return textContent(textEl)
				}
				return ""
			}(),
			Source: func() string {
				if sourceEl != nil {
					return textContent(sourceEl)
				}
				return ""
			}(),
		}
	},
	"big-number": func(bn *etree.Element) contenttree.Node {
		numEl := findChild(bn, "big-number-headline")
		descEl := findChild(bn, "big-number-intro")
		return &contenttree.BigNumber{
			Type: contenttree.BigNumberType,
			Number: func() string {
				if numEl != nil {
					return textContent(numEl)
				}
				return ""
			}(),
			Description: func() string {
				if descEl != nil {
					return textContent(descEl)
				}
				return ""
			}(),
		}
	},
	"img": func(img *etree.Element) contenttree.Node {
		return &contenttree.LayoutImage{
			Type:    contenttree.LayoutImageType,
			ID:      attr(img, "src"),
			Credit:  attr(img, "data-copyright"),
			Alt:     attr(img, "alt"),
			Caption: attr(img, "longdesc"),
		}
	},
	contentType.ImageSet: func(content *etree.Element) contenttree.Node {
		dfrgId := attr(content, "data-fragment-identifier")
		return &contenttree.ImageSet{
			Type:               contenttree.ImageSetType,
			ID:                 attr(content, "id"),
			FragmentIdentifier: dfrgId,
		}
	},
	contentType.Video: func(content *etree.Element) contenttree.Node {
		return &contenttree.Video{
			Type: contenttree.VideoType,
			ID:   attr(content, "id"),
		}
	},
	contentType.Content: func(content *etree.Element) contenttree.Node {
		id := attr(content, "id")
		if attr(content, "data-asset-type") == "flourish" {
			dfrgId := valueOr(attr(content, "data-fragment-identifier"), id)
			return &contenttree.Flourish{
				Type:               contenttree.FlourishType,
				Id:                 id,
				FlourishType:       attr(content, "data-flourish-type"),
				LayoutWidth:        string(toValidLayoutWidth(attr(content, "data-layout-width"))),
				Description:        attr(content, "alt"),
				Timestamp:          attr(content, "data-time-stamp"),
				FragmentIdentifier: dfrgId,
			}
		}
		return &contenttree.Link{
			Type:     contenttree.LinkType,
			URL:      "https://www.ft.com/content/" + id,
			Title:    attr(content, "dataTitle"),
			Children: []*contenttree.Phrasing{},
		}
	},
	contentType.Article: func(content *etree.Element) contenttree.Node {
		return &contenttree.Link{
			Type:     contenttree.LinkType,
			URL:      "https://www.ft.com/content/" + attr(content, "id"),
			Title:    attr(content, "dataTitle"),
			Children: []*contenttree.Phrasing{},
		}
	},
	contentType.CustomCodeComponent: func(content *etree.Element) contenttree.Node {
		return &contenttree.CustomCodeComponent{
			Type:        contenttree.CustomCodeComponentType,
			ID:          attr(content, "id"),
			LayoutWidth: string(toValidLayoutWidth(attr(content, "data-layout-width"))),
		}
	},
	contentType.ClipSet: func(content *etree.Element) contenttree.Node {
		return &contenttree.ClipSet{
			Type:        contenttree.ClipSetType,
			ID:          attr(content, "id"),
			LayoutWidth: string(toValidClipLayoutWidth(attr(content, "data-layout-width"))),
			Autoplay:    attr(content, "autoplay") == "true",
			Loop:        attr(content, "loop") == "true",
			Muted:       attr(content, "muted") == "true",
		}
	},
	"recommended": func(rl *etree.Element) contenttree.Node {
		id := ""
		teaser := ""
		if link := findChild(rl, "content"); link != nil {
			id = attr(link, "id")
			teaser = textContent(link)
		}
		heading := findChild(rl, "recommended-title")
		return &contenttree.Recommended{
			Type: contenttree.RecommendedType,
			ID:   id,
			Heading: func() string {
				if heading != nil {
					return textContent(heading)
				}
				return ""
			}(),
			TeaserTitleOverride: teaser,
		}
	},
	"div": func(div *etree.Element) contenttree.Node {
		switch attr(div, "class") {
		case "n-content-layout":
			return &contenttree.Layout{
				Type:        contenttree.LayoutType,
				LayoutName:  valueOr(attr(div, "data-layout-name"), "auto"),
				LayoutWidth: string(toValidLayoutWidth(attr(div, "data-layout-width"))),
				Children:    []*contenttree.LayoutChild{},
			}
		case "n-content-layout__container":
			return newLiftChildrenNode()
		case "n-content-layout__slot":
			div.Child = flattenedChildren(div)
			return &contenttree.LayoutSlot{
				Type:     contenttree.LayoutSlotType,
				Children: []*contenttree.LayoutSlotChild{},
			}
		default:
			return newUnknownNode(attr(div, "class"), div)
		}
	},
	"section": func(section *etree.Element) contenttree.Node {
		switch attr(section, "data-type") {
		case "timeline":
			{
				timelineTitle := ""
				if h3Element := findChild(section, "h3"); h3Element != nil {
					timelineTitle = textContent(h3Element)
					//extract title but don't treat like a child element
					section.RemoveChild(h3Element)
				}
				return &contenttree.Timeline{
					Type:        contenttree.TimelineType,
					Title:       timelineTitle,
					LayoutWidth: attr(section, "data-layout-width"),
					Children:    []*contenttree.TimelineEvent{},
				}
			}
		}
		return newUnknownNode("", section)
	},
	"experimental": func(_ *etree.Element) contenttree.Node {
		return newLiftChildrenNode()
	},
}
