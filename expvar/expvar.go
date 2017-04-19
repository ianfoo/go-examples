package main

import (
	"context"
	"expvar"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

func main() {
	var (
		s      http.Server
		sigCh  = make(chan os.Signal, 1)
		quitCh = make(chan struct{}, 1)
		wg     sync.WaitGroup
	)

	signal.Notify(sigCh, os.Interrupt)
	go func() {
		<-sigCh
		quitCh <- struct{}{}
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		wg.Wait()
		s.Shutdown(ctx)
	}()

	wg.Add(1)
	go forever(quitCh, &wg)

	addr := ":5050"
	if envAddr := os.Getenv("ADDR"); envAddr != "" {
		addr = envAddr
	}
	s.Addr = addr
	log.Println(s.ListenAndServe())
}

func forever(quitCh <-chan struct{}, wg *sync.WaitGroup) {
	log.Println("starting: adding generated strings to slice")
	i := int64(1)
	ex := expvar.NewInt("sliceLength")
	x := make([]string, 0)
LOOP:
	for ; ; i++ {
		ex.Set(i)
		select {
		case <-quitCh:
			break LOOP
		default:
			label := fmt.Sprintf("member%06d:%d", i, rand.Intn(1e7))
			x = append(x, label)
			if i%100 == 0 {
				log.Printf("added %d entries\n", i)
			}
			time.Sleep(50 * time.Millisecond)
		}
	}
	log.Printf("quitting: added %d entries\n", i)
	wg.Done()
}
