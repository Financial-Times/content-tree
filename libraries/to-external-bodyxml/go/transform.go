package toexternalbodyxml

import (
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"strings"

	contenttree "github.com/Financial-Times/content-tree"
)

// Transform converts content from the content tree format, provided as unmarshalled JSON (json.RawMessage),
// into an "external" XHTML-formatted version of the same content.
//
// The XHTML output is intended for distribution to consumers that only support widely recognized formats like HTML
// or those that should not receive internal-specific details contained in the content tree format.
// Such consumers may be external (non-FT) users, automated systems processing HTML-based content,
// republishing platforms, and more.
func Transform(root json.RawMessage) (string, error) {
	tree := contenttree.Root{}

	err := json.Unmarshal(root, &tree)
	if err != nil {
		return "", fmt.Errorf("failed to instantiate content tree: %w", err)
	}

	return transformNode(&tree)
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
		return "<br>", nil

	case *contenttree.ThematicBreak:
		return "<hr>", nil

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

		return fmt.Sprintf("<%[1]s>%s</%[1]s>", tag, innerXML), nil

	case *contenttree.Strong:
		return fmt.Sprintf("<strong>%s</strong>", innerXML), nil

	case *contenttree.Emphasis:
		return fmt.Sprintf("<em>%s</em>", innerXML), nil

	case *contenttree.Strikethrough:
		return fmt.Sprintf("<s>%s</s>", innerXML), nil

	case *contenttree.Link:
		parts := strings.Split(node.URL, "/")
		if node.Title != "" {
			return fmt.Sprintf("<a href=\"https://ft.com/content/%s\" title=\"%s\">%s</a>", parts[len(parts)-1], node.Title, innerXML), nil
		}

		return fmt.Sprintf("<a href=\"https://ft.com/content/%s\">%s</a>", parts[len(parts)-1], innerXML), nil

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
		return fmt.Sprintf("<pull-quote><pull-quote-text><p>%s</p></pull-quote-text><pull-quote-source>%s</pull-quote-source></pull-quote>", node.Text, node.Source), nil

	case *contenttree.ImageSet:
		return fmt.Sprintf("<ft-content type=\"http://www.ft.com/ontology/content/ImageSet\" url=\"http://api.ft.com/content/%s\" data-embedded=\"true\"></ft-content>", node.ID), nil

	case *contenttree.Flourish:
		return fmt.Sprintf("<ft-content type=\"http://www.ft.com/ontology/content/Content\" url=\"http://api.ft.com/content/%[1]s\" alt=\"%s\" data-asset-type=\"flourish\" data-embedded=\"true\" data-flourish-type=\"%s\" data-layout-width=\"%s\" data-time-stamp=\"%s\" id=\"%[1]s\"></ft-content>", node.Id, node.Description, node.FlourishType, node.LayoutWidth, node.Timestamp), nil

	case *contenttree.TableCaption:
		return fmt.Sprintf("<caption>%s</caption>", innerXML), nil

	case *contenttree.TableCell:
		return fmt.Sprintf("<td>%s</td>", innerXML), nil

	case *contenttree.TableRow:
		return fmt.Sprintf("<tr>%s</tr>", innerXML), nil

	case *contenttree.TableBody:
		return fmt.Sprintf("<tbody>%s</tbody>", innerXML), nil

	case *contenttree.TableFooter:
		return fmt.Sprintf("<tfoot>%s</tfoot>", innerXML), nil

	// TODO: Additional work on table tags will be required as per the resolution of https://github.com/Financial-Times/content-tree/issues/71
	//  The tables have multiple attributes such as
	//  class=\"data-table\"
	//  data-table-collapse-rownum=\"\"
	//  data-table-layout-largescreen=\"auto\"
	//  data-table-layout-smallscreen=\"auto\"
	//  data-table-theme=\"auto\"
	// Is there a match between the Table node and those attributes?
	case *contenttree.Table:
		return fmt.Sprintf("<table>%s</table>", innerXML), nil

	case *contenttree.Video:
		return fmt.Sprintf("<ft-content type=\"http://www.ft.com/ontology/content/Video\" url=\"http://api.ft.com/content/%s\" data-embedded=\"true\"></ft-content>", node.ID), nil

	case *contenttree.YoutubeVideo:
		return fmt.Sprintf("<a data-asset-type=\"video\" data-embedded=\"true\" href=\"%s\"></a>", node.URL), nil

	case *contenttree.Tweet:
		return fmt.Sprintf("<a data-asset-type=\"tweet\" data-embedded=\"true\" href=\"%[1]s\">%[1]s</a>", node.ID), nil

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
	// TODO: What is the equivalent of scrolly copy in the context for the XML tags?
	case *contenttree.ScrollyCopy:
		return "", nil
	// TODO: What is the equivalent of scrolly heading in the context for the XML tags?
	case *contenttree.ScrollyHeading:
		return "", nil
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
		return fmt.Sprintf("<big-number><big-number-headline>%s</big-number-headline><big-number-intro>%s</big-number-intro></big-number>", node.Number, node.Description), nil

	// CCC nodes won't be available in the "external" body XML format.
	case *contenttree.CustomCodeComponent:
		return "", nil

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
	}

	return "", nil
}
