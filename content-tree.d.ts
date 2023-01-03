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
    interface Recommended extends Parent {
        type: "recommended";
        children: [];
    }
    interface ImageSetReference extends Reference {
        referencedType: "image-set";
    }
    interface ImageSet extends Node {
        type: "image-set";
        imageType: "graphic" | "image";
        alt: string;
        caption: string;
        credit: string;
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
    interface ScrollyBlock extends Parent {
        type: "scrolly-block";
        theme: "sans" | "serif";
        children: ScrollySection[];
    }
    interface ScrollySection extends Parent {
        type: "scrolly-section";
        copyTheme: "dark-text" | "light-text" | "dark-text-no-box" | "light-text-no-box";
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
