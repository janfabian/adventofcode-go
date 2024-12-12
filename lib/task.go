package lib

import (
	"runtime"
	"sync"
)

var once sync.Once
var jobs chan func()

func AddTask(fn func()) {
	once.Do(func() {
		poolSize := runtime.NumCPU() * 2
		jobs = make(chan func(), poolSize)

		for i := 0; i < poolSize; i++ {
			go func() {
				for job := range jobs {
					job()
				}
			}()
		}
	})

	jobs <- fn
}

func Close() {
	close(jobs)
}
