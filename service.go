package myworker

import (
	"sync"
)

type Service struct {
	workers *WorkerPool
	jobs    chan Job
	maxJobs int
	wg      sync.WaitGroup
}

func NewService(maxWorkers, maxJobs int) *Service {
	return &Service{
		workers: NewWorkerPool(maxWorkers, maxJobs),
		jobs:    make(chan Job, maxJobs),
	}
}

func (p *Service) Start() {
	p.wg.Add(1)
	p.workers.Start()

	go func() {
		defer p.wg.Done()

		for job := range p.jobs {
			go func(job Job) {
				// 从工作者池取一个工作者
				worker := p.workers.Get()

				// 完成任务后返回给工作者池
				defer p.workers.Put(worker)

				// 提交任务处理(异步)
				worker.AddJob(job)
			}(job)
		}
	}()
}

func (p *Service) Stop() {
	p.workers.Stop()
	close(p.jobs)
	p.wg.Wait()
}

// 提交任务
// 任务管道带较大的缓存, 延缓阻塞的时间
func (p *Service) AddJob(job Job) {
	p.jobs <- job
}
