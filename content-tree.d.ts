export declare namespace ContentTree {
    type FormattingBlock = Paragraph | Heading | List | Blockquote | ThematicBreak | Text;
    type StoryBlock = ImageSet | Flourish | BigNumber | CustomCodeComponent | Layout | Pullquote | ScrollyBlock | Table | Recommended | RecommendedList | Tweet | Video | YoutubeVideo | Timeline | ImagePair | InNumbers | Definition | InfoBox | InfoPair;
    type BodyBlock = FormattingBlock | StoryBlock;
    type LayoutWidth = "auto" | "in-line" | "inset-left" | "inset-right" | "full-bleed" | "full-grid" | "mid-grid" | "full-width";
    type Phrasing = Text | Break | Strong | Emphasis | Strikethrough | Link;
    interface Node {
        type: string;
        data?: any;
    }
    interface Parent extends Node {
        children: Node[];
    }
    interface Root extends Node {
        type: "root";
        body: Body;
    }
    interface Body extends Parent {
        type: "body";
        version: number;
        children: BodyBlock[];
    }
    interface Text extends Node {
        type: "text";
        value: string;
    }
    interface Break extends Node {
        type: "break";
    }
    interface ThematicBreak extends Node {
        type: "thematic-break";
    }
    interface Paragraph extends Parent {
        type: "paragraph";
        children: Phrasing[];
    }
    interface Heading extends Parent {
        type: "heading";
        children: Text[];
        level: "chapter" | "subheading" | "label";
        fragmentIdentifier?: string;
    }
    interface Strong extends Parent {
        type: "strong";
        children: Phrasing[];
    }
    interface Emphasis extends Parent {
        type: "emphasis";
        children: Phrasing[];
    }
    interface Strikethrough extends Parent {
        type: "strikethrough";
        children: Phrasing[];
    }
    interface Link extends Parent {
        type: "link";
        url: string;
        title: string;
        children: Phrasing[];
        styleType?: 'onward-journey';
    }
    interface List extends Parent {
        type: "list";
        ordered: boolean;
        children: ListItem[];
    }
    interface ListItem extends Parent {
        type: "list-item";
        children: (Paragraph | Phrasing)[];
    }
    interface Blockquote extends Parent {
        type: "blockquote";
        children: (Paragraph | Phrasing)[];
    }
    interface Pullquote extends Node {
        type: "pullquote";
        text: string;
        source?: string;
    }
    interface ImageSet extends Node {
        type: "image-set";
        id: string;
        picture: ImageSetPicture;
        fragmentIdentifier?: string;
    }
    type ImageSetPicture = {
        layoutWidth: string;
        imageType: "image" | "graphic";
        alt: string;
        caption: string;
        credit: string;
        images: Image[];
        fallbackImage: Image;
    };
    type Image = {
        id: string;
        width: number;
        height: number;
        format: "desktop" | "mobile" | "square" | "square-ftedit" | "standard" | "wide" | "standard-inline";
        url: string;
        sourceSet?: ImageSource[];
    };
    type ImageSource = {
        url: string;
        width: number;
        dpr: number;
    };
    interface Recommended extends Node {
        type: "recommended";
        id: string;
        heading?: string;
        teaserTitleOverride?: string;
        teaser: Teaser;
    }
    interface RecommendedList extends Node {
        type: "recommended-list";
        heading?: string;
        children: Recommended[];
    }
    type TeaserConcept = {
        apiUrl: string;
        directType: string;
        id: string;
        predicate: string;
        prefLabel: string;
        type: string;
        types: string[];
        url: string;
    };
    type Teaser = {
        id: string;
        url: string;
        type: "article" | "video" | "podcast" | "audio" | "package" | "liveblog" | "promoted-content" | "paid-post";
        title: string;
        publishedDate: string;
        firstPublishedDate: string;
        metaLink?: TeaserConcept;
        metaAltLink?: TeaserConcept;
        metaPrefixText?: string;
        metaSuffixText?: string;
        indicators: {
            accessLevel: "premium" | "subscribed" | "registered" | "free";
            isOpinion?: boolean;
            isColumn?: boolean;
            isPodcast?: boolean;
            isEditorsChoice?: boolean;
            isExclusive?: boolean;
            isScoop?: boolean;
        };
        image: {
            url: string;
            width: number;
            height: number;
        };
        clientName?: string;
    };
    interface Tweet extends Node {
        id: string;
        type: "tweet";
        html: string;
    }
    type FlourishLayoutWidth = Extract<LayoutWidth, "full-grid" | "in-line">;
    interface Flourish extends Node {
        type: "flourish";
        id: string;
        layoutWidth: FlourishLayoutWidth;
        flourishType: string;
        description?: string;
        timestamp?: string;
        fallbackImage?: Image;
        fragmentIdentifier?: string;
    }
    interface BigNumber extends Node {
        type: "big-number";
        number: string;
        description: string;
    }
    interface Video extends Node {
        type: "video";
        id: string;
        title: string;
    }
    interface YoutubeVideo extends Node {
        type: "youtube-video";
        url: string;
    }
    interface ScrollyBlock extends Parent {
        type: "scrolly-block";
        theme: "sans" | "serif";
        children: ScrollySection[];
    }
    interface ScrollySection extends Parent {
        type: "scrolly-section";
        display: "dark-background" | "light-background";
        noBox?: true;
        position: "left" | "center" | "right";
        transition?: "delay-before" | "delay-after";
        children: [ScrollyImage, ...ScrollyCopy[]];
    }
    interface ScrollyImage extends Node {
        type: "scrolly-image";
        id: string;
        picture: ImageSetPicture;
    }
    interface ScrollyCopy extends Parent {
        type: "scrolly-copy";
        children: (ScrollyHeading | Paragraph)[];
    }
    interface ScrollyHeading extends Parent {
        type: "scrolly-heading";
        level: "chapter" | "heading" | "subheading";
        children: Text[];
    }
    interface Layout extends Parent {
        type: "layout";
        layoutName: "auto" | "card" | "timeline";
        layoutWidth: string;
        children: [Heading, LayoutImage, ...LayoutSlot[]] | [Heading, ...LayoutSlot[]] | LayoutSlot[];
    }
    interface LayoutSlot extends Parent {
        type: "layout-slot";
        children: (Heading | Paragraph | LayoutImage)[];
    }
    interface LayoutImage extends Node {
        type: "layout-image";
        id: string;
        alt: string;
        caption: string;
        credit: string;
        picture: ImageSetPicture;
    }
    type TableColumnSettings = {
        hideOnMobile: boolean;
        sortable: boolean;
        sortType: 'text' | 'number' | 'date' | 'currency' | 'percent';
    };
    type TableLayoutWidth = Extract<LayoutWidth, 'auto' | 'full-grid' | 'inset-left' | 'inset-right' | 'full-bleed'>;
    interface TableCaption extends Parent {
        type: 'table-caption';
        children: Phrasing[];
    }
    interface TableCell extends Parent {
        type: 'table-cell';
        heading?: boolean;
        columnSpan?: number;
        rowSpan?: number;
        children: Phrasing[];
    }
    interface TableRow extends Parent {
        type: 'table-row';
        children: TableCell[];
    }
    interface TableBody extends Parent {
        type: 'table-body';
        children: TableRow[];
    }
    interface TableFooter extends Parent {
        type: 'table-footer';
        children: Phrasing[];
    }
    interface Table extends Parent {
        type: 'table';
        stripes: boolean;
        compact: boolean;
        layoutWidth: TableLayoutWidth;
        collapseAfterHowManyRows?: number;
        responsiveStyle: 'overflow' | 'flat' | 'scroll';
        children: [TableCaption, TableBody, TableFooter] | [TableCaption, TableBody] | [TableBody, TableFooter] | [TableBody];
        columnSettings: TableColumnSettings[];
    }
    type CustomCodeComponentAttributes = {
        [key: string]: string | boolean | undefined;
    };
    interface CustomCodeComponent extends Node {
        /** Component type */
        type: "custom-code-component";
        /** Id taken from the CAPI url */
        id: string;
        /** How the component should be presented in the article page according to the column layout system */
        layoutWidth: LayoutWidth;
        /** Repository for the code of the component in the format "[github org]/[github repo]/[component name]". */
        path: string;
        /** Semantic version of the code of the component, e.g. "^0.3.5". */
        versionRange: string;
        /** Last date-time when the attributes for this block were modified, in ISO-8601 format. */
        attributesLastModified: string;
        /** Configuration data to be passed to the component. */
        attributes: CustomCodeComponentAttributes;
    }
    interface ImagePair extends Parent {
        type: 'image-pair';
        children: [ImageSet, ImageSet];
    }
    /**
     * Timeline nodes display a timeline of events in arbitrary order.
     */
    interface Timeline extends Parent {
        type: "timeline";
        /** The title for the timeline */
        title: string;
        children: TimelineEvent[];
    }
    /**
     * TimelineEvent is the representation of a single event in a Timeline.
     */
    interface TimelineEvent extends Parent {
        type: "timeline-event";
        /** The title of the event */
        title: string;
        /** Any combination of paragraphs and image sets */
        children: (Paragraph | ImageSet)[];
    }
    /**
     * A definition has a term and a related description. It is used to describe a term.
     */
    interface Definition extends Node {
        type: "definition";
        term: string;
        description: string;
    }
    /**
     * InNumbers represents a set of numbers with related descriptions.
     */
    interface InNumbers extends Parent {
        type: "in-numbers";
        /** The title for the InNumbers */
        title?: string;
        children: [Definition, Definition, Definition];
    }
    /** Allowed children for a card
    */
    type CardChildren = ImageSet | Exclude<FormattingBlock, Heading>;
    /**
    * A card describes a subject with images and text
    */
    interface Card extends Parent {
        type: "card";
        /** The title of this card */
        title?: string;
        children: CardChildren[];
    }
    /**
    * Allowed layout widths for an InfoBox.
    */
    type InfoBoxLayoutWidth = Extract<LayoutWidth, "full-width" | "inset-left">;
    /**
    * An info box describes a subject via a single card
    */
    interface InfoBox extends Parent {
        type: "info-box";
        /** The layout width supported by this node */
        layoutWidth: InfoBoxLayoutWidth;
        children: [Card];
    }
    /**
    * InfoPair provides exactly two cards.
    */
    interface InfoPair extends Parent {
        type: "info-pair";
        /** The title of the info pair */
        title?: string;
        children: [Card, Card];
    }
    namespace full {
        type FormattingBlock = Paragraph | Heading | List | Blockquote | ThematicBreak | Text;
        type StoryBlock = ImageSet | Flourish | BigNumber | CustomCodeComponent | Layout | Pullquote | ScrollyBlock | Table | Recommended | RecommendedList | Tweet | Video | YoutubeVideo | Timeline | ImagePair | InNumbers | Definition | InfoBox | InfoPair;
        type BodyBlock = FormattingBlock | StoryBlock;
        type LayoutWidth = "auto" | "in-line" | "inset-left" | "inset-right" | "full-bleed" | "full-grid" | "mid-grid" | "full-width";
        type Phrasing = Text | Break | Strong | Emphasis | Strikethrough | Link;
        interface Node {
            type: string;
            data?: any;
        }
        interface Parent extends Node {
            children: Node[];
        }
        interface Root extends Node {
            type: "root";
            body: Body;
        }
        interface Body extends Parent {
            type: "body";
            version: number;
            children: BodyBlock[];
        }
        interface Text extends Node {
            type: "text";
            value: string;
        }
        interface Break extends Node {
            type: "break";
        }
        interface ThematicBreak extends Node {
            type: "thematic-break";
        }
        interface Paragraph extends Parent {
            type: "paragraph";
            children: Phrasing[];
        }
        interface Heading extends Parent {
            type: "heading";
            children: Text[];
            level: "chapter" | "subheading" | "label";
            fragmentIdentifier?: string;
        }
        interface Strong extends Parent {
            type: "strong";
            children: Phrasing[];
        }
        interface Emphasis extends Parent {
            type: "emphasis";
            children: Phrasing[];
        }
        interface Strikethrough extends Parent {
            type: "strikethrough";
            children: Phrasing[];
        }
        interface Link extends Parent {
            type: "link";
            url: string;
            title: string;
            children: Phrasing[];
            styleType?: 'onward-journey';
        }
        interface List extends Parent {
            type: "list";
            ordered: boolean;
            children: ListItem[];
        }
        interface ListItem extends Parent {
            type: "list-item";
            children: (Paragraph | Phrasing)[];
        }
        interface Blockquote extends Parent {
            type: "blockquote";
            children: (Paragraph | Phrasing)[];
        }
        interface Pullquote extends Node {
            type: "pullquote";
            text: string;
            source?: string;
        }
        interface ImageSet extends Node {
            type: "image-set";
            id: string;
            picture: ImageSetPicture;
            fragmentIdentifier?: string;
        }
        type ImageSetPicture = {
            layoutWidth: string;
            imageType: "image" | "graphic";
            alt: string;
            caption: string;
            credit: string;
            images: Image[];
            fallbackImage: Image;
        };
        type Image = {
            id: string;
            width: number;
            height: number;
            format: "desktop" | "mobile" | "square" | "square-ftedit" | "standard" | "wide" | "standard-inline";
            url: string;
            sourceSet?: ImageSource[];
        };
        type ImageSource = {
            url: string;
            width: number;
            dpr: number;
        };
        interface Recommended extends Node {
            type: "recommended";
            id: string;
            heading?: string;
            teaserTitleOverride?: string;
            teaser: Teaser;
        }
        interface RecommendedList extends Node {
            type: "recommended-list";
            heading?: string;
            children: Recommended[];
        }
        type TeaserConcept = {
            apiUrl: string;
            directType: string;
            id: string;
            predicate: string;
            prefLabel: string;
            type: string;
            types: string[];
            url: string;
        };
        type Teaser = {
            id: string;
            url: string;
            type: "article" | "video" | "podcast" | "audio" | "package" | "liveblog" | "promoted-content" | "paid-post";
            title: string;
            publishedDate: string;
            firstPublishedDate: string;
            metaLink?: TeaserConcept;
            metaAltLink?: TeaserConcept;
            metaPrefixText?: string;
            metaSuffixText?: string;
            indicators: {
                accessLevel: "premium" | "subscribed" | "registered" | "free";
                isOpinion?: boolean;
                isColumn?: boolean;
                isPodcast?: boolean;
                isEditorsChoice?: boolean;
                isExclusive?: boolean;
                isScoop?: boolean;
            };
            image: {
                url: string;
                width: number;
                height: number;
            };
            clientName?: string;
        };
        interface Tweet extends Node {
            id: string;
            type: "tweet";
            html: string;
        }
        type FlourishLayoutWidth = Extract<LayoutWidth, "full-grid" | "in-line">;
        interface Flourish extends Node {
            type: "flourish";
            id: string;
            layoutWidth: FlourishLayoutWidth;
            flourishType: string;
            description?: string;
            timestamp?: string;
            fallbackImage?: Image;
            fragmentIdentifier?: string;
        }
        interface BigNumber extends Node {
            type: "big-number";
            number: string;
            description: string;
        }
        interface Video extends Node {
            type: "video";
            id: string;
            title: string;
        }
        interface YoutubeVideo extends Node {
            type: "youtube-video";
            url: string;
        }
        interface ScrollyBlock extends Parent {
            type: "scrolly-block";
            theme: "sans" | "serif";
            children: ScrollySection[];
        }
        interface ScrollySection extends Parent {
            type: "scrolly-section";
            display: "dark-background" | "light-background";
            noBox?: true;
            position: "left" | "center" | "right";
            transition?: "delay-before" | "delay-after";
            children: [ScrollyImage, ...ScrollyCopy[]];
        }
        interface ScrollyImage extends Node {
            type: "scrolly-image";
            id: string;
            picture: ImageSetPicture;
        }
        interface ScrollyCopy extends Parent {
            type: "scrolly-copy";
            children: (ScrollyHeading | Paragraph)[];
        }
        interface ScrollyHeading extends Parent {
            type: "scrolly-heading";
            level: "chapter" | "heading" | "subheading";
            children: Text[];
        }
        interface Layout extends Parent {
            type: "layout";
            layoutName: "auto" | "card" | "timeline";
            layoutWidth: string;
            children: [Heading, LayoutImage, ...LayoutSlot[]] | [Heading, ...LayoutSlot[]] | LayoutSlot[];
        }
        interface LayoutSlot extends Parent {
            type: "layout-slot";
            children: (Heading | Paragraph | LayoutImage)[];
        }
        interface LayoutImage extends Node {
            type: "layout-image";
            id: string;
            alt: string;
            caption: string;
            credit: string;
            picture: ImageSetPicture;
        }
        type TableColumnSettings = {
            hideOnMobile: boolean;
            sortable: boolean;
            sortType: 'text' | 'number' | 'date' | 'currency' | 'percent';
        };
        type TableLayoutWidth = Extract<LayoutWidth, 'auto' | 'full-grid' | 'inset-left' | 'inset-right' | 'full-bleed'>;
        interface TableCaption extends Parent {
            type: 'table-caption';
            children: Phrasing[];
        }
        interface TableCell extends Parent {
            type: 'table-cell';
            heading?: boolean;
            columnSpan?: number;
            rowSpan?: number;
            children: Phrasing[];
        }
        interface TableRow extends Parent {
            type: 'table-row';
            children: TableCell[];
        }
        interface TableBody extends Parent {
            type: 'table-body';
            children: TableRow[];
        }
        interface TableFooter extends Parent {
            type: 'table-footer';
            children: Phrasing[];
        }
        interface Table extends Parent {
            type: 'table';
            stripes: boolean;
            compact: boolean;
            layoutWidth: TableLayoutWidth;
            collapseAfterHowManyRows?: number;
            responsiveStyle: 'overflow' | 'flat' | 'scroll';
            children: [TableCaption, TableBody, TableFooter] | [TableCaption, TableBody] | [TableBody, TableFooter] | [TableBody];
            columnSettings: TableColumnSettings[];
        }
        type CustomCodeComponentAttributes = {
            [key: string]: string | boolean | undefined;
        };
        interface CustomCodeComponent extends Node {
            /** Component type */
            type: "custom-code-component";
            /** Id taken from the CAPI url */
            id: string;
            /** How the component should be presented in the article page according to the column layout system */
            layoutWidth: LayoutWidth;
            /** Repository for the code of the component in the format "[github org]/[github repo]/[component name]". */
            path: string;
            /** Semantic version of the code of the component, e.g. "^0.3.5". */
            versionRange: string;
            /** Last date-time when the attributes for this block were modified, in ISO-8601 format. */
            attributesLastModified: string;
            /** Configuration data to be passed to the component. */
            attributes: CustomCodeComponentAttributes;
        }
        interface ImagePair extends Parent {
            type: 'image-pair';
            children: [ImageSet, ImageSet];
        }
        /**
         * Timeline nodes display a timeline of events in arbitrary order.
         */
        interface Timeline extends Parent {
            type: "timeline";
            /** The title for the timeline */
            title: string;
            children: TimelineEvent[];
        }
        /**
         * TimelineEvent is the representation of a single event in a Timeline.
         */
        interface TimelineEvent extends Parent {
            type: "timeline-event";
            /** The title of the event */
            title: string;
            /** Any combination of paragraphs and image sets */
            children: (Paragraph | ImageSet)[];
        }
        /**
         * A definition has a term and a related description. It is used to describe a term.
         */
        interface Definition extends Node {
            type: "definition";
            term: string;
            description: string;
        }
        /**
         * InNumbers represents a set of numbers with related descriptions.
         */
        interface InNumbers extends Parent {
            type: "in-numbers";
            /** The title for the InNumbers */
            title?: string;
            children: [Definition, Definition, Definition];
        }
        /** Allowed children for a card
        */
        type CardChildren = ImageSet | Exclude<FormattingBlock, Heading>;
        /**
        * A card describes a subject with images and text
        */
        interface Card extends Parent {
            type: "card";
            /** The title of this card */
            title?: string;
            children: CardChildren[];
        }
        /**
        * Allowed layout widths for an InfoBox.
        */
        type InfoBoxLayoutWidth = Extract<LayoutWidth, "full-width" | "inset-left">;
        /**
        * An info box describes a subject via a single card
        */
        interface InfoBox extends Parent {
            type: "info-box";
            /** The layout width supported by this node */
            layoutWidth: InfoBoxLayoutWidth;
            children: [Card];
        }
        /**
        * InfoPair provides exactly two cards.
        */
        interface InfoPair extends Parent {
            type: "info-pair";
            /** The title of the info pair */
            title?: string;
            children: [Card, Card];
        }
    }
    namespace transit {
        type FormattingBlock = Paragraph | Heading | List | Blockquote | ThematicBreak | Text;
        type StoryBlock = ImageSet | Flourish | BigNumber | CustomCodeComponent | Layout | Pullquote | ScrollyBlock | Table | Recommended | RecommendedList | Tweet | Video | YoutubeVideo | Timeline | ImagePair | InNumbers | Definition | InfoBox | InfoPair;
        type BodyBlock = FormattingBlock | StoryBlock;
        type LayoutWidth = "auto" | "in-line" | "inset-left" | "inset-right" | "full-bleed" | "full-grid" | "mid-grid" | "full-width";
        type Phrasing = Text | Break | Strong | Emphasis | Strikethrough | Link;
        interface Node {
            type: string;
            data?: any;
        }
        interface Parent extends Node {
            children: Node[];
        }
        interface Root extends Node {
            type: "root";
            body: Body;
        }
        interface Body extends Parent {
            type: "body";
            version: number;
            children: BodyBlock[];
        }
        interface Text extends Node {
            type: "text";
            value: string;
        }
        interface Break extends Node {
            type: "break";
        }
        interface ThematicBreak extends Node {
            type: "thematic-break";
        }
        interface Paragraph extends Parent {
            type: "paragraph";
            children: Phrasing[];
        }
        interface Heading extends Parent {
            type: "heading";
            children: Text[];
            level: "chapter" | "subheading" | "label";
            fragmentIdentifier?: string;
        }
        interface Strong extends Parent {
            type: "strong";
            children: Phrasing[];
        }
        interface Emphasis extends Parent {
            type: "emphasis";
            children: Phrasing[];
        }
        interface Strikethrough extends Parent {
            type: "strikethrough";
            children: Phrasing[];
        }
        interface Link extends Parent {
            type: "link";
            url: string;
            title: string;
            children: Phrasing[];
            styleType?: 'onward-journey';
        }
        interface List extends Parent {
            type: "list";
            ordered: boolean;
            children: ListItem[];
        }
        interface ListItem extends Parent {
            type: "list-item";
            children: (Paragraph | Phrasing)[];
        }
        interface Blockquote extends Parent {
            type: "blockquote";
            children: (Paragraph | Phrasing)[];
        }
        interface Pullquote extends Node {
            type: "pullquote";
            text: string;
            source?: string;
        }
        interface ImageSet extends Node {
            type: "image-set";
            id: string;
            fragmentIdentifier?: string;
        }
        type ImageSetPicture = {
            layoutWidth: string;
            imageType: "image" | "graphic";
            alt: string;
            caption: string;
            credit: string;
            images: Image[];
            fallbackImage: Image;
        };
        type Image = {
            id: string;
            width: number;
            height: number;
            format: "desktop" | "mobile" | "square" | "square-ftedit" | "standard" | "wide" | "standard-inline";
            url: string;
            sourceSet?: ImageSource[];
        };
        type ImageSource = {
            url: string;
            width: number;
            dpr: number;
        };
        interface Recommended extends Node {
            type: "recommended";
            id: string;
            heading?: string;
            teaserTitleOverride?: string;
        }
        interface RecommendedList extends Node {
            type: "recommended-list";
            heading?: string;
            children: Recommended[];
        }
        type TeaserConcept = {
            apiUrl: string;
            directType: string;
            id: string;
            predicate: string;
            prefLabel: string;
            type: string;
            types: string[];
            url: string;
        };
        type Teaser = {
            id: string;
            url: string;
            type: "article" | "video" | "podcast" | "audio" | "package" | "liveblog" | "promoted-content" | "paid-post";
            title: string;
            publishedDate: string;
            firstPublishedDate: string;
            metaLink?: TeaserConcept;
            metaAltLink?: TeaserConcept;
            metaPrefixText?: string;
            metaSuffixText?: string;
            indicators: {
                accessLevel: "premium" | "subscribed" | "registered" | "free";
                isOpinion?: boolean;
                isColumn?: boolean;
                isPodcast?: boolean;
                isEditorsChoice?: boolean;
                isExclusive?: boolean;
                isScoop?: boolean;
            };
            image: {
                url: string;
                width: number;
                height: number;
            };
            clientName?: string;
        };
        interface Tweet extends Node {
            id: string;
            type: "tweet";
        }
        type FlourishLayoutWidth = Extract<LayoutWidth, "full-grid" | "in-line">;
        interface Flourish extends Node {
            type: "flourish";
            id: string;
            layoutWidth: FlourishLayoutWidth;
            flourishType: string;
            description?: string;
            timestamp?: string;
            fragmentIdentifier?: string;
        }
        interface BigNumber extends Node {
            type: "big-number";
            number: string;
            description: string;
        }
        interface Video extends Node {
            type: "video";
            id: string;
        }
        interface YoutubeVideo extends Node {
            type: "youtube-video";
            url: string;
        }
        interface ScrollyBlock extends Parent {
            type: "scrolly-block";
            theme: "sans" | "serif";
            children: ScrollySection[];
        }
        interface ScrollySection extends Parent {
            type: "scrolly-section";
            display: "dark-background" | "light-background";
            noBox?: true;
            position: "left" | "center" | "right";
            transition?: "delay-before" | "delay-after";
            children: [ScrollyImage, ...ScrollyCopy[]];
        }
        interface ScrollyImage extends Node {
            type: "scrolly-image";
            id: string;
        }
        interface ScrollyCopy extends Parent {
            type: "scrolly-copy";
            children: (ScrollyHeading | Paragraph)[];
        }
        interface ScrollyHeading extends Parent {
            type: "scrolly-heading";
            level: "chapter" | "heading" | "subheading";
            children: Text[];
        }
        interface Layout extends Parent {
            type: "layout";
            layoutName: "auto" | "card" | "timeline";
            layoutWidth: string;
            children: [Heading, LayoutImage, ...LayoutSlot[]] | [Heading, ...LayoutSlot[]] | LayoutSlot[];
        }
        interface LayoutSlot extends Parent {
            type: "layout-slot";
            children: (Heading | Paragraph | LayoutImage)[];
        }
        interface LayoutImage extends Node {
            type: "layout-image";
            id: string;
            alt: string;
            caption: string;
            credit: string;
        }
        type TableColumnSettings = {
            hideOnMobile: boolean;
            sortable: boolean;
            sortType: 'text' | 'number' | 'date' | 'currency' | 'percent';
        };
        type TableLayoutWidth = Extract<LayoutWidth, 'auto' | 'full-grid' | 'inset-left' | 'inset-right' | 'full-bleed'>;
        interface TableCaption extends Parent {
            type: 'table-caption';
            children: Phrasing[];
        }
        interface TableCell extends Parent {
            type: 'table-cell';
            heading?: boolean;
            columnSpan?: number;
            rowSpan?: number;
            children: Phrasing[];
        }
        interface TableRow extends Parent {
            type: 'table-row';
            children: TableCell[];
        }
        interface TableBody extends Parent {
            type: 'table-body';
            children: TableRow[];
        }
        interface TableFooter extends Parent {
            type: 'table-footer';
            children: Phrasing[];
        }
        interface Table extends Parent {
            type: 'table';
            stripes: boolean;
            compact: boolean;
            layoutWidth: TableLayoutWidth;
            collapseAfterHowManyRows?: number;
            responsiveStyle: 'overflow' | 'flat' | 'scroll';
            children: [TableCaption, TableBody, TableFooter] | [TableCaption, TableBody] | [TableBody, TableFooter] | [TableBody];
            columnSettings: TableColumnSettings[];
        }
        type CustomCodeComponentAttributes = {
            [key: string]: string | boolean | undefined;
        };
        interface CustomCodeComponent extends Node {
            /** Component type */
            type: "custom-code-component";
            /** Id taken from the CAPI url */
            id: string;
            /** How the component should be presented in the article page according to the column layout system */
            layoutWidth: LayoutWidth;
        }
        interface ImagePair extends Parent {
            type: 'image-pair';
            children: [ImageSet, ImageSet];
        }
        /**
         * Timeline nodes display a timeline of events in arbitrary order.
         */
        interface Timeline extends Parent {
            type: "timeline";
            /** The title for the timeline */
            title: string;
            children: TimelineEvent[];
        }
        /**
         * TimelineEvent is the representation of a single event in a Timeline.
         */
        interface TimelineEvent extends Parent {
            type: "timeline-event";
            /** The title of the event */
            title: string;
            /** Any combination of paragraphs and image sets */
            children: (Paragraph | ImageSet)[];
        }
        /**
         * A definition has a term and a related description. It is used to describe a term.
         */
        interface Definition extends Node {
            type: "definition";
            term: string;
            description: string;
        }
        /**
         * InNumbers represents a set of numbers with related descriptions.
         */
        interface InNumbers extends Parent {
            type: "in-numbers";
            /** The title for the InNumbers */
            title?: string;
            children: [Definition, Definition, Definition];
        }
        /** Allowed children for a card
        */
        type CardChildren = ImageSet | Exclude<FormattingBlock, Heading>;
        /**
        * A card describes a subject with images and text
        */
        interface Card extends Parent {
            type: "card";
            /** The title of this card */
            title?: string;
            children: CardChildren[];
        }
        /**
        * Allowed layout widths for an InfoBox.
        */
        type InfoBoxLayoutWidth = Extract<LayoutWidth, "full-width" | "inset-left">;
        /**
        * An info box describes a subject via a single card
        */
        interface InfoBox extends Parent {
            type: "info-box";
            /** The layout width supported by this node */
            layoutWidth: InfoBoxLayoutWidth;
            children: [Card];
        }
        /**
        * InfoPair provides exactly two cards.
        */
        interface InfoPair extends Parent {
            type: "info-pair";
            /** The title of the info pair */
            title?: string;
            children: [Card, Card];
        }
    }
    namespace loose {
        type FormattingBlock = Paragraph | Heading | List | Blockquote | ThematicBreak | Text;
        type StoryBlock = ImageSet | Flourish | BigNumber | CustomCodeComponent | Layout | Pullquote | ScrollyBlock | Table | Recommended | RecommendedList | Tweet | Video | YoutubeVideo | Timeline | ImagePair | InNumbers | Definition | InfoBox | InfoPair;
        type BodyBlock = FormattingBlock | StoryBlock;
        type LayoutWidth = "auto" | "in-line" | "inset-left" | "inset-right" | "full-bleed" | "full-grid" | "mid-grid" | "full-width";
        type Phrasing = Text | Break | Strong | Emphasis | Strikethrough | Link;
        interface Node {
            type: string;
            data?: any;
        }
        interface Parent extends Node {
            children: Node[];
        }
        interface Root extends Node {
            type: "root";
            body: Body;
        }
        interface Body extends Parent {
            type: "body";
            version: number;
            children: BodyBlock[];
        }
        interface Text extends Node {
            type: "text";
            value: string;
        }
        interface Break extends Node {
            type: "break";
        }
        interface ThematicBreak extends Node {
            type: "thematic-break";
        }
        interface Paragraph extends Parent {
            type: "paragraph";
            children: Phrasing[];
        }
        interface Heading extends Parent {
            type: "heading";
            children: Text[];
            level: "chapter" | "subheading" | "label";
            fragmentIdentifier?: string;
        }
        interface Strong extends Parent {
            type: "strong";
            children: Phrasing[];
        }
        interface Emphasis extends Parent {
            type: "emphasis";
            children: Phrasing[];
        }
        interface Strikethrough extends Parent {
            type: "strikethrough";
            children: Phrasing[];
        }
        interface Link extends Parent {
            type: "link";
            url: string;
            title: string;
            children: Phrasing[];
            styleType?: 'onward-journey';
        }
        interface List extends Parent {
            type: "list";
            ordered: boolean;
            children: ListItem[];
        }
        interface ListItem extends Parent {
            type: "list-item";
            children: (Paragraph | Phrasing)[];
        }
        interface Blockquote extends Parent {
            type: "blockquote";
            children: (Paragraph | Phrasing)[];
        }
        interface Pullquote extends Node {
            type: "pullquote";
            text: string;
            source?: string;
        }
        interface ImageSet extends Node {
            type: "image-set";
            id: string;
            picture?: ImageSetPicture;
            fragmentIdentifier?: string;
        }
        type ImageSetPicture = {
            layoutWidth: string;
            imageType: "image" | "graphic";
            alt: string;
            caption: string;
            credit: string;
            images: Image[];
            fallbackImage: Image;
        };
        type Image = {
            id: string;
            width: number;
            height: number;
            format: "desktop" | "mobile" | "square" | "square-ftedit" | "standard" | "wide" | "standard-inline";
            url: string;
            sourceSet?: ImageSource[];
        };
        type ImageSource = {
            url: string;
            width: number;
            dpr: number;
        };
        interface Recommended extends Node {
            type: "recommended";
            id: string;
            heading?: string;
            teaserTitleOverride?: string;
            teaser?: Teaser;
        }
        interface RecommendedList extends Node {
            type: "recommended-list";
            heading?: string;
            children: Recommended[];
        }
        type TeaserConcept = {
            apiUrl: string;
            directType: string;
            id: string;
            predicate: string;
            prefLabel: string;
            type: string;
            types: string[];
            url: string;
        };
        type Teaser = {
            id: string;
            url: string;
            type: "article" | "video" | "podcast" | "audio" | "package" | "liveblog" | "promoted-content" | "paid-post";
            title: string;
            publishedDate: string;
            firstPublishedDate: string;
            metaLink?: TeaserConcept;
            metaAltLink?: TeaserConcept;
            metaPrefixText?: string;
            metaSuffixText?: string;
            indicators: {
                accessLevel: "premium" | "subscribed" | "registered" | "free";
                isOpinion?: boolean;
                isColumn?: boolean;
                isPodcast?: boolean;
                isEditorsChoice?: boolean;
                isExclusive?: boolean;
                isScoop?: boolean;
            };
            image: {
                url: string;
                width: number;
                height: number;
            };
            clientName?: string;
        };
        interface Tweet extends Node {
            id: string;
            type: "tweet";
            html?: string;
        }
        type FlourishLayoutWidth = Extract<LayoutWidth, "full-grid" | "in-line">;
        interface Flourish extends Node {
            type: "flourish";
            id: string;
            layoutWidth: FlourishLayoutWidth;
            flourishType: string;
            description?: string;
            timestamp?: string;
            fallbackImage?: Image;
            fragmentIdentifier?: string;
        }
        interface BigNumber extends Node {
            type: "big-number";
            number: string;
            description: string;
        }
        interface Video extends Node {
            type: "video";
            id: string;
            title?: string;
        }
        interface YoutubeVideo extends Node {
            type: "youtube-video";
            url: string;
        }
        interface ScrollyBlock extends Parent {
            type: "scrolly-block";
            theme: "sans" | "serif";
            children: ScrollySection[];
        }
        interface ScrollySection extends Parent {
            type: "scrolly-section";
            display: "dark-background" | "light-background";
            noBox?: true;
            position: "left" | "center" | "right";
            transition?: "delay-before" | "delay-after";
            children: [ScrollyImage, ...ScrollyCopy[]];
        }
        interface ScrollyImage extends Node {
            type: "scrolly-image";
            id: string;
            picture?: ImageSetPicture;
        }
        interface ScrollyCopy extends Parent {
            type: "scrolly-copy";
            children: (ScrollyHeading | Paragraph)[];
        }
        interface ScrollyHeading extends Parent {
            type: "scrolly-heading";
            level: "chapter" | "heading" | "subheading";
            children: Text[];
        }
        interface Layout extends Parent {
            type: "layout";
            layoutName: "auto" | "card" | "timeline";
            layoutWidth: string;
            children: [Heading, LayoutImage, ...LayoutSlot[]] | [Heading, ...LayoutSlot[]] | LayoutSlot[];
        }
        interface LayoutSlot extends Parent {
            type: "layout-slot";
            children: (Heading | Paragraph | LayoutImage)[];
        }
        interface LayoutImage extends Node {
            type: "layout-image";
            id: string;
            alt: string;
            caption: string;
            credit: string;
            picture?: ImageSetPicture;
        }
        type TableColumnSettings = {
            hideOnMobile: boolean;
            sortable: boolean;
            sortType: 'text' | 'number' | 'date' | 'currency' | 'percent';
        };
        type TableLayoutWidth = Extract<LayoutWidth, 'auto' | 'full-grid' | 'inset-left' | 'inset-right' | 'full-bleed'>;
        interface TableCaption extends Parent {
            type: 'table-caption';
            children: Phrasing[];
        }
        interface TableCell extends Parent {
            type: 'table-cell';
            heading?: boolean;
            columnSpan?: number;
            rowSpan?: number;
            children: Phrasing[];
        }
        interface TableRow extends Parent {
            type: 'table-row';
            children: TableCell[];
        }
        interface TableBody extends Parent {
            type: 'table-body';
            children: TableRow[];
        }
        interface TableFooter extends Parent {
            type: 'table-footer';
            children: Phrasing[];
        }
        interface Table extends Parent {
            type: 'table';
            stripes: boolean;
            compact: boolean;
            layoutWidth: TableLayoutWidth;
            collapseAfterHowManyRows?: number;
            responsiveStyle: 'overflow' | 'flat' | 'scroll';
            children: [TableCaption, TableBody, TableFooter] | [TableCaption, TableBody] | [TableBody, TableFooter] | [TableBody];
            columnSettings: TableColumnSettings[];
        }
        type CustomCodeComponentAttributes = {
            [key: string]: string | boolean | undefined;
        };
        interface CustomCodeComponent extends Node {
            /** Component type */
            type: "custom-code-component";
            /** Id taken from the CAPI url */
            id: string;
            /** How the component should be presented in the article page according to the column layout system */
            layoutWidth: LayoutWidth;
            /** Repository for the code of the component in the format "[github org]/[github repo]/[component name]". */
            path?: string;
            /** Semantic version of the code of the component, e.g. "^0.3.5". */
            versionRange?: string;
            /** Last date-time when the attributes for this block were modified, in ISO-8601 format. */
            attributesLastModified?: string;
            /** Configuration data to be passed to the component. */
            attributes?: CustomCodeComponentAttributes;
        }
        interface ImagePair extends Parent {
            type: 'image-pair';
            children: [ImageSet, ImageSet];
        }
        /**
         * Timeline nodes display a timeline of events in arbitrary order.
         */
        interface Timeline extends Parent {
            type: "timeline";
            /** The title for the timeline */
            title: string;
            children: TimelineEvent[];
        }
        /**
         * TimelineEvent is the representation of a single event in a Timeline.
         */
        interface TimelineEvent extends Parent {
            type: "timeline-event";
            /** The title of the event */
            title: string;
            /** Any combination of paragraphs and image sets */
            children: (Paragraph | ImageSet)[];
        }
        /**
         * A definition has a term and a related description. It is used to describe a term.
         */
        interface Definition extends Node {
            type: "definition";
            term: string;
            description: string;
        }
        /**
         * InNumbers represents a set of numbers with related descriptions.
         */
        interface InNumbers extends Parent {
            type: "in-numbers";
            /** The title for the InNumbers */
            title?: string;
            children: [Definition, Definition, Definition];
        }
        /** Allowed children for a card
        */
        type CardChildren = ImageSet | Exclude<FormattingBlock, Heading>;
        /**
        * A card describes a subject with images and text
        */
        interface Card extends Parent {
            type: "card";
            /** The title of this card */
            title?: string;
            children: CardChildren[];
        }
        /**
        * Allowed layout widths for an InfoBox.
        */
        type InfoBoxLayoutWidth = Extract<LayoutWidth, "full-width" | "inset-left">;
        /**
        * An info box describes a subject via a single card
        */
        interface InfoBox extends Parent {
            type: "info-box";
            /** The layout width supported by this node */
            layoutWidth: InfoBoxLayoutWidth;
            children: [Card];
        }
        /**
        * InfoPair provides exactly two cards.
        */
        interface InfoPair extends Parent {
            type: "info-pair";
            /** The title of the info pair */
            title?: string;
            children: [Card, Card];
        }
    }
}
