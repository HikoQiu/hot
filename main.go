package main

import (
	"github.com/fsnotify/fsnotify"
	"log"
	"io/ioutil"
	"path"
	"os"
	"os/exec"
	"bytes"
	"strings"
	"runtime"
	"sync"
	"fmt"
	"hot/model"
	"hot/util"
	"time"
)

var (
	cmd *exec.Cmd
	state sync.Mutex
	paths = make([]string, 0)

	lastUpdateTime int64

/**
 * autobuider cli parameters
 */
	args *model.Args

/**
 * Config file parameters
 */
	conf *model.Conf
)

func main() {
	// Get parameters & Init config from file
	args = model.GetArgs()
	conf = model.NewConf(args)

	// Monitor source files
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if !checkIfWatchExt(event.Name) {
					continue
				}

			// 3 秒内不重复编译
				if lastUpdateTime > time.Now().Unix() - 3 {
					continue
				}

				lastUpdateTime = time.Now().Unix()
				util.ColorPrintln("[INFO]" + event.String(), util.COLOR_INFO)
				go Autobuild(conf.MainFiles)
				if event.Op & fsnotify.Write == fsnotify.Write {
					util.ColorPrintln("[INFO] modified file: " + event.Name + " - " + event.String(), util.COLOR_INFO)
				}
			case err := <-watcher.Errors:
				util.ColorPrintln("[ERROR] error: " + err.Error(), util.COLOR_FAIL)
			}
		}
	}()

	// 定时重刷所有目录
	go func() {
		for {
			paths = make([]string, 0)
			paths = append(paths, conf.MonitorPath)
			readDirs(conf.MonitorPath, &paths)
			for _, p := range paths {
				watcher.Add(p)
			}
			if err != nil {
				log.Fatal(err)
			}
			time.Sleep(time.Second * 10)
		}
	}()

	// Run app when autobuilder starts
	Autobuild(conf.MainFiles)
	<-done
}

func readDirs(directory string, paths *[]string) {
	fileInfos, err := ioutil.ReadDir(directory)
	if err != nil {
		return
	}

	useDirectory := false
	for _, fileInfo := range fileInfos {

		if fileInfo.IsDir() == true && fileInfo.Name()[0] != '.' {
			readDirs(directory + "/" + fileInfo.Name(), paths)
			continue
		}

		if useDirectory == true {
			continue
		}

		if path.Ext(fileInfo.Name()) == ".go" {
			*paths = append(*paths, directory)
			useDirectory = true
		}
	}
	return
}

func Autobuild(files []string) {
	state.Lock()
	defer state.Unlock()

	util.ColorPrintln("[INFO] Start to rebuild [ " + conf.AppName + " ]", util.COLOR_INFO)
	cmdName := "go"
	var err error
	//if conf.GoInstall || conf.Gopm.Install {
	if false {
		icmd := exec.Command("go", "list", "./...")
		buf := bytes.NewBuffer([]byte(""))
		icmd.Stdout = buf
		icmd.Env = append(os.Environ(), "GOGC=off")
		err = icmd.Run()
		if err == nil {
			list := strings.Split(buf.String(), "\n")[1:]
			for _, pkg := range list {
				if len(pkg) == 0 {
					continue
				}
				icmd = exec.Command(cmdName, "install", pkg)
				icmd.Stdout = os.Stdout
				icmd.Stderr = os.Stderr
				icmd.Env = append(os.Environ(), "GOGC=off")
				err = icmd.Run()
				if err != nil {
					break
				}
			}
		}
	}

	if err == nil {
		if runtime.GOOS == "windows" {
			conf.AppName += ".exe"
		}

		args := []string{"build"}
		args = append(args, "-o", conf.AppName)
		//if buildTags != "" {
		//	args = append(args, "-tags", buildTags)
		//}
		args = append(args, files...)
		bcmd := exec.Command(cmdName, args...)
		bcmd.Env = append(os.Environ(), "GOGC=off")
		bcmd.Stdout = os.Stdout
		bcmd.Stderr = os.Stderr
		util.ColorPrintln("[EXEC] " + cmdName + " " + strings.Join(args, " "), util.COLOR_WARNING)
		err = bcmd.Run()
	}

	if err != nil {
		util.ColorPrintln("[ERROR] fail to build " + conf.AppName, util.COLOR_FAIL)
		return
	}


	// mv bin file to BIN_DIR
	err = os.Rename("./" + conf.AppName, conf.BinDir + conf.AppName)
	if err != nil {
		util.ColorPrintln("[ERR] Rename " + conf.AppName + " fail - " + err.Error(), util.COLOR_FAIL)
		os.Exit(-1)
	}
	Restart(conf.BinDir + conf.AppName)
}

func Restart(binFile string) {
	util.ColorPrintln("[INFO] Kill running process of [ " + conf.AppName + " ]", util.COLOR_INFO)
	Kill()
	go Start(binFile)
}

func Kill() {
	defer func() {
		if e := recover(); e != nil {
			util.ColorPrintln(fmt.Sprintf("[ERROR] Kill.recover -> %v", e), util.COLOR_FAIL)
		}
	}()
	if cmd != nil && cmd.Process != nil {
		err := cmd.Process.Kill()
		if err != nil {
			util.ColorPrintln(fmt.Sprintf("[ERROR] Kill -> -> %v", err.Error()), util.COLOR_FAIL)
		}
	}
}

func Start(binFile string) {
	util.ColorPrintln("[INFO] Run [ " + conf.AppName + " ]", util.COLOR_INFO)
	if strings.Index(binFile, "/") != 0 && strings.Index(binFile, "./") == -1 {
		binFile = "./" + binFile
	}

	cmd = exec.Command(binFile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Args = append([]string{binFile, conf.CmdArgs})

	util.ColorPrintln("[EXEC] " + binFile + " " + strings.Join(cmd.Args, " ") + "", util.COLOR_WARNING)
	go cmd.Run()
	util.ColorPrintln("[SUCC] Ok! " + conf.AppName + " started.", util.COLOR_SUCC)
}

func checkIfWatchExt(name string) bool {
	for _, s := range conf.WatchExt {
		if strings.HasSuffix(name, s) {
			return true
		}
	}
	return false
}

