package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// структура  с помощью которой будем управлять воркерами
type WorkerPool struct {
	jobs        chan string
	workerCount int
	workers     map[int]chan struct{}
	mu          sync.Mutex
}

// инициализируем новый WorkerPool
func NewWorkerPool() *WorkerPool {
	return &WorkerPool{
		jobs:    make(chan string),
		workers: make(map[int]chan struct{}),
	}
}

// добавление воркера
func (wp *WorkerPool) AddWorker() {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	wp.workerCount++

	workerID := wp.workerCount
	stopChan := make(chan struct{})

	wp.workers[workerID] = stopChan

	go func(id int, stop <-chan struct{}) {
		for {
			select {
			case job := <-wp.jobs:
				fmt.Printf("воркер %d выполняет задание: %s\n", id, job)
				time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
			case <-stop:
				fmt.Printf("воркер %d удален \n", id)
				return
			}
		}
	}(workerID, stopChan)
}

// удаление воркера
func (wp *WorkerPool) RemoveWorker() {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	if wp.workerCount == 1 {
		fmt.Println("нет других воркеров для выполнения")
		return
	}

	workerID := wp.workerCount
	stopChan := wp.workers[workerID]
	close(stopChan)

	delete(wp.workers, workerID)
	wp.workerCount--
}

func (wp *WorkerPool) AddJob(job string) {
	wp.jobs <- job
}

func main() {
	wp := NewWorkerPool()

	wp.AddWorker()
	wp.AddWorker()
	wp.AddWorker()

	for i := 0; i < 10; i++ {
		wp.AddJob(fmt.Sprintf("задача %d", i+1))
		time.Sleep(100 * time.Millisecond)
	}

	wp.RemoveWorker()
	wp.RemoveWorker()
	wp.RemoveWorker()

	for i := 10; i < 30; i++ {
		wp.AddJob(fmt.Sprintf("задача %d", i+1))
		time.Sleep(100 * time.Millisecond)
	}

	time.Sleep(3 * time.Second)

}
