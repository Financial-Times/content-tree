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

A custom MarshalJSON method is required for union wrapper structs (e.g. BodyBlock, Phrasing, BlockquoteChild, etc.) because they embed
multiple anonymous pointer fields that all export overlapping JSON field names like "type" and "data". The encoding/json package ignores conflicting
fields when marshalling, which results in empty "{}" objects. These MarshalJSON methods ensure only the active
(non-nil) embedded node is serialized.
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
	RecommendedListType     = "recommended-list"
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
	ClipSetType             = "clip-set"

	BodyBlockType           = "body-block"
	BlockquoteChildType     = "blockquote-child"
	LayoutChildType         = "layout-child"
	LayoutSlotChildType     = "layout-slot-child"
	ListItemChildType       = "list-item-child"
	PhrasingType            = "phrasing"
	ScrollyCopyChildType    = "scrolly-copy-child"
	ScrollySectionChildType = "scrolly-section-child"
	TableChildType          = "table-child"

	TimelineType           = "timeline"
	TimelineEventType      = "timeline-event"
	TimelineEventChildType = "timeline-event-child"
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
	Description string      `json:"description"`
	Number      string      `json:"number"`
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

func (n *BigNumber) AppendChild(_ Node) error { return ErrCannotHaveChildren }

type Blockquote struct {
	Type     string             `json:"type"`
	Children []*BlockquoteChild `json:"children"`
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

func (n *Blockquote) AppendChild(child Node) error {
	c, err := makeBlockquoteChild(child)
	if err != nil {
		return err
	}
	n.Children = append(n.Children, c)
	return nil
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

func (n *BlockquoteChild) AppendChild(child Node) error { return ErrCannotHaveChildren }

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

func (n *BlockquoteChild) MarshalJSON() ([]byte, error) {
	switch {
	case n.Paragraph != nil:
		return json.Marshal(n.Paragraph)
	case n.Text != nil:
		return json.Marshal(n.Text)
	case n.Break != nil:
		return json.Marshal(n.Break)
	case n.Strong != nil:
		return json.Marshal(n.Strong)
	case n.Emphasis != nil:
		return json.Marshal(n.Emphasis)
	case n.Strikethrough != nil:
		return json.Marshal(n.Strikethrough)
	case n.Link != nil:
		return json.Marshal(n.Link)
	default:
		return []byte(`{}`), nil
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

type Body struct {
	Type     string       `json:"type"`
	Children []*BodyBlock `json:"children"`
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
	*RecommendedList
	*Tweet
	*Video
	*YoutubeVideo
	*CustomCodeComponent
	*ClipSet
	*Timeline
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
	if n.RecommendedList != nil {
		return n.RecommendedList
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
	if n.ClipSet != nil {
		return n.ClipSet
	}
	if n.Timeline != nil {
		return n.Timeline
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
	if n.RecommendedList != nil {
		return n.RecommendedList.GetChildren()
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
	if n.ClipSet != nil {
		return n.ClipSet.GetChildren()
	}
	if n.Timeline != nil {
		return n.Timeline.GetChildren()
	}
	return nil
}

func (n *BodyBlock) AppendChild(_ Node) error { return ErrCannotHaveChildren }

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
	case RecommendedListType:
		var v RecommendedList
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		n.RecommendedList = &v
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
	case ClipSetType:
		var v ClipSet
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		n.ClipSet = &v
	case TimelineType:
		var v Timeline
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		n.Timeline = &v
	default:
		return fmt.Errorf("failed to unmarshal BodyBlock from %s: %w", data, ErrUnmarshalInvalidNode)
	}
	return nil
}

func (n *BodyBlock) MarshalJSON() ([]byte, error) {
	switch {
	case n.Paragraph != nil:
		return json.Marshal(n.Paragraph)
	case n.Flourish != nil:
		return json.Marshal(n.Flourish)
	case n.Heading != nil:
		return json.Marshal(n.Heading)
	case n.ImageSet != nil:
		return json.Marshal(n.ImageSet)
	case n.BigNumber != nil:
		return json.Marshal(n.BigNumber)
	case n.Layout != nil:
		return json.Marshal(n.Layout)
	case n.List != nil:
		return json.Marshal(n.List)
	case n.Blockquote != nil:
		return json.Marshal(n.Blockquote)
	case n.Pullquote != nil:
		return json.Marshal(n.Pullquote)
	case n.ScrollyBlock != nil:
		return json.Marshal(n.ScrollyBlock)
	case n.ThematicBreak != nil:
		return json.Marshal(n.ThematicBreak)
	case n.Table != nil:
		return json.Marshal(n.Table)
	case n.Text != nil:
		return json.Marshal(n.Text)
	case n.Recommended != nil:
		return json.Marshal(n.Recommended)
	case n.RecommendedList != nil:
		return json.Marshal(n.RecommendedList)
	case n.Tweet != nil:
		return json.Marshal(n.Tweet)
	case n.Video != nil:
		return json.Marshal(n.Video)
	case n.YoutubeVideo != nil:
		return json.Marshal(n.YoutubeVideo)
	case n.CustomCodeComponent != nil:
		return json.Marshal(n.CustomCodeComponent)
	case n.ClipSet != nil:
		return json.Marshal(n.ClipSet)
	case n.Timeline != nil:
		return json.Marshal(n.Timeline)
	default:
		return []byte(`{}`), nil
	}
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
	case RecommendedListType:
		return &BodyBlock{RecommendedList: n.(*RecommendedList)}, nil
	case TweetType:
		return &BodyBlock{Tweet: n.(*Tweet)}, nil
	case VideoType:
		return &BodyBlock{Video: n.(*Video)}, nil
	case YoutubeVideoType:
		return &BodyBlock{YoutubeVideo: n.(*YoutubeVideo)}, nil
	case CustomCodeComponentType:
		return &BodyBlock{CustomCodeComponent: n.(*CustomCodeComponent)}, nil
	case ClipSetType:
		return &BodyBlock{ClipSet: n.(*ClipSet)}, nil
	case TimelineType:
		return &BodyBlock{Timeline: n.(*Timeline)}, nil
	default:
		return nil, ErrInvalidChildType
	}
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

func (n *Break) AppendChild(_ Node) error { return ErrCannotHaveChildren }

type Emphasis struct {
	Type     string      `json:"type"`
	Children []*Phrasing `json:"children"`
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

func (n *Emphasis) AppendChild(child Node) error {
	p, err := makePhrasing(child)
	if err != nil {
		return err
	}
	n.Children = append(n.Children, p)
	return nil
}

type Flourish struct {
	Type               string                 `json:"type"`
	Data               interface{}            `json:"data,omitempty"`
	Description        string                 `json:"description"`
	FallbackImage      *FlourishFallbackImage `json:"fallbackImage,omitempty"`
	FlourishType       string                 `json:"flourishType,omitempty"`
	Id                 string                 `json:"id,omitempty"`
	LayoutWidth        string                 `json:"layoutWidth"`
	Timestamp          string                 `json:"timestamp"`
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

func (n *Flourish) AppendChild(_ Node) error { return ErrCannotHaveChildren }

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
	Children           []*Text     `json:"children"`
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

func (n *Heading) AppendChild(child Node) error {
	if child.GetType() != TextType {
		return ErrInvalidChildType
	}
	n.Children = append(n.Children, child.(*Text))
	return nil
}

type ImageSet struct {
	Type               string      `json:"type"`
	Data               interface{} `json:"data,omitempty"`
	ID                 string      `json:"id"`
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

func (n *ImageSet) AppendChild(_ Node) error { return ErrCannotHaveChildren }

type Layout struct {
	Type        string         `json:"type"`
	Children    []*LayoutChild `json:"children"`
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

func (n *Layout) AppendChild(child Node) error {
	c, err := makeLayoutChild(child)
	if err != nil {
		return err
	}
	n.Children = append(n.Children, c)
	return nil
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

func (n *LayoutChild) AppendChild(_ Node) error { return ErrCannotHaveChildren }

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

func (n *LayoutChild) MarshalJSON() ([]byte, error) {
	switch {
	case n.LayoutSlot != nil:
		return json.Marshal(n.LayoutSlot)
	case n.Heading != nil:
		return json.Marshal(n.Heading)
	case n.LayoutImage != nil:
		return json.Marshal(n.LayoutImage)
	default:
		return []byte(`{}`), nil
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

type LayoutImage struct {
	Type    string      `json:"type"`
	Alt     string      `json:"alt"`
	Caption string      `json:"caption"`
	Credit  string      `json:"credit"`
	Data    interface{} `json:"data,omitempty"`
	ID      string      `json:"id"`
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

func (n *LayoutImage) AppendChild(_ Node) error { return ErrCannotHaveChildren }

type LayoutSlot struct {
	Type     string             `json:"type"`
	Children []*LayoutSlotChild `json:"children"`
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

func (n *LayoutSlot) AppendChild(child Node) error {
	c, err := makeLayoutSlotChild(child)
	if err != nil {
		return err
	}
	n.Children = append(n.Children, c)
	return nil
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

func (n *LayoutSlotChild) AppendChild(_ Node) error { return ErrCannotHaveChildren }

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

func (n *LayoutSlotChild) MarshalJSON() ([]byte, error) {
	switch {
	case n.Paragraph != nil:
		return json.Marshal(n.Paragraph)
	case n.Heading != nil:
		return json.Marshal(n.Heading)
	case n.LayoutImage != nil:
		return json.Marshal(n.LayoutImage)
	default:
		return []byte(`{}`), nil
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

type Link struct {
	Type      string      `json:"type"`
	Children  []*Phrasing `json:"children"`
	Data      interface{} `json:"data,omitempty"`
	Title     string      `json:"title"`
	URL       string      `json:"url"`
	StyleType string      `json:"styleType,omitempty"`
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

func (n *Link) AppendChild(child Node) error {
	p, err := makePhrasing(child)
	if err != nil {
		return err
	}
	n.Children = append(n.Children, p)
	return nil
}

type List struct {
	Type     string      `json:"type"`
	Children []*ListItem `json:"children"`
	Data     interface{} `json:"data,omitempty"`
	Ordered  bool        `json:"ordered"`
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

func (n *List) AppendChild(child Node) error {
	// Keep strict: only accept ListItem
	if child.GetType() != ListItemType {
		return ErrInvalidChildType
	}
	n.Children = append(n.Children, child.(*ListItem))
	return nil
}

type ListItem struct {
	Type     string           `json:"type"`
	Children []*ListItemChild `json:"children"`
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

func (n *ListItem) AppendChild(child Node) error {
	c, err := makeListItemChild(child)
	if err != nil {
		return err
	}
	n.Children = append(n.Children, c)
	return nil
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

func (n *ListItemChild) AppendChild(child Node) error { return ErrCannotHaveChildren }

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

func (n *ListItemChild) MarshalJSON() ([]byte, error) {
	switch {
	case n.Paragraph != nil:
		return json.Marshal(n.Paragraph)
	case n.Text != nil:
		return json.Marshal(n.Text)
	case n.Break != nil:
		return json.Marshal(n.Break)
	case n.Strong != nil:
		return json.Marshal(n.Strong)
	case n.Emphasis != nil:
		return json.Marshal(n.Emphasis)
	case n.Strikethrough != nil:
		return json.Marshal(n.Strikethrough)
	case n.Link != nil:
		return json.Marshal(n.Link)
	default:
		return []byte(`{}`), nil
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

type Paragraph struct {
	Type     string      `json:"type"`
	Children []*Phrasing `json:"children"`
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

func (n *Paragraph) AppendChild(child Node) error {
	p, err := makePhrasing(child)
	if err != nil {
		return err
	}
	n.Children = append(n.Children, p)
	return nil
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

func (n *Phrasing) AppendChild(_ Node) error { return ErrCannotHaveChildren }

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

func (n *Phrasing) MarshalJSON() ([]byte, error) {
	switch {
	case n.Text != nil:
		return json.Marshal(n.Text)
	case n.Break != nil:
		return json.Marshal(n.Break)
	case n.Strong != nil:
		return json.Marshal(n.Strong)
	case n.Emphasis != nil:
		return json.Marshal(n.Emphasis)
	case n.Strikethrough != nil:
		return json.Marshal(n.Strikethrough)
	case n.Link != nil:
		return json.Marshal(n.Link)
	default:
		return []byte(`{}`), nil
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

type Pullquote struct {
	Type   string      `json:"type"`
	Data   interface{} `json:"data,omitempty"`
	Source string      `json:"source"`
	Text   string      `json:"text"`
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

func (n *Pullquote) AppendChild(_ Node) error { return ErrCannotHaveChildren }

type Recommended struct {
	Type                string      `json:"type"`
	Data                interface{} `json:"data,omitempty"`
	Heading             string      `json:"heading"`
	ID                  string      `json:"id"`
	Teaser              *Teaser     `json:"teaser,omitempty"`
	TeaserTitleOverride string      `json:"teaserTitleOverride"`
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

func (n *Recommended) AppendChild(_ Node) error { return ErrCannotHaveChildren }

type RecommendedList struct {
	Type     string         `json:"type"`
	Data     interface{}    `json:"data,omitempty"`
	Heading  string         `json:"heading,omitempty"`
	Children []*Recommended `json:"children"`
}

func (n *RecommendedList) GetType() string {
	return n.Type
}

func (n *RecommendedList) GetEmbedded() Node {
	return nil
}

func (n *RecommendedList) GetChildren() []Node {
	result := make([]Node, len(n.Children))
	for i, v := range n.Children {
		result[i] = v
	}
	return result
}

func (n *RecommendedList) AppendChild(child Node) error {
	if child.GetType() != RecommendedType {
		return ErrInvalidChildType
	}
	n.Children = append(n.Children, child.(*Recommended))
	return nil
}

type ScrollyBlock struct {
	Type     string            `json:"type"`
	Children []*ScrollySection `json:"children"`
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

func (n *ScrollyBlock) AppendChild(child Node) error {
	if child.GetType() != ScrollySectionType {
		return ErrInvalidChildType
	}
	n.Children = append(n.Children, child.(*ScrollySection))
	return nil
}

type ScrollyCopy struct {
	Type     string              `json:"type"`
	Children []*ScrollyCopyChild `json:"children"`
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

func (n *ScrollyCopy) AppendChild(child Node) error {
	c, err := makeScrollyCopyChild(child)
	if err != nil {
		return err
	}
	n.Children = append(n.Children, c)
	return nil
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

func (n *ScrollyCopyChild) AppendChild(child Node) error { return ErrCannotHaveChildren }

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

func (n *ScrollyCopyChild) MarshalJSON() ([]byte, error) {
	switch {
	case n.Paragraph != nil:
		return json.Marshal(n.Paragraph)
	case n.ScrollyHeading != nil:
		return json.Marshal(n.ScrollyHeading)
	default:
		return []byte(`{}`), nil
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

type ScrollyHeading struct {
	Type     string      `json:"type"`
	Children []*Text     `json:"children"`
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

func (n *ScrollyHeading) AppendChild(_ Node) error { return ErrCannotHaveChildren }

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

func (n *ScrollyImage) AppendChild(_ Node) error { return ErrCannotHaveChildren }

type ScrollySection struct {
	Type       string                 `json:"type"`
	Children   []*ScrollySectionChild `json:"children"`
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

func (n *ScrollySection) AppendChild(child Node) error {
	c, err := makeScrollySectionChild(child)
	if err != nil {
		return err
	}
	n.Children = append(n.Children, c)
	return nil
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

func (n *ScrollySectionChild) AppendChild(_ Node) error { return ErrCannotHaveChildren }

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

func (n *ScrollySectionChild) MarshalJSON() ([]byte, error) {
	switch {
	case n.ScrollyCopy != nil:
		return json.Marshal(n.ScrollyCopy)
	case n.ScrollyImage != nil:
		return json.Marshal(n.ScrollyImage)
	default:
		return []byte(`{}`), nil
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

type Strikethrough struct {
	Type     string      `json:"type"`
	Children []*Phrasing `json:"children"`
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

func (n *Strikethrough) AppendChild(child Node) error {
	p, err := makePhrasing(child)
	if err != nil {
		return err
	}
	n.Children = append(n.Children, p)
	return nil
}

type Strong struct {
	Type     string      `json:"type"`
	Children []*Phrasing `json:"children"`
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

func (n *Strong) AppendChild(child Node) error {
	p, err := makePhrasing(child)
	if err != nil {
		return err
	}
	n.Children = append(n.Children, p)
	return nil
}

type Table struct {
	Type                     string                 `json:"type"`
	Children                 []*TableChild          `json:"children"`
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

func (n *Table) AppendChild(child Node) error {
	c, err := makeTableChild(child)
	if err != nil {
		return err
	}
	n.Children = append(n.Children, c)
	return nil
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

func (n *TableChild) AppendChild(_ Node) error { return ErrCannotHaveChildren }

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

func (n *TableChild) MarshalJSON() ([]byte, error) {
	switch {
	case n.TableCaption != nil:
		return json.Marshal(n.TableCaption)
	case n.TableBody != nil:
		return json.Marshal(n.TableBody)
	case n.TableFooter != nil:
		return json.Marshal(n.TableFooter)
	default:
		return []byte(`{}`), nil
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

type TableBody struct {
	Type     string      `json:"type"`
	Children []*TableRow `json:"children"`
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

func (n *TableBody) AppendChild(child Node) error {
	if child.GetType() != TableRowType {
		return ErrInvalidChildType
	}
	n.Children = append(n.Children, child.(*TableRow))
	return nil
}

type TableCaption struct {
	Type     string      `json:"type"`
	Children []*Table    `json:"children"`
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

func (n *TableCaption) AppendChild(child Node) error {
	if child.GetType() != TableType {
		return ErrInvalidChildType
	}
	n.Children = append(n.Children, child.(*Table))
	return nil
}

type TableCell struct {
	Type     string      `json:"type"`
	Children []*Table    `json:"children"`
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

func (n *TableCell) AppendChild(child Node) error {
	if child.GetType() != TableType {
		return ErrInvalidChildType
	}
	n.Children = append(n.Children, child.(*Table))
	return nil
}

type TableFooter struct {
	Type     string      `json:"type"`
	Children []*Table    `json:"children"`
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

func (n *TableFooter) AppendChild(child Node) error {
	if child.GetType() != TableType {
		return ErrInvalidChildType
	}
	n.Children = append(n.Children, child.(*Table))
	return nil
}

type TableRow struct {
	Type     string       `json:"type"`
	Children []*TableCell `json:"children"`
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

func (n *TableRow) AppendChild(child Node) error {
	if child.GetType() != TableCellType {
		return ErrInvalidChildType
	}
	n.Children = append(n.Children, child.(*TableCell))
	return nil
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

func (n *Text) AppendChild(_ Node) error { return ErrCannotHaveChildren }

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

func (n *ThematicBreak) AppendChild(_ Node) error { return ErrCannotHaveChildren }

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

func (n *Tweet) AppendChild(child Node) error { return ErrCannotHaveChildren }

type Video struct {
	Type string      `json:"type"`
	Data interface{} `json:"data,omitempty"`
	ID   string      `json:"id"`
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

func (n *Video) AppendChild(_ Node) error { return ErrCannotHaveChildren }

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

func (n *YoutubeVideo) AppendChild(_ Node) error { return ErrCannotHaveChildren }

type CustomCodeComponent struct {
	Type                   string                 `json:"type"`
	Data                   interface{}            `json:"data,omitempty"`
	ID                     string                 `json:"id"`
	LayoutWidth            string                 `json:"layoutWidth"`
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

func (n *CustomCodeComponent) AppendChild(_ Node) error { return ErrCannotHaveChildren }

type ClipSet struct {
	Type        string      `json:"type"`
	Data        interface{} `json:"data,omitempty"`
	ID          string      `json:"id,omitempty"`
	LayoutWidth string      `json:"layoutWidth,omitempty"`
	Autoplay    bool        `json:"autoplay,omitempty"`
	Loop        bool        `json:"loop,omitempty"`
	Muted       bool        `json:"muted,omitempty"`
}

func (n *ClipSet) GetType() string {
	return n.Type
}

func (n *ClipSet) GetEmbedded() Node {
	return nil
}

func (n *ClipSet) GetChildren() []Node {
	return nil
}

func (n *ClipSet) AppendChild(_ Node) error { return ErrCannotHaveChildren }

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

func (n *Root) AppendChild(_ Node) error { return ErrCannotHaveChildren }

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

type Timeline struct {
	Type        string           `json:"type"`
	Title       string           `json:"title,omitempty"`
	LayoutWidth string           `json:"layoutWidth,omitempty"`
	Children    []*TimelineEvent `json:"children"`
}

func (n *Timeline) GetType() string {
	return n.Type
}

func (n *Timeline) GetEmbedded() Node {
	return nil
}

func (n *Timeline) GetChildren() []Node {
	result := make([]Node, len(n.Children))
	for i, v := range n.Children {
		result[i] = v
	}
	return result
}

func (n *Timeline) AppendChild(child Node) error {
	if child.GetType() != TimelineEventType {
		return ErrInvalidChildType
	}
	n.Children = append(n.Children, child.(*TimelineEvent))
	return nil
}

type TimelineEvent struct {
	Type     string                `json:"type"`
	Title    string                `json:"title,omitempty"`
	Children []*TimelineEventChild `json:"children"`
}

func (n *TimelineEvent) GetType() string {
	return n.Type
}

func (n *TimelineEvent) GetEmbedded() Node {
	return nil
}

func (n *TimelineEvent) GetChildren() []Node {
	result := make([]Node, len(n.Children))
	for i, v := range n.Children {
		result[i] = v
	}
	return result
}

func (n *TimelineEvent) AppendChild(child Node) error {
	c, err := makeTimelineEventChild(child)
	if err != nil {
		return err
	}
	n.Children = append(n.Children, c)
	return nil
}

type TimelineEventChild struct {
	*Paragraph
	*ImageSet
}

func (n *TimelineEventChild) GetType() string {
	return TimelineEventChildType
}

func (n *TimelineEventChild) GetEmbedded() Node {
	if n.Paragraph != nil {
		return n.Paragraph
	}
	if n.ImageSet != nil {
		return n.ImageSet
	}
	return nil
}

func (n *TimelineEventChild) GetChildren() []Node {
	if n.Paragraph != nil {
		return n.Paragraph.GetChildren()
	}
	if n.ImageSet != nil {
		return n.ImageSet.GetChildren()
	}
	return nil
}

func (n *TimelineEventChild) AppendChild(child Node) error { return ErrCannotHaveChildren }

func (n *TimelineEventChild) UnmarshalJSON(data []byte) error {
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
	case ImageSetType:
		var v ImageSet
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		n.ImageSet = &v
	default:
		return fmt.Errorf("failed to unmarshal TimelineEventChild from %s: %w", data, ErrUnmarshalInvalidNode)
	}
	return nil
}

func (n *TimelineEventChild) MarshalJSON() ([]byte, error) {
	switch {
	case n.Paragraph != nil:
		return json.Marshal(n.Paragraph)
	case n.ImageSet != nil:
		return json.Marshal(n.ImageSet)
	default:
		return []byte(`{}`), nil
	}
}

// Build a TimelineEventChild wrapper.
func makeTimelineEventChild(n Node) (*TimelineEventChild, error) {
	switch n.GetType() {
	case ParagraphType:
		return &TimelineEventChild{Paragraph: n.(*Paragraph)}, nil
	case ImageSetType:
		return &TimelineEventChild{ImageSet: n.(*ImageSet)}, nil
	default:
		return nil, ErrInvalidChildType
	}
}
