package spider

import (
	"errors"
	"log"
	"net/url"
	"sync"
	"sync/atomic"
)

var wg sync.WaitGroup

type Engine struct {
	tChan  chan *Task
	done   chan struct{}
	status int32
}

func (e *Engine) Run() {
	for {
		select {
		case <-e.done:
			return
		case t, ok := <-e.tChan:
			if ok {
				log.Println("=============>", len(t.functions))
				go func() {
					wg.Add(1)
					defer wg.Done()
					//finish one task
					if err := t.process(); err == nil {
						//try get next task
						if nextURL, i := t.NextURL(t.Page, t.Selection); nextURL != "" {
							parse, err := url.Parse(nextURL)
							if err != nil {
								log.Println("parse next url failed", err.Error())
							} else {
								//wrap to new task
								t.Request.URL = parse
								newTask, _ := NewTask(t.Request, t.MainSelector)
								newTask.Page = i + 1
								newTask.Pipelines = t.Pipelines
								newTask.functions = t.functions
								newTask.NextURL = t.NextURL
								err = e.AddTask(newTask)
								if err != nil {
									log.Println(err.Error())
								}
							}
						}
					} else {
						log.Println(err.Error())
					}
				}()
			} else {
				log.Println("engine closed")
			}
		}
	}
}

//waiting task done and stop engine
func (e *Engine) Stop() {
	//
	atomic.StoreInt32(&e.status, -1)
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
	if atomic.LoadInt32(&e.status) == -1 {
		return errors.New("engine is closed")
	}
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
		tChan:  make(chan *Task, size),
		done:   make(chan struct{}, 1),
		status: 0,
	}
}
