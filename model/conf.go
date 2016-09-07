package model

import (
	"github.com/larspensjo/config"
	"hot/util"
	"strings"
	"os"
	"fmt"
	"log"
)

/**
 * Config file fields
 */
type Conf struct {
	MonitorPath string
	GoPath      string
	WatchExt    []string
	AppName     string
	MainFiles   []string
	BinDir      string
	CmdArgs     string
}

func NewConf(args *Args) *Conf {
	return readConfFields(args)
}

func readConfFields(args *Args) *Conf {
	c, err := config.ReadDefault(args.ConfigPath)
	if err != nil {
		util.ColorPrintln("[ERROR] " + err.Error(), util.COLOR_FAIL)
		os.Exit(-1)
	}

	// Whether list all apps
	if args.ListAll {
		listAllApps(c)
		os.Exit(0)
	}

	// Check app
	checkApp(c, args.App)

	// read config fields
	conf := new(Conf)
	conf.MonitorPath, _ = c.String(args.App, "MONITOR_PATH")
	exts, _ := c.String(args.App, "WATCH_EXT")
	conf.WatchExt = strings.Split(exts, ",")
	conf.AppName, _ = c.String(args.App, "APP_NAME")
	conf.GoPath, _ = c.String(args.App, "GOPATH")
	if conf.GoPath != "" {
		os.Setenv("GOPATH", conf.GoPath)
	}
	mainFiles, err := c.String(args.App, "MAIN_FILES")
	if err != nil {
		util.ColorPrintln("[ERROR] Please config main file first!", util.COLOR_FAIL)
		os.Exit(-1)
	}

	binDir, _ := c.String(args.App, "BIN_DIR")
	conf.BinDir = strings.TrimRight(binDir, "/") + "/"
	conf.MainFiles = strings.Split(mainFiles, ",")
	conf.CmdArgs, err = c.String(args.App, "CMD_ARGS")

	str := fmt.Sprintf("[INFO] Configs:\n- APP_NAME: %s\n- MONITOR_PATH: %s \n- GOPATH: %s \n- WATCH_EXT: %v \n- MAIN_FILES: %v\n- BIN_DIR: %s\n- CMD_ARGS: %s\n", conf.AppName, conf.MonitorPath, conf.GoPath, conf.WatchExt, conf.MainFiles, conf.BinDir, conf.CmdArgs)
	util.ColorPrintln(str, util.COLOR_INFO)

	return conf
}

func checkApp(c *config.Config, app string) {
	existApp := false
	for _, v := range c.Sections() {
		if app == v {
			existApp = true
			break
		}
	}
	if !existApp {
		util.ColorPrintln("[ERROR] [ " + app + " ] not found in config file", util.COLOR_FAIL)
		os.Exit(-1)
	}
}

func listAllApps(c *config.Config) {
	for i, app := range c.Sections() {
		if i == 0 {
			// DEFAULT config, not app name
			continue
		}
		util.ColorPrintln(fmt.Sprintf("[%d] %s", i - 1, app), util.COLOR_SUCC)
	}
}
