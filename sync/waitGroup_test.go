package sync

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func Test_waitGroup(t *testing.T) {
	var wg sync.WaitGroup

	// 设置等待的 goroutine 数量为 2
	wg.Add(2)

	go func() {
		fmt.Println("Goroutine 1 start")
		defer wg.Done() // 当 goroutine 完成时调用 Done() 函数
		// 执行一些任务...
		time.Sleep(1 * time.Second)
		fmt.Println("Goroutine 1 completed")
	}()

	go func() {
		fmt.Println("Goroutine 2 start")
		defer wg.Done() // 当 goroutine 完成时调用 Done() 函数
		// 执行一些任务...
		time.Sleep(2 * time.Second)
		fmt.Println("Goroutine 2 completed")
	}()

	// 等待所有的 goroutine 完成
	wg.Wait()
	fmt.Println("2 goroutines completed")
	wg.Add(1)
	go func() {
		fmt.Println("Goroutine 3 start")
		defer wg.Done() // 当 goroutine 完成时调用 Done() 函数
		// 执行一些任务...
		time.Sleep(3 * time.Second)
		fmt.Println("Goroutine 3 completed")
	}()
	wg.Wait()
	fmt.Println("All goroutines completed")
	//var mutex sync.Mutex
}
