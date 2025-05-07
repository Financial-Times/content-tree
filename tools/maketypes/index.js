import {createReadStream} from "fs"
import {fromMarkdown} from "mdast-util-from-markdown"
import tangle from "*tangle*"

let filename = process.argv[2]
var input =
	typeof filename == "string" ? createReadStream(filename) : process.stdin
const chunks = []
for await (let chunk of input) {
	chunks.push(chunk)
}
const doc = fromMarkdown(Buffer.concat(chunks))
const code = tangle(doc).replace(
	/^(interface|type)/gm,
	term => `export ${term}`
)

// in the full tree, externals are mandatory
const full = code.replace(/^(\s+)external (.+)$/gm, "$1$2")
// in the transit tree, externals must not be present
const transit = code.replace(/^\s+external (.+:).+$/gm, "")
// in the loose tree, externals are optional
const loose = code
	.replace(/^(\s+)external (.+)\?:(.+)$/gm, "$1$2?:$3")
	.replace(/^(\s+)external (.+):(.+)$/gm, "$1$2?:$3")
process.stdout.write("export namespace ContentTree {\n")
// make content-tree nodes available on the root namespace
process.stdout.write(full.replace(/^/gm, "\t"))

// also make content-tree nodes available on ContentTree.full
process.stdout.write("export namespace full {\n")
process.stdout.write(full.replace(/^/gm, "\t\t"))
process.stdout.write("\n}\n")

// make the transit tree nodes available on ContentTree.transit
process.stdout.write("export namespace transit {\n")
process.stdout.write(transit.replace(/^/gm, "\t\t"))
process.stdout.write("\n}\n")

// make the loose tree nodes available on ContentTree.loose
process.stdout.write("export namespace loose {\n")
process.stdout.write(loose.replace(/^/gm, "\t\t"))
process.stdout.write("\n}\n")

process.stdout.write("\n}\n")
