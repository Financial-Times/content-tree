package toexternalbodyxml

import (
	"encoding/json"
	"errors"
	"fmt"
	"html"
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

var (
	ErrUnknownKind = errors.New("unknown tree kind")
)

// Transform converts content from a content tree representation into an external XHTML-formatted version.
//
// The tree is provided as unmarshalled JSON (json.RawMessage) and must conform to one of the
// supported Schema kinds: TransitTree or BodyTree.
//
// The Schema interface is used to distinguish which type of content tree should be unmarshalled
// and transformed. Implementations of Schema (TransitTree, BodyTree) serve as markers to select
// the appropriate unmarshal/transform logic.
//
// The XHTML output is intended for distribution to consumers that only support widely recognized formats like HTML
// or those that should not receive internal-specific details contained in the content tree format.
// Such consumers may be external (non-FT) users, automated systems processing HTML-based content,
// republishing platforms, and more.
func Transform(tree json.RawMessage, s Schema) (string, error) {
	switch s {
	case TransitTree:
		{
			n := contenttree.Root{}
			return unmarshalAndTransform(tree, &n)
		}
	case BodyTree:
		{
			n := contenttree.Body{}
			return unmarshalAndTransform(tree, &n)
		}
	default:
		return "", fmt.Errorf("%w: %q (expected %q or %q)", ErrUnknownKind, s, TransitTree, BodyTree)
	}
}

func unmarshalAndTransform(tree json.RawMessage, n contenttree.Node) (string, error) {
	if err := json.Unmarshal(tree, n); err != nil {
		return "", fmt.Errorf("failed to instantiate content tree: %w", err)
	}
	return transformNode(n)
}

func transformNode(n contenttree.Node) (string, error) {
	if n == nil {
		return "", errors.New("nil node")
	}

	if n.GetType() == contenttree.RootType {
		root, ok := n.(*contenttree.Root)
		if !ok {
			return "", errors.New("failed to parse node to root")
		}

		return transformNode(root.Body)
	}

	innerXML := ""

	childrenNodes := n.GetChildren()
	if childrenNodes != nil {
		childrenStr := make([]string, 0, len(childrenNodes))
		for _, child := range childrenNodes {
			s, err := transformNode(child)
			if err != nil {
				return "", fmt.Errorf("failed to transform child node to external XML: %w", err)
			}

			childrenStr = append(childrenStr, s)
		}
		innerXML = strings.Join(childrenStr, "")
	}

	switch node := n.(type) {
	case *contenttree.Body:
		return fmt.Sprintf("<body>%s</body>", innerXML), nil

	case *contenttree.Text:
		return html.EscapeString(node.Value), nil

	case *contenttree.Break:
		return "<br/>", nil

	case *contenttree.ThematicBreak:
		return "<hr/>", nil

	case *contenttree.Paragraph:
		return fmt.Sprintf("<p>%s</p>", innerXML), nil

	case *contenttree.Heading:
		tag := ""
		if node.Level == "chapter" {
			tag = "h1"
		}
		if node.Level == "subheading" {
			tag = "h2"
		}
		if node.Level == "label" {
			tag = "h4"
		}

		if tag == "" {
			return "", fmt.Errorf("failed to transform heading with level %s", node.Level)
		}

		if node.FragmentIdentifier != "" {
			return fmt.Sprintf("<%[1]s data-fragment-identifier=\"%[2]s\">%[3]s</%[1]s>", tag, html.EscapeString(node.FragmentIdentifier), innerXML), nil
		}
		return fmt.Sprintf("<%[1]s>%s</%[1]s>", tag, innerXML), nil

	case *contenttree.Strong:
		return fmt.Sprintf("<strong>%s</strong>", innerXML), nil

	case *contenttree.Emphasis:
		return fmt.Sprintf("<em>%s</em>", innerXML), nil

	case *contenttree.Strikethrough:
		return fmt.Sprintf("<s>%s</s>", innerXML), nil

	case *contenttree.Subscript:
		return fmt.Sprintf("<sub>%s</sub>", innerXML), nil

	case *contenttree.Superscript:
		return fmt.Sprintf("<sup>%s</sup>", innerXML), nil

	case *contenttree.Link:
		href := html.EscapeString(node.URL)
		if node.Title != "" {
			return fmt.Sprintf("<a href=\"%s\" title=\"%s\">%s</a>", href, html.EscapeString(node.Title), innerXML), nil
		}

		return fmt.Sprintf("<a href=\"%s\">%s</a>", href, innerXML), nil

	case *contenttree.List:
		tag := "ul"
		if node.Ordered {
			tag = "ol"
		}

		return fmt.Sprintf("<%[1]s>%s</%[1]s>", tag, innerXML), nil

	case *contenttree.ListItem:
		return fmt.Sprintf("<li>%s</li>", innerXML), nil

	case *contenttree.Blockquote:
		return fmt.Sprintf("<blockquote>%s</blockquote>", innerXML), nil

	case *contenttree.Pullquote:
		text := html.EscapeString(node.Text)
		if node.Source != "" {
			return fmt.Sprintf("<pull-quote><pull-quote-text><p>%s</p></pull-quote-text><pull-quote-source>%s</pull-quote-source></pull-quote>", text, html.EscapeString(node.Source)), nil
		}
		return fmt.Sprintf("<pull-quote><pull-quote-text><p>%s</p></pull-quote-text></pull-quote>", text), nil

	case *contenttree.ImageSet:
		{
			if node.FragmentIdentifier != "" {
				return fmt.Sprintf("<ft-content type=\"http://www.ft.com/ontology/content/ImageSet\" url=\"http://api.ft.com/content/%s\" data-embedded=\"true\" data-fragment-identifier=\"%s\"></ft-content>", node.ID, node.FragmentIdentifier), nil
			}
			return fmt.Sprintf("<ft-content type=\"http://www.ft.com/ontology/content/ImageSet\" url=\"http://api.ft.com/content/%s\" data-embedded=\"true\"></ft-content>", node.ID), nil
		}
	case *contenttree.ClipSet:
		attrs := []string{
			"type=\"http://www.ft.com/ontology/content/ClipSet\"",
			fmt.Sprintf("url=\"http://api.ft.com/content/%s\"", node.ID),
			"data-embedded=\"true\"",
		}
		if node.LayoutWidth != "" {
			attrs = append(attrs, fmt.Sprintf("data-layout=\"%s\"", node.LayoutWidth))
		}
		if node.FragmentIdentifier != "" {
			attrs = append(attrs, fmt.Sprintf("data-fragment-identifier=\"%s\"", node.FragmentIdentifier))
		}
		if attr := optionalBoolAttrXML("autoplay", node.Autoplay); attr != "" {
			attrs = append(attrs, attr)
		}
		if attr := optionalBoolAttrXML("loop", node.Loop); attr != "" {
			attrs = append(attrs, attr)
		}
		if attr := optionalBoolAttrXML("muted", node.Muted); attr != "" {
			attrs = append(attrs, attr)
		}
		return fmt.Sprintf("<ft-content %s></ft-content>", strings.Join(attrs, " ")), nil

	case *contenttree.Flourish:
		attrs := []string{
			"type=\"http://www.ft.com/ontology/content/Content\"",
			fmt.Sprintf("url=\"http://api.ft.com/content/%s\"", node.Id),
			"data-asset-type=\"flourish\"",
			"data-embedded=\"true\"",
			fmt.Sprintf("id=\"%s\"", node.Id),
		}
		if node.Description != "" {
			attrs = append(attrs, fmt.Sprintf("alt=\"%s\"", html.EscapeString(node.Description)))
		}
		if node.FlourishType != "" {
			attrs = append(attrs, fmt.Sprintf("data-flourish-type=\"%s\"", node.FlourishType))
		}
		if node.LayoutWidth != "" {
			attrs = append(attrs, fmt.Sprintf("data-layout-width=\"%s\"", node.LayoutWidth))
		}
		if node.Timestamp != "" {
			attrs = append(attrs, fmt.Sprintf("data-time-stamp=\"%s\"", node.Timestamp))
		}
		if node.FragmentIdentifier != "" {
			attrs = append(attrs, fmt.Sprintf("data-fragment-identifier=\"%s\"", node.FragmentIdentifier))
		}
		return fmt.Sprintf("<ft-content %s></ft-content>", strings.Join(attrs, " ")), nil
	case *contenttree.TableCaption:
		return fmt.Sprintf("<caption>%s</caption>", innerXML), nil

	case *contenttree.TableCell:
		tag := "td"
		if node.Heading {
			tag = "th"
		}
		var attrs []string
		if node.ColumnSpan != nil {
			attrs = append(attrs, fmt.Sprintf("colspan=\"%d\"", *node.ColumnSpan))
		}
		if node.RowSpan != nil {
			attrs = append(attrs, fmt.Sprintf("rowspan=\"%d\"", *node.RowSpan))
		}
		if len(attrs) > 0 {
			return fmt.Sprintf("<%s %s>%s</%s>", tag, strings.Join(attrs, " "), innerXML, tag), nil
		}
		return fmt.Sprintf("<%s>%s</%s>", tag, innerXML, tag), nil

	case *contenttree.TableRow:
		return fmt.Sprintf("<tr>%s</tr>", innerXML), nil

	case *contenttree.TableBody:
		return fmt.Sprintf("<tbody>%s</tbody>", innerXML), nil

	case *contenttree.TableFooter:
		return fmt.Sprintf("<tfoot><tr><td>%s</td></tr></tfoot>", innerXML), nil

	case *contenttree.Table:
		var attrs []string
		if theme := tableThemeToExternal(node.Compact, node.Stripes); theme != "" {
			attrs = append(attrs, fmt.Sprintf("data-table-theme=\"%s\"", theme))
		}
		if responseStyle := tableResponsiveStyleToExternal(node.ResponsiveStyle); responseStyle != "" {
			attrs = append(attrs, fmt.Sprintf("data-table-layout-smallscreen=\"%s\"", responseStyle))
		}
		if node.LayoutWidth != "" {
			attrs = append(attrs, fmt.Sprintf("data-table-layout-largescreen=\"%s\"", node.LayoutWidth))
		}
		if node.CollapseAfterHowManyRows != nil {
			attrs = append(attrs, fmt.Sprintf("data-table-collapse-rownum=\"%d\"", *node.CollapseAfterHowManyRows))
		}
		childrenXML := make([]string, 0, len(node.Children))
		for _, child := range node.Children {
			//process <thead> here instead of in a separate case block
			if child.TableHeader != nil {
				rowsXML := make([]string, 0, len(child.TableHeader.Children))
				columnIndex := 0
				for _, row := range child.TableHeader.Children {
					cellsXML := make([]string, 0, len(row.Children))
					for _, cell := range row.Children {
						thXML, err := buildTH(cell, columnIndex, node.ColumnSettings)
						if err != nil {
							return "", fmt.Errorf("failed to transform table header cell to external XML: %w", err)
						}
						cellsXML = append(cellsXML, thXML)
						columnIndex++
					}
					rowsXML = append(rowsXML, fmt.Sprintf("<tr>%s</tr>", strings.Join(cellsXML, "")))
				}
				childrenXML = append(childrenXML, fmt.Sprintf("<thead>%s</thead>", strings.Join(rowsXML, "")))
				continue
			}

			s, err := transformNode(child)
			if err != nil {
				return "", fmt.Errorf("failed to transform table child node to external XML: %w", err)
			}
			childrenXML = append(childrenXML, s)
		}
		if len(attrs) > 0 {
			return fmt.Sprintf("<table %s>%s</table>", strings.Join(attrs, " "), strings.Join(childrenXML, "")), nil
		}
		return fmt.Sprintf("<table>%s</table>", strings.Join(childrenXML, "")), nil

	case *contenttree.Video:
		return fmt.Sprintf("<ft-content type=\"http://www.ft.com/ontology/content/Video\" url=\"http://api.ft.com/content/%s\" data-embedded=\"true\"></ft-content>", node.ID), nil

	case *contenttree.YoutubeVideo:
		return fmt.Sprintf("<a data-asset-type=\"video\" data-embedded=\"true\" href=\"%s\"></a>", html.EscapeString(node.URL)), nil

	case *contenttree.VimeoVideo:
		return fmt.Sprintf("<a data-asset-type=\"video\" data-embedded=\"true\" href=\"%s\"></a>", html.EscapeString(node.URL)), nil

	case *contenttree.AcastPodcast:
		return fmt.Sprintf("<a data-asset-type=\"podcast\" data-embedded=\"true\" href=\"%s\"></a>", html.EscapeString(node.URL)), nil

	case *contenttree.Tweet:
		return fmt.Sprintf("<a data-asset-type=\"tweet\" data-embedded=\"true\" href=\"%[1]s\">%[1]s</a>", html.EscapeString(node.ID)), nil

	// Example from the Native Store to keep the translucent namespace (https://www.ft.com/content/9675cf79-f16d-4132-ab73-8bafa22ee4fc):
	// <tr:scrollable-block theme=\"1\">
	//     <tr:scrollable-section theme-position=\"1\" theme-display=\"1\">
	//         <content type=\"http://www.ft.com/ontology/content/ImageSet\" data-embedded=\"true\" id=\"0184bb0b-1dc8-4501-ade4-d0d49f7dd2e1\"></content>
	//         <tr:scrollable-text>
	//             <p><strong>Founded:</strong> 1946 </p><p><strong>Business:</strong> Electronics </p><p><strong>Headquarters:</strong> Tokyo </p><p><strong>Employees:</strong> 110,000 </p><p><strong>Market value:</strong> ¥12.9tn ($89bn)</p>
	//         </tr:scrollable-text>
	//     </tr:scrollable-section>
	//     <tr:scrollable-section theme-position=\"1\" theme-display=\"1\">
	//         <content type=\"http://www.ft.com/ontology/content/ImageSet\" data-embedded=\"true\" id=\"2065583e-e5f1-4a00-92b0-c3a84fc2e353\"></content>
	//         <tr:scrollable-text>
	//             <p theme-style=\"2\">6-10 per cent</p><p>Price rise for PlayStation 5 consoles in Japan, Europe, China and other key markets this summer due to rising production costs, the global semiconductor shortage and yen weakness</p>
	//         </tr:scrollable-text>
	//     </tr:scrollable-section>
	// </tr:scrollable-block>
	// The scrollable-block and scrollable-section seems to always be published as translucent tags.
	case *contenttree.ScrollyBlock:
		return innerXML, nil
	case *contenttree.ScrollySection:
		return innerXML, nil
	case *contenttree.ScrollyImage:
		return fmt.Sprintf("<ft-content type=\"http://www.ft.com/ontology/content/ImageSet\" url=\"http://api.ft.com/content/%s\" data-embedded=\"true\"></ft-content>", node.ID), nil
	case *contenttree.ScrollyCopy:
		return innerXML, nil
	case *contenttree.ScrollyHeading:
		switch node.Level {
		case "chapter":
			return fmt.Sprintf("<h2>%s</h2>", innerXML), nil
		case "heading", "subheading":
			return fmt.Sprintf("<p>%s</p>", innerXML), nil
		default:
			return "", fmt.Errorf("failed to transform scrolly heading with level %s", node.Level)
		}
	// TODO: Rethink https://github.com/Financial-Times/body-validation-service/pull/80/files (Read the comments)
	//  In the body XML transformation, there is XSLT template that removes all children elements within the h2 tags
	//  in a scrollable-text but leaves their content (text captured between them).
	//  It is not clear how this behaviour is to be replicated.
	//  Additionally, there is XSLT template to match and ignore the @theme-style attributes of h2 and p elements within
	//  translucent:scrollable-text. It is not clear how this behaviour is to be replicated.

	// content tree nodes that were published inside experimental tag and as such are not supported in the "external"
	// body XML format for now
	case *contenttree.Layout:
		return "", nil
	case *contenttree.LayoutSlot:
		return "", nil
	case *contenttree.LayoutImage:
		return "", nil

	case *contenttree.Recommended:
		return "", nil

	case *contenttree.BigNumber:
		return fmt.Sprintf("<big-number><big-number-headline>%s</big-number-headline><big-number-intro>%s</big-number-intro></big-number>", html.EscapeString(node.Number), html.EscapeString(node.Description)), nil

	case *contenttree.CustomCodeComponent:
		attrs := []string{
			"type=\"http://www.ft.com/ontology/content/CustomCodeComponent\"",
			fmt.Sprintf("url=\"http://api.ft.com/content/%s\"", node.ID),
			"data-embedded=\"true\"",
		}
		if node.LayoutWidth != "" {
			attrs = append(attrs, fmt.Sprintf("data-layout-width=\"%s\"", node.LayoutWidth))
		}
		return fmt.Sprintf("<ft-content %s></ft-content>", strings.Join(attrs, " ")), nil
	/*
		TODO: Remove comment to integrate timeline.
		case *contenttree.Timeline:
			{
				titleXML := ""
				if node.Title != "" {
					titleXML = fmt.Sprintf("<h3>%s</h3>", node.Title)
				}
				layoutWidthXML := ""
				if node.LayoutWidth != "" {
					layoutWidthXML = fmt.Sprintf("data-layout-width=\"%s\"", node.LayoutWidth)
				}
				return fmt.Sprintf("<section data-type=\"timeline\" %s>%s<ol data-type=\"timeline_events\">%s</ol></section>", layoutWidthXML, titleXML, innerXML), nil
			}
		case *contenttree.TimelineEvent:
			{
				titleXML := ""
				if node.Title != "" {
					titleXML = fmt.Sprintf("<h4>%s</h4>", node.Title)
				}
				return fmt.Sprintf("<li data-type=\"timeline_event\">%s%s</li>", titleXML, innerXML), nil
			}
	*/
	// content tree nodes which require transformation of their embedded nodes
	case *contenttree.BodyBlock:
		return transformNode(n.GetEmbedded())
	case *contenttree.BlockquoteChild:
		return transformNode(n.GetEmbedded())
	case *contenttree.LayoutChild:
		return transformNode(n.GetEmbedded())
	case *contenttree.LayoutSlotChild:
		return transformNode(n.GetEmbedded())
	case *contenttree.ListItemChild:
		return transformNode(n.GetEmbedded())
	case *contenttree.Phrasing:
		return transformNode(n.GetEmbedded())
	case *contenttree.ScrollyCopyChild:
		return transformNode(n.GetEmbedded())
	case *contenttree.ScrollySectionChild:
		return transformNode(n.GetEmbedded())
	case *contenttree.TableChild:
		return transformNode(n.GetEmbedded())
		/*
			TODO: Remove comment to integrate timeline.
			case *contenttree.TimelineEventChild:
				return transformNode(n.GetEmbedded())
		*/
	}

	return "", nil
}

func optionalBoolAttrXML(name string, value *bool) string {
	if value == nil {
		return ""
	}
	if *value {
		return fmt.Sprintf("%s=\"true\"", name)
	}
	return fmt.Sprintf("%s=\"false\"", name)
}

func tableResponsiveStyleToExternal(v string) string {
	switch v {
	case "flat":
		return "stacked"
	case "scroll":
		return "horizontal-scroll"
	case "overflow":
		return "auto"
	default:
		return ""
	}
}

func tableThemeToExternal(compact, stripes bool) string {
	switch {
	case compact && stripes:
		return "compact-stripes"
	case compact:
		return "compact"
	case stripes:
		return "stripes"
	default:
		return ""
	}
}

func tableHiddenToExternal(hideOnMobile bool) string {
	if hideOnMobile {
		return "small-screen"
	}
	return ""
}

func tableSortTypeToExternal(v string) string {
	switch v {
	case "text":
		return "string"
	case "number":
		return "number"
	case "date":
		return "date"
	case "currency":
		return "currency"
	case "percent":
		return "percent"
	default:
		return ""
	}
}

func buildTH(cell *contenttree.TableCell, columnIndex int, settings []*contenttree.ColumnSettingsItems) (string, error) {
	var cellAttrs []string
	if columnIndex < len(settings) && settings[columnIndex] != nil {
		if hideOnMobile := settings[columnIndex].HideOnMobile; hideOnMobile != nil && *hideOnMobile {
			cellAttrs = append(cellAttrs, `data-column-hidden="small-screen"`)
		}

		if sortable := settings[columnIndex].Sortable; sortable != nil {
			cellAttrs = append(cellAttrs, fmt.Sprintf(`data-column-sortable="%t"`, *sortable))
		}

		if sortType := tableSortTypeToExternal(settings[columnIndex].SortType); sortType != "" {
			cellAttrs = append(cellAttrs, fmt.Sprintf("data-column-type=\"%s\"", sortType))
		}
	}
	if cell.ColumnSpan != nil {
		cellAttrs = append(cellAttrs, fmt.Sprintf("colspan=\"%d\"", *cell.ColumnSpan))
	}
	if cell.RowSpan != nil {
		cellAttrs = append(cellAttrs, fmt.Sprintf("rowspan=\"%d\"", *cell.RowSpan))
	}

	childrenXML := make([]string, 0, len(cell.Children))
	for _, child := range cell.Children {
		s, err := transformNode(child)
		if err != nil {
			return "", err
		}
		childrenXML = append(childrenXML, s)
	}
	innerXML := strings.Join(childrenXML, "")

	if len(cellAttrs) > 0 {
		return fmt.Sprintf("<th %s>%s</th>", strings.Join(cellAttrs, " "), innerXML), nil
	}
	return fmt.Sprintf("<th>%s</th>", innerXML), nil
}
