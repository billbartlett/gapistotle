package main

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Layout constants - fixed chrome sizes
const (
	MenuBarH   = 4 // Menu bar + help line
	StatusBarH = 0 // Reserved for future status bar
	BorderW    = 0 // Visual border/outline width if any
	GutterW    = 1 // Spacing between panes (the separator)
)

// Frame represents the total available window space
type Frame struct {
	Width  int
	Height int
}

// Split represents computed layout dimensions for a two-pane view
type Split struct {
	TotalW     int // Total window width
	LeftW      int // Left pane width
	RightW     int // Right pane width
	ContentW   int // Total content width (minus borders)
	ContentH   int // Total content height (minus menu/status bars)
	MenuBarH   int // Menu bar height
	StatusBarH int // Status bar height
	GutterW    int // Gutter width between panes
	BorderW    int // Border width
}

// NewSplit computes pane sizes from a window frame and desired left pane width
func NewSplit(f Frame, leftW int, minLeft, maxLeft int) Split {
	leftW = Clamp(leftW, minLeft, maxLeft)
	contentH := Max(0, f.Height-MenuBarH-StatusBarH)
	contentW := Max(0, f.Width-(2*BorderW))
	rightW := Max(0, contentW-leftW-GutterW)

	return Split{
		TotalW:     f.Width,
		LeftW:      leftW,
		RightW:     rightW,
		ContentW:   contentW,
		ContentH:   contentH,
		MenuBarH:   MenuBarH,
		StatusBarH: StatusBarH,
		GutterW:    GutterW,
		BorderW:    BorderW,
	}
}

// Clamp constrains n between min and max
func Clamp(n, min, max int) int {
	if n < min {
		return min
	}
	if n > max {
		return max
	}
	return n
}

// Max returns the larger of two integers
func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Min returns the smaller of two integers
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// VSep draws a vertical separator of a given height with a style
func VSep(height int, style lipgloss.Style) string {
	if height <= 0 {
		return ""
	}
	var b strings.Builder
	seg := style.Render("│")
	for i := 0; i < height; i++ {
		b.WriteString(seg)
		if i != height-1 {
			b.WriteByte('\n')
		}
	}
	return b.String()
}

// HR draws a horizontal rule across width using the theme's separator glyph
func HR(width int, style lipgloss.Style, glyph string) string {
	if width <= 0 {
		return ""
	}
	if glyph == "" {
		glyph = "─"
	}
	return style.Render(strings.Repeat(glyph, width))
}

// PadBlock pads a multiline string to exactly height lines (truncate or pad)
func PadBlock(s string, height int) string {
	if height <= 0 {
		return ""
	}
	lines := strings.Split(s, "\n")
	if len(lines) > height {
		lines = lines[:height]
	} else if len(lines) < height {
		padding := make([]string, height-len(lines))
		lines = append(lines, padding...)
	}
	return strings.Join(lines, "\n")
}

// TruncateLines truncates output to maxLines, adding ellipsis if truncated
func TruncateLines(s string, maxLines int) (string, bool) {
	if maxLines <= 0 {
		return "", false
	}
	lines := strings.Split(s, "\n")
	if len(lines) <= maxLines {
		return s, false
	}
	return strings.Join(lines[:maxLines], "\n"), true
}

// RightPanelPageSize calculates the page size for right panel scrolling
// Accounts for borders and padding in the content area
func RightPanelPageSize(split Split) int {
	return Max(0, split.ContentH-2)
}

// HelpScreenPageSize calculates the page size for help screen scrolling
// Accounts for menu bar (4 lines), border (2 lines), and padding (2 lines)
func HelpScreenPageSize(height int) int {
	return Max(0, height-8)
}
