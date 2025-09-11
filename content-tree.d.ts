export declare namespace ContentTree {
    type BodyBlock = Paragraph | Heading | ImageSet | Flourish | BigNumber | CustomCodeComponent | Layout | List | Blockquote | Pullquote | ScrollyBlock | ThematicBreak | Table | Recommended | Tweet | Video | YoutubeVideo | Text | Gallery;
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
        blockIdentifier?: string;
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
        blockIdentifier?: string;
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
    };
    interface Tweet extends Node {
        id: string;
        type: "tweet";
        html: string;
    }
    interface Flourish extends Node {
        type: "flourish";
        id: string;
        layoutWidth: string;
        flourishType: string;
        description?: string;
        timestamp?: string;
        fallbackImage?: Image;
        blockIdentifier?: string;
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
        layoutWidth: 'auto' | 'full-grid' | 'inset-left' | 'inset-right' | 'full-bleed';
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
    type galleryItem = {
        /**
         * @description link for the image
         */
        imageLink?: "text";
        /**
         * @description this is the first Image
         * @default false
         */
        firstImage: boolean;
        /**
         * @description image description
         */
        imageDescription?: string;
        /**
         * @description select or upload image
         */
        picture?: Image;
    };
    /**
     * @sparkGenerateStoryblock true
     */
    interface Gallery extends Node {
        type: "Gallery";
        /**
         * @description gallery description
         * @default default text for the source field
         */
        galleryDescription?: string;
        /**
         * @description autoplay the gallery
         * @default false
         */
        autoPlay?: boolean;
        /**
         * @description each gallery item
         * @maxItems 10
         * @minItems 1
         */
        galleryItems: [galleryItem];
    }
    namespace full {
        type BodyBlock = Paragraph | Heading | ImageSet | Flourish | BigNumber | CustomCodeComponent | Layout | List | Blockquote | Pullquote | ScrollyBlock | ThematicBreak | Table | Recommended | Tweet | Video | YoutubeVideo | Text | Gallery;
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
            blockIdentifier?: string;
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
            blockIdentifier?: string;
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
        };
        interface Tweet extends Node {
            id: string;
            type: "tweet";
            html: string;
        }
        interface Flourish extends Node {
            type: "flourish";
            id: string;
            layoutWidth: string;
            flourishType: string;
            description?: string;
            timestamp?: string;
            fallbackImage?: Image;
            blockIdentifier?: string;
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
            layoutWidth: 'auto' | 'full-grid' | 'inset-left' | 'inset-right' | 'full-bleed';
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
        type galleryItem = {
            /**
             * @description link for the image
             */
            imageLink?: "text";
            /**
             * @description this is the first Image
             * @default false
             */
            firstImage: boolean;
            /**
             * @description image description
             */
            imageDescription?: string;
            /**
             * @description select or upload image
             */
            picture?: Image;
        };
        /**
         * @sparkGenerateStoryblock true
         */
        interface Gallery extends Node {
            type: "Gallery";
            /**
             * @description gallery description
             * @default default text for the source field
             */
            galleryDescription?: string;
            /**
             * @description autoplay the gallery
             * @default false
             */
            autoPlay?: boolean;
            /**
             * @description each gallery item
             * @maxItems 10
             * @minItems 1
             */
            galleryItems: [galleryItem];
        }
    }
    namespace transit {
        type BodyBlock = Paragraph | Heading | ImageSet | Flourish | BigNumber | CustomCodeComponent | Layout | List | Blockquote | Pullquote | ScrollyBlock | ThematicBreak | Table | Recommended | Tweet | Video | YoutubeVideo | Text | Gallery;
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
            blockIdentifier?: string;
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
            blockIdentifier?: string;
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
        };
        interface Tweet extends Node {
            id: string;
            type: "tweet";
        }
        interface Flourish extends Node {
            type: "flourish";
            id: string;
            layoutWidth: string;
            flourishType: string;
            description?: string;
            timestamp?: string;
            blockIdentifier?: string;
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
            layoutWidth: 'auto' | 'full-grid' | 'inset-left' | 'inset-right' | 'full-bleed';
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
        type galleryItem = {
            /**
             * @description link for the image
             */
            imageLink?: "text";
            /**
             * @description this is the first Image
             * @default false
             */
            firstImage: boolean;
            /**
             * @description image description
             */
            imageDescription?: string;
            /**
             * @description select or upload image
             */
            picture?: Image;
        };
        /**
         * @sparkGenerateStoryblock true
         */
        interface Gallery extends Node {
            type: "Gallery";
            /**
             * @description gallery description
             * @default default text for the source field
             */
            galleryDescription?: string;
            /**
             * @description autoplay the gallery
             * @default false
             */
            autoPlay?: boolean;
            /**
             * @description each gallery item
             * @maxItems 10
             * @minItems 1
             */
            galleryItems: [galleryItem];
        }
    }
    namespace loose {
        type BodyBlock = Paragraph | Heading | ImageSet | Flourish | BigNumber | CustomCodeComponent | Layout | List | Blockquote | Pullquote | ScrollyBlock | ThematicBreak | Table | Recommended | Tweet | Video | YoutubeVideo | Text | Gallery;
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
            blockIdentifier?: string;
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
            blockIdentifier?: string;
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
        };
        interface Tweet extends Node {
            id: string;
            type: "tweet";
            html?: string;
        }
        interface Flourish extends Node {
            type: "flourish";
            id: string;
            layoutWidth: string;
            flourishType: string;
            description?: string;
            timestamp?: string;
            fallbackImage?: Image;
            blockIdentifier?: string;
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
            layoutWidth: 'auto' | 'full-grid' | 'inset-left' | 'inset-right' | 'full-bleed';
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
        type galleryItem = {
            /**
             * @description link for the image
             */
            imageLink?: "text";
            /**
             * @description this is the first Image
             * @default false
             */
            firstImage: boolean;
            /**
             * @description image description
             */
            imageDescription?: string;
            /**
             * @description select or upload image
             */
            picture?: Image;
        };
        /**
         * @sparkGenerateStoryblock true
         */
        interface Gallery extends Node {
            type: "Gallery";
            /**
             * @description gallery description
             * @default default text for the source field
             */
            galleryDescription?: string;
            /**
             * @description autoplay the gallery
             * @default false
             */
            autoPlay?: boolean;
            /**
             * @description each gallery item
             * @maxItems 10
             * @minItems 1
             */
            galleryItems: [galleryItem];
        }
    }
}
