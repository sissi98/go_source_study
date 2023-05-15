### noCopy
通过使用noCopy字段，可以在编译时防止结构体被复制，避免了潜在的错误和不一致性。这在一些需要显式控制对象复制的场景中非常有用，例如单例模式或共享资源的引用计数。

需要注意的是，noCopy只提供了静态检查，它并不能完全阻止通过指针或反射等方式对对象进行复制。因此，在使用noCopy时，仍然需要遵循最佳实践，并确保类型的复制行为符合预期

### waitGroup
在sync.WaitGroup结构体中，每个元素的含义如下：

- noCopy：noCopy是一个私有的结构体类型，用于实现防止WaitGroup被拷贝的机制。它的存在可以防止不小心拷贝WaitGroup实例。

- state：state是一个atomic.Uint64类型，用于保存WaitGroup的状态信息。它是一个无符号整数，占用64位，其中高32位用于计数器（counter），低32位用于等待的goroutine计数（waiter count）。

- sema：sema是一个uint32类型，表示一个信号量（semaphore）。它用于控制等待的goroutine的同步和通信。当等待的goroutine数量大于0时，通过对sema的操作来阻塞或唤醒等待的goroutine。 

### mutex 
在上面的代码中，有一些可以借鉴的地方：

1. 使用原子操作：代码中使用`atomic`包提供的原子操作函数来进行对互斥锁状态的读取和更新。这样可以确保对状态的操作是原子的，避免了竞态条件和数据竞争。

2. 基于CAS（Compare-And-Swap）的加锁：在`Lock`方法中，通过使用`atomic.CompareAndSwapInt32`函数来尝试获取锁。这是一种基于CAS的方式，可以避免显式的互斥锁等待。只有在成功获取锁的情况下，才会进入临界区。这种方式可以提供较低的开销，特别是在高并发环境中。

3. 自旋锁和休眠：在`lockSlow`方法中，通过自旋锁（spinlock）的方式尝试获取锁。如果获取锁失败，会进入休眠状态等待唤醒。自旋锁可以减少线程切换的开销，但在自旋一定次数后仍未获取到锁时，会进入休眠状态，避免空转浪费CPU资源。

4. 饥饿模式（starvation mode）：为了防止某些等待时间过长的goroutine长时间无法获取锁（饥饿），代码中引入了饥饿模式。当等待时间超过阈值时，锁会进入饥饿模式，直接将锁的所有权交给队列中的下一个等待者，而不是再进行竞争。这样可以避免长时间等待者无法获得锁的问题，提高公平性和避免尾延迟（tail latency）。

5. 锁的状态管理：代码中使用了位操作来管理互斥锁的状态。通过位运算，可以将不同的状态信息存储在一个整数中，提高了内存效率。此外，还利用了状态的不同位来表示锁是否已锁定、是否处于饥饿模式、是否有等待者等信息。

这些是代码中一些值得借鉴的地方，可以在自己的并发代码实现中考虑使用类似的技术和策略。不过，需要根据具体的应用场景和需求来评估是否适用。

#### 互斥锁的公平性。
//
// 互斥锁可以处于两种操作模式：正常模式和饥饿模式。
// 在正常模式下，等待者按照先进先出（FIFO）的顺序排队，但被唤醒的等待者并不拥有互斥锁，并与新到达的goroutine竞争锁的所有权。新到达的goroutine有一个优势——它们已经在CPU上运行，并且可能有很多个，因此被唤醒的等待者很有可能失败。在这种情况下，它会排在等待队列的前面。如果一个等待者连续超过1毫秒无法获取互斥锁，它会将互斥锁切换到饥饿模式。
//
// 在饥饿模式下，互斥锁的所有权直接从解锁的goroutine传递给等待队列前面的等待者。新到达的goroutine不会尝试获取互斥锁，即使它看起来是解锁的，也不会尝试自旋。相反，它们会排队在等待队列的尾部。
//
// 如果一个等待者获得互斥锁的所有权，并且发现以下情况之一：(1) 它是队列中的最后一个等待者，或者(2) 它等待的时间少于1毫秒，它会将互斥锁切换回正常操作模式。
//
// 正常模式具有更好的性能，因为一个goroutine可以连续多次获取互斥锁，即使有阻塞的等待者。饥饿模式是为了防止尾延迟（tail latency）的病态情况。

#### race
虽然竞争检测工具在生产环境中可能会带来一些性能开销，但它们对于开发和调试阶段非常有价值。
它们帮助开发人员在开发过程中发现潜在的数据竞争问题，并在调试时提供有用的信息。
因此，尽管竞争检测工具在生产环境中可能不是必需的，但在开发和调试阶段，它们通常被视为非常有用的工具

```go

// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package sync provides basic synchronization primitives such as mutual
// exclusion locks. Other than the Once and WaitGroup types, most are intended
// for use by low-level library routines. Higher-level synchronization is
// better done via channels and communication.
//
// Values containing the types defined in this package should not be copied.
package sync

import (
	"internal/race"
	"sync/atomic"
	"unsafe"
)

// Provided by runtime via linkname.
func throw(string)
func fatal(string)

// A Mutex is a mutual exclusion lock.
// The zero value for a Mutex is an unlocked mutex.
//
// A Mutex must not be copied after first use.
//
// In the terminology of the Go memory model,
// the n'th call to Unlock “synchronizes before” the m'th call to Lock
// for any n < m.
// A successful call to TryLock is equivalent to a call to Lock.
// A failed call to TryLock does not establish any “synchronizes before”
// relation at all.-
type Mutex struct {
	state int32
	sema  uint32
}

// A Locker represents an object that can be locked and unlocked.
type Locker interface {
	Lock()
	Unlock()
}

const (
	mutexLocked = 1 << iota // mutex is locked
	mutexWoken
	mutexStarving
	mutexWaiterShift = iota

	// Mutex fairness.
	//
	// Mutex can be in 2 modes of operations: normal and starvation.
	// In normal mode waiters are queued in FIFO order, but a woken up waiter
	// does not own the mutex and competes with new arriving goroutines over
	// the ownership. New arriving goroutines have an advantage -- they are
	// already running on CPU and there can be lots of them, so a woken up
	// waiter has good chances of losing. In such case it is queued at front
	// of the wait queue. If a waiter fails to acquire the mutex for more than 1ms,
	// it switches mutex to the starvation mode.
	//
	// In starvation mode ownership of the mutex is directly handed off from
	// the unlocking goroutine to the waiter at the front of the queue.
	// New arriving goroutines don't try to acquire the mutex even if it appears
	// to be unlocked, and don't try to spin. Instead they queue themselves at
	// the tail of the wait queue.
	//
	// If a waiter receives ownership of the mutex and sees that either
	// (1) it is the last waiter in the queue, or (2) it waited for less than 1 ms,
	// it switches mutex back to normal operation mode.
	//
	// Normal mode has considerably better performance as a goroutine can acquire
	// a mutex several times in a row even if there are blocked waiters.
	// Starvation mode is important to prevent pathological cases of tail latency.
	starvationThresholdNs = 1e6
)

// Lock locks m.
// If the lock is already in use, the calling goroutine
// blocks until the mutex is available.
func (m *Mutex) Lock() {
	// Fast path: grab unlocked mutex.
	if atomic.CompareAndSwapInt32(&m.state, 0, mutexLocked) {
		if race.Enabled {
			race.Acquire(unsafe.Pointer(m))
		}
		return
	}
	// Slow path (outlined so that the fast path can be inlined)
	m.lockSlow()
}

// TryLock tries to lock m and reports whether it succeeded.
//
// Note that while correct uses of TryLock do exist, they are rare,
// and use of TryLock is often a sign of a deeper problem
// in a particular use of mutexes.
func (m *Mutex) TryLock() bool {
	old := m.state
	if old&(mutexLocked|mutexStarving) != 0 {
		return false
	}

	// There may be a goroutine waiting for the mutex, but we are
	// running now and can try to grab the mutex before that
	// goroutine wakes up.
	if !atomic.CompareAndSwapInt32(&m.state, old, old|mutexLocked) {
		return false
	}

	if race.Enabled {
		race.Acquire(unsafe.Pointer(m))
	}
	return true
}

func (m *Mutex) lockSlow() {
	var waitStartTime int64
	starving := false
	awoke := false
	iter := 0
	old := m.state
	for {
		// Don't spin in starvation mode, ownership is handed off to waiters
		// so we won't be able to acquire the mutex anyway.
		if old&(mutexLocked|mutexStarving) == mutexLocked && runtime_canSpin(iter) {
			// Active spinning makes sense.
			// Try to set mutexWoken flag to inform Unlock
			// to not wake other blocked goroutines.
			if !awoke && old&mutexWoken == 0 && old>>mutexWaiterShift != 0 &&
				atomic.CompareAndSwapInt32(&m.state, old, old|mutexWoken) {
				awoke = true
			}
			runtime_doSpin()
			iter++
			old = m.state
			continue
		}
		new := old
		// Don't try to acquire starving mutex, new arriving goroutines must queue.
		if old&mutexStarving == 0 {
			new |= mutexLocked
		}
		if old&(mutexLocked|mutexStarving) != 0 {
			new += 1 << mutexWaiterShift
		}
		// The current goroutine switches mutex to starvation mode.
		// But if the mutex is currently unlocked, don't do the switch.
		// Unlock expects that starving mutex has waiters, which will not
		// be true in this case.
		if starving && old&mutexLocked != 0 {
			new |= mutexStarving
		}
		if awoke {
			// The goroutine has been woken from sleep,
			// so we need to reset the flag in either case.
			if new&mutexWoken == 0 {
				throw("sync: inconsistent mutex state")
			}
			new &^= mutexWoken
		}
		if atomic.CompareAndSwapInt32(&m.state, old, new) {
			if old&(mutexLocked|mutexStarving) == 0 {
				break // locked the mutex with CAS
			}
			// If we were already waiting before, queue at the front of the queue.
			queueLifo := waitStartTime != 0
			if waitStartTime == 0 {
				waitStartTime = runtime_nanotime()
			}
			runtime_SemacquireMutex(&m.sema, queueLifo, 1)
			starving = starving || runtime_nanotime()-waitStartTime > starvationThresholdNs
			old = m.state
			if old&mutexStarving != 0 {
				// If this goroutine was woken and mutex is in starvation mode,
				// ownership was handed off to us but mutex is in somewhat
				// inconsistent state: mutexLocked is not set and we are still
				// accounted as waiter. Fix that.
				if old&(mutexLocked|mutexWoken) != 0 || old>>mutexWaiterShift == 0 {
					throw("sync: inconsistent mutex state")
				}
				delta := int32(mutexLocked - 1<<mutexWaiterShift)
				if !starving || old>>mutexWaiterShift == 1 {
					// Exit starvation mode.
					// Critical to do it here and consider wait time.
					// Starvation mode is so inefficient, that two goroutines
					// can go lock-step infinitely once they switch mutex
					// to starvation mode.
					delta -= mutexStarving
				}
				atomic.AddInt32(&m.state, delta)
				break
			}
			awoke = true
			iter = 0
		} else {
			old = m.state
		}
	}

	if race.Enabled {
		race.Acquire(unsafe.Pointer(m))
	}
}

// Unlock unlocks m.
// It is a run-time error if m is not locked on entry to Unlock.
//
// A locked Mutex is not associated with a particular goroutine.
// It is allowed for one goroutine to lock a Mutex and then
// arrange for another goroutine to unlock it.
func (m *Mutex) Unlock() {
	if race.Enabled {
		_ = m.state
		race.Release(unsafe.Pointer(m))
	}

	// Fast path: drop lock bit.
	new := atomic.AddInt32(&m.state, -mutexLocked)
	if new != 0 {
		// Outlined slow path to allow inlining the fast path.
		// To hide unlockSlow during tracing we skip one extra frame when tracing GoUnblock.
		m.unlockSlow(new)
	}
}

func (m *Mutex) unlockSlow(new int32) {
	if (new+mutexLocked)&mutexLocked == 0 {
		fatal("sync: unlock of unlocked mutex")
	}
	if new&mutexStarving == 0 {
		old := new
		for {
			// If there are no waiters or a goroutine has already
			// been woken or grabbed the lock, no need to wake anyone.
			// In starvation mode ownership is directly handed off from unlocking
			// goroutine to the next waiter. We are not part of this chain,
			// since we did not observe mutexStarving when we unlocked the mutex above.
			// So get off the way.
			if old>>mutexWaiterShift == 0 || old&(mutexLocked|mutexWoken|mutexStarving) != 0 {
				return
			}
			// Grab the right to wake someone.
			new = (old - 1<<mutexWaiterShift) | mutexWoken
			if atomic.CompareAndSwapInt32(&m.state, old, new) {
				runtime_Semrelease(&m.sema, false, 1)
				return
			}
			old = m.state
		}
	} else {
		// Starving mode: handoff mutex ownership to the next waiter, and yield
		// our time slice so that the next waiter can start to run immediately.
		// Note: mutexLocked is not set, the waiter will set it after wakeup.
		// But mutex is still considered locked if mutexStarving is set,
		// so new coming goroutines won't acquire it.
		runtime_Semrelease(&m.sema, true, 1)
	}
}


```