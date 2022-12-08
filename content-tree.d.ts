declare namespace ContentTree {
    type Block = Node;
    type Phrasing = Text | Break | Strong | Emphasis | Strikethrough | Link;
    interface Node {
        type: string;
    }
    interface Parent extends Node {
        children: Node[];
    }
    interface Literal extends Node {
        value: any;
    }
    interface Reference extends Node {
        type: "reference";
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
    interface Text extends Literal {
        type: "text";
        value: string;
    }
    interface Break extends Node {
        type: "break";
    }
    interface ThematicBreak extends Node {
        type: "thematicBreak";
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
        type: "listItem";
        children: Phrasing[];
    }
    interface BlockQuote extends Parent {
        type: "blockquote";
        children: Phrasing[];
    }
    interface PullQuote extends Parent {
        type: "pullQuote";
        children: [PullQuoteText, PullQuoteSource];
    }
    interface PullQuoteText extends Parent {
        type: "pullQuoteText";
        children: Text[];
    }
    interface PullQuoteSource extends Parent {
        type: "pullQuoteSource";
        children: Text[];
    }
    interface Recommended extends Parent {
        type: "recommended";
        children: [];
    }
    interface ImageSetReference extends Reference {
        kind: "imageSet";
        imageType: "Image" | "Graphic";
    }
    interface ImageSet extends Node {
        type: "imageSet";
        alt: string;
        caption?: string;
        imageType: "Image" | "Graphic";
        images: Image[];
    }
    interface Image extends Node {
        type: "image";
    }
    interface TweetReference extends Reference {
        kind: "tweet";
    }
    interface Tweet extends Node {
        type: "tweet";
        id: string;
        children: Phrasing[];
    }
    interface FlourishReference extends Reference {
        kind: "flourish";
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
        type: "bigNumber";
        children: [BigNumberNumber, BigNumberDescription];
    }
    interface BigNumberNumber extends Parent {
        type: "bigNumberNumber";
        children: Phrasing[];
    }
    interface BigNumberDescription extends Parent {
        type: "bigNumberDescription";
        children: Phrasing[];
    }
    interface ScrollableBlock extends Parent {
        type: "scrollableBlock";
        theme: "sans" | "serif";
        children: ScrollableSection[];
    }
    interface ScrollableSection extends Parent {
        type: "scrollableSection";
        display: "dark" | "light";
        position: "left" | "centre" | "right";
        transition?: "delay-before" | "delay-after";
        noBox?: boolean;
        children: Array<ImageSet | ScrollableText>;
    }
    interface ScrollableText extends Parent {
        type: "scrollableText";
        style: "text";
        children: Phrasing[];
    }
    interface ScrollableHeading extends Parent {
        type: "scrollableText";
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
        type: "tableHead";
    }
    interface TableBody {
        type: "tableBody";
    }
}
