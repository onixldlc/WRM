package cmd

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"windows-resolution-manager/config"
	"windows-resolution-manager/display"
)

// InitializeApp initializes the application by parsing command-line arguments and executing commands.
func InitializeApp() {
	// Define the --config-file flag
	configFileFlag := flag.String("config-file", "./config.json", "Path to the configuration file")

	// Parse the flags
	flag.Parse()

	// Remaining arguments after flags
	args := flag.Args()

	// Check if the config file exists; if not, create it with default configurations
	err := config.EnsureConfigFile(*configFileFlag)
	if err != nil {
		fmt.Println("Error ensuring configuration file:", err)
		return
	}

	if len(args) == 0 {
		// Start interactive mode
		StartInteractiveMode(*configFileFlag)
		return
	}

	cmd := strings.ToLower(args[0])
	switch cmd {
	case "help":
		PrintHelp()
	case "list", "ls", "l":
		HandleListCommand(args[1:])
	case "set", "change", "ch", "c", "s":
		HandleSetCommand(args[1:])
	case "config":
		HandleConfigCommand(args[1:], *configFileFlag)
	default:
		fmt.Println("Unknown command:", cmd)
		PrintHelp()
	}
}

// StartInteractiveMode starts the interactive CLI session.
func StartInteractiveMode(configFile string) {
	fmt.Println("Entering interactive mode. Type 'help' for a list of commands.")
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("wrm> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		args := strings.Fields(input)
		if len(args) == 0 {
			continue
		}
		cmd := strings.ToLower(args[0])
		switch cmd {
		case "exit", "quit":
			fmt.Println("Exiting interactive mode.")
			return
		case "help":
			PrintHelp()
		case "list", "ls", "l":
			HandleListCommand(args[1:])
		case "set", "change", "ch", "c", "s":
			HandleSetCommand(args[1:])
		case "config":
			HandleConfigCommand(args[1:], configFile)
		default:
			fmt.Println("Unknown command:", cmd)
			PrintHelp()
		}
	}
}

// HandleConfigCommand processes the 'config' command.
func HandleConfigCommand(args []string, configFile string) {
	config.HandleConfigCommand(args, configFile)
}

// HandleListCommand processes the 'list' command.
func HandleListCommand(args []string) {
	if len(args) == 0 {
		// List monitors
		err := display.PrintMonitors()
		if err != nil {
			fmt.Println("Error listing monitors:", err)
		}
	} else if len(args) >= 1 {
		// Determine if the first argument is a monitor index or a friendly name
		monitorIdentifier := args[0]
		var monitorIndex int
		var err error

		// Try to convert to integer
		monitorIndex, err = strconv.Atoi(monitorIdentifier)
		if err != nil {
			// Not an integer, treat as friendly name
			monitors, listErr := display.ListMonitors()
			if listErr != nil {
				fmt.Println("Error listing monitors:", listErr)
				return
			}
			monitorIndex = -1 // Initialize with invalid index
			for i, mi := range monitors {
				if strings.EqualFold(mi.FriendlyName, monitorIdentifier) {
					monitorIndex = i + 1 // Monitors are 1-indexed
					break
				}
			}
			if monitorIndex == -1 {
				fmt.Printf("Monitor with friendly name '%s' not found.\n", monitorIdentifier)
				return
			}
		}

		// Adjust monitorIndex for 0-based indexing
		zeroBasedIndex := monitorIndex - 1

		// Validate monitorIndex
		monitors, listErr := display.ListMonitors()
		if listErr != nil {
			fmt.Println("Error listing monitors:", listErr)
			return
		}
		if zeroBasedIndex < 0 || zeroBasedIndex >= len(monitors) {
			fmt.Println("Monitor index out of range.")
			return
		}

		if len(args) == 1 {
			// List resolutions for the monitor
			display.ListResolutionsForMonitor(zeroBasedIndex)
		} else if len(args) == 2 {
			// List frequencies for the resolution on the monitor
			resolution := args[1]
			display.ListFrequenciesForResolution(zeroBasedIndex, resolution)
		} else {
			fmt.Println("Invalid list command.")
			PrintHelp()
		}
	} else {
		fmt.Println("Invalid list command.")
		PrintHelp()
	}
}

// HandleSetCommand processes the 'set' command.
func HandleSetCommand(args []string) {
	if len(args) < 1 {
		fmt.Println("Monitor is required for the set command.")
		fmt.Println("Usage: wrm set <monitor> [resolution] [frequency]")
		return
	}
	monitorInput := args[0]
	var monitorIndex int
	var err error

	// Try to convert to integer
	monitorIndex, err = strconv.Atoi(monitorInput)
	if err != nil {
		// Not an integer, treat as friendly name
		monitors, listErr := display.ListMonitors()
		if listErr != nil {
			fmt.Println("Error listing monitors:", listErr)
			return
		}
		monitorIndex = -1 // Initialize with invalid index
		for i, mi := range monitors {
			if strings.EqualFold(mi.FriendlyName, monitorInput) {
				monitorIndex = i + 1 // Monitors are 1-indexed
				break
			}
		}
		if monitorIndex == -1 {
			fmt.Printf("Monitor with friendly name '%s' not found.\n", monitorInput)
			return
		}
	}

	// Adjust monitorIndex for 0-based indexing
	zeroBasedIndex := monitorIndex - 1

	// Retrieve the list of monitors
	monitors, err := display.ListMonitors()
	if err != nil {
		fmt.Println("Error listing monitors:", err)
		return
	}
	if zeroBasedIndex < 0 || zeroBasedIndex >= len(monitors) {
		fmt.Println("Monitor index out of range.")
		return
	}
	mi := monitors[zeroBasedIndex]
	deviceName := mi.DeviceName // Changed from FriendlyName to DeviceName

	if len(args) == 1 {
		// No resolution provided, list resolutions
		display.ListResolutionsForMonitor(zeroBasedIndex)
		fmt.Println("Usage: wrm set <monitor> <resolution> [frequency]")
		return
	}

	resolution := args[1]
	frequency := uint32(0)
	if len(args) >= 3 {
		freqValue, err := strconv.Atoi(args[2])
		if err != nil {
			fmt.Println("Invalid frequency:", args[2])
			return
		}
		frequency = uint32(freqValue)
	}

	err = display.SetResolution(deviceName, resolution, frequency)
	if err != nil {
		fmt.Println("Error setting resolution:", err)
	}
}
