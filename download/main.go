package main

import (
	"flag"
	"log"
)

var chunkSize = 4

func main() {
	var (
		url string
	)
	flag.StringVar(&url, "u", "", "download url")
	flag.IntVar(&chunkSize, "-chunk", 4, "chunk size (must greater than 0)")
	flag.Parse()
	if url == "" || chunkSize <= 0 {
		flag.Usage()
	} else {
		downloader, err := NewHTTPDownloader(url, chunkSize)
		if err != nil {
			log.Fatalf("create downloader error:%v", err)
		}
		err = downloader.Download()
		if err != nil {
			log.Fatalf("downloading error:%v", err)
		}
	}
}
