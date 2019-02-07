package fshelper

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

type DirInfo struct {
	Path  string
	Files []os.FileInfo
	Dirs  []os.FileInfo
}

func ListDir(path string) (*DirInfo, error) {
	fileInfos, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	files := []os.FileInfo{}
	dirs := []os.FileInfo{}
	for _, fileInfo := range fileInfos {
		if fileInfo.IsDir() {
			dirs = append(dirs, fileInfo)
		} else {
			files = append(files, fileInfo)
		}
	}
	return &DirInfo{
		Path:  path,
		Files: files,
		Dirs:  dirs,
	}, err
}

// https://github.com/restic/restic/blob/master/build.go
func DirectoryExists(dirname string) bool {
	stat, err := os.Stat(dirname)
	if err != nil && os.IsNotExist(err) {
		return false
	}

	return stat.IsDir()
}

// CopyFile creates dst from src, preserving file attributes and timestamps.
func CopyFile(src, dst string) error {
	fi, err := os.Stat(src)
	if err != nil {
		return err
	}

	fsrc, err := os.Open(src)
	if err != nil {
		return err
	}

	if err = os.MkdirAll(filepath.Dir(dst), os.ModePerm); err != nil {
		// fmt.Printf("MkdirAll(%v)\n", filepath.Dir(dst))
		return err
	}

	fdst, err := os.Create(dst)
	if err != nil {
		return err
	}

	if _, err = io.Copy(fdst, fsrc); err != nil {
		return err
	}

	if err == nil {
		err = fsrc.Close()
	}

	if err == nil {
		err = fdst.Close()
	}

	if err == nil {
		err = os.Chmod(dst, fi.Mode())
	}

	if err == nil {
		err = os.Chtimes(dst, fi.ModTime(), fi.ModTime())
	}

	return nil
}
