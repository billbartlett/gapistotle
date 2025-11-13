# Gapistotle
> **Know Thy Code**
> Version 0.1.0

terminal UI for running Go tests and viewing coverage. basically, `go test -v -cover` but with an actual interface that doesn't make you want to cry.

## Features

### Test Execution & Coverage
- **package scanning** - finds all your test packages automatically
- **single package execution** - run tests for one package at a time
- **run all tests** - execute everything sequentially (` → Tests → Know It All)
- **directory-specific test modes** - set unit/integration/all per directory and it remembers
- **unit vs integration separation** - separate sections in test details so you can actually see what's what
- **time breakdown** - shows actual test execution time vs setup/overhead time (because testcontainers taking 5 seconds while tests run in 0.8s is confusing without context)
- **slowest-first sorting** - passed tests sorted by duration so you can spot the slow ones immediately
- **coverage analysis** - statement coverage with function-level granularity
- **coverage gap analysis** - identifies untested functions with impact calculations
- **no test caching** - always runs fresh with `-count=1`

### User Interface
- **progressive disclosure** - clean summary by default, details when you need them
- **full-screen mode** - press `f` to view test results without the left panel (easier copy/paste)
- **three-panel layout** - package tree, summary, and drill-down details
- **vim-style navigation** - hjkl keys plus arrow keys because muscle memory
- **resizable panels** - `[` and `]` to adjust
- **help system** - press `?` if you forget the keybindings

### Theming & Customization
- **built-in themes** - 8 color schemes included (gapistotle, dracula, nord, monokai, atom-one-dark, desert, industry, material)
- **fallback theme** - gapistotle theme is hardcoded as a fallback if theme files aren't available
- **custom themes** - full editor with live preview
- **color picker** - 48 curated colors, or enter any hex you want
- **XDG support** - themes stored in `~/.config/gapistotle/themes/`
- **theme sharing** - just copy the .conf files around

### Configuration
- **persistent settings** - UI preferences saved to `~/.config/gapistotle/config.conf`
- **directory-specific test modes** - each directory remembers unit/integration/all setting
- **flexible config paths** - `-c` flag, env var, or default location
- **team workflows** - share configs via custom paths if you're into that
- **structured logging** - JSON logs with configurable levels

## Installation

### from source

```bash
git clone git@github.com:billbartlett/gapistotle.git
cd gapistotle
go build -o gapistotle .
```

move the binary somewhere in your PATH:
```bash
sudo mv gapistotle /usr/local/bin/
```

optionally, copy theme files to your config directory:
```bash
mkdir -p ~/.config/gapistotle/themes
cp themes/*.conf ~/.config/gapistotle/themes/
```

### requirements
- go 1.21 or later
- terminal with 256-color support

## Usage

### basic usage

```bash
# scan and test current directory
gapistotle

# scan specific directory
gapistotle /path/to/project

# show version information
gapistotle --version
gapistotle -v

# use custom config file
gapistotle -c /path/to/config.conf

# use config from environment variable
export GAPISTOTLE_CONFIG=/path/to/config.conf
gapistotle
```

### navigation

**main screen:**
- `↑↓` or `j/k` - navigate package list
- `Enter` - run tests for selected package
- `Tab` - switch between left and right panels
- `[` / `]` - resize left panel
- `t` - cycle through themes
- `?` - show help
- `` ` `` - toggle menu (backtick key)
- `q` - quit

**right panel (when focused):**
- `↑↓` or `j/k` - navigate buttons or scroll content
- `Enter` - select button (TEST DETAILS or COVERAGE GAPS)
- `f` - full-screen mode (shows whichever view is highlighted)
- `g` / `G` - jump to top/bottom
- `PgUp` / `PgDn` - page up/down
- `ESC` - return to summary view

**menu (` - backtick key):**
- settings (placeholder)
- tests → Know It All / Test Mode
- theme → Select Theme / Edit Theme / Reload Themes
- help
- quit

### theme editor

1. press `` ` `` → Theme → Edit Theme
2. navigate properties with `↑↓` or `j/k`
3. press `Enter` to select a property
4. choose a color:
   - navigate palette with arrow keys or `h/j/k/l`, press `Enter` to apply
   - or press `Tab` to switch to hex input, type a hex code (e.g., `#fff`, `#FF5733`), press `Enter` to apply
5. press `s` to save theme
6. press `ESC` to cancel

## Configuration

### config file format

default location: `~/.config/gapistotle/config.conf`

```ini
# theme settings
currentTheme=monokai
themesDirectory=/custom/path/themes  # optional

# panel settings
leftPanelWidth=40
minPanelWidth=15
maxPanelWidthPercent=80
panelResizeIncrement=5

# logging settings
logPath=/tmp/gapistotle.log
logLevel=debug  # debug, info, warn, error
```

### custom themes

**how themes work:**

gapistotle looks for theme files in this order:
1. custom themes directory (if set in config via `themesDirectory`)
2. `$XDG_CONFIG_HOME/gapistotle/themes/`
3. `~/.config/gapistotle/themes/`
4. `./themes/` (when running from source)

themes are loaded from `.conf` files. if no theme files are found, the hardcoded gapistotle theme is used as a fallback.

**creating custom themes:**

create a theme file in `~/.config/gapistotle/themes/` (e.g., `mytheme.conf`):
```ini
name=My Custom Theme
selectedBg=#1e1e1e
selectedFg=#d4d4d4
normalFg=#cccccc
separatorColor=#444444
helpColor=#888888
testCountColor=#888888
borderColor=#444444
treeSymbolColor=#888888
coverageGoodFg=#50fa7b
coverageMediumFg=#f1fa8c
coveragePoorFg=#ff5555
menuNormalFg=#888888
menuActiveFg=#ffffff
menuSelectedBg=#44475a
menuSelectedFg=#ffffff
```

then reload themes via `` ` `` → Theme → Reload Themes.

**installation note:**

when running from source, all 8 themes are available from the `themes/` directory. if you install the binary elsewhere (like `/usr/local/bin/`), you'll need to copy theme files to `~/.config/gapistotle/themes/` or the gapistotle theme will be used as the default.

## Core Philosophy

- **everything accessible via menus** - no hidden hotkeys (press `` ` `` for menu)
- **progressive disclosure** - clean summary by default, detailed output on demand
- **menu-driven interaction** - "Enter = run the thing you selected"

## Architecture

built with:
- [Bubbletea](https://github.com/charmbracelet/bubbletea) - Elm architecture TUI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - terminal styling
- Go's built-in `log/slog` - structured logging

## Contributing

contributions welcome. open an issue or PR if you want to help.

### development

```bash
# build
go build -o gapistotle .

# run from source
go run . /path/to/test/directory

# run tests (when we have them)
go test ./...
```

## License

MIT License - see LICENSE file for details

## Acknowledgments

named after Aristotle, because knowing thy code is the first step to wisdom.

inspired by k9s and the desire for a better Go test coverage workflow.
