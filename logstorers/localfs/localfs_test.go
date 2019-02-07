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

	ls.Write([]byte("dfdbbdbt\n"))
	reader, _ := ls.CreateReader()
	var b = make([]byte, 20)
	s, err := reader.Seek(3, io.SeekStart)
	if s != 3 {
		t.Log("s:", s)
		t.Log("err:", err)
		t.Fatal("s != 3")
	}
	c, err := reader.Seek(0, io.SeekCurrent)
	if c != 3 {
		t.Fatal("c != 3")
	}
	n, err := reader.Read(b)
	if n != 6 {
		t.Log("n:", n)
		t.Log("err:", err)
		t.Fatal("n != 6")
	}
	n, err = reader.Read(b)
	if err != io.EOF {
		t.Fatal("err != io.EOF")
	}
	ls.Write([]byte("dfdbbdbt\n"))
	n, err = reader.ReloadAndRead(b)
	if n != 9 {
		t.Log("n:", n)
		t.Log("err:", err)
		t.Fatal("n != 9")
	}
	ls.Close()
	ls.Write([]byte("dfdbbdbt\n"))
	n, err = reader.ReloadAndRead(b)
	if n != 0 {
		t.Log("n:", n)
		t.Log("err:", err)
		t.Fatal("n != 0")
	}
}
