import test from "node:test"
import assert from "node:assert"
import convert from "@content-tree/from-bodyxml"
import fs from "fs/promises"

let testBase = "../../tests/bodyxml-to-content-tree"
let inputNames = await fs.readdir(`${testBase}/input/`)

for (let inputName of inputNames) {
	let inputPath = `${testBase}/input/${inputName}`
	let outputName = inputName.replace("xml", "json")
	let outputPath = `${testBase}/output/${outputName}`
	let input = await fs.readFile(inputPath, "utf8")
	try {
		let output = await fs.readFile(outputPath, "utf8")
		test(`${inputName} -> ${outputName}`, () => {
			assert.deepStrictEqual(convert(input), JSON.parse(output))
		})
	} catch (error) {
		// couldn't read output, so expecting an error case
		test(`${inputName} -> failure`, () => {
			try {
				convert(input)
				assert.fail("unexpected success")
			} catch (error) {
				assert.ok("threw an error, as expected")
			}
		})
	}
}
