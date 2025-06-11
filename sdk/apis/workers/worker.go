package workers

import "fmt"

type Worker struct {
	ExecutionMethod interface{}
	Name            string
	namespace       string
	maxInProgress   int
}

func NewWorker(ExecutionMethod interface{}, Name string) *Worker {
	worker := &Worker{
		ExecutionMethod: ExecutionMethod,
		Name:            Name,
	}
	worker.SetMaxInProgress(100)
	worker.SetNamespace("default")
	return worker
}

func (worker *Worker) SetMaxInProgress(maxInProgress int) {
	worker.maxInProgress = maxInProgress
}

func (worker *Worker) GetMaxInProgress() int {
	return worker.maxInProgress
}

func (worker *Worker) SetNamespace(namespace string) {
	worker.namespace = namespace
}
func (worker *Worker) GetNamespace() string {
	return worker.namespace
}

func (w *Worker) SetExecutionMethod(ExecutionMethod interface{}) {
	w.ExecutionMethod = ExecutionMethod
}

func (w *Worker) GetExecutionMethod() interface{} {
	return w.ExecutionMethod
}

func (w *Worker) SetName(Name string) {
	w.Name = Name
}

func (w *Worker) GetName() string {
	return w.Name
}

func (w *Worker) String() string {
	return fmt.Sprintf("Worker(name=%s, execution_method=%p)", w.Name, w.ExecutionMethod)
}
