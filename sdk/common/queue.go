package common

type Queue struct {
	channel chan *WorkResponse
	size    int
}

func NewQueue(size int) *Queue {
	return &Queue{
		channel: make(chan *WorkResponse, size),
		size:    size,
	}
}

func (q *Queue) Put(value *WorkResponse) bool {
	select {
	case q.channel <- value:
		return true
	default:
		return false
	}
}

func (q *Queue) Get() (*WorkResponse, bool) {
	select {
	case item := <-q.channel:
		return item, true
	default:
		var zeroValue *WorkResponse
		return zeroValue, false
	}
}

// GetBatch retrieves up to max elements from the queue
func (q *Queue) GetBatch(max int) []*WorkResponse {
	if max <= 0 {
		return []*WorkResponse{}
	}

	var results []*WorkResponse
	for i := 0; i < max; i++ {
		select {
		case item := <-q.channel:
			results = append(results, item)
		default:
			// No more items available
			return results
		}
	}
	return results
}

func (q *Queue) Empty() bool {
	return len(q.channel) == 0
}

// Close closes the channel, preventing further sends
func (q *Queue) Close() {
	close(q.channel)
}

// Size returns the current number of items in the queue
func (q *Queue) Size() int {
	return len(q.channel)
}

// Capacity returns the maximum capacity of the queue
func (q *Queue) Capacity() int {
	return cap(q.channel)
}
