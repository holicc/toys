package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const (
	DELIMIT            = byte(':')
	CONFIG_FILE_SUFFIX = ".json"
	TEMP_FILE_SUFFIX   = ".dtx"
	TEMP_DIR_NAME      = "downloader"
)

var (
	curDir    string
	tempDir   = filepath.Join(os.TempDir(), TEMP_DIR_NAME)
	tempFiles = make(map[string]string, 0)
)

func init() {
	//
	d, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	curDir = d
	//
	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		if os.Mkdir(tempDir, 0666) != nil {
			panic(err)
		}
	}
	//
	filepath.Walk(tempDir, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() && strings.HasSuffix(info.Name(), CONFIG_FILE_SUFFIX) {
			tempFiles[info.Name()] = path
		}
		return err
	})
}

type FileInfo struct {
	FileName     string
	FilePath     string
	Downloaded   int64
	Limit        int64
	TotalSize    int
	TempFile     *os.File `json:"-"`
	OriginalFile *os.File `json:"-"`
	sync.Mutex   `json:"-"`
}

type Downloader interface {
	Download() error
	Pause() error
	Continue() error
	Cancel() error
}

func (f *FileInfo) writeCache(d Downloader) error {
	defer f.Unlock()
	f.Lock()
	if f.TempFile == nil {
		file, err := createTempFile(f.FileName)
		if err != nil {
			return err
		}
		f.TempFile = file
	}
	bytes, err := json.Marshal(d)
	if err != nil {
		return err
	}
	_, err = f.TempFile.WriteAt(bytes, 0)
	return err
}

func getFromTempFile(filename string, d interface{}) error {
	if path := getTempFilePath(filename); path != "" {
		bytes, err := ioutil.ReadFile(path)
		if err == nil {
			err := json.Unmarshal(bytes, d)
			if err != nil {
				return err
			}
			return nil
		}
	}
	return errors.New("temp file does not exists")
}

func getTempFilePath(filename string) string {
	return tempFiles[filename+CONFIG_FILE_SUFFIX]
}

func createOriginalFile(d *HTTPDownloader) (*os.File, error) {
	file, err := os.OpenFile(d.FileInfo.FilePath+TEMP_FILE_SUFFIX, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func createTempFile(filename string) (*os.File, error) {
	return os.Create(filepath.Join(tempDir, filename+CONFIG_FILE_SUFFIX))
}
