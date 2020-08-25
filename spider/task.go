package spider

import (
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"net/url"
	"time"
)

var DefaultPipeline = &ConsolePipeline{}

type Pipeline interface {
	Process(v interface{})
}

type ConsolePipeline struct {
}

func (p *ConsolePipeline) Process(v interface{}) {
	log.Printf("Pipeline recevie :: %d %s", len(v.([]interface{})), v)
}

type Filter func(value interface{}) bool

type Key func(value interface{}) string

type Apply func(selection *goquery.Selection) interface{}

type Sort func(o1, o2 interface{}) int

type Next func(p int, selection *goquery.Selection) (string, int)

type Task struct {
	Request  *http.Request
	Response *http.Response
	//
	MainSelector string
	Selection    *goquery.Selection
	Document     *goquery.Document
	//
	Page      int
	NextURL   Next
	Pipelines []Pipeline
	//
	Values []interface{}
	//
	done bool
	//
	functions []func()
}

func NewTask(req *http.Request, selector string) (*Task, error) {
	return &Task{
		Request:      req,
		Pipelines:    make([]Pipeline, 0),
		functions:    make([]func(), 0),
		MainSelector: selector,
	}, nil
}

func (t *Task) Distinct(key Key) {
	m := make(map[string]interface{}, 0)
	for _, v := range t.Values {
		m[key(v)] = v
	}
	//
	t.Values = t.Values[:0]
	for _, v := range m {
		t.Values = append(t.Values, v)
	}
	//
}

func (t *Task) Filter(filter Filter) {
	t.functions = append(t.functions, func() {
		val := make([]interface{}, 0)
		for _, i := range t.Values {
			if filter(i) {
				val = append(val, i)
			}
		}
		t.Values = val
	})
}

func (t *Task) Map(apply Apply) {
	t.functions = append(t.functions, func() {
		t.Selection.Each(func(i int, selection *goquery.Selection) {
			t.Values = append(t.Values, apply(selection))
		})
	})
}

func (t *Task) Pipeline(p Pipeline) {
	t.Pipelines = append(t.Pipelines, p)
}

func (t *Task) Sort(sort Sort) {
	t.functions = append(t.functions, func() {
		values := t.Values
		//TODO improvement
		for i := range values {
			for j := i + 1; j < len(values); j++ {
				if sort(values[i], values[j]) > 1 {
					t := values[i]
					values[i] = values[j]
					values[j] = t
				}
			}
		}
	})
}

func (t *Task) Collect() []interface{} {
	if t.done {
		return t.Values
	} else {
		//do function
		t.doFunc()
		//active pipeline after all functions
		t.activePipelines()
		//
		t.finish()
		//
		return t.Values
	}
}

func (t *Task) RepeatWithBreak(duration time.Duration, f func(t *Task) bool) {
	if t.done {
		return
	}
	go func() {
		ticker := time.NewTicker(duration)
		for {
			<-ticker.C
			if f(t) {
				log.Println("break.")
				break
			}
			//
			t.process()
			//
		}
	}()
}

func (t *Task) process() bool {
	//
	t.fetchSource()
	//
	t.doFunc()
	//
	go t.activePipelines()
	//
	t.finish()
	//
	return t.done
}

func (t *Task) finish() {
	t.done = true
}

func (t *Task) activePipelines() {
	for i := range t.Pipelines {
		pipeline := t.Pipelines[i]
		//TODO copy values
		go pipeline.Process(t.Values)
	}
}

func (t *Task) doFunc() {
	//
	for _, f := range t.functions {
		f()
	}
	if !t.done {
		if nextURL, page := t.NextURL(t.Page, t.Selection); nextURL != "" {
			parse, err := url.Parse(nextURL)
			if err != nil {
				return
			}
			t.Page = page
			t.Request.URL = parse
		}
	}
}

func (t *Task) fetchSource() {
	response, err := http.DefaultClient.Do(t.Request)
	if err != nil {
		log.Println("http client request error", err.Error())
		return
	}
	defer response.Body.Close()
	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Println("go query document reader failed", err.Error())
		return
	}
	selection := document.Find(t.MainSelector)
	t.Response = response
	t.Selection = selection
}
