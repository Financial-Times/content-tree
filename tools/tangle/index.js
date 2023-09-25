import {createReadStream} from "fs"
import {fromMarkdown} from "mdast-util-from-markdown"
import * as url from "node:url"

export default function tangle(doc) {
	function walk(node) {
		if (node.type == "code") {
			return node.value
		} else if ("children" in node && Array.isArray(node.children)) {
			for (let child of node.children) {
				walk(child)
			}
		}
	}

	let blocks = []
	for (let child of doc.children) {
		let code = walk(child)
		if (code) {
			blocks.push(code)
		}
	}

	return blocks.join("\n")
}

if (import.meta.url.startsWith("file:")) {
	if (process.argv[1] == url.fileURLToPath(import.meta.url)) {
		let filename = process.argv[2]
		var input =
			typeof filename == "string" ? createReadStream(filename) : process.stdin
		const chunks = []
		for await (let chunk of input) {
			chunks.push(chunk)
		}
		const doc = fromMarkdown(Buffer.concat(chunks))
		process.stdout.write(tangle(doc) + "\n")
	}
}
