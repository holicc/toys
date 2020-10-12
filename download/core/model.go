package core

type DownloadConfig struct {
	Link      string
	Path      string
	Protocol  string
	FileName  string
	ChunkSize int
	Limit     int64
}

type DownloadStatus struct {
	TotalSize  int64
	Downloaded int64
	Chunks     []Chunk
	Config     *DownloadConfig
}

type Chunk struct {
	Index int
	From  int64
	To    int64
	Done  bool
}
