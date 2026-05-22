module.exports = function addCustomKeywords(ajv) {
	ajv.addKeyword({
		keyword: "sparkRepeater",
		schemaType: "boolean",
	})
}
