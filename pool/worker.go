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
	w.pool.incRunning()
	go w.work()
}

func (w *worker) work() {
	defer func() {
		w.pool.decRunning()
		w.pool.workerCache.Put(w)
		err := recover()
		if err != nil {

		}
	}()

	for fn := range w.task {
		if fn == nil {
			return
		}

		fn()
		w.pool.putWorker(w)
	}
}
