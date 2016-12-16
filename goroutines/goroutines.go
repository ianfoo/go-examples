// This is an example of using unbuffered channels to coordinate a
// single stream of work with an arbitrary number of workers, and
// allowing the workers to shut down cleanly (without using
// sync.WaitGroup).
package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

const (
	defaultWorkers  = 16
	defaultMessages = 100
	maxSleepMS      = 100
)

var (
	numWorkers  = defaultWorkers
	numMessages = defaultMessages
)

var (
	workCh = make(chan string)
	quitCh = make(chan string)
)

func main() {
	if n := IntFromEnv("NUM_MESSAGES"); n != -1 {
		numMessages = n
	}
	if n := IntFromEnv("NUM_WORKERS"); n != -1 {
		numWorkers = n
	}
	fmt.Printf("[manager] sending %d messages to %d workers\n", numMessages, numWorkers)

	// Launch workers.
	for i := 0; i < numWorkers; i++ {
		go func(id int) {
			Worker{id}.Run(workCh, quitCh)
		}(i + 1)
	}

	// Produce messages for the workers.
	for i := 0; i < numMessages; i++ {
		workCh <- fmt.Sprintf("work message #%d", i+1)
	}

	// Signal workers to stop.
	// All listeners to a channel are notified when the channel is closed.
	close(workCh)

	// Wait for workers to shut down.
	for i := numWorkers; i > 0; i-- {
		msg := <-quitCh
		fmt.Printf("[manager] %s\n", msg)
	}

	fmt.Println("goodbye")
}

// Worker prints messages it receives.
type Worker struct {
	// ID allows each worker to be uniquely identified in its output messages.
	ID int
}

// Run waits for work on the channel and
func (w Worker) Run(work <-chan string, quit chan string) {
	t0 := time.Now()
	for msg := range work {
		time.Sleep(time.Duration(rand.Intn(maxSleepMS)) * time.Millisecond)
		w.log(msg)
	}
	quit <- fmt.Sprintf("worker %d finished in %s", w.ID, time.Since(t0))
}

func (w Worker) log(format string, args ...interface{}) {
	args = append([]interface{}{w.ID}, args...)
	fmt.Printf("[worker %02d] "+format+"\n", args...)
}

// IntFromEnv tries to get an integer value from an environment value.
// If the variable is empty or set to a value that cannot be converted
// to an integer, -1 is returned. So this function will not work in
// cases where -1 is a valid setting.
func IntFromEnv(envvar string) int {
	if s := os.Getenv(envvar); s != "" {
		if n, err := strconv.Atoi(s); err == nil {
			return n
		}
	}
	return -1
}
