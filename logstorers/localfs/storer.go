package localfs

import (
	"errors"
	"github.com/xtforgame/log_mge_utils/fshelper"
	"github.com/xtforgame/log_mge_utils/lmu"
	"os"
	"path/filepath"
	"strings"
)

func StreamNameToNumber(streamName string) uint64 {
	return lmu.IterationNameToNumber(strings.Split(streamName, ".")[0])
}

// ===============================

type LocalFsStorer struct {
	path            string
	iterationNumber uint64
	file            *os.File
	Files           []os.FileInfo
}

// ===============================

func NewLocalFsStorer(path string) (*LocalFsStorer, error) {
	lfs := &LocalFsStorer{path: path}
	if err := lfs.Reinit(); err != nil {
		return nil, err
	}

	return lfs, nil
}

// ==============================

func (lfs *LocalFsStorer) Reinit() error {
	path := lfs.GetPath()
	if !fshelper.DirectoryExists(path) {
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			// fmt.Printf("MkdirAll(%v)\n", path)
			return err
		}
	}
	var iterationNumber uint64
	dirInfo, _ := fshelper.ListDir(path)
	for _, fileInfo := range dirInfo.Files {
		fileNumber := StreamNameToNumber(fileInfo.Name())
		if iterationNumber < fileNumber {
			iterationNumber = fileNumber
		}
	}

	if iterationNumber == 0 {
		iterationNumber = 1
	}
	lfs.iterationNumber = iterationNumber

	file, err := os.OpenFile(lfs.GetStreamName(lfs.GetIteration())+".log", os.O_APPEND|os.O_WRONLY|os.O_CREATE /*|os.O_SYNC*/ /*|os.O_EXCL*/, os.ModePerm)
	if err != nil {
		return err
	}

	lfs.file = file
	lfs.Files = dirInfo.Files
	return nil
}

func (lfs *LocalFsStorer) GetStoreType() string {
	return "local_fs"
}

func (lfs *LocalFsStorer) GetPath() string {
	return lfs.path
}

func (lfs *LocalFsStorer) GetIteration() string {
	return lmu.IterationNumberToName(lfs.iterationNumber)
}

func (lfs *LocalFsStorer) GetStreamName(iteration string) string {
	return filepath.Join(lfs.GetPath(), iteration)
}

func (lfs *LocalFsStorer) CreateReader() (lmu.SReader, error) {
	reader := &LocalFsStoreReader{
		owner:           lfs,
		iterationNumber: lfs.iterationNumber,
	}
	err := reader.Reload()
	if err != nil {
		return nil, err
	}
	return reader, nil
}

func (lfs *LocalFsStorer) SwitchToNextIteration(iteration string) error {
	lfs.Close()
	lfs.iterationNumber += 1
	file, err := os.OpenFile(lfs.GetStreamName(lfs.GetIteration())+".log", os.O_APPEND|os.O_WRONLY|os.O_CREATE /*|os.O_SYNC*/ /*|os.O_EXCL*/, os.ModePerm)
	if err != nil {
		lfs.iterationNumber -= 1
		return err
	}
	lfs.file = file
	return nil
}

func (lfs *LocalFsStorer) Write(p []byte) (int, error) {
	if lfs.file == nil {
		return 0, errors.New("No File Opened")
	}
	n, err := lfs.file.Write(p)
	if err == nil && n > 0 {
		// lfs.file.Sync()
	}
	return n, err
}

func (lfs *LocalFsStorer) RemoveStore() {
	lfs.Close()
	path := lfs.GetPath()
	if fshelper.DirectoryExists(path) {
		os.RemoveAll(path)
	}
}

func (lfs *LocalFsStorer) Close() {
	if lfs.file != nil {
		lfs.file.Sync()
		lfs.file.Close()
		lfs.file = nil
	}
}
