package cli

import (
	"flag"
	"fmt"
)

type InstallOptions struct {
	Production bool
	SaveDev    bool
	Packages   []string
}

func Parse() (InstallOptions, error) {
	prd := flag.Bool("production", false, "production")
	saveDev := flag.Bool("save-dev", false, "save-dev")
	flag.Parse()

	if flag.Arg(0) != "install" {
		return InstallOptions{}, fmt.Errorf("first argument must be 'install'")
	}

	return InstallOptions{
		Production: *prd,
		SaveDev:    *saveDev,
		Packages:   flag.Args()[1:],
	}, nil
}
