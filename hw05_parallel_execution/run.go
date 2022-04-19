package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	wg := sync.WaitGroup{}
	ch := make(chan Task)
	errCount := 0
	mutex := sync.Mutex{}
	for i := 1; i <= n; i++ {
		go func(ci chan Task) {
			defer wg.Done()
			wg.Add(1)
			for task := range ci {
				taskError := task()
				if taskError != nil && errCount < m {
					mutex.Lock()
					errCount++
					mutex.Unlock()
				}
			}
		}(ch)
	}
	for _, t := range tasks {
		if errCount == m {
			close(ch)
			return ErrErrorsLimitExceeded
		}
		ch <- t
	}
	close(ch)
	wg.Wait()
	return nil
}
