package logwatcher

import (
	"fmt"
	"github.com/xtforgame/log_mge_utils/listeners/listenert1"
	"github.com/xtforgame/log_mge_utils/lmu"
	"github.com/xtforgame/log_mge_utils/logbuffers"
	"github.com/xtforgame/log_mge_utils/loggers/loggert1"
	"github.com/xtforgame/log_mge_utils/logstorers/localfs"
	// "os"
	"path"
	"regexp"
	"sync"
)

var localLogWatcherBase = "./tmp/test/log-watcher"
var logNameValidator = regexp.MustCompile(`^[0-9a-zA-Z_-]+$`)

type LoggerHepler struct {
	loggers   map[string]lmu.Logger
	loggersMu sync.Mutex
}

func CreateLoggerHepler() *LoggerHepler {
	// os.RemoveAll(localLogWatcherBase)
	// os.MkdirAll(localLogWatcherBase, os.ModePerm)

	lh := &LoggerHepler{}
	lh.loggers = make(map[string]lmu.Logger)
	return lh
}

func (lh *LoggerHepler) GetLogger(logName string) lmu.Logger {
	validLogName := logNameValidator.FindString(logName)
	if validLogName == "" {
		return nil
	}
	lh.loggersMu.Lock()
	logger, ok := lh.loggers[logName]
	if !ok {
		// os.RemoveAll(localLogWatcherBase)
		logWatcherWorkDir := path.Join(localLogWatcherBase, logName)
		// os.MkdirAll(logWatcherWorkDir, os.ModePerm)

		ls, _ := localfs.NewLocalFsStorer(logWatcherWorkDir)
		lb, _ := logbuffers.NewSimpleBuffer()
		logger, _ = loggert1.NewLoggerT1(ls, lb, listenert1.CreateListenerT1)
		lh.loggers[logName] = logger
	}
	lh.loggersMu.Unlock()
	return logger
}

func (lh *LoggerHepler) RemoveAndCloseLogger(logName string) {
	lh.loggersMu.Lock()
	logger, ok := lh.loggers[logName]
	if ok {
		logger.RemoveAndCloseLogger()
		delete(lh.loggers, logName)
	}
	lh.loggersMu.Unlock()
}

func (lh *LoggerHepler) Close() {
	lh.loggersMu.Lock()
	for k, v := range lh.loggers {
		if v != nil {
			v.Close()
			lh.loggers[k] = nil
		}
	}
	lh.loggers = make(map[string]lmu.Logger)
	fmt.Println("(lh *LoggerHepler) Close()")
	lh.loggersMu.Unlock()
}

var LoggerHeplerInst *LoggerHepler
