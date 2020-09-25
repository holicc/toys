package main

import (
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"spider"
	"strconv"
)

type V2ex struct {
	Title string `selector:""`
	Link  string `selector:""`
	Reply int    `selector:""`
}

func main() {
	engine := spider.NewEngine(16)
	//
	req, _ := http.NewRequest("GET", "https://www.v2ex.com/recent", nil)
	req.AddCookie(&http.Cookie{
		Name:  "PB3_SESSION",
		Value: "2|1:0|10:1598230663|11:PB3_SESSION|40:djJleDoxMTguMTE0LjI1My4xNDE6MzU5NzA1NDI=|abe2637893e203ff4d7cf57ef4894c6bf854a3345065b0b203fe7bb8bf1bb49e",
	})
	req.AddCookie(&http.Cookie{
		Name:  "V2EX_LANG",
		Value: "zhcn",
	})
	req.AddCookie(&http.Cookie{
		Name:  "A2",
		Value: "2|1:0|10:1596016506|2:A2|56:YTMwMDU4ODNlYjFkYmQzODU3MWVkZWIzNzQ5OTkzNmYzN2FjNzViZA==|9adc5fb9f8d7960db311fb11b0ed1cce63efa9fdd71d47850c9698342d994e4f",
	})
	task, _ := spider.NewTask(req, "div.item")
	//
	task.Map(MapToV2ex)
	task.NextURL = NextURL
	//
	err := engine.AddTask(task)
	if err != nil {
		log.Println(err.Error())
	}
	//
	engine.Run()
}

func NextURL(page int, selection *goquery.Selection) (string, int) {
	if page >= 10 {
		return "", -1
	}
	return "https://www.v2ex.com/recent?p=" + strconv.Itoa(page+1), page + 1
}

func MapToV2ex(selection *goquery.Selection) interface{} {
	reply, _ := strconv.Atoi(selection.Find("a.count_livid").Text())
	return &V2ex{
		Title: selection.Find("span.item_title").Text(),
		Link:  selection.Find("span.item_title > a").AttrOr("href", ""),
		Reply: reply,
	}
}
