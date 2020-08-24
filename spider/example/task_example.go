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
	req, _ := http.NewRequest("GET", "https://www.v2ex.com/?tab=hot", nil)
	task, _ := spider.NewTask(req, "div.item")
	task.Map(MapToStruct)
	task.Filter(FilterByReply)
	task.Sort(SortByReply)
	task.Pipeline(spider.DefaultPipeline)
	task.Distinct(Key)
	task.Next(Next)
	//
	task.RepeatWithBreak(1*time.Second, func(t *spider.Task) bool {
		log.Println("repeat do.")
		return len(t.Values) > 30
	})
	//
	//fmt.Println(task.Collect())

	time.Sleep(100 * time.Second)
}

func Next(selection goquery.Selection) string {

	return "https://www.v2ex.com/recent?p=3"
}

func MapToStruct(selection *goquery.Selection) interface{} {
	reply, _ := strconv.Atoi(selection.Find("a.count_livid").Text())
	return &V2ex{
		Title: selection.Find("span.item_title").Text(),
		Link:  selection.Find("span.item_title > a").AttrOr("href", ""),
		Reply: reply,
	}
}

func FilterByReply(v2ex interface{}) bool {
	v := v2ex.(*V2ex)
	return v.Reply >= 20
}

func SortByReply(o1, o2 interface{}) int {
	a := o1.(*V2ex)
	b := o2.(*V2ex)
	if a.Reply > b.Reply {
		return 1
	} else if a.Reply < b.Reply {
		return -1
	} else {
		return 0
	}
}

func Key(o interface{}) string {
	return o.(*V2ex).Title
}
