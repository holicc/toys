package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"spider"
	"strconv"
)

type V2ex struct {
	Title string
	Link  string
	Reply int
}

func main() {
	req, _ := http.NewRequest("GET", "https://www.v2ex.com/?tab=hot", nil)
	task, _ := spider.DoTask(req, "div.item")
	r := task.Map(mapToV2ex).Filter(filterByReply).Sort(sortByReply).Collect()
	//
	for i := range r {
		v := r[i].(*V2ex)
		fmt.Println(*v)
	}
}

func mapToV2ex(selection *goquery.Selection) interface{} {
	reply, _ := strconv.Atoi(selection.Find("a.count_livid").Text())
	return &V2ex{
		Title: selection.Find("span.item_title").Text(),
		Link:  selection.Find("span.item_title > a").AttrOr("href", ""),
		Reply: reply,
	}
}

func filterByReply(v2ex interface{}) bool {
	v := v2ex.(*V2ex)
	return v.Reply >= 20
}

func sortByReply(o1, o2 interface{}) int {
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
