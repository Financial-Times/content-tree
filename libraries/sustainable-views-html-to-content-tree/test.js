import t from "tap"
import convert from "./index.js"
import fs from "fs/promises"

let testBase = "./tests/"
let names = await fs.readdir(`${testBase}/input/`)

for (let name of names) {
	let inputPath = `${testBase}/input/${name}`
	let outputPath = `${testBase}/output/${name}`
	let inputText = await fs.readFile(inputPath, "utf8")
	let input = JSON.parse(inputText)
	console.log(convert(input))
	try {
		let output = await fs.readFile(outputPath, "utf8")
		t.test(`${name.replace("json", "wp")} -> ${name.replace("json", "content-tree")}`, t => {
			t.equal(convert(input).trim(), output.trim())
			t.end()
		})
	} catch (error) {
		// couldn't read output, so expecting an error case
		t.test(`${name} -> failure`, t => {
			try {
				convert(input)
				t.fail("unexpected success")
			} catch (error) {
				t.pass("threw an error, as expected")
			}
			t.end()
		})
	}
}
