# Po translation mcp server

Your code structure should use /example folder as reference. mcp tools are located in `/internal/tools` folder

## Tools

these are the tools you need to implement

1. `listAllPoFiles` will return list of the `po` files in the given dir using `internal/utils/scan.go` for scanning
2. `getUntranslatedTerms` will return n number of untranslated terms. It takes the message.po file path and number of untranslated terms per return using `internal/service/po_service.go`
3. `lookUpTranslation` will search the term key and return the translated value. Take the po file path and number of results return per page. Default to 10
4. `translate` will translate the term by its translation. Takes po file path and map of translations where key is the terms key and value is the translation. After each translation, should write back the translation to the path using `internal/service/po_service.go`

## Testing

Each tool should have its own test

## Logging

Use log pkg for printing and default no log.

## CI

use example/ as reference to setup the auto build github ci script. Take a look at /example/scripts to see how to build macOS pkg. Only macOS for now. No build for other platform.
