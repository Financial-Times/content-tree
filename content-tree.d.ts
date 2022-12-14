export declare namespace ContentTree {
    type Block = Node;
    type Phrasing = Text | Break | Strong | Emphasis | Strikethrough | Link;
    interface ImageSource {
        url: string;
        width: number;
        dpr: number;
    }
    interface TeaserConcept {
        apiUrl: string;
        directType: string;
        id: string;
        predicate: string;
        prefLabel: string;
        type: string;
        types: string[];
        url: string;
    }
    interface TeaserImage {
        url: string;
        width: number;
        height: number;
    }
    interface Indicators {
        accessLevel: "premium" | "subscribed" | "registered" | "free";
        isOpinion?: boolean;
        isColumn?: boolean;
        isPodcast?: boolean;
        isEditorsChoice?: boolean;
        isExclusive?: boolean;
        isScoop?: boolean;
    }
    interface Teaser {
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
        indicators: Indicators;
        image: Image;
    }
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
        children: Block[];
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
        children: Phrasing[];
    }
    interface Blockquote extends Parent {
        type: "blockquote";
        children: Phrasing[];
    }
    interface Pullquote extends Node {
        type: "pullquote";
        text: string;
        source?: string;
    }
    interface Recommended extends Node {
        type: "recommended";
        id: string;
        heading?: string;
        teaserTitleOverride?: string;
        teaser?: Teaser;
    }
    interface ImageSet extends Node {
        type: "image-set";
        id: string;
        layoutWidth: "inline" | "article" | "grid" | "viewport";
        picture?: {
            imageType: "image" | "graphic";
            alt: string;
            caption: string;
            credit: string;
            images: Image[];
            fallbackImage: Image;
        };
    }
    interface Image extends Node {
        type: "image";
        id: string;
        originalWidth: number;
        originalHeight: number;
        format: "desktop" | "mobile" | "square" | "standard" | "wide" | "standard-inline";
        binaryUrl: string;
        sourceSet: ImageSource[];
    }
    interface Tweet extends Node {
        id: string;
        type: "tweet";
        html?: string;
    }
    interface Flourish extends Node {
        type: "flourish";
        id: string;
        layoutWidth: "article" | "grid";
        flourishType: string;
        description?: string;
        timestamp?: string;
        fallbackImage?: Image;
    }
    interface BigNumber extends Parent {
        type: "big-number";
        children: [BigNumberNumber, BigNumberDescription];
    }
    interface BigNumberNumber extends Parent {
        type: "big-number-number";
        children: Phrasing[];
    }
    interface BigNumberDescription extends Parent {
        type: "big-number-description";
        children: Phrasing[];
    }
    interface ScrollyBlock extends Parent {
        type: "scrolly-block";
        theme: "sans" | "serif";
        children: ScrollySection[];
    }
    interface ScrollySection extends Parent {
        type: "scrolly-section";
        theme: "dark-text" | "light-text" | "dark-text-no-box" | "light-text-no-box";
        position: "left" | "center" | "right";
        transition?: "delay-before" | "delay-after";
        children: [ImageSet, ...ScrollyCopy[]];
    }
    interface ScrollyCopy extends Parent {
        type: "scrolly-copy";
        children: ScrollyText[];
    }
    interface ScrollyText extends Parent {
        type: "scrolly-text";
        level: string;
    }
    interface ScrollyHeading extends ScrollyText {
        type: "scrolly-text";
        level: "chapter" | "heading" | "subheading";
        children: Text[];
    }
    interface ScrollyParagraph extends ScrollyText {
        type: "scrolly-text";
        level: "text";
        children: Phrasing[];
    }
    interface Layout extends Parent {
        type: "layout";
        layoutName: "auto" | "card" | "timeline";
        layoutWidth: "inset-left" | "full-width" | "full-grid";
        children: [Heading, ...LayoutSlot[]] | LayoutSlot[];
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
        picture?: Image;
    }
    interface Table extends Parent {
        type: "table";
        children: [Caption | TableHead | TableBody];
    }
    interface Caption {
        type: "caption";
    }
    interface TableHead {
        type: "table-head";
    }
    interface TableBody {
        type: "table-body";
    }
}
