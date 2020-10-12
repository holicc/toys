package core

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	DELIMIT            = byte(':')
	CONFIG_FILE_SUFFIX = ".json"
	TEMP_FILE_SUFFIX   = ".dtx"
	TEMP_DIR_NAME      = "downloader"

	HTTP       = "http"
	BitTorrent = "BT"
)

type Downloader interface {
	download(*Chunk, *os.File, *DownloadStatus) error
}

func Download(cfg *DownloadConfig) error {
	var wg sync.WaitGroup
	wg.Add(cfg.ChunkSize)
	status := &DownloadStatus{
		Config: cfg,
		Chunks: make([]Chunk, cfg.ChunkSize),
	}
	err := refreshStatus(cfg, status)
	if err != nil {
		return err
	}
	pieces := status.TotalSize / int64(cfg.ChunkSize)
	for i := range status.Chunks {
		j := int64(i)
		status.Chunks[i] = Chunk{
			Index: i,
			From:  j * pieces,
			To:    (j * pieces) + pieces,
		}
	}
	// create target temp file
	file, err := getTargetTempFile(cfg.Path, cfg.FileName)
	if err != nil {
		return err
	}
	// ticker record process
	log.Printf("Downloading start,file size: %.2f mb", float32(status.TotalSize)/1024.0/1024.0)
	process := time.NewTicker(1 * time.Second)
	defer process.Stop()
	go func() {
		for {
			<-process.C
			log.Printf("Process %.2f%%", 100*float32(status.Downloaded)/float32(status.TotalSize))
		}
	}()
	// downloading
	downloader := getDownloader(cfg.Protocol)
	for i := range status.Chunks {
		c := &status.Chunks[i]
		if !c.Done {
			go func() {
				defer wg.Done()
				if err := downloader.download(c, file, status); err != nil {
					log.Printf("download chunk error:%v", err)
				}

			}()
		}
	}
	wg.Wait()
	file.Close()
	path := filepath.Join(cfg.Path, cfg.FileName)
	err = os.Rename(path+TEMP_FILE_SUFFIX, path)
	if err != nil {
		return err
	}
	err = os.Remove(getTempFilePath(cfg.FileName))
	if err != nil {
		return err
	}
	return nil
}

func refreshStatus(cfg *DownloadConfig, status *DownloadStatus) error {
	switch cfg.Protocol {
	case HTTP:
		return head(status)
	default:
		return errors.New("unknown protocol")
	}
}

func getDownloader(protocol string) Downloader {
	switch protocol {
	case HTTP:
		return HTTPDownloader(0)
	}
	return nil
}

func getTargetTempFile(path string, name string) (*os.File, error) {
	file, err := os.OpenFile(filepath.Join(path, name)+TEMP_FILE_SUFFIX, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func getTempFilePath(filename string) string {
	return tempFiles[filename+CONFIG_FILE_SUFFIX]
}
