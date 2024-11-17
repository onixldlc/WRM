package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"windows-resolution-manager/display"
)

// Config represents a display configuration
type Config struct {
	Name       string `json:"name"`
	Monitor    int    `json:"monitor"`
	Resolution string `json:"resolution"`
	Frequency  uint32 `json:"frequency"`
}

// Configurations holds a list of Config
type Configurations struct {
	Configs []Config `json:"configurations"`
}

// HandleConfigCommand processes the 'config' command
func HandleConfigCommand(args []string) {
	configFile := "config/wrm_config.json"
	configs, err := LoadConfigurations(configFile)
	if err != nil {
		fmt.Println("Error loading configurations:", err)
		return
	}
	if len(args) == 0 {
		// List configurations
		fmt.Println("Available configurations:")
		for i, cfg := range configs.Configs {
			fmt.Printf("%d. %s: Monitor %d, %s @ %d Hz\n", i+1, cfg.Name, cfg.Monitor, cfg.Resolution, cfg.Frequency)
		}
	} else {
		// Apply a configuration
		cfgIndex, err1 := strconv.Atoi(args[0])
		var cfg *Config
		if err1 == nil && cfgIndex > 0 && cfgIndex <= len(configs.Configs) {
			cfg = &configs.Configs[cfgIndex-1]
		} else {
			// Search by name
			for _, c := range configs.Configs {
				if c.Name == args[0] {
					cfg = &c
					break
				}
			}
		}
		if cfg == nil {
			fmt.Println("Configuration not found.")
			return
		}
		// Apply the configuration
		monitors, err := display.ListMonitors()
		if err != nil {
			fmt.Println("Error listing monitors:", err)
			return
		}
		if cfg.Monitor < 1 || cfg.Monitor > len(monitors) {
			fmt.Println("Monitor index in configuration is out of range.")
			return
		}
		mi := monitors[cfg.Monitor-1]
		deviceName := mi.FriendlyName
		err = display.SetResolution(deviceName, cfg.Resolution, cfg.Frequency)
		if err != nil {
			fmt.Println("Error applying configuration:", err)
		}
	}
}

// LoadConfigurations loads configurations from a JSON file
func LoadConfigurations(filename string) (*Configurations, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	byteValue, _ := ioutil.ReadAll(file)
	var configs Configurations
	json.Unmarshal(byteValue, &configs)
	return &configs, nil
}
