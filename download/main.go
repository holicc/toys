package main

import (
	"flag"
	"log"
)

func main() {
	var (
		url       string
		chunkSize int
		limit     int64
	)
	flag.StringVar(&url, "u", "", "download url")
	flag.IntVar(&chunkSize, "-chunk", 4, "chunk size (must greater than 0)")
	flag.Int64Var(&limit, "-limit", 10*1024, "limit download speed (must greater than 0)")
	flag.Parse()
	if url == "" || chunkSize <= 0 {
		flag.Usage()
	} else {
		downloader, err := NewHTTPDownloader(url, chunkSize, limit)
		if err != nil {
			log.Fatalf("create downloader error:%v", err)
		}
		err = downloader.Download()
		if err != nil {
			log.Fatalf("downloading error:%v", err)
		}
	}
}
