export declare namespace ContentTree {
    type Block = Node;
    type Phrasing = Text | Break | Strong | Emphasis | Strikethrough | Link;
    interface Node {
        type: string;
    }
    interface Parent extends Node {
        children: Node[];
    }
    interface Reference extends Node {
        type: "reference";
        referencedType: string;
        id: string;
        alt?: string;
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
    interface Chapter extends Parent {
        type: "chapter";
        children: Text[];
    }
    interface Heading extends Parent {
        type: "heading";
        children: Text[];
    }
    interface Subheading extends Parent {
        type: "subheading";
        children: Text[];
    }
    interface Label extends Parent {
        type: "label";
        children: Text[];
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
    interface PullQuote extends Parent {
        type: "pull-quote";
        text: string;
        source: string;
    }
    interface Recommended extends Parent {
        type: "recommended";
        children: [];
    }
    interface ImageSetReference extends Reference {
        referencedType: "image-set";
        imageType: "Image" | "Graphic";
    }
    interface ImageSet extends Node {
        type: "image-set";
        alt: string;
        caption?: string;
        imageType: "Image" | "Graphic";
        images: Image[];
    }
    interface Image extends Node {
        type: "image";
    }
    interface TweetReference extends Reference {
        referencedType: "tweet";
    }
    interface Tweet extends Node {
        type: "tweet";
        id: string;
        html: string;
    }
    interface FlourishReference extends Reference {
        referencedType: "flourish";
        flourishType: string;
    }
    interface Flourish extends Node {
        type: "flourish";
        id: string;
        layoutWidth: "" | "full-grid";
        flourishType: string;
        description: string;
        fallbackImage: Image;
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
    interface ScrollableBlock extends Parent {
        type: "scrollable-block";
        theme: "sans" | "serif";
        children: ScrollableSection[];
    }
    interface ScrollableSection extends Parent {
        type: "scrollable-section";
        display: "dark" | "light";
        position: "left" | "centre" | "right";
        transition?: "delay-before" | "delay-after";
        noBox?: boolean;
        children: Array<ImageSet | ScrollableText>;
    }
    interface ScrollableText extends Parent {
        type: "scrollable-text";
        style: "text";
        children: Phrasing[];
    }
    interface ScrollableHeading extends Parent {
        type: "scrollable-text";
        style: "chapter" | "heading" | "subheading";
        children: Text[];
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
