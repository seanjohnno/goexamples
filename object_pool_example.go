package main

import (
	"github.com/seanjohnno/objpool"
	"bytes"
	"fmt"
)

func main() {
	completeChan := make(chan bool, 100)

	pool := objpool.NewTimedExiryPool(3000)
	for i := 0; i < 100; i++ {
		go ReuseFunc(pool, completeChan)
	}

	completeCount := 0
	for {
		<- completeChan
		if completeCount++; completeCount == 100 {
			return
		}
	}
}

func ReuseFunc(pool objpool.ObjectPool, completeChan chan<- bool) {
	if item, present := pool.Retrieve(); present {
		fmt.Println("Woohoo! found an existing buffer...")
		// ...here we'd do something with our buffer...
		pool.Add(item)
	} else {
		fmt.Println("Have to create new buffer...")
		item := bytes.NewBuffer(make([]byte, 50))			// Nothing found in pool so create new
		// ...here we'd do something with our buffer...
		pool.Add(item)										// we're done with it, lets add back into pool so something else can use
	}
	completeChan <- true
}