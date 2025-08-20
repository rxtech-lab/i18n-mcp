# i18n-mcp

PO Translation MCP Server - A Model Context Protocol server for managing .po translation files.

## Features

This MCP server provides tools for working with PO (Portable Object) translation files:

- **listAllPoFiles**: Scan directories to find all .po files
- **getUntranslatedTerms**: Get untranslated terms from a PO file
- **lookUpTranslation**: Search for translations in a PO file
- **translate**: Add or update translations in a PO file

## Installation

### From Source

```bash
git clone https://github.com/rxtech-lab/i18n-mcp.git
cd i18n-mcp
make build
make install-local
```

### macOS Package

Download the latest `.pkg` file from the [releases page](https://github.com/rxtech-lab/i18n-mcp/releases) and double-click to install.

## Configuration

Add the following to your Claude Desktop configuration file (`claude_desktop_config.json`):

```json
{
  "mcpServers": {
    "i18n-mcp": {
      "command": "/usr/local/bin/i18n-mcp"
    }
  }
}
```

## Usage

Once configured, the MCP server will be available in Claude Desktop. You can use the following tools:

### List PO Files
Find all .po files in a directory:
```
Use the listAllPoFiles tool to scan /path/to/translations
```

### Get Untranslated Terms
Get untranslated terms from a PO file:
```
Use getUntranslatedTerms on /path/to/messages.po with limit 10
```

### Search Translations
Search for specific terms in a PO file:
```
Use lookUpTranslation to find "hello" in /path/to/messages.po
```

### Add Translations
Add or update translations:
```
Use translate tool to add translations to /path/to/messages.po:
- "hello": "hola"
- "goodbye": "adiós"
```

## Development

### Requirements

- Go 1.22 or later
- Make

### Building

```bash
make build
```

### Testing

```bash
make test
```

### Running Locally

```bash
make run
```

## License

MIT