package spider

import (
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
)

type Task struct {
	Selection *goquery.Selection
	Response  *http.Response
	values    []interface{}
}

type Filter func(value interface{}) bool

type Apply func(selection *goquery.Selection) interface{}

type Sort func(o1, o2 interface{}) int

func init() {
	log.SetPrefix("[::Spider::]")
	log.SetFlags(log.LstdFlags | log.Llongfile)
}

func DoTask(req *http.Request, selector string) (*Task, error) {
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return nil, err
	}
	selection := document.Find(selector)
	return &Task{
		Response:  response,
		Selection: selection,
	}, nil
}

func (t *Task) Filter(filter Filter) *Task {
	val := make([]interface{}, 0)
	for _, i := range t.values {
		if filter(i) {
			val = append(val, i)
		}
	}
	t.values = val
	return t
}

func (t *Task) Map(apply Apply) *Task {
	t.Selection.Each(func(i int, selection *goquery.Selection) {
		v := apply(selection)
		t.values = append(t.values, v)
	})
	return t
}

func (t *Task) Collect() []interface{} {
	return t.values
}

func (t *Task) Sort(sort Sort) *Task {
	values := t.values
	for i := range values {
		for j := i + 1; j < len(values); j++ {
			if sort(values[i], values[j]) > 1 {
				t := values[i]
				values[i] = values[j]
				values[j] = t
			}
		}
	}
	return t
}
