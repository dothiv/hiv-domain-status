package main

import (
	"fmt"
	"github.com/wsxiaoys/terminal/color"
	"os"
	hivdomainstatus "github.com/dothiv/hiv-domain-status"
)

func error(msg string) {
	color.Fprintln(os.Stderr, "@{!r}ERROR @{|}"+msg)
}

func Help() {
	color.Fprintln(os.Stdout, fmt.Sprintf("Usage: %s %s\n", os.Args[0], "@{g}<command>@{|}"))
	color.Fprintln(os.Stdout, "  @{g}command@{|} may be         help | server\n")
	color.Fprintln(os.Stdout, fmt.Sprintf("Use %s help <command> to get help for a command", os.Args[0]))
}

func main() {
	if len(os.Args) == 1 {
		Help()
		error(fmt.Sprintf("(%s) too few arguments", os.Args[0]))
		os.Exit(1)
	}
	switch os.Args[1] {
	case "help":
		if len(os.Args) == 2 {
			Help()
		} else if len(os.Args) > 2 {
			Help()
			error(fmt.Sprintf("(%s) too many arguments", os.Args[0]))
			os.Exit(1)
		}
		os.Exit(0)
	case "server":
		c, err := hivdomainstatus.NewConfig()
		if err != nil {
			error(err.Error())
			os.Exit(1)
		}
		err = hivdomainstatus.Serve(c)
		if err != nil {
			error(err.Error())
			os.Exit(1)
		}
		os.Exit(0)
	}
}
