package main

import (
	"bufio"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Config struct {
	CurrentTheme         string
	ThemesDirectory      string // Custom path for themes (empty = use XDG default)
	LeftPanelWidth       int
	MinPanelWidth        int
	MaxPanelWidthPercent int
	PanelResizeIncrement int
	LogPath              string // Path to log file (empty = no logging)
	LogLevel             string // Log level: debug, info, warn, error
}

func getConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ".gapistotle.conf" // Fallback to current directory
	}
	return filepath.Join(homeDir, ".gapistotle.conf")
}

// ResolveConfigPath determines the config file path to use
// Priority: 1) -c flag, 2) $GAPISTOTLE_CONFIG, 3) default
func ResolveConfigPath(flagPath string) string {
	// 1. Check command-line flag
	if flagPath != "" {
		return expandPath(flagPath)
	}

	// 2. Check environment variable
	if envPath := os.Getenv("GAPISTOTLE_CONFIG"); envPath != "" {
		return expandPath(envPath)
	}

	// 3. Use default
	return getConfigPath()
}

// LoadConfig loads the config file from the specified path
func LoadConfig(configPath string) Config {
	config := Config{
		CurrentTheme:         "default",
		LeftPanelWidth:       defaultLeftPanelWidth,
		MinPanelWidth:        minPanelWidth,
		MaxPanelWidthPercent: maxPanelWidthPercent,
		PanelResizeIncrement: panelResizeIncrement,
		LogPath:              "/tmp/gapistotle.log",
		LogLevel:             "debug",
	}

	file, err := os.Open(configPath)
	if err != nil {
		return config // File doesn't exist yet, use defaults
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "currentTheme":
			config.CurrentTheme = value
		case "themesDirectory":
			config.ThemesDirectory = value
		case "leftPanelWidth":
			if width, err := strconv.Atoi(value); err == nil {
				config.LeftPanelWidth = width
			}
		case "minPanelWidth":
			if width, err := strconv.Atoi(value); err == nil {
				config.MinPanelWidth = width
			}
		case "maxPanelWidthPercent":
			if percent, err := strconv.Atoi(value); err == nil {
				config.MaxPanelWidthPercent = percent
			}
		case "panelResizeIncrement":
			if increment, err := strconv.Atoi(value); err == nil {
				config.PanelResizeIncrement = increment
			}
		case "logPath":
			config.LogPath = value
		case "logLevel":
			config.LogLevel = value
		}
	}

	return config
}

// SaveConfig saves the config file
func SaveConfig(config Config, configPath string) error {
	file, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	writer.WriteString("# Gapistotle Configuration\n")
	writer.WriteString("# Theme settings\n")
	writer.WriteString("currentTheme=" + config.CurrentTheme + "\n")
	if config.ThemesDirectory != "" {
		writer.WriteString("themesDirectory=" + config.ThemesDirectory + "\n")
	}
	writer.WriteString("\n# Panel settings\n")
	writer.WriteString("leftPanelWidth=" + strconv.Itoa(config.LeftPanelWidth) + "\n")
	writer.WriteString("minPanelWidth=" + strconv.Itoa(config.MinPanelWidth) + "\n")
	writer.WriteString("maxPanelWidthPercent=" + strconv.Itoa(config.MaxPanelWidthPercent) + "\n")
	writer.WriteString("panelResizeIncrement=" + strconv.Itoa(config.PanelResizeIncrement) + "\n")
	writer.WriteString("\n# Logging settings\n")
	writer.WriteString("logPath=" + config.LogPath + "\n")
	writer.WriteString("logLevel=" + config.LogLevel + "\n")

	return writer.Flush()
}

// expandPath expands ~ to home directory
func expandPath(path string) string {
	if path == "" {
		return path
	}
	if strings.HasPrefix(path, "~/") {
		homeDir, err := os.UserHomeDir()
		if err == nil {
			return filepath.Join(homeDir, path[2:])
		}
	}
	return path
}

// GetThemesDir returns the directory for custom themes
// Priority: 1) config override, 2) XDG_CONFIG_HOME, 3) ~/.config, 4) ./themes
func GetThemesDir(config Config) string {
	// 1. Check config file override
	if config.ThemesDirectory != "" {
		return expandPath(config.ThemesDirectory)
	}

	// 2. Check XDG_CONFIG_HOME environment variable
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return filepath.Join(xdg, "gapistotle", "themes")
	}

	// 3. Fallback to ~/.config/gapistotle/themes
	homeDir, err := os.UserHomeDir()
	if err == nil {
		return filepath.Join(homeDir, ".config", "gapistotle", "themes")
	}

	// 4. Last resort: local themes directory (dev mode)
	return "themes"
}
