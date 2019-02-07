package localfs

import (
	"github.com/xtforgame/log_mge_utils/lmu"
	"io"
	"os"
)

// ===============================

type LocalFsStoreReader struct {
	owner           *LocalFsStorer
	iterationNumber uint64
	lastEOFPos      int64
	file            *os.File
}

func (reader *LocalFsStoreReader) Reload() error {
	streamName := reader.GetOwner().GetStreamName(reader.GetIteration())
	file, err := os.Open(streamName + ".log")
	if err != nil {
		return err
	}
	stat, err := file.Stat()
	if err != nil {
		return err
	}
	reader.Close()
	reader.file = file
	reader.lastEOFPos = int64(stat.Size())
	return nil
}

func (reader *LocalFsStoreReader) GetOwner() lmu.Writer {
	return reader.owner
}

func (reader *LocalFsStoreReader) GetIteration() string {
	return lmu.IterationNumberToName(reader.iterationNumber)
}

func (reader *LocalFsStoreReader) GetCurrentPos() int64 {
	pos, _ := reader.file.Seek(0, io.SeekCurrent)
	return pos
}

func (reader *LocalFsStoreReader) GetLastEOFPos() int64 {
	return reader.lastEOFPos
}

func (reader *LocalFsStoreReader) Read(p []byte) (int, error) {
	return reader.file.Read(p)
}

func (reader *LocalFsStoreReader) Seek(off int64, whence int) (int64, error) {
	return reader.file.Seek(off, whence)
}

func (reader *LocalFsStoreReader) ReloadAndRead(p []byte) (int, error) {
	currentPos := reader.GetCurrentPos()
	err := reader.Reload()
	if err != nil {
		return 0, err
	}
	reader.file.Seek(currentPos, io.SeekStart)
	return reader.file.Read(p)
}

func (reader *LocalFsStoreReader) Close() {
	if reader.file != nil {
		reader.file.Close()
		reader.file = nil
	}
}
