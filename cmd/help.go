package cmd

import "fmt"

// PrintHelp displays the help message
func PrintHelp() {
	fmt.Println("Usage:")
	fmt.Println("  wrm help                            Show this help message")
	fmt.Println("  wrm list                            List all monitors")
	fmt.Println("  wrm list <monitor>                  List resolutions for the monitor")
	fmt.Println("  wrm list <monitor> <resolution>     List frequencies for the resolution")
	fmt.Println("  wrm set <monitor> <resolution> [frequency]  Set the monitor resolution")
	fmt.Println("  wrm config                          List pre-configured settings")
	fmt.Println("  wrm config <config_name>            Set resolution based on configuration")
	fmt.Println()
	fmt.Println("Aliases:")
	fmt.Println("  list -> ls, l")
	fmt.Println("  set -> change, ch, c, s")
}
