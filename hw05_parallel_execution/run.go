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
	mutexOne := sync.Mutex{}
	//mutexTwo := sync.Mutex{}
	for i := 1; i <= n; i++ {
		go func(erCnt *int, ci chan Task) {
			defer wg.Done()
			wg.Add(1)
			for task := range ci {
				taskError := task()
				if taskError != nil {
					mutexOne.Lock()
					*erCnt++
					mutexOne.Unlock()
				}
			}
		}(&errCount, ch)
	}
	for _, t := range tasks {
		mutexOne.Lock()
		errCnt := errCount
		mutexOne.Unlock()
		if errCnt == m {
			close(ch)
			return ErrErrorsLimitExceeded
		}
		ch <- t
	}
	close(ch)
	wg.Wait()
	return nil
}
