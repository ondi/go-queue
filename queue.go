//
//
//

package queue

import (
	"sync"

	list "github.com/ondi/go-circular"
)

type List[Value_t any] interface {
	PushFront(Value_t) bool
	PushBack(Value_t) bool
	PopFront() (Value_t, bool)
	PopBack() (Value_t, bool)

	RangeFront(func(Value_t) bool)
	RangeBack(func(Value_t) bool)

	Size() int
}

type Queue[Value_t any] interface {
	PushFront(Value_t) int
	PushBack(Value_t) int
	PopFront() (Value_t, int)
	PopBack() (Value_t, int)

	PushFrontNoLock(Value_t) int
	PushBackNoLock(Value_t) int
	PopFrontNoLock() (Value_t, int)
	PopBackNoLock() (Value_t, int)

	RangeFront(func(Value_t) bool)
	RangeBack(func(Value_t) bool)

	Readers() int
	Writers() int
	Size() int

	Close()
}

func NewOpen[Value_t any](mx sync.Locker, limit int) Queue[Value_t] {
	if limit == 0 {
		return &Empty_t[Value_t]{
			buf:    list.New[Value_t](1),
			reader: sync.NewCond(mx),
			writer: sync.NewCond(mx),
			limit:  limit,
			state:  1,
		}
	} else {
		return &Beffered_t[Value_t]{
			buf:    list.New[Value_t](limit),
			reader: sync.NewCond(mx),
			writer: sync.NewCond(mx),
			limit:  limit,
			state:  1,
		}
	}
}

type QueueSync_t[Value_t any] struct {
	mx sync.Mutex
	q  Queue[Value_t]
}

func NewSync[Value_t any](limit int) Queue[Value_t] {
	self := &QueueSync_t[Value_t]{}
	self.q = NewOpen[Value_t](&self.mx, limit)
	return self
}

func (self *QueueSync_t[Value_t]) PushFront(value Value_t) (ok int) {
	self.mx.Lock()
	ok = self.q.PushFront(value)
	self.mx.Unlock()
	return
}

func (self *QueueSync_t[Value_t]) PushFrontNoLock(value Value_t) (ok int) {
	self.mx.Lock()
	ok = self.q.PushFrontNoLock(value)
	self.mx.Unlock()
	return
}

func (self *QueueSync_t[Value_t]) PushBack(value Value_t) (ok int) {
	self.mx.Lock()
	ok = self.q.PushBack(value)
	self.mx.Unlock()
	return
}

func (self *QueueSync_t[Value_t]) PushBackNoLock(value Value_t) (ok int) {
	self.mx.Lock()
	ok = self.q.PushBackNoLock(value)
	self.mx.Unlock()
	return
}

func (self *QueueSync_t[Value_t]) PopFront() (value Value_t, ok int) {
	self.mx.Lock()
	value, ok = self.q.PopFront()
	self.mx.Unlock()
	return
}

func (self *QueueSync_t[Value_t]) PopFrontNoLock() (value Value_t, ok int) {
	self.mx.Lock()
	value, ok = self.q.PopFrontNoLock()
	self.mx.Unlock()
	return
}

func (self *QueueSync_t[Value_t]) PopBack() (value Value_t, ok int) {
	self.mx.Lock()
	value, ok = self.q.PopBack()
	self.mx.Unlock()
	return
}

func (self *QueueSync_t[Value_t]) PopBackNoLock() (value Value_t, ok int) {
	self.mx.Lock()
	value, ok = self.q.PopBackNoLock()
	self.mx.Unlock()
	return
}

func (self *QueueSync_t[Value_t]) Size() (res int) {
	self.mx.Lock()
	res = self.q.Size()
	self.mx.Unlock()
	return
}

func (self *QueueSync_t[Value_t]) Readers() (res int) {
	self.mx.Lock()
	res = self.q.Readers()
	self.mx.Unlock()
	return
}

func (self *QueueSync_t[Value_t]) Writers() (res int) {
	self.mx.Lock()
	res = self.q.Writers()
	self.mx.Unlock()
	return
}

func (self *QueueSync_t[Value_t]) RangeFront(f func(Value_t) bool) {
	self.mx.Lock()
	self.q.RangeFront(f)
	self.mx.Unlock()
}

func (self *QueueSync_t[Value_t]) RangeBack(f func(Value_t) bool) {
	self.mx.Lock()
	self.q.RangeBack(f)
	self.mx.Unlock()
}

func (self *QueueSync_t[Value_t]) Close() {
	self.mx.Lock()
	self.q.Close()
	self.mx.Unlock()
}
