package main

import (
	"context"
	"errors"
	"sync"
)

type Getter interface {
	Get(ctx context.Context, address, key string) (string, error)
}

// Call `Getter.Get()` for each address in parallel.
// Returns the first successful response.
// If all requests fail, returns an error.
func Get(ctx context.Context, getter Getter, addresses []string, key string) (string, error) {
	if len(addresses) == 0 {
		return "", nil
	}
	result := make(chan string, 1)
	var wg sync.WaitGroup
	for _, address := range addresses {
		wg.Add(1)
		go func() {
			defer wg.Done()
			val, err := getter.Get(ctx, address, key)
			if err == nil {
				select { // check if the channel can still send value, if so send it
				case result <- val:
				default: // if the channel is full, do nothing and exit the goroutine
				}
			}
		}()
	}
	wg.Wait()
	select {
	case res := <-result:
		return res, nil
	default:
		return "", errors.New("All requests failed")
	}
}
