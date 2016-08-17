// Package to abstract finding the current directory, Listing the current directory
// and recursing through given directory
package filesys

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

func Curdir() string {
	curdir, _ := os.Getwd()
	return curdir
}

func Listdir(curdir string) ([]os.FileInfo, error) {
	files, err := ioutil.ReadDir(curdir)
	return files, err
}

func Recurse(curdir string) ([]string, error) {
	fileList := []string{}
	if _, oserr := os.Stat(curdir); oserr != nil {
		return fileList, oserr
	}
	err := filepath.Walk(curdir, func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() {
			fileList = append(fileList, path)
		}
		return err
	})
	return fileList, err
}

func IsFile(file string) (bool, error) {
	f, open_err := os.Open(file)
	if open_err != nil {
		return false, open_err
	}
	fi, stat_err := f.Stat()
	if stat_err != nil {
		return false, stat_err
	}
	return fi.Mode().IsRegular(), nil
}
