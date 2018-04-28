package myworker

import (
	"sync"
)

type Job interface {
	Do() error
}

type Worker struct {
	job  chan Job
	quit chan bool
	wg   sync.WaitGroup
}

// 构造工作者
func NewWorker(maxJobs int) *Worker {
	return &Worker{
		job:  make(chan Job, maxJobs),
		quit: make(chan bool),
	}
}

// 启动任务
func (w *Worker) Start() {
	w.wg.Add(1)

	go func() {
		defer w.wg.Done()

		for {
			// 接收任务
			// 此时工作中已经从工作者池中取出
			select {
			case job := <-w.job:
				// 处理任务
				job.Do()

			case <-w.quit:
				return
			}
		}
	}()
}

// 关闭任务
func (w *Worker) Stop() {
	w.quit <- true
	w.wg.Wait()
}

// 提交任务
func (w *Worker) AddJob(job Job) {
	w.job <- job
}
