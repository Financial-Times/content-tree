package toexternalbodyxml

import (
	"encoding/json"
	"errors"
	"fmt"
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
		return node.Value, nil

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

	// TODO: This implementation is a placeholder. There are different types of links which need to be transformed to
	//  different XHTML tags. For example, there are links that need to be transformed into "<ft-content>" or
	//  "<ft-related>" tags, there are anchors links that shouldn't be transformed at all, and there are regular links
	//  that should be transformed into <a> tags.
	//  This implementation is a placeholder which handles only a link to an FT article.
	//  In seems that the content tree link object at the moment does not provide enough information to distinguish
	//  between different types of links.
	// Example(https://www.ft.com/content/069e537a-ffc2-11e7-9650-9c0ad2d7c5b5):
	// <ft-content type=\"http://www.ft.com/ontology/content/Article\" url=\"http://api.ft.com/content/674697de-fbb5-11e7-9b32-d7d59aace167\" title=\"Apple pledges to invest $30bn and pay $38bn tax bill\">plans to spend $350bn</ft-content>
	case *contenttree.Link:
		parts := strings.Split(node.URL, "/")
		if node.Title != "" {
			return fmt.Sprintf("<ft-content type=\"http://www.ft.com/ontology/content/Article\" url=\"http://api.ft.com/content/%s\" title=\"%s\">%s</ft-content>", parts[len(parts)-1], node.Title, innerXML), nil
		}
		return fmt.Sprintf("<ft-content type=\"http://www.ft.com/ontology/content/Article\" url=\"http://api.ft.com/content/%s\">%s</ftcontent>", parts[len(parts)-1], innerXML), nil

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

	// TODO: The <pull-quote> tag is not a standard HTML tag, it is a custom tag used by the FT. It is worth to
	//  reconsider whether external consumers should receive this tag or it should be transformed into a standard HTML.
	// TODO: The pull-quote node doesn't support children called <pull-quote-image> in the HTML representation.
	//  There is old content with <pull-quote-image> tags.
	//  Example(https://www.ft.com/content/e76980da-3585-11e7-99bd-13beb0903fa3):
	//  <pull-quote>
	//  	<pull-quote-text><p>Norwegian has been very lucky that as they’ve grown, the fuel price has halved. I think without that they wouldn’t be around</p></pull-quote-text>
	//  	<pull-quote-image><ft-content type=\"http://www.ft.com/ontology/content/ImageSet\" url=\"http://api.ft.com/content/1888fefa-3718-11e7-07db-84246ae494ea\" data-embedded=\"true\"></ft-content></pull-quote-image>
	//  	<pull-quote-source>Oliver Sleath, analyst at Barclays</pull-quote-source>
	//  </pull-quote>
	case *contenttree.Pullquote:
		return fmt.Sprintf("<pull-quote><pull-quote-text><p>%s</p></pull-quote-text><pull-quote-source>%s</pull-quote-source></pull-quote>", node.Text, node.Source), nil

	case *contenttree.ImageSet:
		return fmt.Sprintf("<ft-content type=\"http://www.ft.com/ontology/content/ImageSet\" url=\"http://api.ft.com/content/%s\" data-embedded=\"true\"></ft-content>", node.ID), nil

	// TODO: The current content tree definition does include Flourish nodes but the JSON schemas does not.
	// case *contenttree.Flourish:
	// 	return "", nil

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

	// TODO: The tables have multiple attributes such as
	//  class=\"data-table\"
	//  data-table-collapse-rownum=\"\"
	//  data-table-layout-largescreen=\"auto\"
	//  data-table-layout-smallscreen=\"auto\"
	//  data-table-theme=\"auto\"
	// Is there a match between the Table node and those attributes?
	case *contenttree.Table:
		return fmt.Sprintf("<table>%s</table>", innerXML), nil

	// Example(https://www.ft.com/content/9c0516cf-dd12-4665-aa22-712de854fe2f):
	// <ft-content type=\"http://www.ft.com/ontology/content/Video\" url=\"http://api.ft.com/content/1c199563-e2cd-4817-990f-79972f3828fb\" data-embedded=\"true\"></ft-content>
	case *contenttree.Video:
		return fmt.Sprintf("<ft-content type=\"http://www.ft.com/ontology/content/Video\" url=\"http://api.ft.com/content/%s\" data-embedded=\"%t\"></ft-content>", node.ID, node.Embedded), nil

	// TODO: The XHTML representation is more generic, applicable to any video source. Nothing specifies that the video
	//  source is YouTube.
	// Example(https://www.ft.com/content/4d9396e4-cb4b-4937-baa3-97fce6f5cb94):
	// <a data-asset-type=\"video\" data-embedded=\"true\" href=\"https://www.youtube.com/watch?v=Y_uIs_Z9z4w\"></a>
	case *contenttree.YoutubeVideo:
		return fmt.Sprintf("<a data-asset-type=\"video\" data-embedded=\"true\" href=\"%s\"></a>", node.URL), nil

	// TODO: The tweets were represented as anchor tags which require href url to the tweet.
	//  The current content tree definition does not include the url attribute.
	// Example (https://www.ft.com/content/b2899d25-9b16-461d-b406-89cfcadf3afc):
	// <a data-asset-type=\"tweet\" data-embedded=\"true\" href=\"https://x.com/sama/status/1882106524090482701\">https://x.com/sama/status/1882106524090482701</a>
	case *contenttree.Tweet:
		return fmt.Sprintf("<a data-asset-type=\"tweet\" data-embedded=\"true\" href=\"%[1]s\">%[1]s</a>", "unknown url from the tweet"), nil

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
	// body XML format
	case *contenttree.Layout:
		return "", nil
	case *contenttree.LayoutSlot:
		return "", nil
	case *contenttree.LayoutImage:
		return "", nil

	// Example(https://www.ft.com/content/bb94946c-1c76-11e8-aaca-4574d7dabfb6):
	// <recommended>
	//     <recommended-title>Recommended</recommended-title>
	//     <ul>
	//         <li><ft-content type=\"http://www.ft.com/ontology/content/Article\" url=\"http://api.ft.com/content/c46e915c-1bde-11e8-956a-43db76e69936\">Brussels primes political ‘grenades’ in first draft of Brexit treaty</ft-content></li>
	//         <li><ft-content type=\"http://www.ft.com/ontology/content/Article\" url=\"http://api.ft.com/content/63423938-1bde-11e8-aaca-4574d7dabfb6\">Michel Barnier expresses frustration with David Davis over Brexit talks</ft-content></li>
	//         <li><ft-content type=\"http://www.ft.com/ontology/content/Article\" url=\"http://api.ft.com/content/515a2b2c-1c60-11e8-aaca-4574d7dabfb6\">Irish border ‘being used to keep UK in EU’ says Johnson</ft-content></li>
	//     </ul>
	// </recommended>
	// TODO: It seems that <recommended> tags are always published as opaque tags.
	case *contenttree.Recommended:
		return "", nil

	// Example(https://www.ft.com/content/0107f2fa-f75c-11e6-9516-2d969e0d3b65)
	// <big-number>
	//     <big-number-headline>
	//         <p>34%</p>
	//     </big-number-headline>
	//     <big-number-intro>
	//         <p>Gender pay gap in financial services, the largest in any sector according to Pwc</p>
	//     </big-number-intro>
	// </big-number>
	case *contenttree.BigNumber:
		return fmt.Sprintf("<big-number><big-number-headline><p>%s</p></big-number-headline><big-number-intro><p>%s</p></big-number-intro></big-number>", node.Number, node.Description), nil

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

// TODO: The namespaces "opaque" and "translucent" are not part of the content tree definition.
//  However they are important for the transformation to the "external" XHTML. By definition each opaque tag
//  is stripped along with all its children from the XHTML. The translucent tags are stripped but their children
//  are kept in the XHTML.
//  Example of opaque tag(https://www.ft.com/content/6a858bff-476a-44f7-91af-da636f0d6b93):
//  <opaque:recommended>
//  	<recommended-title>Recommended</recommended-title>
//  	<ul>
//  		<li><content type=\"http://www.ft.com/ontology/content/Content\" id=\"52946dd2-7316-420b-aa0b-2ac13ea5ea68\">Syria caught up in Lebanon fallout</content></li>
//  	</ul>
//  </opaque:recommended>
//  Example of translucent is the scrollable block example above.
//  The recent implementation of anchor tags relies on the "translucent" namespace as well.

// TODO: The content tree definition lacks "concept" nodes. However, there is a lot of content pieces published
//  with <concept>/<ft-concept> tags in the past.
// Example(https://www.ft.com/content/a2868e64-4e37-11e4-bfda-00144feab7de):
// <ft-concept type=\"http://www.ft.com/ontology/company/PublicCompany\" url=\"http://api.ft.com/organisations/897610dc-fc82-3257-a4f7-c26abfca3bb6\">edX</ft-concept>

// TODO: The content tree definition lacks "related" nodes. There very very few old content pieces (2) published
//  with <related>/<ft-related> tags.
//  Example(https://www.ft.com/content/fa6de70c-e9b8-11e6-893c-082c54a7f539 and https://www.ft.com/content/8885c026-8a1b-11e6-8cb7-e7ada1d123b1):
//  <ft-related type=\"http://www.ft.com/ontology/content/Article\" url=\"http://api.ft.com/content/cc5795be-bb25-11e6-8b45-b8b81dd5d080\">
//  	<title>Analysis</title>
//  	<headline>Venezuela struggles to tame triple-digit inflation\n</headline>
//  	<media><ft-content type=\"http://www.ft.com/ontology/content/ImageSet\" url=\"http://api.ft.com/content/aee2c57e-dc0b-11e6-18ca-65c901033a8f\" data-embedded=\"true\"></ft-content></media>
//  	<intro><p>Shopkeepers resort to weighing banknotes in echo of Weimar Germany</p></intro>
//  </ft-related>

// TODO: There is a tag <promo-box> which doesn't seem to be used anymore but there is still old content using it.
//  Example(https://www.ft.com/content/50c06966-277c-11e3-ae16-00144feab7de):
//  <promo-box>
//      <promo-title><p>In depth</p></promo-title>
//      <promo-headline><p><a href=\"http://www.ft.com/indepth/libor-scandal\" title=\"Libor scandal in depth - FT.com\">Libor scandal</a>\n</p></promo-headline>
//      <promo-image><ft-content type=\"http://www.ft.com/ontology/content/ImageSet\" url=\"http://api.ft.com/content/ecec1d70-6b95-11e1-3243-978e959e1fd3\" data-embedded=\"true\"></ft-content></promo-image>
//      <promo-intro><p>Regulators across the globe probe alleged manipulation by US and European banks of the London interbank offered rate and other key benchmark lending rates</p></promo-intro>
//  </promo-box>

// TODO: We introduced a new tag for CCC - <fallback> which is not part of the content tree definition.

// TODO: There are many node types which have property "Data". It is not clear how it is utilised and whether it
//  should be taken in consideration when transforming to the "external" XHTML version of the content.
//  Such nodes are Blockquote, Link, Table, Video, YoutubeVideo and more.

// TODO: The bodyXML of many articles contain "\n" which is not represented in the content tree. The transformer won't
//  be able to produce exactly the same bodyXML as the one published with "\n" in it.
