package main

import (
	"fmt"
	"sync"
	"time"
)

type Bucket struct {
	sync.RWMutex
	// 请求总数
	ReqTotal int

	Timestamp time.Time
}

func NewBucket() *Bucket {
	return &Bucket{
		Timestamp: time.Now(),
	}
}

func (b *Bucket) Record() {
	b.Lock()
	defer b.Unlock()
	b.ReqTotal++
}

type SlidingWindow struct {
	sync.RWMutex
	limit   bool
	size    int
	buckets []*Bucket
	// 触发限流的请求总数阈值
	reqThreshold int
	//上次限流时间
	lastLimitTime time.Time
	limitTimeGap  time.Duration
}

func NewSlidingWindow(size, reqThreshold int, limitTimeGap time.Duration) *SlidingWindow {
	return &SlidingWindow{
		size:         size,
		buckets:      make([]*Bucket, 0, size),
		reqThreshold: reqThreshold,
		limitTimeGap: limitTimeGap,
	}
}

func (sw *SlidingWindow) AddBucket() {
	sw.Lock()
	defer sw.Unlock()
	sw.buckets = append(sw.buckets, NewBucket())
	if !(len(sw.buckets) < sw.size+1) {
		sw.buckets = sw.buckets[1:]
	}
}
func (sw *SlidingWindow) GetBucket() *Bucket {
	if len(sw.buckets) == 0 {
		sw.AddBucket()
	}
	return sw.buckets[len(sw.buckets)-1]
}

func (sw *SlidingWindow) RecordReqResult() {
	sw.GetBucket().Record()
}

// 根据当前滑动窗口判断是否需要触发限流
func (sw *SlidingWindow) ReqLimit() bool {
	sw.RLock()
	defer sw.RUnlock()
	reqTotal := 0
	for _, v := range sw.buckets {
		reqTotal += v.ReqTotal
	}
	if reqTotal > sw.reqThreshold {
		return true
	}
	return false
}

func (sw *SlidingWindow) Monitor() {
	go func() {
		for {
			if sw.limit {
				if sw.OverLimitTimeGap() {
					sw.Lock()
					sw.limit = false
					sw.Unlock()
				}
				continue
			}
			if sw.ReqLimit() {
				sw.Lock()
				sw.limit = true
				sw.lastLimitTime = time.Now()
				sw.Unlock()
			}
		}
	}()
}

func (sw *SlidingWindow) OverLimitTimeGap() bool {
	return time.Since(sw.lastLimitTime) > sw.limitTimeGap
}

func (sw *SlidingWindow) IsLimit() bool {
	return sw.limit
}

//启动
func (sw *SlidingWindow) Start() {
	go func() {
		for {
			//每隔一段时间添加一个bucket
			sw.AddBucket()
			time.Sleep(time.Millisecond * 100)
		}
	}()
}

func main() {
	//可以添加到gin 框架的中间件

	//简单模拟
	limitTimeGap := time.Millisecond * 50
	sw := NewSlidingWindow(5, 15, limitTimeGap)
	sw.Start()
	sw.Monitor()

	for i := 0; i < 50; i++ {
		//添加一个请求
		sw.RecordReqResult()
		time.Sleep(time.Millisecond * 30)
		if sw.IsLimit() {
			fmt.Println("这个请求被限流啦。。。", i)
		} else {
			fmt.Println("这个请求通过。。。", i)
		}
	}
	time.Sleep(3 * time.Second)
}
