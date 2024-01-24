/**
 * @template {ContentTree.Node} Node
 * @param {Node} node
 * @returns {node is ContentTree.Text}
 */
function isText(node) {
	return node.type == "text"
}

/**
 * @template {ContentTree.Node} Node
 * @param {Node} node
 * @returns {node is ContentTree.Parent}
 */
function isParent(node) {
	return "children" in node && Array.isArray(node.children)
}

/**
 * @template {ContentTree.Node} Node
 * @param {Node} node
 * @returns {node is ContentTree.Root}
 */
function isRoot(node) {
	return node.type == "root"
}

/**
 * @typedef {Object} Options
 * @prop {Record<string, (node: ContentTree.node) => string>?} transformers?
 */

let separate = ["heading", "paragraph"]

/**
 * @param {ContentTree.Node} node
 * @param {Options} options
 * @returns {string}
 */
export default function stringify(node, options = {}) {
	options = Object.assign(
		{},
		{
			transformers: {},
		},
		options
	)

	if (node.type in options.transformers) {
		return options.transformers[node.type](node)
	}

	if (isText(node)) {
		return node.value
	}

	if (isRoot(node)) {
		return stringify(node.body, options)
	}

	if (isParent(node)) {
		let content = node.children.map(stringify).join("")
		if (separate.includes(node.type) && content.length) {
			return (content + " ").replace(/\s+/g, " ")
		} else {
			return content
		}
	}

	return ""
}
