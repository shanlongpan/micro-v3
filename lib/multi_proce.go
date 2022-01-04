/**
* @Author:Tristan
* @Date: 2021/12/31 6:15 下午
 */

package lib

import (
	"log"
	"sync"
	"sync/atomic"
	"time"
)

var DefaultAsyncWithBufferHelper *AsyncCallHelper = NewAsyncCallHelper(256, 64, true)

type AsyncCallHelper struct {
	callChan   chan func() error
	callCount  int64
	withBuffer bool
	callBuffer []func() error
	mutex      sync.Mutex
}

func NewAsyncCallHelper(maxChanLen, maxprocs int, withBuffer bool) *AsyncCallHelper {
	m := &AsyncCallHelper{
		callChan:   make(chan func() error, maxChanLen),
		withBuffer: withBuffer,
	}

	for i := 0; i < maxprocs; i++ {
		go m.run()
	}

	go m.qps()

	return m
}

func (m *AsyncCallHelper) AsyncCall(call func() error) {
	//fun := "AsyncCallHelper.AsyncCall -->"

	atomic.AddInt64(&m.callCount, 1)

	if m.withBuffer == false {
		m.callChan <- call

	} else {

		select {
		case m.callChan <- call:

		default:
			m.mutex.Lock()
			m.callBuffer = append(m.callBuffer, call)
			m.mutex.Unlock()

		}
	}
}

func (m *AsyncCallHelper) run() {
	fun := "AsyncCallHelper.run-->"

	ticker := time.NewTicker(time.Second * 30)
	defer ticker.Stop()

	for {
		select {
		case call := <-m.callChan:
			err := call()
			if err != nil {
				log.Printf("%s err:%s", fun, err)
			}

		case <-ticker.C:
			if m.withBuffer {
				m.handleCallBuffer()
			}
		}
	}
}

func (m *AsyncCallHelper) handleCallBuffer() {
	fun := "AsyncCallHelper.handleCallBuffer-->"

	m.mutex.Lock()
	buffer := m.callBuffer
	m.callBuffer = nil
	m.mutex.Unlock()

	for _, call := range buffer {
		err := call()
		if err != nil {
			log.Printf("%s err:%s", fun, err)
		}
	}
}

func (m *AsyncCallHelper) qps() {
	fun := "AsyncCallHelper.qps -->"

	ticker := time.NewTicker(30 * time.Second)
	for {
		select {
		case <-ticker.C:
			count := atomic.SwapInt64(&m.callCount, 0)
			log.Printf("%s qps:%d", fun, count/30)
		}
	}
}

