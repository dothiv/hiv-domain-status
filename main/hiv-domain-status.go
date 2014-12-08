package main

import (
	"fmt"
	"database/sql"
	"github.com/wsxiaoys/terminal/color"
	"os"
	hivdomainstatus "github.com/dothiv/hiv-domain-status"
)

func error(msg string) {
	color.Fprintln(os.Stderr, "@{!r}ERROR @{|}"+msg)
}

func Help() {
	color.Fprintln(os.Stdout, fmt.Sprintf("Usage: %s %s\n", os.Args[0], "@{g}<command>@{|}"))
	color.Fprintln(os.Stdout, "  @{g}command@{|} may be         help | server | check\n")
	color.Fprintln(os.Stdout, fmt.Sprintf("Use %s help <command> to get help for a command", os.Args[0]))
}

func main() {
	if len(os.Args) == 1 {
		Help()
		error(fmt.Sprintf("(%s) too few arguments", os.Args[0]))
		os.Exit(1)
	}
	c, err := hivdomainstatus.NewConfig()
	if err != nil {
		error(err.Error())
		os.Exit(1)
	}

	switch os.Args[1] {
	case "help":
		if len(os.Args) == 2 {
			Help()
		} else if len(os.Args) > 3 {
			Help()
			error(fmt.Sprintf("(%s) too many arguments", os.Args[0]))
			os.Exit(1)
		}
		switch(os.Args[2]) {
		case "check":
			color.Fprintln(os.Stdout, fmt.Sprintf("Usage: %s %s\n", os.Args[0], " check @{g}<hiv-domain>@{|}"))
			os.Stdout.WriteString("Check hiv domains.\n")
			os.Stdout.WriteString("\n")
			color.Fprintln(os.Stdout, "  @{g}hiv-domain@{|}           the .hiv domain to check")
			color.Fprintln(os.Stdout, "                       check all registered domains if not set")
		}
		os.Exit(0)
	case "server":
		err = hivdomainstatus.Serve(c)
		if err != nil {
			error(err.Error())
			os.Exit(1)
		}
		os.Exit(0)
	case "check":
		// Open DB
		db, err := sql.Open("postgres", c.DSN())
		if err != nil {
			error(err.Error())
			os.Exit(1)
		}
		domainRepo := hivdomainstatus.NewDomainRepository(db)
		domainCheckRepo := hivdomainstatus.NewDomainCheckRepository(db)
		manager := hivdomainstatus.NewManager(domainRepo, domainCheckRepo)

		if len(os.Args) > 2 {
			var result *hivdomainstatus.DomainCheckResult
			result, err = hivdomainstatus.CheckDomain(c, os.Args[2])
			manager.OnCheckDomainResult(result)
		} else {
			domains, findAllErr := domainRepo.FindAll()
			if findAllErr != nil {
				error(findAllErr.Error())
				os.Exit(1)
			}
			for _, domain := range domains {
				result, _ := hivdomainstatus.CheckDomain(c, domain.Name)
				manager.OnCheckDomainResult(result)
			}
		}
		os.Exit(0)
	}
}
