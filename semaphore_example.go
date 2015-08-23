package main

import (
	"github.com/seanjohnno/semaphore"
	"sync"
	"container/list"
	"fmt"
	"math/rand"
	"time"
)

const (
	// ConsumerCount is the amount of consumer goroutines we want to run
	ConsumerCount = 50

	// PublisherCount is the amount of publisher goroutines we want to run
	PublisherCount = 25

	// MinSleepDurationMilli is used to generate a random sleep interval for publishers (so we can see things out or order)
	MinSleepDurationMilli = 0

	// MaxSleepDurationMilli is used to generate a random sleep interval for publishers (so we can see things out or order)
	MaxSleepDurationMilli = 30
)

var (
	// SyncMutex is used to synchronise/lock the publish/Consume queue
	SyncMutex = sync.Mutex{}

	// CountingSem is our semaphore used to wait for published items
	CountingSem = semaphore.New()

	// PubsubQueue is our publisher/consumer queue where items are added and consumed from
	PubsubQueue = list.New()

	// ConsumeCount is the total amount of items that have been consumed (so we know when to end the example)
	ConsumeCount = 0

	// WaitChan is used to pause the main routine until we're done
	WaitChan = make(chan bool)
)
	
// main - entry point
func main() {	
	// Run our consumer routines - they'll wait until something has been published
	for i := 0; i < ConsumerCount; i++ {
		go Consume(WaitChan)
	}

	// Run our publisher routines
	for i := 0; i < PublisherCount; i++ {
		go Publish(i)
	}

	// Wait on the channel, Consume sends to the channel when we've consumed everything
	<-WaitChan
}

// Publish waits a random interval, adds its value to the queue and then signals the semaphore
func Publish(val int) {
	RandSleep(val)							// Sleep for random interval

	SyncMutex.Lock()						// Sync method as we're changing state on the queue
	defer SyncMutex.Unlock()				
	
	fmt.Printf("Publishing: %d\n", val)
	PubsubQueue.PushFront(val)				// Add value to the queue

	// Signal semaphore - this will incrememnt the semaphores internal Count. If a single or multiple
	// goroutines are currently waiting on it, it'll unblock one of them
	CountingSem.Signal()					
}

// Consume waits on the semaphore until there is something for it to grab from the queue
func Consume(waitChan chan<- bool) {
	// Forever loop
	for {

		// Wait on semaphore - If calls to signal are outstripping calls to Wait then we'll return
		// immediately here without actually blocking
		CountingSem.Wait()


		SyncMutex.Lock()									// [Start sync]
		
		// Pop the oldest element in the queue (FIFO)
		elem := PubsubQueue.Back()							
		PubsubQueue.Remove(elem)							
		
		fmt.Printf("Consumed: %d\n", elem.Value.(int))

		// Increment count and break out of loop if we've consumed everything
		ConsumeCount++
		if ConsumeCount == PublisherCount {
			SyncMutex.Unlock()
			break
		}

		SyncMutex.Unlock()									// [End sync]
	}

	// Send to channel so we can exit
	WaitChan <- true
}

// RandSleep pauses the goroutine for a random interval between MaxSleepDurationMilli - MinSleepDurationMilli
func RandSleep(seed int) {
	rand.Seed(time.Now().Unix() + int64(seed))
    time.Sleep(time.Duration(rand.Intn(MaxSleepDurationMilli - MinSleepDurationMilli) + MinSleepDurationMilli) * time.Millisecond)
}