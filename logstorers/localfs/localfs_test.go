package localfs

import (
	"testing"
	// "errors"
	// "bufio"
	// "github.com/xtforgame/log_mge_utils/lmu"
	"io"
	"os"
)

var localFsTestFolder = "../../tmp/test/fstest"
var localFsTestFile = localFsTestFolder + "/x.x"

// func TestLocalFsWrite1(t *testing.T) {
// 	os.RemoveAll(localFsTestFolder)
// 	os.MkdirAll(localFsTestFolder, os.ModePerm)
// 	ls, _ := logstorers.NewLocalFsStorer(localFsTestFolder)

// 	file, err := os.OpenFile("../tmp/test/fstest/eee", os.O_WRONLY|os.O_CREATE|os.O_EXCL, os.ModePerm)
// 	if err != nil {
// 		t.Error("open file FAIL")
// 	}
// 	file.Close()

// 	t.Log(ls)
// }

func TestLocalFsWrite2(t *testing.T) {
	os.RemoveAll(localFsTestFolder)
	os.MkdirAll(localFsTestFolder, os.ModePerm)
	// file, err := os.OpenFile(localFsTestFolder+"/ff.log", os.O_WRONLY|os.O_CREATE, os.ModePerm)
	// if err != nil {
	// 	t.Error("open file FAIL")
	// }
	// file.Close()

	ls, _ := NewLocalFsStorer(localFsTestFolder)
	defer ls.Close()

	AssertInt64Equal := func(varName string, _var int64, expectedVal int64, err error) {
		if _var != expectedVal {
			t.Log("err:", err)
			t.Fatal("expect "+varName+": ", expectedVal, ", actual:", _var)
		}
	}

	AssertIntEqual := func(varName string, _var int, expectedVal int, err error) {
		if _var != expectedVal {
			t.Log("err:", err)
			t.Fatal("expect "+varName+": ", expectedVal, ", actual:", _var)
		}
	}

	ls.Write([]byte("dfdbbdbt\n"))
	reader, _ := ls.CreateReader()
	var b = make([]byte, 20)
	s, err := reader.Seek(3, io.SeekStart)
	AssertInt64Equal("seek", s, 3, err)

	c, err := reader.Seek(0, io.SeekCurrent)
	AssertInt64Equal("current", c, 3, err)

	n, err := reader.Read(b)
	AssertIntEqual("current", n, 6, err)

	n, err = reader.Read(b)
	if err != io.EOF {
		t.Fatal("err != io.EOF")
	}
	ls.Write([]byte("dfdbbdbt\n"))
	n, err = reader.ReloadAndRead(b)
	AssertIntEqual("current", n, 9, err)

	ls.Close()
	ls.Write([]byte("dfdbbdbt\n"))
	n, err = reader.ReloadAndRead(b)
	AssertIntEqual("current", n, 0, err)
}
