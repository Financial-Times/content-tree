import t from "tap"
import {u} from "unist-builder"
import stringify from "@content-tree/content-tree-to-string"
import fs from "fs/promises"

let testBase = "../../tests/content-tree-to-string"
let inputNames = await fs.readdir(`${testBase}/input/`)

for (let inputName of inputNames) {
	let inputPath = `${testBase}/input/${inputName}`
	let outputName = inputName.replace(/json$/, "text")
	let outputPath = `${testBase}/output/${outputName}`
	let inputText = await fs.readFile(inputPath, "utf8")
	let input = JSON.parse(inputText)
	try {
		let output = await fs.readFile(outputPath, "utf8")
		t.test(`${inputName} -> ${outputName}`, t => {
			t.equal(stringify(input).trim(), output.trim())
			t.end()
		})
	} catch (error) {
		// couldn't read output, so expecting an error case
		t.test(`${inputName} -> failure`, t => {
			try {
				stringify(input)
				t.fail("unexpected success")
			} catch (error) {
				t.pass("threw an error, as expected")
			}
			t.end()
		})
	}
}

t.test("Supports a custom transformer", t => {
	let result = stringify(
		{
			type: "media",
			alt: "hello :)",
			credit: "chee",
			caption: "cool beans",
		},
		{
			transformers: {
				media: node => {
					return `${node.alt} ${node.caption} © ${node.credit}`
				},
			},
		}
	)
	t.equal("hello :) cool beans © chee", result)
	t.end()
})
