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
	go func() {
		for {
			select {
			case <-e.done:
				return
			case t := <-e.tChan:
				wg.Add(1)
				go func() {
					defer wg.Done()
					//TODO return error
					//finish one task
					if t.process() {
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
					}
				}()
			}
		}
	}()
}

//waiting task done and stop engine
func (e *Engine) Stop() {
	//
	close(e.tChan)
	//
	wg.Wait()
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
