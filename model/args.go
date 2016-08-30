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
	ListAll    bool
}

func GetArgs() *Args {
	args := new(Args)
	h := flag.Bool("h", false, "Get help informations")
	c := flag.String("c", "./src/hot/conf/config.ini", "Hot's config file path")
	a := flag.String("a", "app", "Assign target app to autobuild (configed in config file)")
	l := flag.Bool("l", false, "List all apps configed in Hot's config file")
	flag.Parse()
	if *h {
		flag.Usage()
		os.Exit(0)
	}

	args.ConfigPath = *c
	args.App = *a
	args.ListAll = *l
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

