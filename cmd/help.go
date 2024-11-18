package cmd

import "fmt"

// PrintHelp displays the combined help message
func PrintHelp() {
	helpMessage := `
Usage: wrm [--config-file <path>] <command> [arguments]

Commands:
  help                                Show this help message
  list                                List all monitors
  list <monitor>                      List resolutions for the specified monitor
  list <monitor> <resolution>         List frequencies for the specified resolution on the monitor
  set <monitor> <resolution> [freq]    Set the resolution and frequency for the specified monitor
  config                              List pre-configured settings
  config <config_name/index>          Apply a saved configuration by name or index

Aliases:
  list -> ls, l
  set -> change, ch, c, s

Flags:
  --config-file <path>                Specify a custom configuration file path (default: ./config.json)

Examples:
  wrm list
  wrm ls 1
  wrm ls 27G2G5
  wrm l 2 1920x1080
  wrm set 1 1280x720 60
  wrm set 27G2G5 1280x720 60
  wrm config
  wrm config "Gaming Setup"
  wrm config 2
`
	fmt.Println(helpMessage)
}
