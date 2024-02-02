// Part of the code is generated using schema-generate (https://github.com/a-h/generate.).
// The schema used for the code generations is modified version of schemas/content-tree.schema.json.
// The modification applied was dropping blocks of the schema describing arrays of heterogeneous elements.
// The result Go structs doesn't capture fully the data definition in the JSON schemas used for its genteration.
// There no rules for enums, number of elements in array, etc.
package main

type ColumnSettingsItems struct {
	HideOnMobile bool   `json:"hideOnMobile,omitempty"`
	SortType     string `json:"sortType,omitempty"`
	Sortable     bool   `json:"sortable,omitempty"`
}

type ContentTreeFullBigNumber struct {
	Data        interface{} `json:"data,omitempty"`
	Description string      `json:"description,omitempty"`
	Number      string      `json:"number,omitempty"`
	Type        string      `json:"type,omitempty"`
}

type ContentTreeFullBlockquote struct {
	Children []*ContentTreeFullBlockquoteChild `json:"children,omitempty"`
	Data     interface{}                       `json:"data,omitempty"`
	Type     string                            `json:"type,omitempty"`
}

type ContentTreeFullBlockquoteChild struct {
	ContentTreeFullParagraph
	ContentTreeFullText
	ContentTreeFullBreak
	ContentTreeFullStrong
	ContentTreeFullEmphasis
	ContentTreeFullStrikethrough
	ContentTreeFullLink
}

type ContentTreeFullBody struct {
	Children []*ContentTreeFullBodyBlock `json:"children,omitempty"`
	Data     interface{}                 `json:"data,omitempty"`
	Type     string                      `json:"type,omitempty"`
	Version  float64                     `json:"version,omitempty"`
}

type ContentTreeFullBodyBlock struct {
	ContentTreeFullParagraph
	ContentTreeFullHeading
	ContentTreeFullImageSet
	ContentTreeFullBigNumber
	ContentTreeFullLayout
	ContentTreeFullList
	ContentTreeFullBlockquote
	ContentTreeFullPullquote
	ContentTreeFullScrollyBlock
	ContentTreeFullThematicBreak
	ContentTreeFullTable
	ContentTreeFullRecommended
	ContentTreeFullTweet
	ContentTreeFullVideo
	ContentTreeFullYoutubeVideo
}

type ContentTreeFullBreak struct {
	Data interface{} `json:"data,omitempty"`
	Type string      `json:"type,omitempty"`
}

type ContentTreeFullEmphasis struct {
	Children []*ContentTreeFullPhrasing `json:"children,omitempty"`
	Data     interface{}                `json:"data,omitempty"`
	Type     string                     `json:"type,omitempty"`
}

type ContentTreeFullHeading struct {
	Children []*ContentTreeFullText `json:"children,omitempty"`
	Data     interface{}            `json:"data,omitempty"`
	Level    string                 `json:"level,omitempty"`
	Type     string                 `json:"type,omitempty"`
}

type ContentTreeFullImageSet struct {
	Data    interface{} `json:"data,omitempty"`
	Id      string      `json:"id,omitempty"`
	Picture *Picture    `json:"picture,omitempty"`
	Type    string      `json:"type,omitempty"`
}

type ContentTreeFullLayout struct {
	Children    []*ContentTreeFullLayoutChild `json:"children,omitempty"`
	Data        interface{}                   `json:"data,omitempty"`
	LayoutName  string                        `json:"layoutName,omitempty"`
	LayoutWidth string                        `json:"layoutWidth,omitempty"`
	Type        string                        `json:"type,omitempty"`
}

// ContentTreeFullLayoutChild definition is significantly simplified form the definition of the type
// in the JSON schema.
type ContentTreeFullLayoutChild struct {
	ContentTreeFullLayoutSlot
	ContentTreeFullHeading
	ContentTreeFullLayoutImage
}

type ContentTreeFullLayoutImage struct {
	Alt     string      `json:"alt,omitempty"`
	Caption string      `json:"caption,omitempty"`
	Credit  string      `json:"credit,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Id      string      `json:"id,omitempty"`
	Picture *Picture    `json:"picture,omitempty"`
	Type    string      `json:"type,omitempty"`
}

type ContentTreeFullLayoutSlot struct {
	Children []*ContentTreeFullLayoutSlotChild `json:"children,omitempty"`
	Data     interface{}                       `json:"data,omitempty"`
	Type     string                            `json:"type,omitempty"`
}

type ContentTreeFullLayoutSlotChild struct {
	ContentTreeFullParagraph
	ContentTreeFullHeading
	ContentTreeFullLayoutImage
}

type ContentTreeFullLink struct {
	Children []*ContentTreeFullPhrasing `json:"children,omitempty"`
	Data     interface{}                `json:"data,omitempty"`
	Title    string                     `json:"title,omitempty"`
	Type     string                     `json:"type,omitempty"`
	Url      string                     `json:"url,omitempty"`
}

type ContentTreeFullList struct {
	Children []*ContentTreeFullListItem `json:"children,omitempty"`
	Data     interface{}                `json:"data,omitempty"`
	Ordered  bool                       `json:"ordered,omitempty"`
	Type     string                     `json:"type,omitempty"`
}

type ContentTreeFullListItem struct {
	Children []*ContentTreeFullListItemChild `json:"children,omitempty"`
	Data     interface{}                     `json:"data,omitempty"`
	Type     string                          `json:"type,omitempty"`
}

type ContentTreeFullListItemChild struct {
	ContentTreeFullParagraph
	ContentTreeFullText
	ContentTreeFullBreak
	ContentTreeFullStrong
	ContentTreeFullEmphasis
	ContentTreeFullStrikethrough
	ContentTreeFullLink
}

type ContentTreeFullParagraph struct {
	Children []*ContentTreeFullPhrasing `json:"children,omitempty"`
	Data     interface{}                `json:"data,omitempty"`
	Type     string                     `json:"type,omitempty"`
}

type ContentTreeFullPhrasing struct {
	ContentTreeFullText
	ContentTreeFullBreak
	ContentTreeFullStrong
	ContentTreeFullEmphasis
	ContentTreeFullStrikethrough
	ContentTreeFullLink
}

type ContentTreeFullPullquote struct {
	Data   interface{} `json:"data,omitempty"`
	Source string      `json:"source,omitempty"`
	Text   string      `json:"text,omitempty"`
	Type   string      `json:"type,omitempty"`
}

type ContentTreeFullRecommended struct {
	Data                interface{} `json:"data,omitempty"`
	Heading             string      `json:"heading,omitempty"`
	Id                  string      `json:"id,omitempty"`
	Teaser              *Teaser     `json:"teaser,omitempty"`
	TeaserTitleOverride string      `json:"teaserTitleOverride,omitempty"`
	Type                string      `json:"type,omitempty"`
}

type ContentTreeFullScrollyBlock struct {
	Children []*ContentTreeFullScrollySection `json:"children,omitempty"`
	Data     interface{}                      `json:"data,omitempty"`
	Theme    string                           `json:"theme,omitempty"`
	Type     string                           `json:"type,omitempty"`
}

type ContentTreeFullScrollyCopy struct {
	Children []*ContentTreeFullScrollyCopyChild `json:"children,omitempty"`
	Data     interface{}                        `json:"data,omitempty"`
	Type     string                             `json:"type,omitempty"`
}

type ContentTreeFullScrollyCopyChild struct {
	ContentTreeFullParagraph
	ContentTreeFullScrollyHeading
}

type ContentTreeFullScrollyHeading struct {
	Children []*ContentTreeFullText `json:"children,omitempty"`
	Data     interface{}            `json:"data,omitempty"`
	Level    string                 `json:"level,omitempty"`
	Type     string                 `json:"type,omitempty"`
}

type ContentTreeFullScrollyImage struct {
	Data    interface{} `json:"data,omitempty"`
	Id      string      `json:"id,omitempty"`
	Picture *Picture    `json:"picture,omitempty"`
	Type    string      `json:"type,omitempty"`
}

type ContentTreeFullScrollySection struct {
	Children   []*ContentTreeFullScrollySectionChild `json:"children,omitempty"`
	Data       interface{}                           `json:"data,omitempty"`
	Display    string                                `json:"display,omitempty"`
	NoBox      bool                                  `json:"noBox,omitempty"`
	Position   string                                `json:"position,omitempty"`
	Transition string                                `json:"transition,omitempty"`
	Type       string                                `json:"type,omitempty"`
}

// ContentTreeFullScrollySectionChild is example of a type much simpler than the definition in the JSON schema.
type ContentTreeFullScrollySectionChild struct {
	ContentTreeFullScrollyCopy
	ContentTreeFullScrollyImage
}

type ContentTreeFullStrikethrough struct {
	Children []*ContentTreeFullPhrasing `json:"children,omitempty"`
	Data     interface{}                `json:"data,omitempty"`
	Type     string                     `json:"type,omitempty"`
}

type ContentTreeFullStrong struct {
	Children []*ContentTreeFullPhrasing `json:"children,omitempty"`
	Data     interface{}                `json:"data,omitempty"`
	Type     string                     `json:"type,omitempty"`
}

type ContentTreeFullTable struct {
	Children                 []*ContentTreeFullTableChild `json:"children,omitempty"`
	CollapseAfterHowManyRows float64                      `json:"collapseAfterHowManyRows,omitempty"`
	ColumnSettings           []*ColumnSettingsItems       `json:"columnSettings,omitempty"`
	Compact                  bool                         `json:"compact,omitempty"`
	Data                     interface{}                  `json:"data,omitempty"`
	LayoutWidth              string                       `json:"layoutWidth,omitempty"`
	ResponsiveStyle          string                       `json:"responsiveStyle,omitempty"`
	Stripes                  bool                         `json:"stripes,omitempty"`
	Type                     string                       `json:"type,omitempty"`
}

type ContentTreeFullTableChild struct {
	ContentTreeFullTableCaption
	ContentTreeFullTableBody
	ContentTreeFullTableFooter
}

type ContentTreeFullTableBody struct {
	Children []*ContentTreeFullTableRow `json:"children,omitempty"`
	Data     interface{}                `json:"data,omitempty"`
	Type     string                     `json:"type,omitempty"`
}

type ContentTreeFullTableCaption struct {
	Children []*ContentTreeFullTable `json:"children,omitempty"`
	Data     interface{}             `json:"data,omitempty"`
	Type     string                  `json:"type,omitempty"`
}

type ContentTreeFullTableCell struct {
	Children []*ContentTreeFullTable `json:"children,omitempty"`
	Data     interface{}             `json:"data,omitempty"`
	Heading  bool                    `json:"heading,omitempty"`
	Type     string                  `json:"type,omitempty"`
}

type ContentTreeFullTableFooter struct {
	Children []*ContentTreeFullTable `json:"children,omitempty"`
	Data     interface{}             `json:"data,omitempty"`
	Type     string                  `json:"type,omitempty"`
}

type ContentTreeFullTableRow struct {
	Children []*ContentTreeFullTableCell `json:"children,omitempty"`
	Data     interface{}                 `json:"data,omitempty"`
	Type     string                      `json:"type,omitempty"`
}

type ContentTreeFullText struct {
	Data  interface{} `json:"data,omitempty"`
	Type  string      `json:"type,omitempty"`
	Value string      `json:"value,omitempty"`
}

type ContentTreeFullThematicBreak struct {
	Data interface{} `json:"data,omitempty"`
	Type string      `json:"type,omitempty"`
}

type ContentTreeFullTweet struct {
	Data interface{} `json:"data,omitempty"`
	Html string      `json:"html,omitempty"`
	Id   string      `json:"id,omitempty"`
	Type string      `json:"type,omitempty"`
}

type ContentTreeFullVideo struct {
	Data     interface{} `json:"data,omitempty"`
	Embedded bool        `json:"embedded,omitempty"`
	Id       string      `json:"id,omitempty"`
	Type     string      `json:"type,omitempty"`
}

type ContentTreeFullYoutubeVideo struct {
	Data interface{} `json:"data,omitempty"`
	Type string      `json:"type,omitempty"`
	Url  string      `json:"url,omitempty"`
}

type FallbackImage struct {
	Format    string            `json:"format,omitempty"`
	Height    float64           `json:"height,omitempty"`
	Id        string            `json:"id,omitempty"`
	SourceSet []*SourceSetItems `json:"sourceSet,omitempty"`
	Url       string            `json:"url,omitempty"`
	Width     float64           `json:"width,omitempty"`
}

type Image struct {
	Format    string            `json:"format,omitempty"`
	Height    float64           `json:"height,omitempty"`
	Id        string            `json:"id,omitempty"`
	SourceSet []*SourceSetItems `json:"sourceSet,omitempty"`
	Url       string            `json:"url,omitempty"`
	Width     float64           `json:"width,omitempty"`
}

type ImagesItems struct {
	Format    string            `json:"format,omitempty"`
	Height    float64           `json:"height,omitempty"`
	Id        string            `json:"id,omitempty"`
	SourceSet []*SourceSetItems `json:"sourceSet,omitempty"`
	Url       string            `json:"url,omitempty"`
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
	ApiUrl     string   `json:"apiUrl,omitempty"`
	DirectType string   `json:"directType,omitempty"`
	Id         string   `json:"id,omitempty"`
	Predicate  string   `json:"predicate,omitempty"`
	PrefLabel  string   `json:"prefLabel,omitempty"`
	Type       string   `json:"type,omitempty"`
	Types      []string `json:"types,omitempty"`
	Url        string   `json:"url,omitempty"`
}

type MetaLink struct {
	ApiUrl     string   `json:"apiUrl,omitempty"`
	DirectType string   `json:"directType,omitempty"`
	Id         string   `json:"id,omitempty"`
	Predicate  string   `json:"predicate,omitempty"`
	PrefLabel  string   `json:"prefLabel,omitempty"`
	Type       string   `json:"type,omitempty"`
	Types      []string `json:"types,omitempty"`
	Url        string   `json:"url,omitempty"`
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
	Body *ContentTreeFullBody `json:"body,omitempty"`
	Data interface{}          `json:"data,omitempty"`
	Type string               `json:"type,omitempty"`
}

type SourceSetItems struct {
	Dpr   float64 `json:"dpr,omitempty"`
	Url   string  `json:"url,omitempty"`
	Width float64 `json:"width,omitempty"`
}

type Teaser struct {
	FirstPublishedDate string       `json:"firstPublishedDate,omitempty"`
	Id                 string       `json:"id,omitempty"`
	Image              *Image       `json:"image,omitempty"`
	Indicators         *Indicators  `json:"indicators,omitempty"`
	MetaAltLink        *MetaAltLink `json:"metaAltLink,omitempty"`
	MetaLink           *MetaLink    `json:"metaLink,omitempty"`
	MetaPrefixText     string       `json:"metaPrefixText,omitempty"`
	MetaSuffixText     string       `json:"metaSuffixText,omitempty"`
	PublishedDate      string       `json:"publishedDate,omitempty"`
	Title              string       `json:"title,omitempty"`
	Type               string       `json:"type,omitempty"`
	Url                string       `json:"url,omitempty"`
}

func main() {}
