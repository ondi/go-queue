//
//
//

package queue

import "sync"

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

	PushFrontNoWait(Value_t) int
	PushBackNoWait(Value_t) int
	PopFrontNoWait() (Value_t, int)
	PopBackNoWait() (Value_t, int)

	RangeFront(func(Value_t) bool)
	RangeBack(func(Value_t) bool)

	Readers() int
	Writers() int
	Size() int

	Close() List[Value_t]
}

type Queue_t[Value_t any] struct {
	mx sync.Mutex
	q  Queue[Value_t]
}

func New[Value_t any](limit int) Queue[Value_t] {
	self := &Queue_t[Value_t]{}
	self.q = NewOpen[Value_t](&self.mx, limit)
	return self
}

func (self *Queue_t[Value_t]) PushFront(value Value_t) (ok int) {
	self.mx.Lock()
	ok = self.q.PushFront(value)
	self.mx.Unlock()
	return
}

func (self *Queue_t[Value_t]) PushFrontNoWait(value Value_t) (ok int) {
	self.mx.Lock()
	ok = self.q.PushFrontNoWait(value)
	self.mx.Unlock()
	return
}

func (self *Queue_t[Value_t]) PushBack(value Value_t) (ok int) {
	self.mx.Lock()
	ok = self.q.PushBack(value)
	self.mx.Unlock()
	return
}

func (self *Queue_t[Value_t]) PushBackNoWait(value Value_t) (ok int) {
	self.mx.Lock()
	ok = self.q.PushBackNoWait(value)
	self.mx.Unlock()
	return
}

func (self *Queue_t[Value_t]) PopFront() (value Value_t, ok int) {
	self.mx.Lock()
	value, ok = self.q.PopFront()
	self.mx.Unlock()
	return
}

func (self *Queue_t[Value_t]) PopFrontNoWait() (value Value_t, ok int) {
	self.mx.Lock()
	value, ok = self.q.PopFrontNoWait()
	self.mx.Unlock()
	return
}

func (self *Queue_t[Value_t]) PopBack() (value Value_t, ok int) {
	self.mx.Lock()
	value, ok = self.q.PopBack()
	self.mx.Unlock()
	return
}

func (self *Queue_t[Value_t]) PopBackNoWait() (value Value_t, ok int) {
	self.mx.Lock()
	value, ok = self.q.PopBackNoWait()
	self.mx.Unlock()
	return
}

func (self *Queue_t[Value_t]) Size() (res int) {
	self.mx.Lock()
	res = self.q.Size()
	self.mx.Unlock()
	return
}

func (self *Queue_t[Value_t]) Readers() (res int) {
	self.mx.Lock()
	res = self.q.Readers()
	self.mx.Unlock()
	return
}

func (self *Queue_t[Value_t]) Writers() (res int) {
	self.mx.Lock()
	res = self.q.Writers()
	self.mx.Unlock()
	return
}

func (self *Queue_t[Value_t]) RangeFront(f func(Value_t) bool) {
	self.mx.Lock()
	self.q.RangeFront(f)
	self.mx.Unlock()
}

func (self *Queue_t[Value_t]) RangeBack(f func(Value_t) bool) {
	self.mx.Lock()
	self.q.RangeBack(f)
	self.mx.Unlock()
}

func (self *Queue_t[Value_t]) Close() List[Value_t] {
	self.mx.Lock()
	self.q = NewClosed(self.q.Close())
	self.mx.Unlock()
	return nil
}
