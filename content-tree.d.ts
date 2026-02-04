export declare namespace ContentTree {
    export type LayoutWidth = "auto" | "in-line" | "inset-left" | "inset-right" | "full-bleed" | "full-grid" | "mid-grid" | "full-width";
    export interface Node {
        type: string;
        data?: any;
    }
    export interface Parent extends Node {
        children: Node[];
    }
    export interface Root extends Node {
        type: "root";
        body: Body;
    }
    export interface Body extends Parent {
        type: "body";
        version: number;
        children: BodyBlock[];
    }
    export type BodyBlock = FormattingBlock | StoryBlock;
    export type FormattingBlock = Paragraph | Heading | List | Blockquote | ThematicBreak | Text;
    export interface Text extends Node {
        type: "text";
        value: string;
    }
    export type Phrasing = Text | Break | Strong | Emphasis | Strikethrough | Link | FindOutMoreLink;
    export interface Break extends Node {
        type: "break";
    }
    export interface ThematicBreak extends Node {
        type: "thematic-break";
    }
    export interface Paragraph extends Parent {
        type: "paragraph";
        children: Phrasing[];
    }
    export interface Heading extends Parent {
        type: "heading";
        children: Text[];
        level: "chapter" | "subheading" | "label";
        fragmentIdentifier?: string;
    }
    export interface Strong extends Parent {
        type: "strong";
        children: Phrasing[];
    }
    export interface Emphasis extends Parent {
        type: "emphasis";
        children: Phrasing[];
    }
    export interface Strikethrough extends Parent {
        type: "strikethrough";
        children: Phrasing[];
    }
    export interface Link extends Parent {
        type: "link";
        url: string;
        title: string;
        children: Phrasing[];
    }
    export interface FindOutMoreLink extends Parent {
        type: "find-out-more-link";
        url: string;
        title: string;
        children: Phrasing[];
    }
    export interface List extends Parent {
        type: "list";
        ordered: boolean;
        children: ListItem[];
    }
    export interface ListItem extends Parent {
        type: "list-item";
        children: (Paragraph | Phrasing)[];
    }
    export interface Blockquote extends Parent {
        type: "blockquote";
        children: (Paragraph | Phrasing)[];
    }
    export type StoryBlock = ImageSet | Flourish | BigNumber | CustomCodeComponent | Layout | Pullquote | ScrollyBlock | Table | Recommended | RecommendedList | Tweet | Video | YoutubeVideo | Timeline | ImagePair | InNumbers | Definition | InfoBox | InfoPair | AudioPlayer;
    export interface Pullquote extends Node {
        type: "pullquote";
        text: string;
        source?: string;
    }
    export interface ImageSet extends Node {
        type: "image-set";
        id: string;
        picture: ImageSetPicture;
        fragmentIdentifier?: string;
    }
    export type ImageSetPicture = {
        layoutWidth: string;
        imageType: "image" | "graphic";
        alt: string;
        caption: string;
        credit: string;
        images: Image[];
        fallbackImage: Image;
    };
    export type Image = {
        id: string;
        width: number;
        height: number;
        format: "desktop" | "mobile" | "square" | "square-ftedit" | "standard" | "wide" | "standard-inline";
        url: string;
        sourceSet?: ImageSource[];
    };
    export type ImageSource = {
        url: string;
        width: number;
        dpr: number;
    };
    export interface Recommended extends Node {
        type: "recommended";
        id: string;
        heading?: string;
        teaserTitleOverride?: string;
        teaser: Teaser;
    }
    export interface RecommendedList extends Node {
        type: "recommended-list";
        heading?: string;
        children: Recommended[];
    }
    export type TeaserConcept = {
        apiUrl: string;
        directType: string;
        id: string;
        predicate: string;
        prefLabel: string;
        type: string;
        types: string[];
        url: string;
    };
    export type Teaser = {
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
    export interface Tweet extends Node {
        id: string;
        type: "tweet";
        html: string;
    }
    export type FlourishLayoutWidth = Extract<LayoutWidth, "full-grid" | "in-line">;
    export interface Flourish extends Node {
        type: "flourish";
        id: string;
        layoutWidth: FlourishLayoutWidth;
        flourishType: string;
        description?: string;
        timestamp?: string;
        fallbackImage?: Image;
        fragmentIdentifier?: string;
    }
    export interface BigNumber extends Node {
        type: "big-number";
        number: string;
        description: string;
    }
    export interface Video extends Node {
        type: "video";
        id: string;
        title: string;
    }
    export interface YoutubeVideo extends Node {
        type: "youtube-video";
        url: string;
    }
    export interface ScrollyBlock extends Parent {
        type: "scrolly-block";
        theme: "sans" | "serif";
        children: ScrollySection[];
    }
    export interface ScrollySection extends Parent {
        type: "scrolly-section";
        display: "dark-background" | "light-background";
        noBox?: true;
        position: "left" | "center" | "right";
        transition?: "delay-before" | "delay-after";
        children: [ScrollyImage, ...ScrollyCopy[]];
    }
    export interface ScrollyImage extends Node {
        type: "scrolly-image";
        id: string;
        picture: ImageSetPicture;
    }
    export interface ScrollyCopy extends Parent {
        type: "scrolly-copy";
        children: (ScrollyHeading | Paragraph)[];
    }
    export interface ScrollyHeading extends Parent {
        type: "scrolly-heading";
        level: "chapter" | "heading" | "subheading";
        children: Text[];
    }
    export interface Layout extends Parent {
        type: "layout";
        layoutName: "auto" | "card" | "timeline";
        layoutWidth: string;
        children: [Heading, LayoutImage, ...LayoutSlot[]] | [Heading, ...LayoutSlot[]] | LayoutSlot[];
    }
    export interface LayoutSlot extends Parent {
        type: "layout-slot";
        children: (Heading | Paragraph | LayoutImage)[];
    }
    export interface LayoutImage extends Node {
        type: "layout-image";
        id: string;
        alt: string;
        caption: string;
        credit: string;
        picture: ImageSetPicture;
    }
    export type TableColumnSettings = {
        hideOnMobile: boolean;
        sortable: boolean;
        sortType: 'text' | 'number' | 'date' | 'currency' | 'percent';
    };
    export type TableLayoutWidth = Extract<LayoutWidth, 'auto' | 'full-grid' | 'inset-left' | 'inset-right' | 'full-bleed'>;
    export interface TableCaption extends Parent {
        type: 'table-caption';
        children: Phrasing[];
    }
    export interface TableCell extends Parent {
        type: 'table-cell';
        heading?: boolean;
        columnSpan?: number;
        rowSpan?: number;
        children: Phrasing[];
    }
    export interface TableRow extends Parent {
        type: 'table-row';
        children: TableCell[];
    }
    export interface TableBody extends Parent {
        type: 'table-body';
        children: TableRow[];
    }
    export interface TableFooter extends Parent {
        type: 'table-footer';
        children: Phrasing[];
    }
    export interface Table extends Parent {
        type: 'table';
        stripes: boolean;
        compact: boolean;
        layoutWidth: TableLayoutWidth;
        collapseAfterHowManyRows?: number;
        responsiveStyle: 'overflow' | 'flat' | 'scroll';
        children: [TableCaption, TableBody, TableFooter] | [TableCaption, TableBody] | [TableBody, TableFooter] | [TableBody];
        columnSettings: TableColumnSettings[];
    }
    export type CustomCodeComponentAttributes = {
        [key: string]: string | boolean | undefined;
    };
    export interface CustomCodeComponent extends Node {
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
    export interface ImagePair extends Parent {
        type: 'image-pair';
        children: [ImageSet, ImageSet];
    }
    /**
     * Timeline nodes display a timeline of events in arbitrary order.
     */
    export interface Timeline extends Parent {
        type: "timeline";
        /** The title for the timeline */
        title: string;
        children: TimelineEvent[];
    }
    /**
     * TimelineEvent is the representation of a single event in a Timeline.
     */
    export interface TimelineEvent extends Parent {
        type: "timeline-event";
        /** The title of the event */
        title: string;
        /** Any combination of paragraphs and image sets */
        children: (Paragraph | ImageSet)[];
    }
    /**
     * A definition has a term and a related description. It is used to describe a term.
     */
    export interface Definition extends Node {
        type: "definition";
        term: string;
        description: string;
    }
    /**
     * InNumbers represents a set of numbers with related descriptions.
     */
    export interface InNumbers extends Parent {
        type: "in-numbers";
        /** The title for the InNumbers */
        title?: string;
        children: [Definition, Definition, Definition];
    }
    /** Allowed children for a card
    */
    export type CardChildren = ImageSet | Exclude<FormattingBlock, Heading>;
    /**
    * A card describes a subject with images and text
    */
    export interface Card extends Parent {
        type: "card";
        /** The title of this card */
        title?: string;
        children: CardChildren[];
    }
    /**
    * Allowed layout widths for an InfoBox.
    */
    export type InfoBoxLayoutWidth = Extract<LayoutWidth, "in-line" | "inset-left">;
    /**
    * An info box describes a subject via a single card
    */
    export interface InfoBox extends Parent {
        type: "info-box";
        /** The layout width supported by this node */
        layoutWidth: InfoBoxLayoutWidth;
        children: [Card];
    }
    /**
    * InfoPair provides exactly two cards.
    */
    export interface InfoPair extends Parent {
        type: "info-pair";
        /** The title of the info pair */
        title?: string;
        children: [Card, Card];
    }
    /**
       * @sparkGenerateStoryblock true
       **/
    type AudioPlayer = AudioPlayerV1 | AudioPlayerV2 | AudioPlayerV3;
    /** @support deprecated */
    export interface AudioPlayerV1 extends Node {
        type: "audio-player";
        version: 1;
        title: string;
        audioUrl: string;
    }
    export interface AudioPlayerV2 extends Node {
        type: "audio-player";
        version: 2;
        title: string;
        audioId: string;
        audio: AudioSet;
    }
    /** @support prerelease */
    export interface AudioPlayerV3 extends Node {
        type: "audio-player";
        version: 3;
        title: string;
        audioId: string;
        transcriptionId: string;
        audio: AudioSet;
        transcription: Transcription;
    }
    /**
     * Demo placeholders so the AudioPlayer versioning example compiles.
     */
    export interface AudioSet extends Node {
        url: string;
    }
    export interface Transcription extends Node {
        text: string;
    }
    export namespace full {
        export type LayoutWidth = "auto" | "in-line" | "inset-left" | "inset-right" | "full-bleed" | "full-grid" | "mid-grid" | "full-width";
        export interface Node {
            type: string;
            data?: any;
        }
        export interface Parent extends Node {
            children: Node[];
        }
        export interface Root extends Node {
            type: "root";
            body: Body;
        }
        export interface Body extends Parent {
            type: "body";
            version: number;
            children: BodyBlock[];
        }
        export type BodyBlock = FormattingBlock | StoryBlock;
        export type FormattingBlock = Paragraph | Heading | List | Blockquote | ThematicBreak | Text;
        export interface Text extends Node {
            type: "text";
            value: string;
        }
        export type Phrasing = Text | Break | Strong | Emphasis | Strikethrough | Link | FindOutMoreLink;
        export interface Break extends Node {
            type: "break";
        }
        export interface ThematicBreak extends Node {
            type: "thematic-break";
        }
        export interface Paragraph extends Parent {
            type: "paragraph";
            children: Phrasing[];
        }
        export interface Heading extends Parent {
            type: "heading";
            children: Text[];
            level: "chapter" | "subheading" | "label";
            fragmentIdentifier?: string;
        }
        export interface Strong extends Parent {
            type: "strong";
            children: Phrasing[];
        }
        export interface Emphasis extends Parent {
            type: "emphasis";
            children: Phrasing[];
        }
        export interface Strikethrough extends Parent {
            type: "strikethrough";
            children: Phrasing[];
        }
        export interface Link extends Parent {
            type: "link";
            url: string;
            title: string;
            children: Phrasing[];
        }
        export interface FindOutMoreLink extends Parent {
            type: "find-out-more-link";
            url: string;
            title: string;
            children: Phrasing[];
        }
        export interface List extends Parent {
            type: "list";
            ordered: boolean;
            children: ListItem[];
        }
        export interface ListItem extends Parent {
            type: "list-item";
            children: (Paragraph | Phrasing)[];
        }
        export interface Blockquote extends Parent {
            type: "blockquote";
            children: (Paragraph | Phrasing)[];
        }
        export type StoryBlock = ImageSet | Flourish | BigNumber | CustomCodeComponent | Layout | Pullquote | ScrollyBlock | Table | Recommended | RecommendedList | Tweet | Video | YoutubeVideo | Timeline | ImagePair | InNumbers | Definition | InfoBox | InfoPair | AudioPlayer;
        export interface Pullquote extends Node {
            type: "pullquote";
            text: string;
            source?: string;
        }
        export interface ImageSet extends Node {
            type: "image-set";
            id: string;
            picture: ImageSetPicture;
            fragmentIdentifier?: string;
        }
        export type ImageSetPicture = {
            layoutWidth: string;
            imageType: "image" | "graphic";
            alt: string;
            caption: string;
            credit: string;
            images: Image[];
            fallbackImage: Image;
        };
        export type Image = {
            id: string;
            width: number;
            height: number;
            format: "desktop" | "mobile" | "square" | "square-ftedit" | "standard" | "wide" | "standard-inline";
            url: string;
            sourceSet?: ImageSource[];
        };
        export type ImageSource = {
            url: string;
            width: number;
            dpr: number;
        };
        export interface Recommended extends Node {
            type: "recommended";
            id: string;
            heading?: string;
            teaserTitleOverride?: string;
            teaser: Teaser;
        }
        export interface RecommendedList extends Node {
            type: "recommended-list";
            heading?: string;
            children: Recommended[];
        }
        export type TeaserConcept = {
            apiUrl: string;
            directType: string;
            id: string;
            predicate: string;
            prefLabel: string;
            type: string;
            types: string[];
            url: string;
        };
        export type Teaser = {
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
        export interface Tweet extends Node {
            id: string;
            type: "tweet";
            html: string;
        }
        export type FlourishLayoutWidth = Extract<LayoutWidth, "full-grid" | "in-line">;
        export interface Flourish extends Node {
            type: "flourish";
            id: string;
            layoutWidth: FlourishLayoutWidth;
            flourishType: string;
            description?: string;
            timestamp?: string;
            fallbackImage?: Image;
            fragmentIdentifier?: string;
        }
        export interface BigNumber extends Node {
            type: "big-number";
            number: string;
            description: string;
        }
        export interface Video extends Node {
            type: "video";
            id: string;
            title: string;
        }
        export interface YoutubeVideo extends Node {
            type: "youtube-video";
            url: string;
        }
        export interface ScrollyBlock extends Parent {
            type: "scrolly-block";
            theme: "sans" | "serif";
            children: ScrollySection[];
        }
        export interface ScrollySection extends Parent {
            type: "scrolly-section";
            display: "dark-background" | "light-background";
            noBox?: true;
            position: "left" | "center" | "right";
            transition?: "delay-before" | "delay-after";
            children: [ScrollyImage, ...ScrollyCopy[]];
        }
        export interface ScrollyImage extends Node {
            type: "scrolly-image";
            id: string;
            picture: ImageSetPicture;
        }
        export interface ScrollyCopy extends Parent {
            type: "scrolly-copy";
            children: (ScrollyHeading | Paragraph)[];
        }
        export interface ScrollyHeading extends Parent {
            type: "scrolly-heading";
            level: "chapter" | "heading" | "subheading";
            children: Text[];
        }
        export interface Layout extends Parent {
            type: "layout";
            layoutName: "auto" | "card" | "timeline";
            layoutWidth: string;
            children: [Heading, LayoutImage, ...LayoutSlot[]] | [Heading, ...LayoutSlot[]] | LayoutSlot[];
        }
        export interface LayoutSlot extends Parent {
            type: "layout-slot";
            children: (Heading | Paragraph | LayoutImage)[];
        }
        export interface LayoutImage extends Node {
            type: "layout-image";
            id: string;
            alt: string;
            caption: string;
            credit: string;
            picture: ImageSetPicture;
        }
        export type TableColumnSettings = {
            hideOnMobile: boolean;
            sortable: boolean;
            sortType: 'text' | 'number' | 'date' | 'currency' | 'percent';
        };
        export type TableLayoutWidth = Extract<LayoutWidth, 'auto' | 'full-grid' | 'inset-left' | 'inset-right' | 'full-bleed'>;
        export interface TableCaption extends Parent {
            type: 'table-caption';
            children: Phrasing[];
        }
        export interface TableCell extends Parent {
            type: 'table-cell';
            heading?: boolean;
            columnSpan?: number;
            rowSpan?: number;
            children: Phrasing[];
        }
        export interface TableRow extends Parent {
            type: 'table-row';
            children: TableCell[];
        }
        export interface TableBody extends Parent {
            type: 'table-body';
            children: TableRow[];
        }
        export interface TableFooter extends Parent {
            type: 'table-footer';
            children: Phrasing[];
        }
        export interface Table extends Parent {
            type: 'table';
            stripes: boolean;
            compact: boolean;
            layoutWidth: TableLayoutWidth;
            collapseAfterHowManyRows?: number;
            responsiveStyle: 'overflow' | 'flat' | 'scroll';
            children: [TableCaption, TableBody, TableFooter] | [TableCaption, TableBody] | [TableBody, TableFooter] | [TableBody];
            columnSettings: TableColumnSettings[];
        }
        export type CustomCodeComponentAttributes = {
            [key: string]: string | boolean | undefined;
        };
        export interface CustomCodeComponent extends Node {
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
        export interface ImagePair extends Parent {
            type: 'image-pair';
            children: [ImageSet, ImageSet];
        }
        /**
         * Timeline nodes display a timeline of events in arbitrary order.
         */
        export interface Timeline extends Parent {
            type: "timeline";
            /** The title for the timeline */
            title: string;
            children: TimelineEvent[];
        }
        /**
         * TimelineEvent is the representation of a single event in a Timeline.
         */
        export interface TimelineEvent extends Parent {
            type: "timeline-event";
            /** The title of the event */
            title: string;
            /** Any combination of paragraphs and image sets */
            children: (Paragraph | ImageSet)[];
        }
        /**
         * A definition has a term and a related description. It is used to describe a term.
         */
        export interface Definition extends Node {
            type: "definition";
            term: string;
            description: string;
        }
        /**
         * InNumbers represents a set of numbers with related descriptions.
         */
        export interface InNumbers extends Parent {
            type: "in-numbers";
            /** The title for the InNumbers */
            title?: string;
            children: [Definition, Definition, Definition];
        }
        /** Allowed children for a card
        */
        export type CardChildren = ImageSet | Exclude<FormattingBlock, Heading>;
        /**
        * A card describes a subject with images and text
        */
        export interface Card extends Parent {
            type: "card";
            /** The title of this card */
            title?: string;
            children: CardChildren[];
        }
        /**
        * Allowed layout widths for an InfoBox.
        */
        export type InfoBoxLayoutWidth = Extract<LayoutWidth, "in-line" | "inset-left">;
        /**
        * An info box describes a subject via a single card
        */
        export interface InfoBox extends Parent {
            type: "info-box";
            /** The layout width supported by this node */
            layoutWidth: InfoBoxLayoutWidth;
            children: [Card];
        }
        /**
        * InfoPair provides exactly two cards.
        */
        export interface InfoPair extends Parent {
            type: "info-pair";
            /** The title of the info pair */
            title?: string;
            children: [Card, Card];
        }
        /**
           * @sparkGenerateStoryblock true
           **/
        type AudioPlayer = AudioPlayerV1 | AudioPlayerV2 | AudioPlayerV3;
        /** @support deprecated */
        export interface AudioPlayerV1 extends Node {
            type: "audio-player";
            version: 1;
            title: string;
            audioUrl: string;
        }
        export interface AudioPlayerV2 extends Node {
            type: "audio-player";
            version: 2;
            title: string;
            audioId: string;
            audio: AudioSet;
        }
        /** @support prerelease */
        export interface AudioPlayerV3 extends Node {
            type: "audio-player";
            version: 3;
            title: string;
            audioId: string;
            transcriptionId: string;
            audio: AudioSet;
            transcription: Transcription;
        }
        /**
         * Demo placeholders so the AudioPlayer versioning example compiles.
         */
        export interface AudioSet extends Node {
            url: string;
        }
        export interface Transcription extends Node {
            text: string;
        }
        export {};
    }
    export namespace transit {
        export type LayoutWidth = "auto" | "in-line" | "inset-left" | "inset-right" | "full-bleed" | "full-grid" | "mid-grid" | "full-width";
        export interface Node {
            type: string;
            data?: any;
        }
        export interface Parent extends Node {
            children: Node[];
        }
        export interface Root extends Node {
            type: "root";
            body: Body;
        }
        export interface Body extends Parent {
            type: "body";
            version: number;
            children: BodyBlock[];
        }
        export type BodyBlock = FormattingBlock | StoryBlock;
        export type FormattingBlock = Paragraph | Heading | List | Blockquote | ThematicBreak | Text;
        export interface Text extends Node {
            type: "text";
            value: string;
        }
        export type Phrasing = Text | Break | Strong | Emphasis | Strikethrough | Link | FindOutMoreLink;
        export interface Break extends Node {
            type: "break";
        }
        export interface ThematicBreak extends Node {
            type: "thematic-break";
        }
        export interface Paragraph extends Parent {
            type: "paragraph";
            children: Phrasing[];
        }
        export interface Heading extends Parent {
            type: "heading";
            children: Text[];
            level: "chapter" | "subheading" | "label";
            fragmentIdentifier?: string;
        }
        export interface Strong extends Parent {
            type: "strong";
            children: Phrasing[];
        }
        export interface Emphasis extends Parent {
            type: "emphasis";
            children: Phrasing[];
        }
        export interface Strikethrough extends Parent {
            type: "strikethrough";
            children: Phrasing[];
        }
        export interface Link extends Parent {
            type: "link";
            url: string;
            title: string;
            children: Phrasing[];
        }
        export interface FindOutMoreLink extends Parent {
            type: "find-out-more-link";
            url: string;
            title: string;
            children: Phrasing[];
        }
        export interface List extends Parent {
            type: "list";
            ordered: boolean;
            children: ListItem[];
        }
        export interface ListItem extends Parent {
            type: "list-item";
            children: (Paragraph | Phrasing)[];
        }
        export interface Blockquote extends Parent {
            type: "blockquote";
            children: (Paragraph | Phrasing)[];
        }
        export type StoryBlock = ImageSet | Flourish | BigNumber | CustomCodeComponent | Layout | Pullquote | ScrollyBlock | Table | Recommended | RecommendedList | Tweet | Video | YoutubeVideo | Timeline | ImagePair | InNumbers | Definition | InfoBox | InfoPair | AudioPlayer;
        export interface Pullquote extends Node {
            type: "pullquote";
            text: string;
            source?: string;
        }
        export interface ImageSet extends Node {
            type: "image-set";
            id: string;
            fragmentIdentifier?: string;
        }
        export type ImageSetPicture = {
            layoutWidth: string;
            imageType: "image" | "graphic";
            alt: string;
            caption: string;
            credit: string;
            images: Image[];
            fallbackImage: Image;
        };
        export type Image = {
            id: string;
            width: number;
            height: number;
            format: "desktop" | "mobile" | "square" | "square-ftedit" | "standard" | "wide" | "standard-inline";
            url: string;
            sourceSet?: ImageSource[];
        };
        export type ImageSource = {
            url: string;
            width: number;
            dpr: number;
        };
        export interface Recommended extends Node {
            type: "recommended";
            id: string;
            heading?: string;
            teaserTitleOverride?: string;
        }
        export interface RecommendedList extends Node {
            type: "recommended-list";
            heading?: string;
            children: Recommended[];
        }
        export type TeaserConcept = {
            apiUrl: string;
            directType: string;
            id: string;
            predicate: string;
            prefLabel: string;
            type: string;
            types: string[];
            url: string;
        };
        export type Teaser = {
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
        export interface Tweet extends Node {
            id: string;
            type: "tweet";
        }
        export type FlourishLayoutWidth = Extract<LayoutWidth, "full-grid" | "in-line">;
        export interface Flourish extends Node {
            type: "flourish";
            id: string;
            layoutWidth: FlourishLayoutWidth;
            flourishType: string;
            description?: string;
            timestamp?: string;
            fragmentIdentifier?: string;
        }
        export interface BigNumber extends Node {
            type: "big-number";
            number: string;
            description: string;
        }
        export interface Video extends Node {
            type: "video";
            id: string;
        }
        export interface YoutubeVideo extends Node {
            type: "youtube-video";
            url: string;
        }
        export interface ScrollyBlock extends Parent {
            type: "scrolly-block";
            theme: "sans" | "serif";
            children: ScrollySection[];
        }
        export interface ScrollySection extends Parent {
            type: "scrolly-section";
            display: "dark-background" | "light-background";
            noBox?: true;
            position: "left" | "center" | "right";
            transition?: "delay-before" | "delay-after";
            children: [ScrollyImage, ...ScrollyCopy[]];
        }
        export interface ScrollyImage extends Node {
            type: "scrolly-image";
            id: string;
        }
        export interface ScrollyCopy extends Parent {
            type: "scrolly-copy";
            children: (ScrollyHeading | Paragraph)[];
        }
        export interface ScrollyHeading extends Parent {
            type: "scrolly-heading";
            level: "chapter" | "heading" | "subheading";
            children: Text[];
        }
        export interface Layout extends Parent {
            type: "layout";
            layoutName: "auto" | "card" | "timeline";
            layoutWidth: string;
            children: [Heading, LayoutImage, ...LayoutSlot[]] | [Heading, ...LayoutSlot[]] | LayoutSlot[];
        }
        export interface LayoutSlot extends Parent {
            type: "layout-slot";
            children: (Heading | Paragraph | LayoutImage)[];
        }
        export interface LayoutImage extends Node {
            type: "layout-image";
            id: string;
            alt: string;
            caption: string;
            credit: string;
        }
        export type TableColumnSettings = {
            hideOnMobile: boolean;
            sortable: boolean;
            sortType: 'text' | 'number' | 'date' | 'currency' | 'percent';
        };
        export type TableLayoutWidth = Extract<LayoutWidth, 'auto' | 'full-grid' | 'inset-left' | 'inset-right' | 'full-bleed'>;
        export interface TableCaption extends Parent {
            type: 'table-caption';
            children: Phrasing[];
        }
        export interface TableCell extends Parent {
            type: 'table-cell';
            heading?: boolean;
            columnSpan?: number;
            rowSpan?: number;
            children: Phrasing[];
        }
        export interface TableRow extends Parent {
            type: 'table-row';
            children: TableCell[];
        }
        export interface TableBody extends Parent {
            type: 'table-body';
            children: TableRow[];
        }
        export interface TableFooter extends Parent {
            type: 'table-footer';
            children: Phrasing[];
        }
        export interface Table extends Parent {
            type: 'table';
            stripes: boolean;
            compact: boolean;
            layoutWidth: TableLayoutWidth;
            collapseAfterHowManyRows?: number;
            responsiveStyle: 'overflow' | 'flat' | 'scroll';
            children: [TableCaption, TableBody, TableFooter] | [TableCaption, TableBody] | [TableBody, TableFooter] | [TableBody];
            columnSettings: TableColumnSettings[];
        }
        export type CustomCodeComponentAttributes = {
            [key: string]: string | boolean | undefined;
        };
        export interface CustomCodeComponent extends Node {
            /** Component type */
            type: "custom-code-component";
            /** Id taken from the CAPI url */
            id: string;
            /** How the component should be presented in the article page according to the column layout system */
            layoutWidth: LayoutWidth;
        }
        export interface ImagePair extends Parent {
            type: 'image-pair';
            children: [ImageSet, ImageSet];
        }
        /**
         * Timeline nodes display a timeline of events in arbitrary order.
         */
        export interface Timeline extends Parent {
            type: "timeline";
            /** The title for the timeline */
            title: string;
            children: TimelineEvent[];
        }
        /**
         * TimelineEvent is the representation of a single event in a Timeline.
         */
        export interface TimelineEvent extends Parent {
            type: "timeline-event";
            /** The title of the event */
            title: string;
            /** Any combination of paragraphs and image sets */
            children: (Paragraph | ImageSet)[];
        }
        /**
         * A definition has a term and a related description. It is used to describe a term.
         */
        export interface Definition extends Node {
            type: "definition";
            term: string;
            description: string;
        }
        /**
         * InNumbers represents a set of numbers with related descriptions.
         */
        export interface InNumbers extends Parent {
            type: "in-numbers";
            /** The title for the InNumbers */
            title?: string;
            children: [Definition, Definition, Definition];
        }
        /** Allowed children for a card
        */
        export type CardChildren = ImageSet | Exclude<FormattingBlock, Heading>;
        /**
        * A card describes a subject with images and text
        */
        export interface Card extends Parent {
            type: "card";
            /** The title of this card */
            title?: string;
            children: CardChildren[];
        }
        /**
        * Allowed layout widths for an InfoBox.
        */
        export type InfoBoxLayoutWidth = Extract<LayoutWidth, "in-line" | "inset-left">;
        /**
        * An info box describes a subject via a single card
        */
        export interface InfoBox extends Parent {
            type: "info-box";
            /** The layout width supported by this node */
            layoutWidth: InfoBoxLayoutWidth;
            children: [Card];
        }
        /**
        * InfoPair provides exactly two cards.
        */
        export interface InfoPair extends Parent {
            type: "info-pair";
            /** The title of the info pair */
            title?: string;
            children: [Card, Card];
        }
        /**
           * @sparkGenerateStoryblock true
           **/
        type AudioPlayer = AudioPlayerV1 | AudioPlayerV2 | AudioPlayerV3;
        /** @support deprecated */
        export interface AudioPlayerV1 extends Node {
            type: "audio-player";
            version: 1;
            title: string;
            audioUrl: string;
        }
        export interface AudioPlayerV2 extends Node {
            type: "audio-player";
            version: 2;
            title: string;
            audioId: string;
        }
        /** @support prerelease */
        export interface AudioPlayerV3 extends Node {
            type: "audio-player";
            version: 3;
            title: string;
            audioId: string;
            transcriptionId: string;
        }
        /**
         * Demo placeholders so the AudioPlayer versioning example compiles.
         */
        export interface AudioSet extends Node {
            url: string;
        }
        export interface Transcription extends Node {
            text: string;
        }
        export {};
    }
    export namespace loose {
        export type LayoutWidth = "auto" | "in-line" | "inset-left" | "inset-right" | "full-bleed" | "full-grid" | "mid-grid" | "full-width";
        export interface Node {
            type: string;
            data?: any;
        }
        export interface Parent extends Node {
            children: Node[];
        }
        export interface Root extends Node {
            type: "root";
            body: Body;
        }
        export interface Body extends Parent {
            type: "body";
            version: number;
            children: BodyBlock[];
        }
        export type BodyBlock = FormattingBlock | StoryBlock;
        export type FormattingBlock = Paragraph | Heading | List | Blockquote | ThematicBreak | Text;
        export interface Text extends Node {
            type: "text";
            value: string;
        }
        export type Phrasing = Text | Break | Strong | Emphasis | Strikethrough | Link | FindOutMoreLink;
        export interface Break extends Node {
            type: "break";
        }
        export interface ThematicBreak extends Node {
            type: "thematic-break";
        }
        export interface Paragraph extends Parent {
            type: "paragraph";
            children: Phrasing[];
        }
        export interface Heading extends Parent {
            type: "heading";
            children: Text[];
            level: "chapter" | "subheading" | "label";
            fragmentIdentifier?: string;
        }
        export interface Strong extends Parent {
            type: "strong";
            children: Phrasing[];
        }
        export interface Emphasis extends Parent {
            type: "emphasis";
            children: Phrasing[];
        }
        export interface Strikethrough extends Parent {
            type: "strikethrough";
            children: Phrasing[];
        }
        export interface Link extends Parent {
            type: "link";
            url: string;
            title: string;
            children: Phrasing[];
        }
        export interface FindOutMoreLink extends Parent {
            type: "find-out-more-link";
            url: string;
            title: string;
            children: Phrasing[];
        }
        export interface List extends Parent {
            type: "list";
            ordered: boolean;
            children: ListItem[];
        }
        export interface ListItem extends Parent {
            type: "list-item";
            children: (Paragraph | Phrasing)[];
        }
        export interface Blockquote extends Parent {
            type: "blockquote";
            children: (Paragraph | Phrasing)[];
        }
        export type StoryBlock = ImageSet | Flourish | BigNumber | CustomCodeComponent | Layout | Pullquote | ScrollyBlock | Table | Recommended | RecommendedList | Tweet | Video | YoutubeVideo | Timeline | ImagePair | InNumbers | Definition | InfoBox | InfoPair | AudioPlayer;
        export interface Pullquote extends Node {
            type: "pullquote";
            text: string;
            source?: string;
        }
        export interface ImageSet extends Node {
            type: "image-set";
            id: string;
            picture?: ImageSetPicture;
            fragmentIdentifier?: string;
        }
        export type ImageSetPicture = {
            layoutWidth: string;
            imageType: "image" | "graphic";
            alt: string;
            caption: string;
            credit: string;
            images: Image[];
            fallbackImage: Image;
        };
        export type Image = {
            id: string;
            width: number;
            height: number;
            format: "desktop" | "mobile" | "square" | "square-ftedit" | "standard" | "wide" | "standard-inline";
            url: string;
            sourceSet?: ImageSource[];
        };
        export type ImageSource = {
            url: string;
            width: number;
            dpr: number;
        };
        export interface Recommended extends Node {
            type: "recommended";
            id: string;
            heading?: string;
            teaserTitleOverride?: string;
            teaser?: Teaser;
        }
        export interface RecommendedList extends Node {
            type: "recommended-list";
            heading?: string;
            children: Recommended[];
        }
        export type TeaserConcept = {
            apiUrl: string;
            directType: string;
            id: string;
            predicate: string;
            prefLabel: string;
            type: string;
            types: string[];
            url: string;
        };
        export type Teaser = {
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
        export interface Tweet extends Node {
            id: string;
            type: "tweet";
            html?: string;
        }
        export type FlourishLayoutWidth = Extract<LayoutWidth, "full-grid" | "in-line">;
        export interface Flourish extends Node {
            type: "flourish";
            id: string;
            layoutWidth: FlourishLayoutWidth;
            flourishType: string;
            description?: string;
            timestamp?: string;
            fallbackImage?: Image;
            fragmentIdentifier?: string;
        }
        export interface BigNumber extends Node {
            type: "big-number";
            number: string;
            description: string;
        }
        export interface Video extends Node {
            type: "video";
            id: string;
            title?: string;
        }
        export interface YoutubeVideo extends Node {
            type: "youtube-video";
            url: string;
        }
        export interface ScrollyBlock extends Parent {
            type: "scrolly-block";
            theme: "sans" | "serif";
            children: ScrollySection[];
        }
        export interface ScrollySection extends Parent {
            type: "scrolly-section";
            display: "dark-background" | "light-background";
            noBox?: true;
            position: "left" | "center" | "right";
            transition?: "delay-before" | "delay-after";
            children: [ScrollyImage, ...ScrollyCopy[]];
        }
        export interface ScrollyImage extends Node {
            type: "scrolly-image";
            id: string;
            picture?: ImageSetPicture;
        }
        export interface ScrollyCopy extends Parent {
            type: "scrolly-copy";
            children: (ScrollyHeading | Paragraph)[];
        }
        export interface ScrollyHeading extends Parent {
            type: "scrolly-heading";
            level: "chapter" | "heading" | "subheading";
            children: Text[];
        }
        export interface Layout extends Parent {
            type: "layout";
            layoutName: "auto" | "card" | "timeline";
            layoutWidth: string;
            children: [Heading, LayoutImage, ...LayoutSlot[]] | [Heading, ...LayoutSlot[]] | LayoutSlot[];
        }
        export interface LayoutSlot extends Parent {
            type: "layout-slot";
            children: (Heading | Paragraph | LayoutImage)[];
        }
        export interface LayoutImage extends Node {
            type: "layout-image";
            id: string;
            alt: string;
            caption: string;
            credit: string;
            picture?: ImageSetPicture;
        }
        export type TableColumnSettings = {
            hideOnMobile: boolean;
            sortable: boolean;
            sortType: 'text' | 'number' | 'date' | 'currency' | 'percent';
        };
        export type TableLayoutWidth = Extract<LayoutWidth, 'auto' | 'full-grid' | 'inset-left' | 'inset-right' | 'full-bleed'>;
        export interface TableCaption extends Parent {
            type: 'table-caption';
            children: Phrasing[];
        }
        export interface TableCell extends Parent {
            type: 'table-cell';
            heading?: boolean;
            columnSpan?: number;
            rowSpan?: number;
            children: Phrasing[];
        }
        export interface TableRow extends Parent {
            type: 'table-row';
            children: TableCell[];
        }
        export interface TableBody extends Parent {
            type: 'table-body';
            children: TableRow[];
        }
        export interface TableFooter extends Parent {
            type: 'table-footer';
            children: Phrasing[];
        }
        export interface Table extends Parent {
            type: 'table';
            stripes: boolean;
            compact: boolean;
            layoutWidth: TableLayoutWidth;
            collapseAfterHowManyRows?: number;
            responsiveStyle: 'overflow' | 'flat' | 'scroll';
            children: [TableCaption, TableBody, TableFooter] | [TableCaption, TableBody] | [TableBody, TableFooter] | [TableBody];
            columnSettings: TableColumnSettings[];
        }
        export type CustomCodeComponentAttributes = {
            [key: string]: string | boolean | undefined;
        };
        export interface CustomCodeComponent extends Node {
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
        export interface ImagePair extends Parent {
            type: 'image-pair';
            children: [ImageSet, ImageSet];
        }
        /**
         * Timeline nodes display a timeline of events in arbitrary order.
         */
        export interface Timeline extends Parent {
            type: "timeline";
            /** The title for the timeline */
            title: string;
            children: TimelineEvent[];
        }
        /**
         * TimelineEvent is the representation of a single event in a Timeline.
         */
        export interface TimelineEvent extends Parent {
            type: "timeline-event";
            /** The title of the event */
            title: string;
            /** Any combination of paragraphs and image sets */
            children: (Paragraph | ImageSet)[];
        }
        /**
         * A definition has a term and a related description. It is used to describe a term.
         */
        export interface Definition extends Node {
            type: "definition";
            term: string;
            description: string;
        }
        /**
         * InNumbers represents a set of numbers with related descriptions.
         */
        export interface InNumbers extends Parent {
            type: "in-numbers";
            /** The title for the InNumbers */
            title?: string;
            children: [Definition, Definition, Definition];
        }
        /** Allowed children for a card
        */
        export type CardChildren = ImageSet | Exclude<FormattingBlock, Heading>;
        /**
        * A card describes a subject with images and text
        */
        export interface Card extends Parent {
            type: "card";
            /** The title of this card */
            title?: string;
            children: CardChildren[];
        }
        /**
        * Allowed layout widths for an InfoBox.
        */
        export type InfoBoxLayoutWidth = Extract<LayoutWidth, "in-line" | "inset-left">;
        /**
        * An info box describes a subject via a single card
        */
        export interface InfoBox extends Parent {
            type: "info-box";
            /** The layout width supported by this node */
            layoutWidth: InfoBoxLayoutWidth;
            children: [Card];
        }
        /**
        * InfoPair provides exactly two cards.
        */
        export interface InfoPair extends Parent {
            type: "info-pair";
            /** The title of the info pair */
            title?: string;
            children: [Card, Card];
        }
        /**
           * @sparkGenerateStoryblock true
           **/
        type AudioPlayer = AudioPlayerV1 | AudioPlayerV2 | AudioPlayerV3;
        /** @support deprecated */
        export interface AudioPlayerV1 extends Node {
            type: "audio-player";
            version: 1;
            title: string;
            audioUrl: string;
        }
        export interface AudioPlayerV2 extends Node {
            type: "audio-player";
            version: 2;
            title: string;
            audioId: string;
            audio?: AudioSet;
        }
        /** @support prerelease */
        export interface AudioPlayerV3 extends Node {
            type: "audio-player";
            version: 3;
            title: string;
            audioId: string;
            transcriptionId: string;
            audio?: AudioSet;
            transcription?: Transcription;
        }
        /**
         * Demo placeholders so the AudioPlayer versioning example compiles.
         */
        export interface AudioSet extends Node {
            url: string;
        }
        export interface Transcription extends Node {
            text: string;
        }
        export {};
    }
    export {};
}
