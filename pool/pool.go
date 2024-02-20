package pool

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

// TODO 优化或重构协程池，目前无法满足百万级并发

type signal struct{}

type Pool struct {
	// 最大运行worker数量
	cap int64

	// 正在运行worker数量
	running int64

	// 空闲worker的过期时间
	expire time.Duration

	// pool销毁信号
	release chan signal

	// 资源互斥锁
	lock sync.Mutex

	// 限制单次操作
	once sync.Once

	// 阻塞唤醒
	cond *sync.Cond

	// 空闲worker队列
	workers []*worker

	// worker缓存
	workerCache sync.Pool
}

func NewPool(cap int64, expire int64) (*Pool, error) {
	if cap < 0 {
		return nil, errors.New("cap can not be less than 0")
	}
	if expire < 0 {
		return nil, errors.New("expire can not be less than 0")
	}

	pool := &Pool{}
	pool.cap = cap
	pool.expire = time.Duration(expire) * time.Second
	pool.release = make(chan signal, 1)
	pool.cond = sync.NewCond(&pool.lock)
	pool.workerCache.New = func() any {
		return &worker{
			pool: pool,
			task: make(chan func(), 1),
		}
	}

	go pool.expireWorker()
	return pool, nil
}

func (p *Pool) Submit(task func()) error {
	if len(p.release) > 0 {
		return errors.New("pool has been destroyed")
	}

	w := p.getWorker()
	w.task <- task

	return nil
}

func (p *Pool) Release() {
	p.once.Do(func() {
		p.lock.Lock()
		workers := p.workers
		for i, w := range workers {
			w.pool = nil
			w.task = nil
			workers[i] = nil
		}
		p.workers = nil
		p.lock.Unlock()

		p.release <- signal{}
	})
}

func (p *Pool) RunningWorkers() int {
	return int(atomic.LoadInt64(&p.running))
}

func (p *Pool) expireWorker() {
	ticker := time.NewTicker(p.expire)
	for range ticker.C {
		//fmt.Printf("%v\n", p.workers)

		p.lock.Lock()
		var index = -1
		freeWorkers := p.workers
		if len(freeWorkers) > 0 {
			for i, w := range freeWorkers {
				if time.Now().Sub(w.last) <= p.expire {
					break
				} else {
					w.task <- nil
					index = i
					break
				}
			}

			freeWorkers = freeWorkers[index+1:]
			p.workers = freeWorkers
		}
		p.lock.Unlock()
	}
}

func (p *Pool) getWorker() *worker {
	var n int
	var freeWorkers []*worker

	freeWorkers = p.workers
	n = len(freeWorkers)

	if n > 0 {
		p.lock.Lock()
		w := freeWorkers[0]
		freeWorkers = freeWorkers[1:]
		p.workers = freeWorkers
		p.lock.Unlock()

		return w
	}

	if p.running < p.cap {
		var w *worker
		tw := p.workerCache.Get()
		if tw != nil {
			w = tw.(*worker)
		} else {
			w = &worker{
				pool: p,
				task: make(chan func(), 1),
			}
		}
		w.run()

		return w
	} else {
		for {
			p.lock.Lock()
			p.cond.Wait()
			freeWorkers = p.workers
			n = len(freeWorkers)

			if n <= 0 {
				p.lock.Unlock()
				continue
			}

			w := freeWorkers[0]
			freeWorkers = freeWorkers[1:]
			p.workers = freeWorkers
			p.lock.Unlock()

			return w
		}
	}
}

func (p *Pool) putWorker(w *worker) {
	w.last = time.Now()
	p.lock.Lock()
	p.workers = append(p.workers, w)
	p.cond.Signal()
	p.lock.Unlock()
}

func (p *Pool) incRunning() {
	atomic.AddInt64(&p.running, 1)
}

func (p *Pool) decRunning() {
	atomic.AddInt64(&p.running, -1)
}
