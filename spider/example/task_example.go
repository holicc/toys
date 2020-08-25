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
		Value: "",
	})
	req.AddCookie(&http.Cookie{
		Name:  "V2EX_LANG",
		Value: "zhcn",
	})
	req.AddCookie(&http.Cookie{
		Name:  "A2",
		Value: "",
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
