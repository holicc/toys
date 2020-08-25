package main

import (
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"spider"
	"strconv"
	"time"
)

type V2ex struct {
	Title string
	Link  string
	Reply int
}

func main() {
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
	task.NextURL = Next
	//
	task.Map(MapToStruct)
	task.Pipeline(spider.DefaultPipeline)
	//
	task.RepeatWithBreak(1*time.Second, func(t *spider.Task) bool {
		log.Println("repeat do ", t.Request.URL)
		return t.Page == 20
	})

	task.Wait()
}

func Next(page int, selection *goquery.Selection) (string, int) {
	return "https://www.v2ex.com/recent?p=" + strconv.Itoa(page+1), page + 1
}

func MapToStruct(selection *goquery.Selection) interface{} {
	reply, _ := strconv.Atoi(selection.Find("a.count_livid").Text())
	return &V2ex{
		Title: selection.Find("span.item_title").Text(),
		Link:  selection.Find("span.item_title > a").AttrOr("href", ""),
		Reply: reply,
	}
}
