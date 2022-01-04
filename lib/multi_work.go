/**
* @Author:Tristan
* @Date: 2022/1/4 11:22 上午
 */

package lib

import (
	"fmt"
	"time"
	"context"
)

type Worker struct {
	Num int
	ch1 chan interface{}
}

func (w *Worker) Run(ctx context.Context) {
	for {
		select {
		case ch1 := <-w.ch1:
			fmt.Println(ch1)
		//  doing
		case <-ctx.Done():
			return
		}
	}
}
func (w *Worker) Send(num int) {
	w.ch1 <- num
}
func NewWork(workNum int, ctx context.Context) *Worker {
	w := &Worker{Num: workNum}
	for i := 0; i < w.Num; i++ {
		go func() {
			w.Run(ctx)
		}()
	}
	return w
}

func testWork() {
	ctx, cancel := context.WithCancel(context.TODO())
	w := NewWork(10, ctx)

	timer := time.NewTicker(25 * time.Second)
	for {
		select {
		case <-timer.C:
			cancel()
			return
		default:
			w.Send(int(time.Now().UnixMilli()))
			time.Sleep(1000 * time.Millisecond)
		}
	}
}

