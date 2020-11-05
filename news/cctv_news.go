package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

const URL = "https://news.cctv.com/2019/07/gaiban/cmsdatainterface/page/news_1.jsonp?cb=news"

type Response struct {
	Data CCTVNews `json:"data"`
}

type CCTVNews struct {
	Total int    `json:"total"`
	List  []News `json:"list"`
}

type News struct {
	ID        string `json:"id"`
	Image2    string `json:"image2"`
	Title     string `json:"title"`
	Keywords  string `json:"keywords"`
	Count     string `json:"count"`
	ExtField  string `json:"ext_field"`
	Image     string `json:"image"`
	FocusDate string `json:"focus_date"`
	Image3    string `json:"image3"`
	Brief     string `json:"brief"`
	URL       string `json:"url"`
}

func GetNews(c *gin.Context) {
	news, err := news()
	if err != nil {
		c.JSON(500, err.Error())
	}
	c.JSON(200, news)
}

func news() (*CCTVNews, error) {
	resp, err := http.Get(URL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	content := string(d)
	var r Response
	err = json.Unmarshal([]byte(content[5:len(content)-1]), &r)
	if err != nil {
		return nil, err
	}
	return &r.Data, nil
}
