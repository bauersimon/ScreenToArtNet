package ambilight

type workerPool struct {
	inputChannel chan job
	size         int
}

type queue []func()

func (q *queue) enqueue(task func()) {
	*q = append(*q, task)
}

func (q *queue) dequeue() func() {
	job, newQueue := (*q)[0], (*q)[1:]
	*q = newQueue
	return job
}

func (q *queue) isEmpty() bool {
	return len(*q) == 0
}

type job struct {
	doneChannel chan interface{}
	task        func()
}

func newWorkerPool(size int) workerPool {
	inputChannel := make(chan job)
	wp := workerPool{inputChannel: inputChannel, size: size}
	for i := 0; i < size; i++ {
		go spawnWorker(inputChannel)
	}
	return wp
}

func (w *workerPool) workOn(taskQueue queue) {
	doneChan := make(chan interface{})
	defer close(doneChan)

	for {
		jobsTaken := 0
		for i := 0; i < w.size; i++ {
			if taskQueue.isEmpty() {
				break
			}
			w.inputChannel <- job{task: taskQueue.dequeue(), doneChannel: doneChan}
			jobsTaken++
		}
		for i := 0; i < jobsTaken; i++ {
			<-doneChan

		}
		if jobsTaken == 0 {
			break
		}
	}

}

func spawnWorker(inputChannel chan job) {
	for job := range inputChannel {
		job.task()
		job.doneChannel <- struct{}{}
	}
}
