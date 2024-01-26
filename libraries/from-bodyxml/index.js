import {fromXml as xastFromXml} from "xast-util-from-xml"
import {toString as xastToString} from "xast-util-to-string"
import {find} from "unist-util-find"

let ContentType = {
	imageset: "http://www.ft.com/ontology/content/ImageSet",
	video: "http://www.ft.com/ontology/content/Video",
	content: "http://www.ft.com/ontology/content/Content",
	article: "http://www.ft.com/ontology/content/Article",
}

/**
 * @typedef {import("unist").Parent} UParent
 * @typedef {import("unist").Node} UNode
 */

/**
 * @typedef {import("xast").Node} XNode
 */

/**
 * @template {UNode & UParent} Node
 * @typedef {Omit<Node, "children"> & (Node extends UParent ? {children?: Node["children"]} : {children?: null})} TransNode
 */

export let defaultTransformers = {
	/**
	 * @param {import("xast").Element} p
	 * @returns {TransNode<ContentTree.transit.Paragraph>}
	 */
	p(p) {
		return {
			type: "paragraph",
		}
	},
	/**
	 * @param {import("xast").Element} h1
	 * @returns {TransNode<ContentTree.transit.Heading>}
	 */
	h1(h1) {
		return {
			type: "heading",
			level: "chapter",
		}
	},
	/**
	 * @param {import("xast").Element} h2
	 * @returns {TransNode<ContentTree.transit.Heading>}
	 */
	h2(h2) {
		return {
			type: "heading",
			level: "subheading",
		}
	},
	/**
	 * @param {import("xast").Element} h4
	 * @returns {TransNode<ContentTree.transit.Heading>}
	 */
	h4(h4) {
		return {
			type: "heading",
			level: "label",
		}
	},
	/**
	 * @param {import("xast").Element} pq
	 * @returns {TransNode<ContentTree.transit.Pullquote>}
	 */
	["pull-quote"](pq) {
		// find children pull-quote-text
		// find children pull-quote-source
		let text = find(pq, {name: "pull-quote-text"})
		let source = find(pq, {name: "pull-quote-source"})
		return {
			type: "pullquote",
			text: text ? xastToString(text) : "",
			source: source ? xastToString(source) : "",
			children: null,
		}
	},
	/**
	 * @param {import("xast").Element} img
	 * @returns {TransNode<ContentTree.transit.LayoutImage>}
	 */
	img(img) {
		return {
			type: "layout-image",
			id: img.attributes.src ?? "",
			credit: img.attributes["data-copyright"]?.replace(/^© /, "") ?? "",
			// todo this can't be right
			alt: img.attributes.longdesc ?? "",
			caption: img.attributes.longdesc ?? "",
		}
	},
	/**
	 * @param {import("xast").Element} em
	 * @returns {TransNode<ContentTree.transit.Emphasis>}
	 */
	em(em) {
		return {
			type: "emphasis",
		}
	},
	/**
	 * @param {import("xast").Element} strong
	 * @returns {TransNode<ContentTree.transit.Strong>}
	 */
	strong(strong) {
		return {
			type: "strong",
		}
	},
	/**
	 * @param {import("xast").Element} s
	 * @returns {TransNode<ContentTree.transit.Strikethrough>}
	 */
	s(s) {
		return {
			type: "strikethrough",
		}
	},
	/**
	 * @param {import("xast").Element} content
	 * @returns {TransNode<ContentTree.transit.ImageSet>}
	 */
	[ContentType.imageset](content) {
		return {
			type: "image-set",
			id: content.attributes.id ?? "",
			children: null,
		}
	},
	/**
	 * @param {import("xast").Element} content
	 * @returns {TransNode<ContentTree.transit.Video>}
	 */
	[ContentType.video](content) {
		return {
			type: "video",
			id: content.attributes.id ?? "",
			embedded: content.attributes["data-embedded"] == "true" ? true : false,
			children: null,
		}
	},
	/**
	 * @param {import("xast").Element} content
	 * @returns {TransNode<ContentTree.transit.Link>}
	 */
	[ContentType.content](content) {
		return {
			type: "link",
			url: `https://www.ft.com/content/${content.attributes.id}`,
			title: content.attributes.dataTitle ?? "",
		}
	},
	/**
	 * @param {import("xast").Element} content
	 * @returns {TransNode<ContentTree.transit.Link>}
	 */
	[ContentType.article](content) {
		return {
			type: "link",
			url: `https://www.ft.com/content/${content.attributes.id}`,
			title: content.attributes.dataTitle ?? "",
		}
	},
}

/**
 * @param {import("xast").Node} node
 * @returns {node is import("xast").Element}
 */
function isElement(node) {
	return node.type == "element"
}

/**
 * @param {import("xast").Node} node
 * @returns {node is import("xast").Text}
 */
function isText(node) {
	return node.type == "text"
}

/**
 * @param {import("xast").Node} node
 * @returns {node is import("xast").Root}
 */
function isRoot(node) {
	return node.type == "root"
}

/**
 * @param {import("xast").Node} bodyxast
 * @returns {ContentTree.transit.Root}
 */
export function fromXast(bodyxast, transformers = defaultTransformers) {
	return (function walk(xmlnode) {
		if (isRoot(xmlnode)) {
			return {
				type: "root",
				body: {
					type: "body",
					children: xmlnode.children[0].children.map(walk),
				},
			}
		} else if (isElement(xmlnode)) {
			// i thought about this solution for no more than 5 seconds
			if (xmlnode.name == "experimental") {
				return xmlnode.children.map(walk)
			}
			let transformer =
				xmlnode.name == "content"
					? String(xmlnode.attributes.type)
					: xmlnode.name

			if (transformer in transformers) {
				let ctnode = transformers[transformer](xmlnode)
				if ("children" in ctnode && ctnode.children === null) {
					// this is how we indicate we shouldn't iterate, but this thing
					// shouldn't have any children
					delete ctnode.children
					return ctnode
				} else if ("children" in ctnode && Array.isArray(ctnode.children)) {
					return ctnode
				} else if ("children" in xmlnode) {
					return {
						...ctnode,
						// this is a flatmap because of <experimental/>
						children: xmlnode.children.flatMap(walk),
					}
				}
				return ctnode
			} else {
				return {type: "__UNKNOWN__"}
			}
		} else if (isText(xmlnode)) {
			return {
				type: "text",
				value: xmlnode.value,
			}
		} else {
			return {type: "__UNKNOWN__"}
		}
	})(bodyxast)
}

/** @param {string} bodyxml */
export function fromXML(bodyxml) {
	return fromXast(xastFromXml(bodyxml))
}

export default fromXML
