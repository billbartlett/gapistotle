# Gapistotle
> **Know Thy Code**

A terminal UI for running Go tests and viewing coverage with rich detail and progressive disclosure.

## Features

### Test Execution & Coverage
- **Package scanning** - Automatically discovers test packages in your project
- **Single package execution** - Run tests for individual packages
- **Run all tests** - Execute all packages at once via menu (F1 → Tests → Run All)
- **Real-time results** - Pass/fail/skip counts with execution timing
- **Coverage analysis** - Statement coverage with function-level granularity
- **Coverage gap analysis** - Identifies untested functions with impact calculations
- **No test caching** - Always runs fresh tests with `-count=1`

### User Interface
- **Progressive disclosure** - Clean summary view, detailed output on demand
- **Three-panel layout** - Package tree, summary, and drill-down details
- **Vim-style navigation** - h/j/k/l keys plus arrow keys
- **Resizable panels** - Adjust layout with `[` and `]` keys
- **Help system** - Press `?` for full keybinding reference

### Theming & Customization
- **Built-in themes** - 7 carefully crafted color schemes
- **Custom themes** - Full theme editor with live preview
- **Color picker** - 48 curated colors to choose from, or enter any hex color
- **XDG support** - Themes stored in `~/.config/gapistotle/themes/`
- **Theme sharing** - Import/export themes via config files

### Configuration
- **Persistent settings** - UI preferences saved to `~/.config/gapistotle/config.conf`
- **Flexible config paths** - CLI flag `-c`, env var, or default location
- **Team workflows** - Share configs via custom paths
- **Structured logging** - JSON logs with configurable levels

## Installation

### From Source

```bash
git clone git@github.com:billbartlett/gapistotle.git
cd gapistotle
go build -o gapistotle .
```

Move the binary to your PATH:
```bash
sudo mv gapistotle /usr/local/bin/
```

### Requirements
- Go 1.21 or later
- Terminal with 256-color support

## Usage

### Basic Usage

```bash
# Scan and test current directory
gapistotle

# Scan specific directory
gapistotle /path/to/project

# Use custom config file
gapistotle -c /path/to/config.conf

# Use config from environment variable
export GAPISTOTLE_CONFIG=/path/to/config.conf
gapistotle
```

### Navigation

**Main Screen:**
- `↑↓` or `j/k` - Navigate package list
- `Enter` - Run tests for selected package
- `Tab` - Switch between left and right panels
- `[` / `]` - Resize left panel
- `t` - Cycle through themes
- `?` - Show help
- `F1` - Toggle menu
- `q` - Quit

**Right Panel (when focused):**
- `↑↓` or `j/k` - Navigate buttons or scroll content
- `Enter` - Select button (TEST DETAILS or COVERAGE GAPS)
- `g` / `G` - Jump to top/bottom
- `PgUp` / `PgDn` - Page up/down
- `ESC` - Return to summary view

**Menu (F1):**
- Settings (placeholder)
- Tests → Run All Tests
- Theme → Select Theme / Edit Theme / Reload Themes
- Help
- Quit

### Theme Editor

1. Press `F1` → Theme → Edit Theme
2. Navigate properties with `↑↓` or `j/k`
3. Press `Enter` to select a property
4. Choose a color:
   - Navigate palette with arrow keys or `h/j/k/l`, press `Enter` to apply
   - Or press `Tab` to switch to hex input, type a hex code (e.g., `#fff`, `#FF5733`), press `Enter` to apply
5. Press `s` to save theme
6. Press `ESC` to cancel

## Configuration

### Config File Format

Default location: `~/.config/gapistotle/config.conf`

```ini
# Theme settings
currentTheme=monokai
themesDirectory=/custom/path/themes  # Optional

# Panel settings
leftPanelWidth=40
minPanelWidth=15
maxPanelWidthPercent=80
panelResizeIncrement=5

# Logging settings
logPath=/tmp/gapistotle.log
logLevel=debug  # debug, info, warn, error
```

### Custom Themes

Themes are stored in `~/.config/gapistotle/themes/` (or `$XDG_CONFIG_HOME/gapistotle/themes/`).

Create a theme file (e.g., `mytheme.conf`):
```ini
name=My Custom Theme
bgColor=#1e1e1e
fgColor=#d4d4d4
menuBgColor=#2d2d30
menuFgColor=#cccccc
# ... more properties
```

Then reload themes via `F1` → Theme → Reload Themes.

## Core Philosophy

- **Everything accessible via menus** - No hidden hotkeys (press F1)
- **Progressive disclosure** - Clean summary by default, detailed output on demand
- **Menu-driven interaction** - "Enter = run the thing you selected"

## Architecture

Built with:
- [Bubbletea](https://github.com/charmbracelet/bubbletea) - Elm architecture TUI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Terminal styling
- Go's built-in `log/slog` - Structured logging

## Contributing

Contributions welcome! Please open an issue or PR.

### Development

```bash
# Build
go build -o gapistotle .

# Run from source
go run . /path/to/test/directory

# Run tests (when we have them)
go test ./...
```

## License

MIT License - see LICENSE file for details

## Acknowledgments

Named after Aristotle, because knowing thy code is the first step to wisdom.

Inspired by k9s and the desire for a better Go test coverage workflow.
