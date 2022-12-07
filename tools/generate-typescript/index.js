import {unified} from "unified"
import {stream} from "unified-stream"
import remarkParse from "remark-parse"
import remarkGfm from "remark-gfm"
import {createReadStream} from "fs"

let filename = process.argv[2]
var input =
	typeof filename == "string" ? createReadStream(filename) : process.stdin

process.stdout.write("namespace ContentTree {\n")
input
	.pipe(
		stream(
			unified()
				.use(remarkParse)
				.use(remarkGfm)
				.use(function extract() {
					Object.assign(this, {Compiler: walk})
					function walk(node) {
						if (node.type == "code") {
							process.stdout.write(
								node.value
									.replace(/^(interface|type)/gm, term => `export ${term}`)
									.replace(/^/gm, "\t") + "\n\n"
							)
						} else if ("children" in node && Array.isArray(node.children)) {
							for (let child of node.children) {
								walk(child)
							}
						}
					}
				})
		)
	)
	.on("end", () => {
		process.stdout.write("}\n")
		process.exit()
	})
