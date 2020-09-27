package main

import "C"
import (
	"fmt"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var wg sync.WaitGroup
var client = http.DefaultClient

type HTTPDownloader struct {
	URL          string
	ETag         string
	LastModify   time.Time
	ChunkSize    int
	Chunks       []Chunk
	SupportRange bool
	FileInfo     *FileInfo
}

type Chunk struct {
	Index int
	From  int64
	To    int64
	Done  bool
}

func NewHTTPDownloader(url string, chunkSize int, limit int64) (*HTTPDownloader, error) {
	d := &HTTPDownloader{
		URL:       url,
		ChunkSize: chunkSize,
		Chunks:    make([]Chunk, chunkSize),
		FileInfo: &FileInfo{
			Limit: limit,
		},
	}
	if err := getFileInfo(d); err != nil {
		return nil, err
	}
	if err := getFromTempFile(d.FileInfo.FileName, d); err == nil {
		return d, nil
	}
	d.initChunks(chunkSize)

	return d, d.FileInfo.writeCache(d)
}

func (d *HTTPDownloader) Download() error {
	file, err := createOriginalFile(d)
	if err != nil {
		return err
	}
	d.FileInfo.OriginalFile = file
	log.Printf("Downloading start,file size: %.2f mb", float32(d.FileInfo.TotalSize)/1024.0/1024.0)
	process := time.NewTicker(1 * time.Second)
	update := time.NewTicker(10 * time.Second)
	defer process.Stop()
	defer update.Stop()
	go func() {
		for {
			<-process.C
			log.Printf("Process %.2f%%", 100*float32(d.FileInfo.Downloaded)/float32(d.FileInfo.TotalSize))
		}
	}()
	go func() {
		for {
			<-process.C
			d.FileInfo.writeCache(d)
		}
	}()
	if d.SupportRange {
		for i := range d.Chunks {
			c := &d.Chunks[i]
			if !c.Done {
				wg.Add(1)
				go func() {
					defer wg.Done()
					if err := d.downloadChunk(c); err != nil {
						log.Printf("download chunk error:%v", err)
					}
				}()
			}
		}
		wg.Wait()

		d.FileInfo.OriginalFile.Close()
		d.FileInfo.TempFile.Close()

		err := os.Rename(d.FileInfo.FilePath+TEMP_FILE_SUFFIX, d.FileInfo.FilePath)
		if err != nil {
			return err
		}
		err = os.Remove(getTempFilePath(d.FileInfo.FileName))
		if err != nil {
			return err
		}
	}
	return d.normalDownload()
}

func (d *HTTPDownloader) Pause() error {
	panic("implement me")
}

func (d *HTTPDownloader) Continue() error {
	panic("implement me")
}

func (d *HTTPDownloader) Cancel() error {
	panic("implement me")
}

func (d *HTTPDownloader) normalDownload() error {
	resp, err := client.Get(d.URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	file, err := os.Create(d.FileInfo.FilePath)
	if err != nil {
		return err
	}
	_, err = file.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func (d *HTTPDownloader) initChunks(chunkSize int) {
	pieces := d.FileInfo.TotalSize / chunkSize
	for i := range d.Chunks {
		size := (i * pieces) + pieces
		chunk := Chunk{
			Index: i,
			From:  int64(i * pieces),
			To:    int64(size),
			Done:  false,
		}
		d.Chunks[i] = chunk
	}
}

func (d *HTTPDownloader) downloadChunk(c *Chunk) error {
	req, err := http.NewRequest(http.MethodGet, d.URL, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Range", fmt.Sprintf("bytes=%v-%v", c.From, c.To))
	response, err := client.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	index := c.From
	for {
		bytes := make([]byte, d.FileInfo.Limit)
		read, err := response.Body.Read(bytes)
		if err != nil || read <= 0 {
			log.Printf("[%d] Chunk Done...", c.Index)
			break
		}
		write, err := d.FileInfo.OriginalFile.WriteAt(bytes, index)
		if err != nil || write <= 0 {
			log.Printf("write file error:%v", err)
			return err
		}
		index += int64(write)
		atomic.AddInt64(&d.FileInfo.Downloaded, int64(write))
	}
	c.Done = true
	return nil
}

func getFileInfo(d *HTTPDownloader) (err error) {
	req, err := http.NewRequest(http.MethodHead, d.URL, nil)
	if err != nil {
		return
	}
	r, err := client.Do(req)
	if err != nil {
		return
	}

	d.SupportRange = r.Header.Get("Accept-Ranges") == "bytes"
	d.FileInfo.TotalSize, _ = strconv.Atoi(r.Header.Get("Content-Length"))
	d.ETag = r.Header.Get("ETag")
	d.LastModify, _ = time.Parse(http.TimeFormat, r.Header.Get("Last-Modified"))

	// resolve duplicate file name
	filename, extension := parseFileInfoFrom(r)
	path := filepath.Join(curDir, filename+extension)
	for {
		var i int
		if _, err := os.Stat(path); os.IsNotExist(err) {
			d.FileInfo.FileName = filename + extension
			d.FileInfo.FilePath = path
			break
		} else {
			i++
			filename += fmt.Sprintf(" (%d) ", i+1)
			path = filepath.Join(curDir, filename+fmt.Sprintf(" (%d) ", i+1)+extension)
		}
	}

	return
}

func parseFileInfoFrom(resp *http.Response) (string, string) {
	var filename string
	contentDisposition := resp.Header.Get("Content-Disposition")
	if contentDisposition != "" {
		_, params, err := mime.ParseMediaType(contentDisposition)
		if err != nil {
			panic(err)
		}
		filename = params["filename"]
	} else {
		filename = filepath.Base(resp.Request.URL.Path)
	}
	index := strings.LastIndex(filename, ".")
	if index != -1 {
		return filename[:index], filename[index:]
	} else {
		return filename, ""
	}
}
