package spider

import (
	"errors"
	"log"
	"net/url"
	"sync"
)

var wg sync.WaitGroup

type Engine struct {
	tChan chan *Task
	done  chan struct{}
}

func (e *Engine) Run() {
	for {
		select {
		case <-e.done:
			return
		case t, ok := <-e.tChan:
			if ok {
				wg.Add(1)
				go func() {
					defer wg.Done()
					//TODO return error
					//finish one task
					if err := t.process(); err == nil {
						//try get next task
						if nextURL, i := t.NextURL(t.Page, t.Selection); nextURL != "" && i > t.Page {
							parse, err := url.Parse(nextURL)
							if err != nil {
								log.Println("parse next url failed", err.Error())
							} else {
								//wrap to new task
								t.Request.URL = parse
								e.tChan <- &Task{
									Request:      t.Request,
									MainSelector: t.MainSelector,
									functions:    t.functions,
									NextURL:      t.NextURL,
									Pipelines:    t.Pipelines,
									Page:         t.Page,
								}
							}
						}
					} else {
						log.Println(err.Error())
					}
				}()
			}
		}
	}
}

//waiting task done and stop engine
func (e *Engine) Stop() {
	//
	wg.Wait()
	close(e.tChan)
	//
	e.done <- struct{}{}
	close(e.done)
}

//shutdown engine immediately
func (e *Engine) Shutdown() {

}

func (e *Engine) AddTask(task *Task) error {
	if len(task.functions) == 0 {
		return errors.New("add task failed,because task should have at least one function")
	} else {
		if len(task.Pipelines) == 0 {
			task.Pipelines = append(task.Pipelines, DefaultPipeline)
		}
		e.tChan <- task
		return nil
	}
}

func NewEngine(size int) *Engine {
	return &Engine{
		tChan: make(chan *Task, size),
		done:  make(chan struct{}, 1),
	}
}
