a library (and cli) for creating a (transit) content-tree from a bodyXML

## Testing

### Generate a bodyTree from a content api URL

```
curl "$CONTENT_API_HOST/internalcontent/9e9192ba-02f1-4034-aaf0-200270d0f6d7?apiKey=$CONTENT_API_READ_KEY" | jq '.bodyXML' -r | node libraries/from-bodyxml/cli.js
```

### Validate from-bodyxml's output against the schema

Ensure you have the environment variables:

- CONTENT_API_HOST
- CONTENT_API_READ_KEY

```
 node libraries/from-bodyxml/validate.js  9e9192ba-02f1-4034-aaf0-200270d0f6d7
```

### Validate articles published in the last 5 minutes

Ensure you have the environment variables:

- CONTENT_API_HOST
- CONTENT_API_READ_KEY

```
node --test libraries/from-bodyxml/smoke-test.js
```
