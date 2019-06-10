package main

import (
	"flag"
	"fmt"
	"github.com/xtforgame/log_mge_utils/tests/logwatcher"
	"github.com/xtforgame/log_mge_utils/utils"
	"os"
	"path/filepath"
)

func NormalizePath(path string) (string, error) {
	// if path == "" {
	// 	return "", errors.New("Invalid path : " + path)
	// }
	var err error = nil
	if !filepath.IsAbs(path) {
		var base string
		base, err = os.Getwd()
		if err != nil {
			return "", err
		}
		path = filepath.Join(base, path)
	}

	path, err = filepath.Abs(path)
	if err != nil {
		return "", err
	}
	return path, nil
}

// https://github.com/restic/restic/blob/master/build.go
func DirectoryExists(dirname string) bool {
	stat, err := os.Stat(dirname)
	if err != nil && os.IsNotExist(err) {
		return false
	}

	return stat.IsDir()
}

func main() {
	flag.Parse()
	args := flag.Args()
	logPath, _ := NormalizePath("./logs")
	webPath, _ := NormalizePath("./web")
	if len(args) > 0 {
		var err error
		logPath, err = NormalizePath(args[0])
		if err != nil {
			fmt.Println("wrong logPath :", logPath)
			return
		}
	}

	if len(args) > 1 {
		var err error
		webPath, err = NormalizePath(args[1])
		if err != nil {
			fmt.Println("wrong webPath :", webPath)
			return
		}
	}

	defer utils.FinalReport()
	defer func() {
		if logwatcher.LoggerHeplerInst != nil {
			logwatcher.LoggerHeplerInst.Close()
		}
	}()
	// os.Exit(0)

	hs := logwatcher.NewHttpServer(logPath, webPath)
	hs.Init()
	hs.Start()
}
