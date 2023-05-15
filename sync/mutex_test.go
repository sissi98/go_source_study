package sync

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
)

func Test_Cas(t *testing.T) {
	var counter int64
	var wg sync.WaitGroup

	wg.Add(2)

	// 使用 CAS 增加计数器的值
	go func() {
		defer wg.Done()

		for i := 0; i < 1000; i++ {
			atomic.AddInt64(&counter, 1)
		}
	}()

	// 使用 CAS 减少计数器的值
	go func() {
		defer wg.Done()

		for i := 0; i < 1000; i++ {
			atomic.AddInt64(&counter, -1)
		}
	}()

	wg.Wait()

	// 输出最终计数器的值
	fmt.Println("Counter:", atomic.LoadInt64(&counter))
}
