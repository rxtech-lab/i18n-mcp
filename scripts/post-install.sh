#!/bin/bash

# Post-installation script for i18n-mcp
echo "PO Translation MCP Server has been installed successfully!"
echo ""
echo "To use it with Claude Desktop, add the following to your claude_desktop_config.json:"
echo ""
echo '  "mcpServers": {'
echo '    "i18n-mcp": {'
echo '      "command": "/usr/local/bin/i18n-mcp"'
echo '    }'
echo '  }'
echo ""
echo "For more information, visit: https://github.com/rxtech-lab/i18n-mcp"