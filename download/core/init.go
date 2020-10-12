package core

import (
	"os"
	"path/filepath"
	"strings"
)

var (
	tempDir   = filepath.Join(os.TempDir(), TEMP_DIR_NAME)
	tempFiles = make(map[string]string, 0)
)

func init() {
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
