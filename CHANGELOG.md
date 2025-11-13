# Changelog

All notable changes to gapistotle will (should...) be documented in this file.

## [0.1.0] - 12 Nov 2025

### Test Execution
- Test package scanning with automatic discovery
- Single package and sequential "Run All" execution
- Directory-specific test modes (unit/integration/all) with persistent config
- Integration test detection via build tags
- Separate display sections for unit vs integration tests
- Time breakdown showing actual test execution vs setup/overhead time
- Slowest-first sorting for passed tests to identify performance issues
- Coverage analysis with statement-level granularity
- Coverage gap analysis with function-level impact calculations

### UI & Navigation
- Context-aware full-screen mode ('f' key) for test results and coverage gaps
- Progressive disclosure with summary/details/gaps views
- Three-panel layout with resizable panels ('[' and ']' keys)
- Vim-style navigation (hjkl + arrow keys)
- Comprehensive help screen with keybinding reference
- Test mode indicator showing current directory's mode
- Version display in help screen and via --version flag

### Themes & Customization
- 8 built-in themes (gapistotle, dracula, nord, monokai, atom-one-dark, desert, industry, material)
- Full theme editor with live preview and 48-color palette
- Custom theme support via ~/.config/gapistotle/themes/
- Theme file-based loading with XDG Base Directory support

### Configuration
- Persistent settings in ~/.config/gapistotle/config.conf
- XDG Base Directory support for config and themes
- Structured logging with configurable levels
- Flexible config paths via -c flag or environment variable

### Documentation
- Comprehensive README with installation and usage guide
- Theme system documentation with examples
- MIT License
