package tests

import (
	"testing"
	// "errors"
	// "bufio"
	// "github.com/xtforgame/log_mge_utils/lmu"
	"github.com/xtforgame/log_mge_utils/listeners/listenert1"
	"github.com/xtforgame/log_mge_utils/logbuffers"
	"github.com/xtforgame/log_mge_utils/loggers/loggert1"
	"github.com/xtforgame/log_mge_utils/logstorers/localfs"
	"os"
)

var localLoggerFsTestFolder = "../tmp/test/logger-fstest"
var localLoggerFsTestFile = localLoggerFsTestFolder + "/x.x"

// func TestLoggerLocalFsWrite(t *testing.T) {
// 	os.RemoveAll(localLoggerFsTestFolder)
// 	os.MkdirAll(localLoggerFsTestFolder, os.ModePerm)

// 	ls, _ := localfs.NewLocalFsStorer(localLoggerFsTestFolder)
// 	lb, _ := logbuffers.NewSimpleBuffer()
// 	logger, _ := lmu.NewLoggerT1(ls, lb, lmu.CreateListenerT1)
// 	defer logger.Close()

// 	logger.Write([]byte("dfdbbdbt\n"))
// 	logger.SwitchToNextIteration("")
// 	logger.Write([]byte("dfdbbdbt\n"))
// 	logger.Write([]byte("dfdbbdbt\n"))
// 	logger.Write([]byte("dfdbbdbt\n"))
// 	logger.SwitchToNextIteration("")
// 	logger.Close()
// }

func TestLoggerLocalFsWrite2(t *testing.T) {
	os.RemoveAll(localLoggerFsTestFolder)
	os.MkdirAll(localLoggerFsTestFolder, os.ModePerm)

	ls, _ := localfs.NewLocalFsStorer(localLoggerFsTestFolder)
	lb, _ := logbuffers.NewSimpleBuffer()
	logger, _ := loggert1.NewLoggerT1(ls, lb, listenert1.CreateListenerT1)
	defer logger.Close()

	listener1, _ := logger.CreateListener(nil)
	lh1 := CreateListenerHelperT1(t, "listener 1", listener1)

	listener2, _ := logger.CreateListener(nil)
	lh2 := CreateListenerHelperT1(t, "listener 2", listener2)

	logger.Write([]byte("dfdbbdbt\n"))
	err := lh1.Listener.Restore()
	if err != nil {
		// t.Fatal(err)
	}

	lh1.Listener.Listen()

	lh1.AssertRestoreEventCounter(1)

	lh2.AssertRestoreEventCounter(0)
	lh2.Listener.Restore()

	logger.Write([]byte("dfdbbdbt\n"))
	lh1.AssertWriteEventCounter(1)
	lh2.AssertWriteEventCounter(0)
	lh2.Listener.Listen()
	logger.Write([]byte("dfdbbdbt\n"))
	lh1.AssertWriteEventCounter(2)
	lh2.AssertRestoreEventCounter(2)

	logger.Write([]byte("dfdbbdbt\n"))
	lh2.AssertWriteEventCounter(1)
}
