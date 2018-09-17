package main

import (
    "sync"
    "sync/atomic"
    "time"
)

type Manager struct {
    limit                int
    jobs                 chan *Job
    activeCount          int32
    completionCount        int64
    wg                   *sync.WaitGroup
    started              time.Time
}


// job result
type Future struct {
    result interface{}
    signal chan bool // if read channel returned, result is ready
}

// Represent a job that returns something for future retrieval. Could be nil
type JobFunc func() interface{}

// A function that is no args, no returns
type VoidFunc func()

// Internal concept of a Job and its produced result
type Job struct {
    jobf   JobFunc
    result *Future
}

// Create a new Manager.
// max_pending_jobs：最大等待job 数
// limit: 启动协程数

func NewManager(limit int, max_pending_jobs int) *Manager{
    return &Manager{
        limit: limit,
        jobs: make(chan *Job, max_pending_jobs),
        activeCount: 0,
        completionCount: 0,
        wg: new(sync.WaitGroup)}
}

// Starts
func (v *Manager) Start() {
    for i := 0; i < v.limit; i++ {
        v.wg.Add(1)
        go func() {
            defer v.wg.Done()
            for next := range v.jobs {
                atomic.AddInt32(&v.activeCount, 1)
                result := next.jobf()
                next.result.updateResult(result)
                atomic.AddInt32(&v.activeCount, -1)
                atomic.AddInt64(&v.completionCount, 1)
            }
        }()
    }
    v.started = time.Now()
}

// 返回执行任务数
func (v *Manager) ActiveCount() int {
    return int(v.activeCount)
}

// 所有任务开始时间
func (v *Manager) StartedTime() time.Time {
    return v.started
}

//　等待任务数
func (v *Manager) PendingCount() int {
    return len(v.jobs)
}

// How many jobs has been completed
func (v *Manager) CompletedCount() int64 {
    return v.completionCount
}

// Stop accepting new jobs. After this call is called, future calls to Submit will panic
// You can't shutdown more than once, sorry
func (v *Manager) Shutdown() {
    close(v.jobs) //now submission will panic
}

// Wait until all jobs are processed. after this, All previously returned future should be ready for retrieval
// Must call Shutdown() first or Wait() will block forever
func (v *Manager) Wait() {
    v.wg.Wait()
}

// Submit a job and return a Future value can be retrieved later sync or async
func (v *Manager) Submit(j JobFunc) *Future {
    if j == nil {
        panic("Can't submit nil function")
    }
    result := &Future{nil, make(chan bool, 1)}
    nj := &Job{j, result}
    v.jobs <- nj
    return result
}

func (v *Future) updateResult(result interface{}) {
    v.result = result
    v.signal <- true
    close(v.signal)
}

// Get the future value without wait. bool value is whether this retrieve did retrieve something, the interface{} value
// is the actual future result
func (v *Future) GetNoWait() (bool, interface{}) {
    return v.GetWaitTimeout(0 * time.Second)
}

// Synchronously retrieve the future's value. It will block until the value is available
func (v *Future) GetWait() interface{} {
    <-v.signal
    return v.result
}

// Retrieve the futures value, with a timeout. The bool value represent whether this retrieval did succeed
func (v *Future) GetWaitTimeout(t time.Duration) (bool, interface{}) {
    select {
    case <-v.signal:
        return true, v.result
    case <-time.After(t):
        return false, nil
    }
}

type FutureGroup struct {
    Futures []*Future
}

func (v FutureGroup) WaitAll() []interface{} {
    result := make([]interface{}, len(v.Futures))
    for idx, nf := range v.Futures {
        result[idx] = nf.GetWait()
    }
    return result
}


func (v FutureGroup) WaitAllTimeout(t time.Duration) []interface{} {
    result := make([]interface{}, len(v.Futures))
    for idx, nf := range v.Futures {
        _, result[idx] = nf.GetWaitTimeout(t)
    }
    return result
}

// 并行处理jobs ，并行数为：job的数量
func ParallelDo(jobs []func() interface{}) FutureGroup {
    return ParallelDoWithLimit(jobs, len(jobs))
}

// nTthreads 个并行　
func ParallelDoWithLimit(jobs []func() interface{}, nThreads int) FutureGroup {
    if nThreads > len(jobs) {
        nThreads = len(jobs)
    }
    tp :=    NewManager(nThreads, len(jobs))
    tp.Start()
    defer func() {
        tp.Shutdown()
        tp.Wait()
    }()
    result := make([]*Future, len(jobs))
    for idx, nj := range(jobs) {
        result[idx] = tp.Submit(nj)
    }
    return FutureGroup{result}
}
// Run a func and get its result as a futur immediately
// Note this is unmanaged, it is as good as your own go func(){}()
// Just that it is wrapped with a nice Future object, and you can
// retrieve it as many times as you want, and you can retrieve with timeout
func FutureOf(f JobFunc) *Future {
    if f == nil {
        panic("Can't create Future of nil function")
    }
    result := &Future{nil, make(chan bool, 1)}
    go func() {
        resultVal := f()
        result.updateResult(resultVal)
    }()
    return result
}