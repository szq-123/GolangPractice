package golang

import (
	"GolangPractice/utils/logger"
	"errors"
	"fmt"
	"github.com/panjf2000/ants/v2"
	"golang.org/x/sync/errgroup"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// concurrency-safe Map
// theories, refer to：https://blog.csdn.net/u011957758/article/details/96633984
func TestSyncMap(t *testing.T) {
	key, subKey := "", 50
	m := &sync.Map{}
	setMap(m, key, subKey)
}

func setMap(m *sync.Map, key string, subKey int) {
	if subM, ok := m.Load(key); ok {
		e, ok1 := subM.(*sync.Map).Load(subKey)
		if ok1 {
			subM.(*sync.Map).Store(subKey, e.(int)+1)
		} else {
			subM.(*sync.Map).Store(subKey, 1)
		}
	} else {
		newSubM := &sync.Map{}
		newSubM.Store(subKey, 1)
		m.Store(key, newSubM)
	}
}

// synchronization
// mutually exclusive
func TestMutex(t *testing.T) {
	var m sync.Mutex

	m.Lock()
	go func() {
		defer m.Unlock()
		println("hello world")
	}()

	m.Lock()
}

func TestChannelSyncNoBuf(t *testing.T) {
	done := make(chan int)
	go func() {
		time.Sleep(time.Second)
		println("hello world")
		done <- 1 // < done. using `<- done` as an alternative also works, but it can be better to let the recipient block main threading.
	}()

	<-done // done <- 1
}

func TestChannelSyncBuf(t *testing.T) {
	done := make(chan int, 10)

	for i := 0; i < 10; i++ {
		go func(i int) {
			fmt.Printf("goroutine: %d\n", i)
			done <- 1
		}(i)
	}

	for i := 0; i < 10; i++ {
		<-done
	}
}

// sync.Pool an object pool.
// It is safe for use by multiple goroutines simultaneously.
func TestSyncPool(t *testing.T) {
	var numCalcsCreated int32

	objectPool := &sync.Pool{
		New: func() any {
			atomic.AddInt32(&numCalcsCreated, 1)
			buffer := make([]byte, 1024)
			return &buffer
		},
	}

	const numWorkers int = 1e6

	var wg sync.WaitGroup
	wg.Add(numWorkers)

	for i := 0; i < numWorkers; i++ {
		go func() {
			defer wg.Done()
			// acquire an instance of buffer
			buffer := objectPool.Get()
			_ = buffer.(*[]byte)
			// return the instance of buffer
			defer objectPool.Put(buffer)
		}()
	}
	wg.Wait()

	logger.Infoln(numCalcsCreated)
}

func TestSyncCond(t *testing.T) {
	var sharedRsc = make(map[int]interface{})

	var wg sync.WaitGroup
	wg.Add(2)
	m := sync.Mutex{}
	c := sync.NewCond(&m)

	for i := 0; i < 2; i++ {
		go func(index int) {
			// this go routine wait for changes to the sharedRsc
			c.L.Lock()
			for len(sharedRsc) == 0 {
				c.Wait()
			}
			fmt.Println(sharedRsc[index])
			c.L.Unlock()
			wg.Done()
		}(i)
	}

	// this one writes changes to sharedRsc
	c.L.Lock()
	sharedRsc[1] = "foo"
	sharedRsc[2] = "bar"
	c.Broadcast()
	c.L.Unlock()
	wg.Wait()
}

// channel is actually a queue. It's easy(convenient) to implement `Observer`(also known(referred to) as producer-consumer) model using it.
func TestProduceConsume(t *testing.T) {
	ch := make(chan int, 2)
	go Producer(5, ch)
	go Producer(2, ch)
	go Consumer(ch)

	time.Sleep(4 * time.Second)
}

func Producer(factor int, out chan<- int) {
	for i := 0; i < 10; i++ {
		out <- i * factor
	}
}

func Consumer(in <-chan int) {
	for v := range in {
		fmt.Printf("consume: %d\n", v)
	}
}

// limit(restrain,restrict,control) parallel count, by channel with buffer.
// TryLock，can implement it by channel with a buffer size of 1. it's like semaphore.

// goroutine
func TestWaitGroup(t *testing.T) {
	const goroutineNum = 2

	wg := &sync.WaitGroup{}
	wg.Add(2)

	for i := 0; i < goroutineNum; i++ {
		// reserved word `go`, it creates a coroutine.
		// the created coroutine will be added to a wait queue. so goroutine is not totally/completely real-time, but has the characteristic of delayed execution(delayed-execution featured), like `defer`.
		// so it will copy needed parameters while being established.
		go func(index int) {
			defer wg.Done()
			<-time.After(time.Second)
			println("done", index)
		}(i)
	}

	// `wg.Wait()` can be called in multiple places.
	// it's a for-loop which checks a variable if it hits particular conditions.
	go func() {
		wg.Wait()
		println("goroutine exits")
	}()

	wg.Wait()
	println("main exits\n")

	// os.Exit() exit the program directly(straightly, instantly, immediately). no `defer` shall be run.

	// Using Factory Pattern to bind goroutine to a channel.
	// some goroutines are permanently waiting for receiving messages from or sending messages to a channel due to certain bug in our program. it causes memory leaks.
}

// https://golang.org/x/sync/errgroup
// Slightly different from `go func()...`, it handles errors.
// ErrGroup does NOT offer the functionality of recovering from panic.
func TestErrGroup(t *testing.T) {
	g := new(errgroup.Group)

	g.Go(func() error {
		return nil
	})
	g.Go(func() error {
		return errors.New("very good")
	})

	err := g.Wait()
	if err != nil {
		logger.Errorln(err)
	}
}

// https://github.com/panjf2000/ants/v2 Goroutine Pool.
func TestGoroutinePool(t *testing.T) {
	const TaskNum = 1e3

	pool, err := ants.NewPool(TaskNum / 10)
	if err != nil {
		logger.Errorln(err)
		return
	}
	defer pool.Release()

	waitGroup := new(sync.WaitGroup)

	for i := 0; i < TaskNum; i++ {
		_ = pool.Submit(func() {
			waitGroup.Add(1)
			defer waitGroup.Done()
			fmt.Println("Good task")
		})
	}

	waitGroup.Wait()
}
