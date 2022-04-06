package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var (
	all, force bool
	mode       int
)

func main() {
	// ls command and its own parameters
	lsCmd := flag.NewFlagSet("ls", flag.ExitOnError)
	lsCmd.BoolVar(&all, "a", false, "List all files including hidden ones")
	// rm command and its own parameters
	rmCmd := flag.NewFlagSet("rm", flag.ExitOnError)
	rmCmd.BoolVar(&force, "f", false, "Force remove")
	// touch command and its own parameters
	touchCmd := flag.NewFlagSet("touch", flag.ExitOnError)
	touchCmd.IntVar(&mode, "m", 0644, "Default mode")

	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		fmt.Println("ls/rm/touch or count subcommand is required")
		return
	}

	switch args[0] {
	case lsCmd.Name():
		// Test with: ./flag-subcommands ls -al .
		lsCmd.Parse(args[1:])
		if all {
			log.Printf("listing all files includin excluded ones")
		}
		log.Printf("Called ls with args %v", lsCmd.Args())
	case rmCmd.Name():
		// Test with: ./flag-subcommands rm -f .
		rmCmd.Parse(args[1:])
		log.Printf("Called rm %v", rmCmd.Args())
	case touchCmd.Name():
		touchCmd.Parse(args[1:])
		// Test with: ./flag-subcommands touch -m 0666 asdasd
		if len(touchCmd.Args()) == 0 {
			log.Printf("Filename is required")
			return
		}
		log.Printf("Called touch %v", touchCmd.Args())
	default:
		log.Printf("unkown command %s", args[0])
		os.Exit(1)
	}
}
