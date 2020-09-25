package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

const (
	DELIMIT            = byte(':')
	CONFIG_FILE_SUFFIX = ".json"
	TEMP_FILE_SUFFIX   = ".dtx"
)

var dir string

func init() {
	d, err := os.Getwd()
	if err != nil {
		log.Fatalf("get dir error:%v", err)
	}
	dir = d
}

type FileInfo struct {
	FileName     string
	FilePath     string
	Downloaded   int64
	Limit        int
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
		file, err := createTempFile(f.FilePath + CONFIG_FILE_SUFFIX)
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

func getFromTempFile(path string, d interface{}) error {
	bytes, err := ioutil.ReadFile(path)
	if err == nil {
		err := json.Unmarshal(bytes, d)
		if err != nil {
			return err
		}
		return nil
	}
	return err
}

func createOriginalFile(d *HTTPDownloader) (*os.File, error) {
	file, err := os.OpenFile(d.FileInfo.FilePath+TEMP_FILE_SUFFIX, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func createTempFile(path string) (*os.File, error) {
	return os.Create(path)
}
