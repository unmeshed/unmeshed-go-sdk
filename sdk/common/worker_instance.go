package common

import apis "github.com/unmeshed/unmeshed-go-sdk/sdk/apis/workers"

type WorkerInstance struct {
	worker   *apis.Worker
	ioThread bool
}

func NewWorkerInstance(worker *apis.Worker, ioThread bool) *WorkerInstance {
	return &WorkerInstance{
		worker:   worker,
		ioThread: ioThread,
	}
}

func (wi *WorkerInstance) GetWorker() *apis.Worker {
	return wi.worker
}

func (wi *WorkerInstance) SetWorker(worker *apis.Worker) {
	wi.worker = worker
}

func (wi *WorkerInstance) IsIOThread() bool {
	return wi.ioThread
}

func (wi *WorkerInstance) SetIOThread(ioThread bool) {
	wi.ioThread = ioThread
}
