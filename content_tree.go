/*
Package contenttree defines a content tree object and provides functionality for unmarshalling it into memory.
Each type of node in the tree is represented by a struct that implements the Node interface.

The node structs that embed other structs, such as ListItemChild, are artificially created structs and do not have
a corresponding node definition in the content tree. These structs are used to represent child nodes of a parent node
that can have heterogeneous children — children of different types. This approach enables the representation of
a heterogeneous list of objects in a strongly typed language like Go.

A custom UnmarshalJSON method is implemented for these artificially created types. It handles a specific challenge
when unmarshalling objects with embedded structs that implement the same interface. The customization is necessary
because the embedded structs contain a field with the same name, the field "Type".
According to the official Go documentation (https://pkg.go.dev/encoding/json), when multiple fields with the same name
exist, during unmarshalling they are all ignored, and no error is returned. As a result, the objects are not
unmarshalled correctly unless custom unmarshalling logic is applied.
*/
package contenttree

import (
	"encoding/json"
	"errors"
	"fmt"
)

const (
	RootType                = "root"
	BodyType                = "body"
	TextType                = "text"
	BreakType               = "break"
	ThematicBreakType       = "thematic-break"
	ParagraphType           = "paragraph"
	HeadingType             = "heading"
	StrongType              = "strong"
	EmphasisType            = "emphasis"
	StrikethroughType       = "strikethrough"
	LinkType                = "link"
	ListType                = "list"
	ListItemType            = "list-item"
	BlockquoteType          = "blockquote"
	PullquoteType           = "pullquote"
	ImageSetType            = "image-set"
	RecommendedType         = "recommended"
	TweetType               = "tweet"
	FlourishType            = "flourish"
	BigNumberType           = "big-number"
	VideoType               = "video"
	YoutubeVideoType        = "youtube-video"
	ScrollyBlockType        = "scrolly-block"
	ScrollySectionType      = "scrolly-section"
	ScrollyImageType        = "scrolly-image"
	ScrollyCopyType         = "scrolly-copy"
	ScrollyHeadingType      = "scrolly-heading"
	LayoutType              = "layout"
	LayoutSlotType          = "layout-slot"
	LayoutImageType         = "layout-image"
	TableCaptionType        = "table-caption"
	TableCellType           = "table-cell"
	TableRowType            = "table-row"
	TableBodyType           = "table-body"
	TableFooterType         = "table-footer"
	TableType               = "table"
	CustomCodeComponentType = "custom-code-component"

	BodyBlockType           = "body-block"
	BlockquoteChildType     = "blockquote-child"
	LayoutChildType         = "layout-child"
	LayoutSlotChildType     = "layout-slot-child"
	ListItemChildType       = "list-item-child"
	PhrasingType            = "phrasing"
	ScrollyCopyChildType    = "scrolly-copy-child"
	ScrollySectionChildType = "scrolly-section-child"
	TableChildType          = "table-child"
)

var (
	// returned when calling AppendChild on a node that doesn't own a Children slice
	ErrCannotHaveChildren = errors.New("node cannot have children")
	// returned when a child is not one of the allowed types for a parent
	ErrInvalidChildType = errors.New("invalid child type for this parent")
)

// Node represents a unified interface for different types of content tree nodes.
// It facilitates easy traversal of the tree structure without requiring type casting.
type Node interface {
	// GetType returns the type of the node as a string.
	GetType() string
	// GetChildren retrieves a list of child nodes, enabling hierarchical traversal.
	GetChildren() []Node
	// GetEmbedded returns the embedded node, if applicable.
	// It is useful for traversing node structs which embed other node structs.
	GetEmbedded() Node
	// AppendChild attempts to append a child node, returning an error if not allowed.
	AppendChild(child Node) error
}

// typed() is a small utility to read a node's type without full unmarshal.
func typed(v any) string {
	if n, ok := v.(Node); ok && n != nil {
		return n.GetType()
	}
	return ""
}

// typedNode is a lightweight struct that holds only the type information of a content tree node.
// It is primarily used for unmarshalling when only the node type is required.
type typedNode struct {
	Type string `json:"type"`
}

// ErrUnmarshalInvalidNode is returned when a node fails to unmarshal due to an invalid type.
// This occurs when the node type is not among the allowed child types defined by its parent node.
// It is primarily used in custom UnmarshalJSON methods for structs containing embedded structs.
var ErrUnmarshalInvalidNode = errors.New("unmarshalling node with invalid type")

type ColumnSettingsItems struct {
	HideOnMobile bool   `json:"hideOnMobile,omitempty"`
	SortType     string `json:"sortType,omitempty"`
	Sortable     bool   `json:"sortable,omitempty"`
}

type BigNumber struct {
	Type        string      `json:"type"`
	Data        interface{} `json:"data,omitempty"`
	Description string      `json:"description,omitempty"`
	Number      string      `json:"number,omitempty"`
}

func (n *BigNumber) GetType() string {
	return n.Type
}

func (n *BigNumber) GetEmbedded() Node {
	return nil
}

func (n *BigNumber) GetChildren() []Node {
	return nil
}

type Blockquote struct {
	Type     string             `json:"type"`
	Children []*BlockquoteChild `json:"children,omitempty"`
	Data     interface{}        `json:"data,omitempty"`
}

func (n *Blockquote) GetType() string {
	return n.Type
}

func (n *Blockquote) GetEmbedded() Node {
	return nil
}

func (n *Blockquote) GetChildren() []Node {
	result := make([]Node, len(n.Children))
	for i, v := range n.Children {
		result[i] = v
	}
	return result
}

type BlockquoteChild struct {
	*Paragraph
	*Text
	*Break
	*Strong
	*Emphasis
	*Strikethrough
	*Link
}

func (n *BlockquoteChild) GetType() string {
	return BlockquoteChildType
}

func (n *BlockquoteChild) GetEmbedded() Node {
	if n.Paragraph != nil {
		return n.Paragraph
	}
	if n.Text != nil {
		return n.Text
	}
	if n.Break != nil {
		return n.Break
	}
	if n.Strong != nil {
		return n.Strong
	}
	if n.Emphasis != nil {
		return n.Emphasis
	}
	if n.Strikethrough != nil {
		return n.Strikethrough
	}
	if n.Link != nil {
		return n.Link
	}
	return nil
}

func (n *BlockquoteChild) GetChildren() []Node {
	if n.Paragraph != nil {
		return n.Paragraph.GetChildren()
	}
	if n.Text != nil {
		return n.Text.GetChildren()
	}
	if n.Break != nil {
		return n.Break.GetChildren()
	}
	if n.Strong != nil {
		return n.Strong.GetChildren()
	}
	if n.Emphasis != nil {
		return n.Emphasis.GetChildren()
	}
	if n.Strikethrough != nil {
		return n.Strikethrough.GetChildren()
	}
	if n.Link != nil {
		return n.Link.GetChildren()
	}
	return nil
}

func (n *BlockquoteChild) UnmarshalJSON(data []byte) error {
	var tn typedNode
	if err := json.Unmarshal(data, &tn); err != nil {
		return err
	}

	switch tn.Type {
	case ParagraphType:
		var v Paragraph
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		n.Paragraph = &v
	case TextType:
		var v Text
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		n.Text = &v
	case BreakType:
		var v Break
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		n.Break = &v
	case StrongType:
		var v Strong
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		n.Strong = &v
	case EmphasisType:
		var v Emphasis
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		n.Emphasis = &v
	case StrikethroughType:
		var v Strikethrough
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		n.Strikethrough = &v
	case LinkType:
		var v Link
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		n.Link = &v
	default:
		return fmt.Errorf("failed to unmarshal BlockquoteChild from %s: %w", data, ErrUnmarshalInvalidNode)
	}
	return nil
}

type Body struct {
	Type     string       `json:"type"`
	Children []*BodyBlock `json:"children,omitempty"`
	Data     interface{}  `json:"data,omitempty"`
	Version  float64      `json:"version,omitempty"`
}

func (n *Body) GetType() string {
	return n.Type
}

func (n *Body) GetEmbedded() Node {
	return nil
}

func (n *Body) GetChildren() []Node {
	result := make([]Node, len(n.Children))
	for i, v := range n.Children {
		result[i] = v
	}
	return result
}

type BodyBlock struct {
	*Paragraph
	*Flourish
	*Heading
	*ImageSet
	*BigNumber
	*Layout
	*List
	*Blockquote
	*Pullquote
	*ScrollyBlock
	*ThematicBreak
	*Table
	*Text
	*Recommended
	*Tweet
	*Video
	*YoutubeVideo
	*CustomCodeComponent
}

func (n *BodyBlock) GetType() string {
	return BodyBlockType
}

func (n *BodyBlock) GetEmbedded() Node {
	if n.Paragraph != nil {
		return n.Paragraph
	}
	if n.Flourish != nil {
		return n.Flourish
	}
	if n.Heading != nil {
		return n.Heading
	}
	if n.ImageSet != nil {
		return n.ImageSet
	}
	if n.BigNumber != nil {
		return n.BigNumber
	}
	if n.Layout != nil {
		return n.Layout
	}
	if n.List != nil {
		return n.List
	}
	if n.Blockquote != nil {
		return n.Blockquote
	}
	if n.Pullquote != nil {
		return n.Pullquote
	}
	if n.ScrollyBlock != nil {
		return n.ScrollyBlock
	}
	if n.ThematicBreak != nil {
		return n.ThematicBreak
	}
	if n.Table != nil {
		return n.Table
	}
	if n.Text != nil {
		return n.Text
	}
	if n.Recommended != nil {
		return n.Recommended
	}
	if n.Tweet != nil {
		return n.Tweet
	}
	if n.Video != nil {
		return n.Video
	}
	if n.YoutubeVideo != nil {
		return n.YoutubeVideo
	}
	if n.CustomCodeComponent != nil {
		return n.CustomCodeComponent
	}
	return nil
}

func (n *BodyBlock) GetChildren() []Node {
	if n.Paragraph != nil {
		return n.Paragraph.GetChildren()
	}
	if n.Flourish != nil {
		return n.Flourish.GetChildren()
	}
	if n.Heading != nil {
		return n.Heading.GetChildren()
	}
	if n.ImageSet != nil {
		return n.ImageSet.GetChildren()
	}
	if n.BigNumber != nil {
		return n.BigNumber.GetChildren()
	}
	if n.Layout != nil {
		return n.Layout.GetChildren()
	}
	if n.List != nil {
		return n.List.GetChildren()
	}
	if n.Blockquote != nil {
		return n.Blockquote.GetChildren()
	}
	if n.Pullquote != nil {
		return n.Pullquote.GetChildren()
	}
	if n.ScrollyBlock != nil {
		return n.ScrollyBlock.GetChildren()
	}
	if n.ThematicBreak != nil {
		return n.ThematicBreak.GetChildren()
	}
	if n.Table != nil {
		return n.Table.GetChildren()
	}
	if n.Text != nil {
		return n.Text.GetChildren()
	}
	if n.Recommended != nil {
		return n.Recommended.GetChildren()
	}
	if n.Tweet != nil {
		return n.Tweet.GetChildren()
	}
	if n.Video != nil {
		return n.Video.GetChildren()
	}
	if n.YoutubeVideo != nil {
		return n.YoutubeVideo.GetChildren()
	}
	if n.CustomCodeComponent != nil {
		return n.CustomCodeComponent.GetChildren()
	}
	return nil
}

func (n *BodyBlock) UnmarshalJSON(data []byte) error {
	var tn typedNode
	if err := json.Unmarshal(data, &tn); err != nil {
		return err
	}

	switch tn.Type {
	case ParagraphType:
		var v Paragraph
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		n.Paragraph = &v
	case FlourishType:
		var v Flourish
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		n.Flourish = &v
	case HeadingType:
		var v Heading
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		n.Heading = &v
	case ImageSetType:
		var v ImageSet
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		n.ImageSet = &v
	case BigNumberType:
		var v BigNumber
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		n.BigNumber = &v
	case LayoutType:
		var v Layout
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		n.Layout = &v
	case ListType:
		var v List
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		n.List = &v
	case BlockquoteType:
		var v Blockquote
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		n.Blockquote = &v
	case PullquoteType:
		var v Pullquote
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		n.Pullquote = &v
	case ScrollyBlockType:
		var v ScrollyBlock
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		n.ScrollyBlock = &v
	case ThematicBreakType:
		var v ThematicBreak
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		n.ThematicBreak = &v
	case TableType:
		var v Table
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		n.Table = &v
	case TextType:
		var v Text
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		n.Text = &v
	case RecommendedType:
		var v Recommended
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		n.Recommended = &v
	case TweetType:
		var v Tweet
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		n.Tweet = &v
	case VideoType:
		var v Video
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		n.Video = &v
	case YoutubeVideoType:
		var v YoutubeVideo
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		n.YoutubeVideo = &v
	case CustomCodeComponentType:
		var v CustomCodeComponent
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		n.CustomCodeComponent = &v
	default:
		return fmt.Errorf("failed to unmarshal BodyBlock from %s: %w", data, ErrUnmarshalInvalidNode)
	}
	return nil
}

type Break struct {
	Type string      `json:"type"`
	Data interface{} `json:"data,omitempty"`
}

func (n *Break) GetType() string {
	return n.Type
}

func (n *Break) GetEmbedded() Node {
	return nil
}

func (n *Break) GetChildren() []Node {
	return nil
}

type Emphasis struct {
	Type     string      `json:"type"`
	Children []*Phrasing `json:"children,omitempty"`
	Data     interface{} `json:"data,omitempty"`
}

func (n *Emphasis) GetType() string {
	return n.Type
}

func (n *Emphasis) GetEmbedded() Node {
	return nil
}

func (n *Emphasis) GetChildren() []Node {
	result := make([]Node, len(n.Children))
	for i, v := range n.Children {
		result[i] = v
	}
	return result
}

type Flourish struct {
	Type               string                 `json:"type"`
	Data               interface{}            `json:"data,omitempty"`
	Description        string                 `json:"description,omitempty"`
	FallbackImage      *FlourishFallbackImage `json:"fallbackImage,omitempty"`
	FlourishType       string                 `json:"flourishType,omitempty"`
	Id                 string                 `json:"id,omitempty"`
	LayoutWidth        string                 `json:"layoutWidth,omitempty"`
	Timestamp          string                 `json:"timestamp,omitempty"`
	FragmentIdentifier string                 `json:"fragmentIdentifier,omitempty"`
}

func (n *Flourish) GetType() string {
	return n.Type
}

func (n *Flourish) GetEmbedded() Node {
	return nil
}

func (n *Flourish) GetChildren() []Node {
	return nil
}

type FlourishFallbackImage struct {
	Format    string                               `json:"format,omitempty"`
	Height    float64                              `json:"height,omitempty"`
	Id        string                               `json:"id,omitempty"`
	SourceSet []FlourishFallbackImageSourceSetElem `json:"sourceSet,omitempty"`
	Url       string                               `json:"url,omitempty"`
	Width     float64                              `json:"width,omitempty"`
}

type FlourishFallbackImageSourceSetElem struct {
	Dpr   float64 `json:"dpr,omitempty"`
	Url   string  `json:"url,omitempty"`
	Width float64 `json:"width,omitempty"`
}

type Heading struct {
	Type               string      `json:"type"`
	Children           []*Text     `json:"children,omitempty"`
	Data               interface{} `json:"data,omitempty"`
	Level              string      `json:"level,omitempty"`
	FragmentIdentifier string      `json:"fragmentIdentifier,omitempty"`
}

func (n *Heading) GetType() string {
	return n.Type
}

func (n *Heading) GetEmbedded() Node {
	return nil
}

func (n *Heading) GetChildren() []Node {
	result := make([]Node, len(n.Children))
	for i, v := range n.Children {
		result[i] = v
	}
	return result
}

type ImageSet struct {
	Type               string      `json:"type"`
	Data               interface{} `json:"data,omitempty"`
	ID                 string      `json:"id,omitempty"`
	Picture            *Picture    `json:"picture,omitempty"`
	FragmentIdentifier string      `json:"fragmentIdentifier,omitempty"`
}

func (n *ImageSet) GetType() string {
	return n.Type
}

func (n *ImageSet) GetEmbedded() Node {
	return nil
}

func (n *ImageSet) GetChildren() []Node {
	return nil
}

type Layout struct {
	Type        string         `json:"type"`
	Children    []*LayoutChild `json:"children,omitempty"`
	Data        interface{}    `json:"data,omitempty"`
	LayoutName  string         `json:"layoutName,omitempty"`
	LayoutWidth string         `json:"layoutWidth,omitempty"`
}

func (n *Layout) GetType() string {
	return n.Type
}

func (n *Layout) GetEmbedded() Node {
	return nil
}

func (n *Layout) GetChildren() []Node {
	result := make([]Node, len(n.Children))
	for i, v := range n.Children {
		result[i] = v
	}
	return result
}

type LayoutChild struct {
	*LayoutSlot
	*Heading
	*LayoutImage
}

func (n *LayoutChild) GetType() string {
	return LayoutSlotChildType
}

func (n *LayoutChild) GetEmbedded() Node {
	if n.LayoutSlot != nil {
		return n.LayoutSlot
	}
	if n.Heading != nil {
		return n.Heading
	}
	if n.LayoutImage != nil {
		return n.LayoutImage
	}
	return nil
}

func (n *LayoutChild) GetChildren() []Node {
	if n.LayoutSlot != nil {
		return n.LayoutSlot.GetChildren()
	}
	if n.Heading != nil {
		return n.Heading.GetChildren()
	}
	if n.LayoutImage != nil {
		return n.LayoutImage.GetChildren()
	}
	return nil
}

func (n *LayoutChild) UnmarshalJSON(data []byte) error {
	var tn typedNode
	if err := json.Unmarshal(data, &tn); err != nil {
		return err
	}

	switch tn.Type {
	case LayoutSlotType:
		var v LayoutSlot
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		n.LayoutSlot = &v
	case HeadingType:
		var v Heading
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		n.Heading = &v
	case LayoutImageType:
		var v LayoutImage
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		n.LayoutImage = &v
	default:
		return fmt.Errorf("failed to unmarshal LayoutChild from %s: %w", data, ErrUnmarshalInvalidNode)
	}
	return nil
}

type LayoutImage struct {
	Type    string      `json:"type"`
	Alt     string      `json:"alt,omitempty"`
	Caption string      `json:"caption,omitempty"`
	Credit  string      `json:"credit,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	ID      string      `json:"id,omitempty"`
	Picture *Picture    `json:"picture,omitempty"`
}

func (n *LayoutImage) GetType() string {
	return n.Type
}

func (n *LayoutImage) GetEmbedded() Node {
	return nil
}

func (n *LayoutImage) GetChildren() []Node {
	return nil
}

type LayoutSlot struct {
	Type     string             `json:"type"`
	Children []*LayoutSlotChild `json:"children,omitempty"`
	Data     interface{}        `json:"data,omitempty"`
}

func (n *LayoutSlot) GetType() string {
	return n.Type
}

func (n *LayoutSlot) GetEmbedded() Node {
	return nil
}

func (n *LayoutSlot) GetChildren() []Node {
	result := make([]Node, len(n.Children))
	for i, v := range n.Children {
		result[i] = v
	}
	return result
}

type LayoutSlotChild struct {
	*Paragraph
	*Heading
	*LayoutImage
}

func (n *LayoutSlotChild) GetType() string {
	return LayoutSlotChildType
}

func (n *LayoutSlotChild) GetEmbedded() Node {
	if n.Paragraph != nil {
		return n.Paragraph
	}
	if n.Heading != nil {
		return n.Heading
	}
	if n.LayoutImage != nil {
		return n.LayoutImage
	}
	return nil
}

func (n *LayoutSlotChild) GetChildren() []Node {
	if n.Paragraph != nil {
		return n.Paragraph.GetChildren()
	}
	if n.Heading != nil {
		return n.Heading.GetChildren()
	}
	if n.LayoutImage != nil {
		return n.LayoutImage.GetChildren()
	}
	return nil
}

func (n *LayoutSlotChild) UnmarshalJSON(data []byte) error {
	var tn typedNode
	if err := json.Unmarshal(data, &tn); err != nil {
		return err
	}

	switch tn.Type {
	case ParagraphType:
		var v Paragraph
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		n.Paragraph = &v
	case HeadingType:
		var v Heading
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		n.Heading = &v
	case LayoutImageType:
		var v LayoutImage
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		n.LayoutImage = &v
	default:
		return fmt.Errorf("failed to unmarshal LayoutSlotChild from %s: %w", data, ErrUnmarshalInvalidNode)
	}
	return nil
}

type Link struct {
	Type     string      `json:"type"`
	Children []*Phrasing `json:"children,omitempty"`
	Data     interface{} `json:"data,omitempty"`
	Title    string      `json:"title,omitempty"`
	URL      string      `json:"url,omitempty"`
}

func (n *Link) GetType() string {
	return n.Type
}

func (n *Link) GetEmbedded() Node {
	return nil
}

func (n *Link) GetChildren() []Node {
	result := make([]Node, len(n.Children))
	for i, v := range n.Children {
		result[i] = v
	}
	return result
}

type List struct {
	Type     string      `json:"type"`
	Children []*ListItem `json:"children,omitempty"`
	Data     interface{} `json:"data,omitempty"`
	Ordered  bool        `json:"ordered,omitempty"`
}

func (n *List) GetType() string {
	return n.Type
}

func (n *List) GetEmbedded() Node {
	return nil
}

func (n *List) GetChildren() []Node {
	result := make([]Node, len(n.Children))
	for i, v := range n.Children {
		result[i] = v
	}
	return result
}

type ListItem struct {
	Type     string           `json:"type"`
	Children []*ListItemChild `json:"children,omitempty"`
	Data     interface{}      `json:"data,omitempty"`
}

func (n *ListItem) GetType() string {
	return n.Type
}

func (n *ListItem) GetEmbedded() Node {
	return nil
}

func (n *ListItem) GetChildren() []Node {
	result := make([]Node, len(n.Children))
	for i, v := range n.Children {
		result[i] = v
	}
	return result
}

type ListItemChild struct {
	*Paragraph
	*Text
	*Break
	*Strong
	*Emphasis
	*Strikethrough
	*Link
}

func (n *ListItemChild) GetType() string {
	return ListItemChildType
}

func (n *ListItemChild) GetEmbedded() Node {
	if n.Paragraph != nil {
		return n.Paragraph
	}
	if n.Text != nil {
		return n.Text
	}
	if n.Break != nil {
		return n.Break
	}
	if n.Strong != nil {
		return n.Strong
	}
	if n.Emphasis != nil {
		return n.Emphasis
	}
	if n.Strikethrough != nil {
		return n.Strikethrough
	}
	if n.Link != nil {
		return n.Link
	}
	return nil
}

func (n *ListItemChild) GetChildren() []Node {
	if n.Paragraph != nil {
		return n.Paragraph.GetChildren()
	}
	if n.Text != nil {
		return n.Text.GetChildren()
	}
	if n.Break != nil {
		return n.Break.GetChildren()
	}
	if n.Strong != nil {
		return n.Strong.GetChildren()
	}
	if n.Emphasis != nil {
		return n.Emphasis.GetChildren()
	}
	if n.Strikethrough != nil {
		return n.Strikethrough.GetChildren()
	}
	if n.Link != nil {
		return n.Link.GetChildren()
	}
	return nil
}

func (n *ListItemChild) UnmarshalJSON(data []byte) error {
	type node struct {
		Type string `json:"type"`
	}
	var tn node
	if err := json.Unmarshal(data, &tn); err != nil {
		return err
	}

	switch tn.Type {
	case ParagraphType:
		var v Paragraph
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		n.Paragraph = &v
	case TextType:
		var v Text
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		n.Text = &v
	case BreakType:
		var v Break
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		n.Break = &v
	case StrongType:
		var v Strong
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		n.Strong = &v
	case EmphasisType:
		var v Emphasis
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		n.Emphasis = &v
	case StrikethroughType:
		var v Strikethrough
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		n.Strikethrough = &v
	case LinkType:
		var v Link
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		n.Link = &v
	default:
		return fmt.Errorf("failed to unmarshal ListItemChild from %s: %w", data, ErrUnmarshalInvalidNode)
	}
	return nil
}

type Paragraph struct {
	Type     string      `json:"type"`
	Children []*Phrasing `json:"children,omitempty"`
	Data     interface{} `json:"data,omitempty"`
}

func (n *Paragraph) GetType() string {
	return n.Type
}

func (n *Paragraph) GetEmbedded() Node {
	return nil
}

func (n *Paragraph) GetChildren() []Node {
	result := make([]Node, len(n.Children))
	for i, v := range n.Children {
		result[i] = v
	}
	return result
}

type Phrasing struct {
	*Text
	*Break
	*Strong
	*Emphasis
	*Strikethrough
	*Link
}

func (n *Phrasing) GetType() string {
	return PhrasingType
}

func (n *Phrasing) GetEmbedded() Node {
	if n.Text != nil {
		return n.Text
	}
	if n.Break != nil {
		return n.Break
	}
	if n.Strong != nil {
		return n.Strong
	}
	if n.Emphasis != nil {
		return n.Emphasis
	}
	if n.Strikethrough != nil {
		return n.Strikethrough
	}
	if n.Link != nil {
		return n.Link
	}
	return nil
}

func (n *Phrasing) GetChildren() []Node {
	if n.Text != nil {
		return n.Text.GetChildren()
	}
	if n.Break != nil {
		return n.Break.GetChildren()
	}
	if n.Strong != nil {
		return n.Strong.GetChildren()
	}
	if n.Emphasis != nil {
		return n.Emphasis.GetChildren()
	}
	if n.Strikethrough != nil {
		return n.Strikethrough.GetChildren()
	}
	if n.Link != nil {
		return n.Link.GetChildren()
	}
	return nil
}

func (n *Phrasing) UnmarshalJSON(data []byte) error {
	var tn typedNode
	if err := json.Unmarshal(data, &tn); err != nil {
		return err
	}

	switch tn.Type {
	case TextType:
		var v Text
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		n.Text = &v
	case BreakType:
		var v Break
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		n.Break = &v
	case StrongType:
		var v Strong
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		n.Strong = &v
	case EmphasisType:
		var v Emphasis
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		n.Emphasis = &v
	case StrikethroughType:
		var v Strikethrough
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		n.Strikethrough = &v
	case LinkType:
		var v Link
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		n.Link = &v
	default:
		return fmt.Errorf("failed to unmarshal Phrasing from %s: %w", data, ErrUnmarshalInvalidNode)
	}
	return nil
}

type Pullquote struct {
	Type   string      `json:"type"`
	Data   interface{} `json:"data,omitempty"`
	Source string      `json:"source,omitempty"`
	Text   string      `json:"text,omitempty"`
}

func (n *Pullquote) GetType() string {
	return n.Type
}

func (n *Pullquote) GetEmbedded() Node {
	return nil
}

func (n *Pullquote) GetChildren() []Node {
	return nil
}

type Recommended struct {
	Type                string      `json:"type"`
	Data                interface{} `json:"data,omitempty"`
	Heading             string      `json:"heading,omitempty"`
	ID                  string      `json:"id,omitempty"`
	Teaser              *Teaser     `json:"teaser,omitempty"`
	TeaserTitleOverride string      `json:"teaserTitleOverride,omitempty"`
}

func (n *Recommended) GetType() string {
	return n.Type
}

func (n *Recommended) GetEmbedded() Node {
	return nil
}

func (n *Recommended) GetChildren() []Node {
	return nil
}

type ScrollyBlock struct {
	Type     string            `json:"type"`
	Children []*ScrollySection `json:"children,omitempty"`
	Data     interface{}       `json:"data,omitempty"`
	Theme    string            `json:"theme,omitempty"`
}

func (n *ScrollyBlock) GetType() string {
	return n.Type
}

func (n *ScrollyBlock) GetEmbedded() Node {
	return nil
}

func (n *ScrollyBlock) GetChildren() []Node {
	result := make([]Node, len(n.Children))
	for i, v := range n.Children {
		result[i] = v
	}
	return result
}

type ScrollyCopy struct {
	Type     string              `json:"type"`
	Children []*ScrollyCopyChild `json:"children,omitempty"`
	Data     interface{}         `json:"data,omitempty"`
}

func (n *ScrollyCopy) GetType() string {
	return n.Type
}

func (n *ScrollyCopy) GetEmbedded() Node {
	return nil
}

func (n *ScrollyCopy) GetChildren() []Node {
	result := make([]Node, len(n.Children))
	for i, v := range n.Children {
		result[i] = v
	}
	return result
}

type ScrollyCopyChild struct {
	*Paragraph
	*ScrollyHeading
}

func (n *ScrollyCopyChild) GetType() string {
	return ScrollyCopyChildType
}

func (n *ScrollyCopyChild) GetEmbedded() Node {
	if n.Paragraph != nil {
		return n.Paragraph
	}
	if n.ScrollyHeading != nil {
		return n.ScrollyHeading
	}
	return nil
}

func (n *ScrollyCopyChild) GetChildren() []Node {
	if n.Paragraph != nil {
		return n.Paragraph.GetChildren()
	}
	if n.ScrollyHeading != nil {
		return n.ScrollyHeading.GetChildren()
	}
	return nil
}

func (n *ScrollyCopyChild) UnmarshalJSON(data []byte) error {
	var tn typedNode
	if err := json.Unmarshal(data, &tn); err != nil {
		return err
	}

	switch tn.Type {
	case ParagraphType:
		var v Paragraph
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		n.Paragraph = &v
	case ScrollyHeadingType:
		var v ScrollyHeading
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		n.ScrollyHeading = &v
	default:
		return fmt.Errorf("failed to unmarshal ScrollyCopyChild from %s: %w", data, ErrUnmarshalInvalidNode)
	}
	return nil
}

type ScrollyHeading struct {
	Type     string      `json:"type"`
	Children []*Text     `json:"children,omitempty"`
	Data     interface{} `json:"data,omitempty"`
	Level    string      `json:"level,omitempty"`
}

func (n *ScrollyHeading) GetType() string {
	return n.Type
}

func (n *ScrollyHeading) GetEmbedded() Node {
	return nil
}

func (n *ScrollyHeading) GetChildren() []Node {
	result := make([]Node, len(n.Children))
	for i, v := range n.Children {
		result[i] = v
	}
	return result
}

type ScrollyImage struct {
	Type    string      `json:"type"`
	Data    interface{} `json:"data,omitempty"`
	ID      string      `json:"id,omitempty"`
	Picture *Picture    `json:"picture,omitempty"`
}

func (n *ScrollyImage) GetType() string {
	return n.Type
}

func (n *ScrollyImage) GetEmbedded() Node {
	return nil
}

func (n *ScrollyImage) GetChildren() []Node {
	return nil
}

type ScrollySection struct {
	Type       string                 `json:"type"`
	Children   []*ScrollySectionChild `json:"children,omitempty"`
	Data       interface{}            `json:"data,omitempty"`
	Display    string                 `json:"display,omitempty"`
	NoBox      bool                   `json:"noBox,omitempty"`
	Position   string                 `json:"position,omitempty"`
	Transition string                 `json:"transition,omitempty"`
}

func (n *ScrollySection) GetType() string {
	return n.Type
}

func (n *ScrollySection) GetEmbedded() Node {
	return nil
}

func (n *ScrollySection) GetChildren() []Node {
	result := make([]Node, len(n.Children))
	for i, v := range n.Children {
		result[i] = v
	}
	return result
}

type ScrollySectionChild struct {
	*ScrollyCopy
	*ScrollyImage
}

func (n *ScrollySectionChild) GetType() string {
	return ScrollySectionChildType
}

func (n *ScrollySectionChild) GetEmbedded() Node {
	if n.ScrollyCopy != nil {
		return n.ScrollyCopy
	}
	if n.ScrollyImage != nil {
		return n.ScrollyImage
	}
	return nil
}

func (n *ScrollySectionChild) GetChildren() []Node {
	if n.ScrollyCopy != nil {
		return n.ScrollyCopy.GetChildren()
	}
	if n.ScrollyImage != nil {
		return n.ScrollyImage.GetChildren()
	}
	return nil
}

func (n *ScrollySectionChild) UnmarshalJSON(data []byte) error {
	var tn typedNode
	if err := json.Unmarshal(data, &tn); err != nil {
		return err
	}

	switch tn.Type {
	case ScrollyCopyType:
		var v ScrollyCopy
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		n.ScrollyCopy = &v
	case ScrollyImageType:
		var v ScrollyImage
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		n.ScrollyImage = &v
	default:
		return fmt.Errorf("failed to unmarshal ScrollySectionChild from %s: %w", data, ErrUnmarshalInvalidNode)
	}
	return nil
}

type Strikethrough struct {
	Type     string      `json:"type"`
	Children []*Phrasing `json:"children,omitempty"`
	Data     interface{} `json:"data,omitempty"`
}

func (n *Strikethrough) GetType() string {
	return n.Type
}

func (n *Strikethrough) GetEmbedded() Node {
	return nil
}

func (n *Strikethrough) GetChildren() []Node {
	result := make([]Node, len(n.Children))
	for i, v := range n.Children {
		result[i] = v
	}
	return result
}

type Strong struct {
	Type     string      `json:"type"`
	Children []*Phrasing `json:"children,omitempty"`
	Data     interface{} `json:"data,omitempty"`
}

func (n *Strong) GetType() string {
	return n.Type
}

func (n *Strong) GetEmbedded() Node {
	return nil
}

func (n *Strong) GetChildren() []Node {
	result := make([]Node, len(n.Children))
	for i, v := range n.Children {
		result[i] = v
	}
	return result
}

type Table struct {
	Type                     string                 `json:"type"`
	Children                 []*TableChild          `json:"children,omitempty"`
	CollapseAfterHowManyRows float64                `json:"collapseAfterHowManyRows,omitempty"`
	ColumnSettings           []*ColumnSettingsItems `json:"columnSettings,omitempty"`
	Compact                  bool                   `json:"compact,omitempty"`
	Data                     interface{}            `json:"data,omitempty"`
	LayoutWidth              string                 `json:"layoutWidth,omitempty"`
	ResponsiveStyle          string                 `json:"responsiveStyle,omitempty"`
	Stripes                  bool                   `json:"stripes,omitempty"`
}

func (n *Table) GetType() string {
	return n.Type
}

func (n *Table) GetEmbedded() Node {
	return nil
}

func (n *Table) GetChildren() []Node {
	result := make([]Node, len(n.Children))
	for i, v := range n.Children {
		result[i] = v
	}
	return result
}

type TableChild struct {
	*TableCaption
	*TableBody
	*TableFooter
}

func (n *TableChild) GetType() string {
	return TableChildType
}

func (n *TableChild) GetEmbedded() Node {
	if n.TableCaption != nil {
		return n.TableCaption
	}
	if n.TableBody != nil {
		return n.TableCaption
	}
	if n.TableBody != nil {
		return n.TableCaption
	}
	return nil
}

func (n *TableChild) GetChildren() []Node {
	if n.TableCaption != nil {
		return n.TableCaption.GetChildren()
	}
	if n.TableBody != nil {
		return n.TableCaption.GetChildren()
	}
	if n.TableBody != nil {
		return n.TableCaption.GetChildren()
	}
	return nil
}

func (n *TableChild) UnmarshalJSON(data []byte) error {
	var tn typedNode
	if err := json.Unmarshal(data, &tn); err != nil {
		return err
	}

	switch tn.Type {
	case TableCaptionType:
		var v TableCaption
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		n.TableCaption = &v
	case TableBodyType:
		var v TableBody
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		n.TableBody = &v
	case TableFooterType:
		var v TableFooter
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		n.TableFooter = &v
	default:
		return fmt.Errorf("failed to unmarshal TableChild from %s: %w", data, ErrUnmarshalInvalidNode)
	}
	return nil
}

type TableBody struct {
	Type     string      `json:"type"`
	Children []*TableRow `json:"children,omitempty"`
	Data     interface{} `json:"data,omitempty"`
}

func (n *TableBody) GetType() string {
	return n.Type
}

func (n *TableBody) GetEmbedded() Node {
	return nil
}

func (n *TableBody) GetChildren() []Node {
	result := make([]Node, len(n.Children))
	for i, v := range n.Children {
		result[i] = v
	}
	return result
}

type TableCaption struct {
	Type     string      `json:"type"`
	Children []*Table    `json:"children,omitempty"`
	Data     interface{} `json:"data,omitempty"`
}

func (n *TableCaption) GetType() string {
	return n.Type
}

func (n *TableCaption) GetEmbedded() Node {
	return nil
}

func (n *TableCaption) GetChildren() []Node {
	result := make([]Node, len(n.Children))
	for i, v := range n.Children {
		result[i] = v
	}
	return result
}

type TableCell struct {
	Type     string      `json:"type"`
	Children []*Table    `json:"children,omitempty"`
	Data     interface{} `json:"data,omitempty"`
	Heading  bool        `json:"heading,omitempty"`
}

func (n *TableCell) GetType() string {
	return n.Type
}

func (n *TableCell) GetEmbedded() Node {
	return nil
}

func (n *TableCell) GetChildren() []Node {
	result := make([]Node, len(n.Children))
	for i, v := range n.Children {
		result[i] = v
	}
	return result
}

type TableFooter struct {
	Type     string      `json:"type"`
	Children []*Table    `json:"children,omitempty"`
	Data     interface{} `json:"data,omitempty"`
}

func (n *TableFooter) GetType() string {
	return n.Type
}

func (n *TableFooter) GetEmbedded() Node {
	return nil
}

func (n *TableFooter) GetChildren() []Node {
	result := make([]Node, len(n.Children))
	for i, v := range n.Children {
		result[i] = v
	}
	return result
}

type TableRow struct {
	Type     string       `json:"type"`
	Children []*TableCell `json:"children,omitempty"`
	Data     interface{}  `json:"data,omitempty"`
}

func (n *TableRow) GetType() string {
	return n.Type
}

func (n *TableRow) GetEmbedded() Node {
	return nil
}

func (n *TableRow) GetChildren() []Node {
	result := make([]Node, len(n.Children))
	for i, v := range n.Children {
		result[i] = v
	}
	return result
}

type Text struct {
	Type  string      `json:"type"`
	Data  interface{} `json:"data,omitempty"`
	Value string      `json:"value,omitempty"`
}

func (n *Text) GetType() string {
	return n.Type
}

func (n *Text) GetEmbedded() Node {
	return nil
}

func (n *Text) GetChildren() []Node {
	return nil
}

type ThematicBreak struct {
	Type string      `json:"type"`
	Data interface{} `json:"data,omitempty"`
}

func (n *ThematicBreak) GetType() string {
	return n.Type
}

func (n *ThematicBreak) GetEmbedded() Node {
	return nil
}

func (n *ThematicBreak) GetChildren() []Node {
	return nil
}

type Tweet struct {
	Type string      `json:"type"`
	Data interface{} `json:"data,omitempty"`
	HTML string      `json:"html,omitempty"`
	ID   string      `json:"id,omitempty"`
}

func (n *Tweet) GetType() string {
	return n.Type
}

func (n *Tweet) GetEmbedded() Node {
	return nil
}

func (n *Tweet) GetChildren() []Node {
	return nil
}

type Video struct {
	Type string      `json:"type"`
	Data interface{} `json:"data,omitempty"`
	ID   string      `json:"id,omitempty"`
}

func (n *Video) GetType() string {
	return n.Type
}

func (n *Video) GetEmbedded() Node {
	return nil
}

func (n *Video) GetChildren() []Node {
	return nil
}

type YoutubeVideo struct {
	Type string      `json:"type"`
	Data interface{} `json:"data,omitempty"`
	URL  string      `json:"url,omitempty"`
}

func (n *YoutubeVideo) GetType() string {
	return n.Type
}

func (n *YoutubeVideo) GetEmbedded() Node {
	return nil
}

func (n *YoutubeVideo) GetChildren() []Node {
	return nil
}

type CustomCodeComponent struct {
	Type                   string                 `json:"type"`
	Data                   interface{}            `json:"data,omitempty"`
	ID                     string                 `json:"id,omitempty"`
	LayoutWidth            string                 `json:"layoutWidth,omitempty"`
	Attributes             map[string]interface{} `json:"attributes,omitempty"`
	AttributesLastModified string                 `json:"attributesLastModified,omitempty"`
	Path                   string                 `json:"path,omitempty"`
	VersionRange           string                 `json:"versionRange,omitempty"`
}

func (n *CustomCodeComponent) GetType() string {
	return n.Type
}

func (n *CustomCodeComponent) GetEmbedded() Node {
	return nil
}

func (n *CustomCodeComponent) GetChildren() []Node {
	return nil
}

type FallbackImage struct {
	Format    string            `json:"format,omitempty"`
	Height    float64           `json:"height,omitempty"`
	ID        string            `json:"id,omitempty"`
	SourceSet []*SourceSetItems `json:"sourceSet,omitempty"`
	URL       string            `json:"url,omitempty"`
	Width     float64           `json:"width,omitempty"`
}

type Image struct {
	Format    string            `json:"format,omitempty"`
	Height    float64           `json:"height,omitempty"`
	ID        string            `json:"id,omitempty"`
	SourceSet []*SourceSetItems `json:"sourceSet,omitempty"`
	URL       string            `json:"url,omitempty"`
	Width     float64           `json:"width,omitempty"`
}

type ImagesItems struct {
	Format    string            `json:"format,omitempty"`
	Height    float64           `json:"height,omitempty"`
	ID        string            `json:"id,omitempty"`
	SourceSet []*SourceSetItems `json:"sourceSet,omitempty"`
	URL       string            `json:"url,omitempty"`
	Width     float64           `json:"width,omitempty"`
}

type Indicators struct {
	AccessLevel     string `json:"accessLevel,omitempty"`
	IsColumn        bool   `json:"isColumn,omitempty"`
	IsEditorsChoice bool   `json:"isEditorsChoice,omitempty"`
	IsExclusive     bool   `json:"isExclusive,omitempty"`
	IsOpinion       bool   `json:"isOpinion,omitempty"`
	IsPodcast       bool   `json:"isPodcast,omitempty"`
	IsScoop         bool   `json:"isScoop,omitempty"`
}

type MetaAltLink struct {
	Type       string   `json:"type"`
	APIUrl     string   `json:"apiUrl,omitempty"`
	DirectType string   `json:"directType,omitempty"`
	ID         string   `json:"id,omitempty"`
	Predicate  string   `json:"predicate,omitempty"`
	PrefLabel  string   `json:"prefLabel,omitempty"`
	Types      []string `json:"types,omitempty"`
	URL        string   `json:"url,omitempty"`
}

func (n *MetaAltLink) GetType() string {
	return n.Type
}

func (n *MetaAltLink) GetEmbedded() Node {
	return nil
}

func (n *MetaAltLink) GetChildren() []Node {
	return nil
}

type MetaLink struct {
	Type       string   `json:"type"`
	APIUrl     string   `json:"apiUrl,omitempty"`
	DirectType string   `json:"directType,omitempty"`
	ID         string   `json:"id,omitempty"`
	Predicate  string   `json:"predicate,omitempty"`
	PrefLabel  string   `json:"prefLabel,omitempty"`
	Types      []string `json:"types,omitempty"`
	URL        string   `json:"url,omitempty"`
}

func (n *MetaLink) GetType() string {
	return n.Type
}

func (n *MetaLink) GetEmbedded() Node {
	return nil
}

func (n *MetaLink) GetChildren() []Node {
	return nil
}

type Picture struct {
	Alt           string         `json:"alt,omitempty"`
	Caption       string         `json:"caption,omitempty"`
	Credit        string         `json:"credit,omitempty"`
	FallbackImage *FallbackImage `json:"fallbackImage,omitempty"`
	ImageType     string         `json:"imageType,omitempty"`
	Images        []*ImagesItems `json:"images,omitempty"`
	LayoutWidth   string         `json:"layoutWidth,omitempty"`
}

type Root struct {
	Type string      `json:"type"`
	Body *Body       `json:"body,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

func (n *Root) GetType() string {
	return n.Type
}

func (n *Root) GetEmbedded() Node {
	return nil
}

func (n *Root) GetChildren() []Node {
	return nil
}

type SourceSetItems struct {
	Dpr   float64 `json:"dpr,omitempty"`
	URL   string  `json:"url,omitempty"`
	Width float64 `json:"width,omitempty"`
}

type Teaser struct {
	Type               string       `json:"type"`
	FirstPublishedDate string       `json:"firstPublishedDate,omitempty"`
	ID                 string       `json:"id,omitempty"`
	Image              *Image       `json:"image,omitempty"`
	Indicators         *Indicators  `json:"indicators,omitempty"`
	MetaAltLink        *MetaAltLink `json:"metaAltLink,omitempty"`
	MetaLink           *MetaLink    `json:"metaLink,omitempty"`
	MetaPrefixText     string       `json:"metaPrefixText,omitempty"`
	MetaSuffixText     string       `json:"metaSuffixText,omitempty"`
	PublishedDate      string       `json:"publishedDate,omitempty"`
	Title              string       `json:"title,omitempty"`
	URL                string       `json:"url,omitempty"`
}

func (n *Teaser) GetType() string {
	return n.Type
}

func (n *Teaser) GetEmbedded() Node {
	return nil
}

func (n *Teaser) GetChildren() []Node {
	return nil
}

// Build a BodyBlock wrapper from any allowed top-level block node.
func makeBodyBlock(n Node) (*BodyBlock, error) {
	switch n.GetType() {
	case ParagraphType:
		return &BodyBlock{Paragraph: n.(*Paragraph)}, nil
	case FlourishType:
		return &BodyBlock{Flourish: n.(*Flourish)}, nil
	case HeadingType:
		return &BodyBlock{Heading: n.(*Heading)}, nil
	case ImageSetType:
		return &BodyBlock{ImageSet: n.(*ImageSet)}, nil
	case BigNumberType:
		return &BodyBlock{BigNumber: n.(*BigNumber)}, nil
	case LayoutType:
		return &BodyBlock{Layout: n.(*Layout)}, nil
	case ListType:
		return &BodyBlock{List: n.(*List)}, nil
	case BlockquoteType:
		return &BodyBlock{Blockquote: n.(*Blockquote)}, nil
	case PullquoteType:
		return &BodyBlock{Pullquote: n.(*Pullquote)}, nil
	case ScrollyBlockType:
		return &BodyBlock{ScrollyBlock: n.(*ScrollyBlock)}, nil
	case ThematicBreakType:
		return &BodyBlock{ThematicBreak: n.(*ThematicBreak)}, nil
	case TableType:
		return &BodyBlock{Table: n.(*Table)}, nil
	case TextType:
		return &BodyBlock{Text: n.(*Text)}, nil
	case RecommendedType:
		return &BodyBlock{Recommended: n.(*Recommended)}, nil
	case TweetType:
		return &BodyBlock{Tweet: n.(*Tweet)}, nil
	case VideoType:
		return &BodyBlock{Video: n.(*Video)}, nil
	case YoutubeVideoType:
		return &BodyBlock{YoutubeVideo: n.(*YoutubeVideo)}, nil
	case CustomCodeComponentType:
		return &BodyBlock{CustomCodeComponent: n.(*CustomCodeComponent)}, nil
	default:
		return nil, ErrInvalidChildType
	}
}

// Build a Phrasing wrapper for paragraph/phrasing-bearing parents.
func makePhrasing(n Node) (*Phrasing, error) {
	switch n.GetType() {
	case TextType:
		return &Phrasing{Text: n.(*Text)}, nil
	case BreakType:
		return &Phrasing{Break: n.(*Break)}, nil
	case StrongType:
		return &Phrasing{Strong: n.(*Strong)}, nil
	case EmphasisType:
		return &Phrasing{Emphasis: n.(*Emphasis)}, nil
	case StrikethroughType:
		return &Phrasing{Strikethrough: n.(*Strikethrough)}, nil
	case LinkType:
		return &Phrasing{Link: n.(*Link)}, nil
	default:
		return nil, ErrInvalidChildType
	}
}

// Build a BlockquoteChild wrapper.
func makeBlockquoteChild(n Node) (*BlockquoteChild, error) {
	switch n.GetType() {
	case ParagraphType:
		return &BlockquoteChild{Paragraph: n.(*Paragraph)}, nil
	case TextType:
		return &BlockquoteChild{Text: n.(*Text)}, nil
	case BreakType:
		return &BlockquoteChild{Break: n.(*Break)}, nil
	case StrongType:
		return &BlockquoteChild{Strong: n.(*Strong)}, nil
	case EmphasisType:
		return &BlockquoteChild{Emphasis: n.(*Emphasis)}, nil
	case StrikethroughType:
		return &BlockquoteChild{Strikethrough: n.(*Strikethrough)}, nil
	case LinkType:
		return &BlockquoteChild{Link: n.(*Link)}, nil
	default:
		return nil, ErrInvalidChildType
	}
}

// Build a ListItemChild wrapper.
func makeListItemChild(n Node) (*ListItemChild, error) {
	switch n.GetType() {
	case ParagraphType:
		return &ListItemChild{Paragraph: n.(*Paragraph)}, nil
	case TextType:
		return &ListItemChild{Text: n.(*Text)}, nil
	case BreakType:
		return &ListItemChild{Break: n.(*Break)}, nil
	case StrongType:
		return &ListItemChild{Strong: n.(*Strong)}, nil
	case EmphasisType:
		return &ListItemChild{Emphasis: n.(*Emphasis)}, nil
	case StrikethroughType:
		return &ListItemChild{Strikethrough: n.(*Strikethrough)}, nil
	case LinkType:
		return &ListItemChild{Link: n.(*Link)}, nil
	default:
		return nil, ErrInvalidChildType
	}
}

// Build LayoutChild wrapper.
func makeLayoutChild(n Node) (*LayoutChild, error) {
	switch n.GetType() {
	case LayoutSlotType:
		return &LayoutChild{LayoutSlot: n.(*LayoutSlot)}, nil
	case HeadingType:
		return &LayoutChild{Heading: n.(*Heading)}, nil
	case LayoutImageType:
		return &LayoutChild{LayoutImage: n.(*LayoutImage)}, nil
	default:
		return nil, ErrInvalidChildType
	}
}

// Build LayoutSlotChild wrapper.
func makeLayoutSlotChild(n Node) (*LayoutSlotChild, error) {
	switch n.GetType() {
	case ParagraphType:
		return &LayoutSlotChild{Paragraph: n.(*Paragraph)}, nil
	case HeadingType:
		return &LayoutSlotChild{Heading: n.(*Heading)}, nil
	case LayoutImageType:
		return &LayoutSlotChild{LayoutImage: n.(*LayoutImage)}, nil
	default:
		return nil, ErrInvalidChildType
	}
}

// Build ScrollySectionChild wrapper.
func makeScrollySectionChild(n Node) (*ScrollySectionChild, error) {
	switch n.GetType() {
	case ScrollyCopyType:
		return &ScrollySectionChild{ScrollyCopy: n.(*ScrollyCopy)}, nil
	case ScrollyImageType:
		return &ScrollySectionChild{ScrollyImage: n.(*ScrollyImage)}, nil
	default:
		return nil, ErrInvalidChildType
	}
}

// Build ScrollyCopyChild wrapper.
func makeScrollyCopyChild(n Node) (*ScrollyCopyChild, error) {
	switch n.GetType() {
	case ParagraphType:
		return &ScrollyCopyChild{Paragraph: n.(*Paragraph)}, nil
	case ScrollyHeadingType:
		return &ScrollyCopyChild{ScrollyHeading: n.(*ScrollyHeading)}, nil
	default:
		return nil, ErrInvalidChildType
	}
}

// Build TableChild wrapper.
func makeTableChild(n Node) (*TableChild, error) {
	switch n.GetType() {
	case TableCaptionType:
		return &TableChild{TableCaption: n.(*TableCaption)}, nil
	case TableBodyType:
		return &TableChild{TableBody: n.(*TableBody)}, nil
	case TableFooterType:
		return &TableChild{TableFooter: n.(*TableFooter)}, nil
	default:
		return nil, ErrInvalidChildType
	}
}

// ===== Nodes WITH children =====

func (n *Body) AppendChild(child Node) error {
	if n == nil {
		return fmt.Errorf("nil Body: %w", ErrCannotHaveChildren)
	}
	bb, err := makeBodyBlock(child)
	if err != nil {
		return err
	}
	n.Children = append(n.Children, bb)
	return nil
}

func (n *Paragraph) AppendChild(child Node) error {
	p, err := makePhrasing(child)
	if err != nil {
		return err
	}
	n.Children = append(n.Children, p)
	return nil
}

func (n *Emphasis) AppendChild(child Node) error {
	p, err := makePhrasing(child)
	if err != nil {
		return err
	}
	n.Children = append(n.Children, p)
	return nil
}

func (n *Strong) AppendChild(child Node) error {
	p, err := makePhrasing(child)
	if err != nil {
		return err
	}
	n.Children = append(n.Children, p)
	return nil
}

func (n *Strikethrough) AppendChild(child Node) error {
	p, err := makePhrasing(child)
	if err != nil {
		return err
	}
	n.Children = append(n.Children, p)
	return nil
}

func (n *Link) AppendChild(child Node) error {
	p, err := makePhrasing(child)
	if err != nil {
		return err
	}
	n.Children = append(n.Children, p)
	return nil
}

func (n *Blockquote) AppendChild(child Node) error {
	c, err := makeBlockquoteChild(child)
	if err != nil {
		return err
	}
	n.Children = append(n.Children, c)
	return nil
}

func (n *List) AppendChild(child Node) error {
	// Keep strict: only accept ListItem
	if child.GetType() != ListItemType {
		return ErrInvalidChildType
	}
	n.Children = append(n.Children, child.(*ListItem))
	return nil
}

func (n *ListItem) AppendChild(child Node) error {
	c, err := makeListItemChild(child)
	if err != nil {
		return err
	}
	n.Children = append(n.Children, c)
	return nil
}

func (n *Layout) AppendChild(child Node) error {
	c, err := makeLayoutChild(child)
	if err != nil {
		return err
	}
	n.Children = append(n.Children, c)
	return nil
}

func (n *LayoutSlot) AppendChild(child Node) error {
	c, err := makeLayoutSlotChild(child)
	if err != nil {
		return err
	}
	n.Children = append(n.Children, c)
	return nil
}

func (n *Heading) AppendChild(child Node) error {
	// Heading only allows Text nodes
	if child.GetType() != TextType {
		return ErrInvalidChildType
	}
	n.Children = append(n.Children, child.(*Text))
	return nil
}

func (n *ScrollyBlock) AppendChild(child Node) error {
	if child.GetType() != ScrollySectionType {
		return ErrInvalidChildType
	}
	n.Children = append(n.Children, child.(*ScrollySection))
	return nil
}

func (n *ScrollySection) AppendChild(child Node) error {
	c, err := makeScrollySectionChild(child)
	if err != nil {
		return err
	}
	n.Children = append(n.Children, c)
	return nil
}

func (n *ScrollyCopy) AppendChild(child Node) error {
	c, err := makeScrollyCopyChild(child)
	if err != nil {
		return err
	}
	n.Children = append(n.Children, c)
	return nil
}

func (n *Table) AppendChild(child Node) error {
	c, err := makeTableChild(child)
	if err != nil {
		return err
	}
	n.Children = append(n.Children, c)
	return nil
}

func (n *TableBody) AppendChild(child Node) error {
	if child.GetType() != TableRowType {
		return ErrInvalidChildType
	}
	n.Children = append(n.Children, child.(*TableRow))
	return nil
}

func (n *TableRow) AppendChild(child Node) error {
	if child.GetType() != TableCellType {
		return ErrInvalidChildType
	}
	n.Children = append(n.Children, child.(*TableCell))
	return nil
}

func (n *TableCaption) AppendChild(child Node) error {
	if child.GetType() != TableType {
		return ErrInvalidChildType
	}
	n.Children = append(n.Children, child.(*Table))
	return nil
}

func (n *TableFooter) AppendChild(child Node) error {
	if child.GetType() != TableType {
		return ErrInvalidChildType
	}
	n.Children = append(n.Children, child.(*Table))
	return nil
}

func (n *TableCell) AppendChild(child Node) error {
	if child.GetType() != TableType {
		return ErrInvalidChildType
	}
	n.Children = append(n.Children, child.(*Table))
	return nil
}

func (n *BigNumber) AppendChild(child Node) error           { return ErrCannotHaveChildren }
func (n *Break) AppendChild(child Node) error               { return ErrCannotHaveChildren }
func (n *Flourish) AppendChild(child Node) error            { return ErrCannotHaveChildren }
func (n *ImageSet) AppendChild(child Node) error            { return ErrCannotHaveChildren }
func (n *LayoutImage) AppendChild(child Node) error         { return ErrCannotHaveChildren }
func (n *Pullquote) AppendChild(child Node) error           { return ErrCannotHaveChildren }
func (n *Recommended) AppendChild(child Node) error         { return ErrCannotHaveChildren }
func (n *ScrollyImage) AppendChild(child Node) error        { return ErrCannotHaveChildren }
func (n *ScrollyHeading) AppendChild(child Node) error      { return ErrCannotHaveChildren }
func (n *Text) AppendChild(child Node) error                { return ErrCannotHaveChildren }
func (n *ThematicBreak) AppendChild(child Node) error       { return ErrCannotHaveChildren }
func (n *Tweet) AppendChild(child Node) error               { return ErrCannotHaveChildren }
func (n *Video) AppendChild(child Node) error               { return ErrCannotHaveChildren }
func (n *YoutubeVideo) AppendChild(child Node) error        { return ErrCannotHaveChildren }
func (n *CustomCodeComponent) AppendChild(child Node) error { return ErrCannotHaveChildren }
func (n *Root) AppendChild(child Node) error                { return ErrCannotHaveChildren }

// union wrappers (no own Children slice)
func (n *BodyBlock) AppendChild(child Node) error           { return ErrCannotHaveChildren }
func (n *Phrasing) AppendChild(child Node) error            { return ErrCannotHaveChildren }
func (n *BlockquoteChild) AppendChild(child Node) error     { return ErrCannotHaveChildren }
func (n *ListItemChild) AppendChild(child Node) error       { return ErrCannotHaveChildren }
func (n *LayoutChild) AppendChild(child Node) error         { return ErrCannotHaveChildren }
func (n *LayoutSlotChild) AppendChild(child Node) error     { return ErrCannotHaveChildren }
func (n *ScrollySectionChild) AppendChild(child Node) error { return ErrCannotHaveChildren }
func (n *ScrollyCopyChild) AppendChild(child Node) error    { return ErrCannotHaveChildren }
func (n *TableChild) AppendChild(child Node) error          { return ErrCannotHaveChildren }
