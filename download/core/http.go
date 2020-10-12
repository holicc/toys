package core

import "C"
import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync/atomic"
)

var client = http.DefaultClient

type HTTPDownloader int

func (d HTTPDownloader) download(c *Chunk, file *os.File, status *DownloadStatus) error {
	req, err := http.NewRequest(http.MethodGet, status.Config.Link, nil)
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
		data := make([]byte, status.Config.Limit)
		read, err := response.Body.Read(data)
		if err != nil || read <= 0 {
			log.Printf("[%d] Chunk Done...", c.Index)
			break
		}
		write, err := file.WriteAt(data, index)
		index += int64(write)
		if err != nil || write <= 0 {
			log.Printf("write file error:%v", err)
		}
		atomic.AddInt64(&status.Downloaded, int64(write))
	}
	c.Done = true
	return nil
}

func head(status *DownloadStatus) (err error) {
	req, err := http.NewRequest(http.MethodHead, status.Config.Link, nil)
	if err != nil {
		return
	}
	r, err := client.Do(req)
	if err != nil {
		return
	}

	t, _ := strconv.Atoi(r.Header.Get("Content-Length"))
	status.TotalSize = int64(t)

	return
}
