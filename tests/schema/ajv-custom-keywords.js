module.exports = function addCustomKeywords(ajv) {
	ajv.addKeyword({
		keyword: "sparkRepeater",
		schemaType: "boolean",
	})
	ajv.addKeyword({
		keyword: "sparkGenerateStoryblock",
		schemaType: "boolean",
	})
	ajv.addKeyword({
		keyword: "sparkMapNodeType",
		schemaType: "string",
	})
	ajv.addKeyword({
		keyword: "propertyOrder",
		schemaType: "array",
	})
}
