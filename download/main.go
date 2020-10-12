package main

import (
	"download/core"
	"flag"
	"log"
	"os"
)

func main() {
	var (
		url       string
		protocol  string
		chunkSize int
		filename  string
		path      string
		limit     int64
	)
	flag.StringVar(&url, "u", "", "download url")
	flag.StringVar(&path, "o", "", "output dir")
	flag.StringVar(&filename, "filename", "", "filename")
	flag.StringVar(&protocol, "-protocol", core.HTTP, "download protocol")
	flag.IntVar(&chunkSize, "-chunk", 4, "chunk size (must greater than 0)")
	flag.Int64Var(&limit, "-limit", 10*1024, "limit download speed (must greater than 0)")
	flag.Parse()
	if url == "" || filename == "" || chunkSize <= 0 {
		flag.Usage()
	} else {
		if path == "" {
			path, _ = os.Getwd()
		}
		err := core.Download(&core.DownloadConfig{
			Link:      url,
			Protocol:  protocol,
			ChunkSize: chunkSize,
			FileName:  filename,
			Path:      path,
			Limit:     limit,
		})
		if err != nil {
			log.Fatalf("downloading error:%v", err)
		}
	}
}
