package common

type Queue struct {
	items []*WorkResponse
	size  int
}

func NewQueue(size int) *Queue {
	return &Queue{
		items: []*WorkResponse{},
		size:  size,
	}
}

func (q *Queue) Put(value *WorkResponse) bool {

	if len(q.items) >= q.size {
		return false
	}

	q.items = append(q.items, value)
	return true
}

func (q *Queue) Get() (*WorkResponse, bool) {
	if len(q.items) == 0 {
		var zeroValue *WorkResponse
		return zeroValue, false
	}

	item := q.items[0]
	q.items = q.items[1:]
	return item, true
}

func (q *Queue) Empty() bool {
	return len(q.items) == 0
}
