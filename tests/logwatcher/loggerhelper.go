package logwatcher

import (
	"fmt"
	"github.com/xtforgame/log_mge_utils/listeners/listenert1"
	"github.com/xtforgame/log_mge_utils/lmu"
	"github.com/xtforgame/log_mge_utils/logbuffers"
	"github.com/xtforgame/log_mge_utils/loggers/loggert1"
	"github.com/xtforgame/log_mge_utils/logstorers/localfs"
	"os"
)

var localLoggWatcherFolder = "./tmp/test/log-watcher"

type LoggerHepler struct {
	Logger lmu.Logger
}

func CreateLoggerHepler() *LoggerHepler {
	os.RemoveAll(localLoggWatcherFolder)
	os.MkdirAll(localLoggWatcherFolder, os.ModePerm)

	ls, _ := localfs.NewLocalFsStorer(localLoggWatcherFolder)
	lb, _ := logbuffers.NewSimpleBuffer()
	logger, _ := loggert1.NewLoggerT1(ls, lb, listenert1.CreateListenerT1)

	lh := &LoggerHepler{
		Logger: logger,
	}
	return lh
}

func (lh *LoggerHepler) Close() {
	if lh.Logger != nil {
		lh.Logger.Close()
		lh.Logger = nil
	}
	fmt.Println("(lh *LoggerHepler) Close()")
}

var LoggerHeplerInst *LoggerHepler
