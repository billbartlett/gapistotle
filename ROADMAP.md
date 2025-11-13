# Gapistotle Roadmap

This document outlines planned features and enhancements for Gapistotle.

## Current Status

Gapistotle is feature-complete for core test execution and coverage analysis. The focus now is on expanding visualization capabilities and developer workflow integration.

---

## High Priority

### Full Test Output View
**Status:** Planned
**Goal:** Show complete logs and error messages for failing tests

When tests fail, users need to see the full output including:
- Stack traces
- Error messages
- Test-specific output (fmt.Println, t.Log, etc.)
- Colored diff output for assertion failures

**Implementation:**
- Add dedicated failure detail view
- Syntax highlighting for stack traces
- Smart wrapping for long lines
- Searchable output

---

### Enhanced Failure Handling
**Status:** Planned
**Goal:** Prominently display failure details with actionable information

Make it easy to understand why tests fail:
- Quick navigation to first failure
- Visual indicators for error locations
- Links to source files (when terminal supports it)
- Failure grouping by package/type

---

## Medium Priority

### Detailed Coverage View
**Status:** Under consideration
**Goal:** Per-file breakdown with line-by-line coverage visualization

Drill down from package → file → line coverage:
- File list with coverage percentages
- Red/green bars showing covered vs uncovered lines
- Source code viewer with coverage overlay
- Jump to uncovered sections

**Proposal:**
- Parse coverage profile for line-level data
- Use lipgloss to render colorized source
- Navigate with vim-style keys
- Export annotated HTML coverage reports

---

### Filter/Search Tests
**Status:** Planned
**Goal:** Find tests by name, status, or package

Essential for large codebases:
- Fuzzy search by test name
- Filter by status (pass/fail/skip)
- Filter by package pattern
- Save filter presets

---

### Test History & Trends
**Status:** Exploring
**Goal:** Track coverage changes over time

Help teams monitor quality:
- Store test results in SQLite or JSONL
- Sparkline/ASCII charts of coverage trends
- Compare current run to previous runs
- Highlight regressions

**Possible UI:**
```
Coverage Trend (last 10 runs):
█▇▆▇█▇▆▅▄█  87.1%
```

---

## Low Priority

### Export Results
**Status:** Planned
**Goal:** Export test results in multiple formats

- JSON export for CI/CD integration
- HTML coverage reports
- Markdown summary for PR comments
- JUnit XML for compatibility

---

### Browser Integration
**Status:** Exploring
**Goal:** Open browser-based coverage viewer from TUI

Launch Go's built-in coverage tool:
```bash
# From gapistotle, press 'b' to open in browser
go tool cover -html=coverage.out
```

---

### Parallel Test Execution
**Status:** Research needed
**Goal:** Run multiple packages concurrently with live updates

Technical challenges:
- Managing concurrent go test processes
- Real-time streaming of results
- Handling stdout/stderr interleaving
- Progress indicators for running tests

---

### Settings Screen
**Status:** Placeholder exists
**Goal:** Edit configuration without leaving TUI

Currently settings are in config file. Add UI for:
- Panel width preferences
- Theme selection
- Log level configuration
- Test execution flags

---

## Nice to Have

### Watch Mode
**Status:** Detailed proposal exists
**Goal:** Auto-rerun tests when files change

**Implementation details:**
- Use `github.com/fsnotify/fsnotify` for file watching
- Watch `*_test.go`, `*.go`, and `go.mod` files
- Debounce 300-500ms for editor save bursts
- Knight Rider style "Cylon" scanner indicator:
  ```
  `: menu | [*       ] watching | ↑↓/jk: navigate
  `: menu | [ *      ] watching | ↑↓/jk: navigate
  `: menu | [  *     ] watching | ↑↓/jk: navigate
  ```
- Toggle with `w` key
- Config option: `watchDebounceMs=500`

---

### CI/CD Integration
**Status:** Ideas only
**Goal:** Generate outputs suitable for CI pipelines

- Exit codes based on test success/failure
- Structured output for parsing
- Coverage threshold enforcement
- Integration with GitHub Actions, GitLab CI, etc.

---

### Benchmark Results
**Status:** Ideas only
**Goal:** Display benchmark timing and allocations

Parse `go test -bench` output:
- Bar charts for benchmark comparisons
- Memory allocation statistics
- Historical trend tracking
- Regression detection

---

### Test Profiling
**Status:** Ideas only
**Goal:** Visualize CPU/memory profiles from tests

Integration with `pprof`:
- Run tests with `-cpuprofile`/`-memprofile`
- ASCII flame graphs
- Hot path identification
- Profile comparison

---

### Configurable Keymaps
**Status:** Exploring
**Goal:** Let users customize keyboard shortcuts

- Define keybindings in config file
- Support vim/emacs modes
- Visual keymap editor
- Export/import keymap presets

---

### Theme Sharing
**Status:** Planned
**Goal:** Import/export themes easily

- Import theme from URL
  ```
  https://github.com/billbartlett/gapistotle-themes/solarized.conf
  ```
- Export current theme to file
- Community theme repository
- One-command theme installation

---

### ASCII Art Splash Screen
**Status:** Fun idea
**Goal:** Startup branding with philosophical quotes

```
     ____             _     _        _   _
    / ___| __ _ _ __ (_)___| |_ ___ | |_| | ___
   | |  _ / _` | '_ \| / __| __/ _ \| __| |/ _ \
   | |_| | (_| | |_) | \__ \ || (_) | |_| |  __/
    \____|\__,_| .__/|_|___/\__\___/ \__|_|\___|
               |_|

   "Know thy code."
   - Gapistotle

   [Random philosophical quote about testing]
```

---

### Demo Mode Soundtrack
**Status:** Wild idea
**Goal:** Because no other TUI has ever done this

Ultimate demoscene tribute:
- Embed .mod/.xm tracker files (50-200KB)
- Pure Go playback with `hajimehoshi/oto`
- Toggle with `m` key
- VU meter indicator: `♫ [||||::||||||:::||]`
- Auto-disable on SSH/headless environments
- Community-submitted tracks
- Different moods: testing, passing, debugging

**Why:**
- Nostalgia factor
- Conference demo material
- Memorable experience
- Love letter to demoscene culture

---

## Completed Features

See [Recent Changes](#recent-changes) in the main documentation for recently completed work.

---

## Contributing Ideas

Have a feature idea? Open an issue on GitHub with:
- Use case description
- Mockup or ASCII art of proposed UI
- Technical approach (if applicable)
- Priority level (your opinion)

The best features start as detailed proposals with clear user value.
