package main

import (
	"container/list"
	"fmt"
	"sync"
	"time"
)

// workq to hold worker reqests, this is decoupling
// and serializing the execution, this prevents to many
// threads in case src traversal is too fast (which is probably
// normal)
var (
	workq      list.List
	workqmutex sync.Mutex
)

// CopySlice is the RPC routine that does the copying, it will
// be called on several servers
func (*PafimiServerT) CopySlice(arg []string, result *string) error {
	fmt.Println("pushing", arg)
	workqmutex.Lock()
	workq.PushBack(arg)
	workqmutex.Unlock()
	return nil
}

// CopyWorker is the go routine checking workq and processing
// the data
func CopyWorker() {
	for true {
		time.Sleep(100 * time.Millisecond)
		// we have only one consumer, so no need to lock the length check
		for workq.Len() > 0 {
			workqmutex.Lock()
			element := workq.Remove(workq.Front())
			// DO WORK HERE
			workqmutex.Unlock()
			fmt.Println("picked:", element)
		}
	}
}
