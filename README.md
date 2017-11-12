# MDN

An API for accessing the MDN Javascript documentation. Exposes a single endpoint `/search` which requires one query parameter `q` set to a search query. Returns JSON.

## Example query

### Request
`/search?q=regex.test`

### Result
```json
{
    "ID": 5278,
    "Label": "RegExp.prototype.test()",
    "Locale": "en-US",
    "Modified": "2017-11-06T15:15:06.150712",
    "Slug": "Web/JavaScript/Reference/Global_Objects/RegExp/test",
    "Subpages": [],
    "Summary": "The \u003cstrong\u003e\u003ccode\u003etest()\u003c/code\u003e\u003c/strong\u003e method executes a search for a match between a regular expression and a specified string. Returns \u003ccode\u003etrue\u003c/code\u003e or \u003ccode\u003efalse\u003c/code\u003e.",
    "Tags": ["Reference", "RegExp", "Prototype", "Regular Expressions", "JavaScript", "Method"],
    "Title": "RegExp.prototype.test()",
    "Translations": [],
    "UUID": "383f2015-768e-4e6a-b1ec-7380ed6c17c3",
    "URL": "/en-US/docs/Web/JavaScript/Reference/Global_Objects/RegExp/test"
}
```
