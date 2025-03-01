package concurrency_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/andrejacobs/go-aj/concurrency"
	"github.com/andrejacobs/go-aj/random"
	"github.com/stretchr/testify/assert"
)

func TestDuplicateOutputs(t *testing.T) {
	expectedCount := 10000
	producer := make(chan int, 1000)

	consumerCount := 100
	consumers := make([]chan int, consumerCount)
	for i := 0; i < consumerCount; i++ {
		consumers[i] = make(chan int, 100)
	}

	// Start producing
	go func() {
		for i := 0; i < expectedCount; i++ {
			producer <- i
		}
		close(producer)
	}()

	// Duplicate the producer to output to multiple channels
	go concurrency.DuplicateOutputs[int](context.Background(), producer, consumers...)

	// Consume from all the duplicate producers
	wg := sync.WaitGroup{}
	for i := 0; i < consumerCount; i++ {
		wg.Add(1)
		go func(consumer chan int) {
			received := make([]int, 0, expectedCount)
			for v := range consumer {
				received = append(received, v)
			}
			wg.Done()
			// Verify
			assert.Equal(t, expectedCount, len(received))
			for i := 0; i < len(received); i++ {
				assert.Equal(t, i, received[i])
			}
		}(consumers[i])
	}

	wg.Wait()
}

func TestDuplicateTransformedOutputs(t *testing.T) {
	expectedCount := 10000
	producer := make(chan int, 1000)

	consumerCount := 100
	consumers := make([]chan int, consumerCount)
	for i := 0; i < consumerCount; i++ {
		consumers[i] = make(chan int, 100)
	}

	// Start producing
	go func() {
		for i := 0; i < expectedCount; i++ {
			producer <- i
		}
		close(producer)
	}()

	// Duplicate the producer to output to multiple channels
	go concurrency.DuplicateTransformedOutputs[int](context.Background(),
		func(in int) int {
			return in * 2
		},
		producer, consumers...)

	// Consume from all the duplicate producers
	wg := sync.WaitGroup{}
	for i := 0; i < consumerCount; i++ {
		wg.Add(1)
		go func(consumer chan int) {
			received := make([]int, 0, expectedCount)
			for v := range consumer {
				received = append(received, v)
			}
			wg.Done()
			// Verify
			assert.Equal(t, expectedCount, len(received))
			for i := 0; i < len(received); i++ {
				assert.Equal(t, i*2, received[i])
			}
		}(consumers[i])
	}

	wg.Wait()
}

func TestDuplicateOutputsWithTimeout(t *testing.T) {
	expectedCount := 1000
	producer := make(chan int)

	consumerCount := 10
	consumers := make([]chan int, consumerCount)
	for i := 0; i < consumerCount; i++ {
		consumers[i] = make(chan int)
	}

	// Start producing
	go func() {
		for i := 0; i < expectedCount; i++ {
			producer <- i
		}
		// Never close the channel so the timeout will kick in first (which then closes the channel)
	}()

	// Duplicate the producer to output to multiple channels
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	defer cancel()
	go concurrency.DuplicateOutputs[int](ctx, producer, consumers...)

	// Consume from all the duplicate producers
	wg := sync.WaitGroup{}
	for i := 0; i < consumerCount; i++ {
		wg.Add(1)
		go func(consumer chan int) {
			count := 0
			for range consumer {
				count++
			}
			wg.Done()
		}(consumers[i])
	}

	wg.Wait()
}

func TestDuplicateOutputsDifferentRates(t *testing.T) {
	expectedCount := 100
	producer := make(chan int, 100)

	consumerCount := 10
	consumers := make([]chan int, consumerCount)
	for i := 0; i < consumerCount; i++ {
		consumers[i] = make(chan int, random.Int(4, 20))
	}

	// Start producing
	go func() {
		for i := 0; i < expectedCount; i++ {
			producer <- i
			time.Sleep(time.Millisecond)
		}
		close(producer)
		// fmt.Printf("AJ### Producer is finished\n")
	}()

	// Duplicate the producer to output to multiple channels
	go concurrency.DuplicateOutputs[int](context.Background(), producer, consumers...)

	// Consume from all the duplicate producers
	wg := sync.WaitGroup{}
	for i := 0; i < consumerCount; i++ {
		wg.Add(1)
		go func(consumer chan int, delay time.Duration) {
			received := make([]int, 0, expectedCount)
			for v := range consumer {
				received = append(received, v)
				// fmt.Printf("[%v] %d\n", delay, v)
				time.Sleep(delay)
			}
			wg.Done()
			// Verify
			assert.Equal(t, expectedCount, len(received))
			for i := 0; i < len(received); i++ {
				assert.Equal(t, i, received[i])
			}
		}(consumers[i], time.Millisecond*time.Duration(i))
	}

	wg.Wait()
}
