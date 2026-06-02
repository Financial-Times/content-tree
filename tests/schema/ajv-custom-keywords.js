module.exports = function addCustomKeywords(ajv) {
	ajv.addKeyword({
		keyword: "sparkRepeater",
		schemaType: "boolean",
	})
	ajv.addKeyword({
		keyword: "sparkGenerateStoryblock",
		schemaType: "boolean",
	})
}
