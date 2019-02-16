package main

import (
	"github.com/xtforgame/log_mge_utils/tests/logwatcher"
	"github.com/xtforgame/log_mge_utils/utils"
)

func main() {
	defer utils.FinalReport()
	defer func() {
		if logwatcher.LoggerHeplerInst != nil {
			logwatcher.LoggerHeplerInst.Close()
		}
	}()
	// os.Exit(0)

	hs := logwatcher.NewHttpServer()
	hs.Init()
	hs.Start()
}
