package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"windows-resolution-manager/display"
)

// Config represents a display configuration
type Config struct {
	Name        string `json:"name"`
	Monitor     int    `json:"monitor"`      // Optional if MonitorName is used
	MonitorName string `json:"monitor_name"` // Optional if Monitor is used
	Resolution  string `json:"resolution"`
	Frequency   uint32 `json:"frequency"`
}

// Configurations holds a list of Config
type Configurations struct {
	Configs []Config `json:"configurations"`
}

// HandleConfigCommand processes the 'config' command with the provided config file path.
func HandleConfigCommand(args []string, configFile string) {
	// Ensure the configuration file exists; create with defaults if it doesn't
	err := EnsureConfigFile(configFile)
	if err != nil {
		fmt.Println("Error ensuring configuration file:", err)
		return
	}

	// Load configurations from the file
	configs, err := LoadConfigurations(configFile)
	if err != nil {
		fmt.Println("Error loading configurations:", err)
		return
	}

	if len(args) == 0 {
		// List configurations
		fmt.Println("Available configurations:")
		for i, cfg := range configs.Configs {
			monitorIdentifier := ""
			if cfg.MonitorName != "" {
				monitorIdentifier = fmt.Sprintf("(%s)", cfg.MonitorName)
			} else {
				monitorIdentifier = fmt.Sprintf("(Monitor %d)", cfg.Monitor)
			}
			fmt.Printf("%d. %s: %s, %s @ %d Hz\n", i+1, cfg.Name, monitorIdentifier, cfg.Resolution, cfg.Frequency)
		}
	} else {
		// Apply a configuration
		cfgIndex, err1 := strconv.Atoi(args[0])
		var cfg *Config
		if err1 == nil && cfgIndex > 0 && cfgIndex <= len(configs.Configs) {
			cfg = &configs.Configs[cfgIndex-1]
		} else {
			// Search by configuration name
			for _, c := range configs.Configs {
				if strings.EqualFold(c.Name, args[0]) {
					cfg = &c
					break
				}
			}
		}
		if cfg == nil {
			fmt.Println("Configuration not found.")
			return
		}

		// Retrieve the list of monitors
		monitors, err := display.ListMonitors()
		if err != nil {
			fmt.Println("Error listing monitors:", err)
			return
		}

		var targetMonitor display.MonitorInfo
		if cfg.MonitorName != "" {
			// Find monitor by friendly name
			found := false
			for _, mi := range monitors {
				if strings.EqualFold(mi.FriendlyName, cfg.MonitorName) {
					targetMonitor = mi
					found = true
					break
				}
			}
			if !found {
				fmt.Printf("Monitor with friendly name '%s' not found.\n", cfg.MonitorName)
				return
			}
		} else {
			// Find monitor by index
			if cfg.Monitor < 1 || cfg.Monitor > len(monitors) {
				fmt.Println("Monitor index in configuration is out of range.")
				return
			}
			targetMonitor = monitors[cfg.Monitor-1]
		}

		// Apply the configuration
		err = display.SetResolution(targetMonitor.DeviceName, cfg.Resolution, cfg.Frequency)
		if err != nil {
			fmt.Println("Error applying configuration:", err)
		} else {
			fmt.Printf("Configuration '%s' applied successfully to %s.\n", cfg.Name, targetMonitor.FriendlyName)
		}
	}
}

// EnsureConfigFile checks if the config file exists. If not, it creates one with default configurations.
func EnsureConfigFile(filename string) error {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		// Retrieve connected monitors to validate monitor indices
		monitors, err := display.ListMonitors()
		if err != nil {
			return fmt.Errorf("error listing monitors for default configurations: %v", err)
		}
		if len(monitors) < 2 {
			return fmt.Errorf("requires at least 2 monitors to create default configurations")
		}

		// Create default configurations based on connected monitors
		defaultConfigs := Configurations{
			Configs: []Config{
				{
					Name:       "Gaming Setup",
					Monitor:    1,
					Resolution: "1920x1080",
					Frequency:  180,
				},
				{
					Name:        "Work Setup",
					MonitorName: monitors[0].FriendlyName,
					Resolution:  "2560x1440",
					Frequency:   60,
				},
			},
		}

		// Serialize to JSON
		data, err := json.MarshalIndent(defaultConfigs, "", "  ")
		if err != nil {
			return fmt.Errorf("error marshaling default configurations: %v", err)
		}

		// Ensure the directory exists
		dir := getDir(filename)
		if dir != "" {
			if err := os.MkdirAll(dir, os.ModePerm); err != nil {
				return fmt.Errorf("error creating directories for config file: %v", err)
			}
		}

		// Write to file
		err = ioutil.WriteFile(filename, data, 0644)
		if err != nil {
			return fmt.Errorf("error writing default configuration to file: %v", err)
		}

		fmt.Printf("Configuration file '%s' created with default settings.\n", filename)
	}
	return nil
}

// getDir extracts the directory from a file path.
func getDir(filePath string) string {
	lastSlash := strings.LastIndex(filePath, "/")
	if lastSlash == -1 {
		lastSlash = strings.LastIndex(filePath, "\\")
	}
	if lastSlash == -1 {
		return ""
	}
	return filePath[:lastSlash]
}

// LoadConfigurations loads configurations from a JSON file.
func LoadConfigurations(filename string) (*Configurations, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("could not open config file '%s': %v", filename, err)
	}
	defer file.Close()

	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("could not read config file '%s': %v", filename, err)
	}

	var configs Configurations
	err = json.Unmarshal(byteValue, &configs)
	if err != nil {
		return nil, fmt.Errorf("invalid JSON in config file '%s': %v", filename, err)
	}
	return &configs, nil
}
