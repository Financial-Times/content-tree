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
 * @template {UNode} NodeType
 * @typedef {(el: import("xast").Element) => TransNode<NodeType>} Transformer
 */

/**
 * @template {UNode | UParent} Node
 * @typedef {Omit<Node, "children"> & (Node extends UParent ? {children?: Node["children"]} : {children?: null})} TransNode
 */

export let defaultTransformers = {
	/**
	 * @type {Transformer<ContentTree.transit.Heading>}
	 */
	h1(h1) {
		return {
			type: "heading",
			level: "chapter",
		}
	},
	/**
	 * @type {Transformer<ContentTree.transit.Heading>}
	 */
	h2(h2) {
		return {
			type: "heading",
			level: "subheading",
		}
	},
	/**
	 * @type {Transformer<ContentTree.transit.Heading>}
	 */
	h4(h4) {
		return {
			type: "heading",
			level: "label",
		}
	},
	/**
	 * @type {Transformer<ContentTree.transit.Paragraph>}
	 */
	p(p) {
		return {
			type: "paragraph",
		}
	},
	/**
	 * @type {Transformer<ContentTree.transit.Emphasis>}
	 */
	em(em) {
		return {
			type: "emphasis",
		}
	},
	/**
	 * @type {Transformer<ContentTree.transit.Strong>}
	 */
	strong(strong) {
		return {
			type: "strong",
		}
	},
	/**
	 * @type {Transformer<ContentTree.transit.Strikethrough>}
	 */
	s(s) {
		return {
			type: "strikethrough",
		}
	},
	/**
	 * @type {Transformer<ContentTree.transit.Break>}
	 */
	br(br) {
		return {
			type: "break",
		}
	},
	/**
	 * @type {Transformer<ContentTree.transit.ThematicBreak>}
	 */
	hr(hr) {
		return {
			type: "thematic-break",
		}
	},
	/**
	 * @type {Transformer<ContentTree.transit.Link | ContentTree.transit.YoutubeVideo>}
	 */
	a(a) {
		if(a.attributes['data-asset-type'] === 'video') {
			const url = a.attributes.href ?? '';
			if(url.includes('youtube.com')) {
				return /** @type {ContentTree.transit.YoutubeVideo} */({
					type: "youtube-video",
					url: url,
					children: null
				})
			}
			//TODO: specialist support Vimeo, but this isn't in the Content Tree spec yet
		}
		return /** @type {ContentTree.transit.Link} */({
			type: "link",
			title: a.attributes.title ?? "",
			url: a.attributes.href ?? "",
		})
	},
	/**
	 * @type {Transformer<ContentTree.transit.List>}
	 */
	ol(ol) {
		return {
			type: "list",
			ordered: true,
		}
	},
	/**
	 * @type {Transformer<ContentTree.transit.List>}
	 */
	ul(ul) {
		return {
			type: "list",
			ordered: false,
		}
	},
	/**
	 * @type {Transformer<ContentTree.transit.ListItem>}
	 */
	li(li) {
		return {
			type: "list-item",
		}
	},
	/**
	 * @type {Transformer<ContentTree.transit.Blockquote>}
	 */
	blockquote(blockquote) {
		return {
			type: "blockquote",
		}
	},
	/**
	 * @type {Transformer<ContentTree.transit.Pullquote>}
	 */
	["pull-quote"](pq) {
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
	 * @type {Transformer<ContentTree.transit.BigNumber>}
	 */
	["big-number"](bn) {
		let number = find(bn, {name: "big-number-headline"})
		let description = find(bn, {name: "big-number-intro"})
		return {
			type: "big-number",
			number: number ? xastToString(number) : "",
			description: description ? xastToString(description) : "",
			children: null,
		}
	},
	/**
	 * @type {Transformer<ContentTree.transit.LayoutImage>}
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
	 * @type {Transformer<ContentTree.transit.ImageSet>}
	 */
	[ContentType.imageset](content) {
		return {
			type: "image-set",
			id: content.attributes.url ?? "",
			children: null,
		}
	},
	/**
	 * @type {Transformer<ContentTree.transit.Video>}
	 */
	[ContentType.video](content) {
		return {
			type: "video",
			id: content.attributes.url ?? "",
			embedded: content.attributes["data-embedded"] == "true" ? true : false,
			children: null,
		}
	},
	// TODO these two Link transforms may be wrong. what is a "content" or an "article"?
	/**
	 * @type {Transformer<ContentTree.transit.Flourish | ContentTree.transit.Link>}
	 */
	[ContentType.content](content) {
		if (content.attributes["data-asset-type"] == "flourish") {
			return /** @type {ContentTree.transit.Flourish} */ ({
				type: "flourish",
				flourishType: content.attributes["data-flourish-type"] || "",
				layoutWidth: content.attributes["data-layout-width"] || "",
				description: content.attributes["alt"] || "",
				timestamp: content.attributes["data-time-stamp"] || "",
				// fallbackImage -- TODO should this be external in content-tree?
			})
		}
		const id = content.attributes.url ?? "";
		const uuid = id.split('/').pop();
		return /** @type {ContentTree.transit.Link} */({
			type: "link",
			url: `https://www.ft.com/content/${uuid}`,
			title: content.attributes.dataTitle ?? "",
		})
	},
	/**
	 * @type {Transformer<ContentTree.transit.Link>}
	 */
	[ContentType.article](content) {
		const id = content.attributes.url ?? "";
		const uuid = id.split('/').pop();
		return {
			type: "link",
			url: `https://www.ft.com/content/${uuid}`,
			title: content.attributes.dataTitle ?? "",
		}
	},
	/**
	 * @type {Transformer<ContentTree.transit.Recommended>}
	 */
	recommended(rl) {
		const link = find(rl, { name: 'ft-content'});
		const heading = find(rl, { name: 'recommended-title'});
		return {
			type: "recommended",
			id: link?.attributes?.url ?? "",
			heading: heading ? xastToString(heading) : "",
			teaserTitleOverride: link ? xastToString(link) : "",
			children: null
		}
	},
}

/**
 * @param {import("xast").Node} node
 * @returns {node is import("xast").Element}
 */
function isXElement(node) {
	return node.type == "element"
}

/**
 * @param {import("xast").Node} node
 * @returns {node is import("xast").Text}
 */
function isXText(node) {
	return node.type == "text"
}

/**
 * @param {import("xast").Node} node
 * @returns {node is import("xast").Root}
 */
function isXRoot(node) {
	return node.type == "root"
}

/**
 * @param {import("xast").Node} bodyxast
 * @returns {ContentTree.transit.Root}
 */
export function fromXast(bodyxast, transformers = defaultTransformers) {
	return (function walk(xmlnode) {
		if (isXRoot(xmlnode)) {
			return {
				type: "root",
				body: {
					type: "body",
					children: xmlnode.children[0].children.map(walk),
				},
			}
		} else if (isXElement(xmlnode)) {
			// i thought about this solution for no more than 5 seconds
			if (xmlnode.name == "experimental") {
				return xmlnode.children.map(walk)
			}
			let transformer =
				(xmlnode.name == "content" || xmlnode.name == "ft-content")
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
		} else if (isXText(xmlnode)) {
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
