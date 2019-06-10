package logwatcher

import (
	"fmt"
	"github.com/xtforgame/log_mge_utils/listeners/listenert1"
	"github.com/xtforgame/log_mge_utils/lmu"
	"github.com/xtforgame/log_mge_utils/logbuffers"
	"github.com/xtforgame/log_mge_utils/loggers/loggert1"
	"github.com/xtforgame/log_mge_utils/logstorers/localfs"
	"path/filepath"
	// "os"
	"path"
	"regexp"
	"sync"
)

var logNameValidator = regexp.MustCompile(`^[0-9a-zA-Z_-]+$`)

type LoggerHepler struct {
	logPath             string
	localLogWatcherBase string
	loggers             map[string]lmu.Logger
	loggersMu           sync.Mutex
}

func CreateLoggerHepler(logPath string) *LoggerHepler {
	lh := &LoggerHepler{
		logPath:             logPath,
		localLogWatcherBase: filepath.Join(logPath, "log-watcher"),
	}
	// os.RemoveAll(lh.localLogWatcherBase)
	// os.MkdirAll(lh.localLogWatcherBase, os.ModePerm)

	lh.loggers = make(map[string]lmu.Logger)
	return lh
}

func (lh *LoggerHepler) FindLogger(logName string) lmu.Logger {
	validLogName := logNameValidator.FindString(logName)
	if validLogName == "" {
		return nil
	}
	lh.loggersMu.Lock()
	logger, _ := lh.loggers[logName]
	lh.loggersMu.Unlock()
	return logger
}

func (lh *LoggerHepler) CreateOrGetLogger(logName string) lmu.Logger {
	validLogName := logNameValidator.FindString(logName)
	if validLogName == "" {
		return nil
	}
	lh.loggersMu.Lock()
	logger, ok := lh.loggers[logName]
	if !ok {
		// os.RemoveAll(lh.localLogWatcherBase)
		logWatcherWorkDir := path.Join(lh.localLogWatcherBase, logName)
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
