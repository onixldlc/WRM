package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"windows-resolution-manager/config"
	"windows-resolution-manager/display"
)

func InitializeApp() {
	args := os.Args[1:]
	if len(args) == 0 {
		// Start interactive mode
		StartInteractiveMode()
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
		HandleConfigCommand(args[1:])
	default:
		fmt.Println("Unknown command:", cmd)
		PrintHelp()
	}
}

// StartInteractiveMode starts the interactive CLI session
func StartInteractiveMode() {
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
			HandleConfigCommand(args[1:])
		default:
			fmt.Println("Unknown command:", cmd)
			PrintHelp()
		}
	}
}

// HandleListCommand processes the 'list' command
func HandleListCommand(args []string) {
	if len(args) == 0 {
		// List monitors
		err := display.PrintMonitors()
		if err != nil {
			fmt.Println("Error listing monitors:", err)
		}
	} else if len(args) == 1 {
		// List resolutions for a monitor
		monitorIndex, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Invalid monitor index:", args[0])
			return
		}
		display.ListResolutionsForMonitor(monitorIndex - 1)
	} else if len(args) == 2 {
		// List frequencies for a resolution
		monitorIndex, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Invalid monitor index:", args[0])
			return
		}
		resolution := args[1]
		display.ListFrequenciesForResolution(monitorIndex-1, resolution)
	} else {
		fmt.Println("Invalid list command.")
		PrintHelp()
	}
}

// HandleSetCommand processes the 'set' command
func HandleSetCommand(args []string) {
	if len(args) < 1 {
		fmt.Println("Monitor is required for the set command.")
		fmt.Println("Usage: wrm set <monitor> [resolution] [frequency]")
		return
	}
	monitorInput := args[0]
	monitorIndex, err := strconv.Atoi(monitorInput)
	if err != nil {
		fmt.Println("Invalid monitor index:", monitorInput)
		return
	}
	monitors, err := display.ListMonitors()
	if err != nil {
		fmt.Println("Error listing monitors:", err)
		return
	}
	if monitorIndex < 1 || monitorIndex > len(monitors) {
		fmt.Println("Monitor index out of range")
		return
	}
	mi := monitors[monitorIndex-1]
	deviceName := mi.FriendlyName

	if len(args) == 1 {
		// No resolution provided, list resolutions
		display.ListResolutionsForMonitor(monitorIndex - 1)
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

// HandleConfigCommand processes the 'config' command
func HandleConfigCommand(args []string) {
	config.HandleConfigCommand(args)
}
