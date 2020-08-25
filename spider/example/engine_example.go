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
	engine := spider.NewEngine(16)
	//
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
	task.Map(MapToV2ex)
	task.NextURL = NextURL
	//
	err := engine.AddTask(task)
	if err != nil {
		log.Println(err.Error())
	}
	//
	engine.Run()
	time.Sleep(3 * time.Second)
	engine.Stop()
}

func NextURL(page int, selection *goquery.Selection) (string, int) {
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
