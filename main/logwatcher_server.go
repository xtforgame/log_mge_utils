package main

import (
	"github.com/xtforgame/log_mge_utils/tests/logwatcher"
	"github.com/xtforgame/log_mge_utils/utils"
)

func main() {
	defer utils.FinalReport()
	// os.Exit(0)

	hs := logwatcher.NewHttpServer()
	hs.Init()
	hs.Start()
}
