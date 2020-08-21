package main

import (
	"fmt"
	"github.com/reujab/wallpaper"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"regexp"
	"runtime"
)

const IMAGE_API = "https://bing.ioliu.cn/v1?p=1&size=1&callback=z"
const UHD_API = "https://cn.bing.com/th?id=OHR.%s_UHD.jpg"

func main() {
	resp, err := http.Get(IMAGE_API)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	all, _ := ioutil.ReadAll(resp.Body)
	rgx := regexp.MustCompile("rb/(.*?)_1920x1080.jpg")
	imageName := rgx.FindSubmatch(all)[1]
	//
	sPath := storePath(string(imageName) + "_1920x1080.jpg")
	if _, err := os.Stat(sPath); os.IsNotExist(err) {
		response, _ := http.Get(fmt.Sprintf(UHD_API, imageName))
		defer response.Body.Close()
		readAll, _ := ioutil.ReadAll(response.Body)
		file, _ := os.Create(sPath)
		_, err := file.Write(readAll)
		if err != nil {
			panic(err)
		}
		err = file.Close()
		if err != nil {
			panic(err)
		}
	} else {
		fmt.Println("image exists!")
	}
	//
	wallpaper.SetFromFile(sPath)
}

func storePath(imageName string) string {
	homeDir, _ := os.UserHomeDir()
	storePath := path.Join(homeDir, "Pictures")
	goos := runtime.GOOS
	if goos == "windows" {
		return path.Join(storePath, "Saved Pictures", imageName)
	} else {
		return path.Join(storePath, imageName)
	}
}
