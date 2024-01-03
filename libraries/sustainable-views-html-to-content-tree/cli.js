// todo use package name
import convert from "./index.js"

import {createReadStream} from "fs"
let path = process.argv[2]

let stream = path && path != "-" ? createReadStream(path) : process.stdin

// goodbye memory
let json = await new Promise((yay, boo) => {
	let string = ""
	stream.on("data", data => string += data.toString("utf8"))
	stream.on("end", () => yay(string))
	stream.on("error", boo)
})

process.stdout.write(JSON.stringify(convert(JSON.parse(json)), null, "\t"))
