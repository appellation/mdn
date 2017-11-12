# MDN

An API for accessing the MDN Javascript documentation. Exposes a single endpoint `/search` which requires one query parameter `q` set to a search query. Returns JSON.

## Example query

### Request
`/search?q=string.prototype.match`

### Result
```json
{
    "Locale": "en-US",
    "Slug": "Web/JavaScript/Reference/Global_Objects/String/match",
    "Title": "String.prototype.match()",
    "URL": "/en-US/docs/Web/JavaScript/Reference/Global_Objects/String/match",
    "Subpages": []
}
```
