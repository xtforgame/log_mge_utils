package tests

import (
	"testing"
	// "errors"
	// "bufio"
	"github.com/xtforgame/log_mge_utils/lmu"
	"github.com/xtforgame/log_mge_utils/logbuffers"
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
// 	logger, _ := lmu.NewLoggerT1(ls, lb)
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
	// file, err := os.OpenFile(localFsTestFolder+"/ff.log", os.O_WRONLY|os.O_CREATE, os.ModePerm)
	// if err != nil {
	// 	t.Error("open file FAIL")
	// }
	// file.Close()

	ls, _ := localfs.NewLocalFsStorer(localLoggerFsTestFolder)
	lb, _ := logbuffers.NewSimpleBuffer()
	logger, _ := lmu.NewLoggerT1(ls, lb)
	defer logger.Close()

	restoreCounter := 0
	writeDataCounter := 0

	AssertRestoreCounter := func(counter int) {
		if restoreCounter != counter {
			t.Fatal("expect restoreCounter: ", counter, ", actual:", writeDataCounter)
		}
	}

	AssertWriteCounter := func(counter int) {
		if writeDataCounter != counter {
			t.Fatal("expect writeDataCounter: ", counter, ", actual:", writeDataCounter)
		}
	}

	listener, _ := logger.CreateListener()
	listener.OnEvent(func(event *lmu.LoggerEvent) {
		data, ok := event.Data.(*lmu.DataEventPayload)
		if ok {
			if data.IsFromRestoring {
				restoreCounter++
			} else {
				writeDataCounter++
			}
			t.Log("cc", event.Position)
		}

	})
	logger.Write([]byte("dfdbbdbt\n"))
	err := listener.StartRestore()
	if err != nil {
		// t.Fatal(err)
	}

	AssertRestoreCounter(1)

	logger.Write([]byte("dfdbbdbt\n"))
	AssertWriteCounter(1)
	logger.Write([]byte("dfdbbdbt\n"))
	AssertWriteCounter(2)
}
