package pool

import (
	"time"
)

type worker struct {
	// worker所属的pool
	pool *Pool

	// 最近一次执行任务的时间
	last time.Time

	// 绑定至worker的任务
	task chan func()
}

func (w *worker) run() {
	go w.work()
}

func (w *worker) work() {
	for fn := range w.task {
		fn()
		w.pool.putWorker(w)
		w.pool.decRunning()
	}
}
