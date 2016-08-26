package model

import (
	"flag"
	"hot/util"
	"os"
	"strings"
	"errors"
)

/**
 * autobuider cli parameters
 */
type Args struct {
	ConfigPath string
	App        string
}

func GetArgs() *Args {
	args := new(Args)
	c := flag.String("c", "./conf/config.ini", "Autobuilder's config file path")
	a := flag.String("a", "app", "Assign target app to autobuild (configed in config file)")
	flag.Parse()

	args.ConfigPath = *c
	args.App = *a
	err := checkArgs(args)
	if err != nil {
		util.ColorPrintln("[ERROR] Invalid args: " + err.Error(), util.COLOR_FAIL)
		os.Exit(-1)
	}
	return args
}

func checkArgs(args *Args) error {
	// Check whether config file exists
	if strings.Index(args.ConfigPath, "./") != 0 && strings.Index(args.ConfigPath, "/") != 0 {
		args.ConfigPath = "./" + args.ConfigPath
	}
	_, err := os.Stat(args.ConfigPath)
	if err != nil {
		return errors.New("Config file not found")
	}

	return nil
}

