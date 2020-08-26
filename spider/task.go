package spider

import (
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
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
	functions []func(t *Task)
}

func NewTask(req *http.Request, selector string) (*Task, error) {
	return &Task{
		Request:      req,
		Pipelines:    make([]Pipeline, 0),
		functions:    make([]func(t *Task), 0),
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
	t.functions = append(t.functions, func(self *Task) {
		val := make([]interface{}, 0)
		for _, i := range self.Values {
			if filter(i) {
				val = append(val, i)
			}
		}
		self.Values = val
	})
}

func (t *Task) Map(apply Apply) {
	t.functions = append(t.functions, func(self *Task) {
		self.Selection.Each(func(i int, selection *goquery.Selection) {
			self.Values = append(self.Values, apply(selection))
		})
	})
}

func (t *Task) Pipeline(p Pipeline) {
	t.Pipelines = append(t.Pipelines, p)
}

func (t *Task) Sort(sort Sort) {
	t.functions = append(t.functions, func(self *Task) {
		values := self.Values
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
		t.process()
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

func (t *Task) process() error {
	//
	err := t.fetchSource()
	if err != nil {
		return err
	}
	//
	err = t.doFunc()
	if err != nil {
		return err
	}
	//
	t.activePipelines()
	//
	t.finish()
	//
	return nil
}

func (t *Task) finish() {
	t.done = true
}

func (t *Task) activePipelines() {
	for i := range t.Pipelines {
		pipeline := t.Pipelines[i]
		//TODO copy values
		pipeline.Process(t.Values)
	}
}

func (t *Task) doFunc() error {
	//
	for _, f := range t.functions {
		f(t)
	}
	return nil
}

func (t *Task) fetchSource() error {
	log.Println("URL===>", t.Request.URL)
	response, err := http.DefaultClient.Do(t.Request)
	if err != nil {
		log.Println("http client request error", err.Error())
		return err
	}
	defer response.Body.Close()
	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Println("go query document reader failed", err.Error())
		return err
	}
	selection := document.Find(t.MainSelector)
	t.Response = response
	t.Selection = selection
	return nil
}
