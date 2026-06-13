#!/bin/bash

cat <<'EOF'
# Project README

## Overview

This is a **demonstration** of the `markdown` output type in `tui-workflow`.

## Features

- Beautiful rendering with *Glamour*
- Support for `code blocks`
- Tables, lists, and more

## Example Table

| Feature | Status |
|---------|--------|
| Markdown | ✅ Working |
| Text     | ✅ Working |
| Styles   | ✅ Dark theme |

## Code Example

```bash
#!/bin/bash
echo "Hello, world!"
```

> This output was rendered from a shell script using the `output_type: markdown` setting in the workflow JSON.

EOF