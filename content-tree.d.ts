export declare namespace ContentTree {
    type Block = Node;
    type Phrasing = Text | Break | Strong | Emphasis | Strikethrough | Link;
    interface ImageSource {
        url: string;
        width: number;
        dpr: number;
    }
    interface External<ExternalData> {
        id: string;
        external: true;
        resolved: ExternalData;
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
    interface Root extends Parent {
        type: "root";
        children: [Body];
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
    interface Pullquote extends Parent {
        type: "pullquote";
        text: string;
        source: string;
    }
    interface Recommended extends External<Teaser>, Node {
        type: "recommended";
        heading?: string;
        teaserTitleOverride?: string;
    }
    interface ImageSetContent {
        imageType: "image" | "graphic";
        alt: string;
        caption: string;
        credit: string;
        images: Image[];
        fallbackImage: Image;
    }
    interface ImageSet extends External<ImageSetContent>, Node {
        type: "image-set";
        layoutWidth: "inline" | "article" | "grid" | "viewport";
    }
    interface Image extends Node {
        type: "image";
        id: string;
        originalWidth: number;
        originalHeight: number;
        format: "desktop" | "mobile" | "square" | "standard" | "wide" | "standard-inline";
        binaryUrl: "string";
        sourceSet: ImageSource[];
    }
    interface TweetContent {
        html: string;
    }
    interface Tweet extends External<TweetContent>, Node {
        type: "tweet";
    }
    interface FlourishContent {
        fallbackImage: Image;
    }
    interface Flourish extends External<FlourishContent>, Node {
        type: "flourish";
        layoutWidth: "article" | "grid";
        flourishType: string;
        description: string;
        timestamp: string;
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
