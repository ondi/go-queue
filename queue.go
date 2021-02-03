//
//
//

package queue

import "sync"

type List interface {
	PushFront(interface{}) bool
	PushBack(interface{}) bool
	PopFront() (interface{}, bool)
	PopBack() (interface{}, bool)

	RangeFront(func(interface{}) bool)
	RangeBack(func(interface{}) bool)

	Size() int
}

type Queue interface {
	PushFront(interface{}) int
	PushBack(interface{}) int
	PopFront() (interface{}, int)
	PopBack() (interface{}, int)

	PushFrontNoWait(interface{}) int
	PushBackNoWait(interface{}) int
	PopFrontNoWait() (interface{}, int)
	PopBackNoWait() (interface{}, int)

	RangeFront(func(interface{}) bool)
	RangeBack(func(interface{}) bool)

	Readers() int
	Writers() int
	Size() int

	Close() List
}

type Queue_t struct {
	mx sync.Mutex
	q  Queue
}

func New(limit int) Queue {
	self := &Queue_t{}
	self.q = NewOpen(&self.mx, limit)
	return self
}

func (self *Queue_t) PushFront(value interface{}) (ok int) {
	self.mx.Lock()
	ok = self.q.PushFront(value)
	self.mx.Unlock()
	return
}

func (self *Queue_t) PushFrontNoWait(value interface{}) (ok int) {
	self.mx.Lock()
	ok = self.q.PushFrontNoWait(value)
	self.mx.Unlock()
	return
}

func (self *Queue_t) PushBack(value interface{}) (ok int) {
	self.mx.Lock()
	ok = self.q.PushBack(value)
	self.mx.Unlock()
	return
}

func (self *Queue_t) PushBackNoWait(value interface{}) (ok int) {
	self.mx.Lock()
	ok = self.q.PushBackNoWait(value)
	self.mx.Unlock()
	return
}

func (self *Queue_t) PopFront() (value interface{}, ok int) {
	self.mx.Lock()
	value, ok = self.q.PopFront()
	self.mx.Unlock()
	return
}

func (self *Queue_t) PopFrontNoWait() (value interface{}, ok int) {
	self.mx.Lock()
	value, ok = self.q.PopFrontNoWait()
	self.mx.Unlock()
	return
}

func (self *Queue_t) PopBack() (value interface{}, ok int) {
	self.mx.Lock()
	value, ok = self.q.PopBack()
	self.mx.Unlock()
	return
}

func (self *Queue_t) PopBackNoWait() (value interface{}, ok int) {
	self.mx.Lock()
	value, ok = self.q.PopBackNoWait()
	self.mx.Unlock()
	return
}

func (self *Queue_t) Size() (res int) {
	self.mx.Lock()
	res = self.q.Size()
	self.mx.Unlock()
	return
}

func (self *Queue_t) Readers() (res int) {
	self.mx.Lock()
	res = self.q.Readers()
	self.mx.Unlock()
	return
}

func (self *Queue_t) Writers() (res int) {
	self.mx.Lock()
	res = self.q.Writers()
	self.mx.Unlock()
	return
}

func (self *Queue_t) RangeFront(f func(interface{}) bool) {
	self.mx.Lock()
	self.q.RangeFront(f)
	self.mx.Unlock()
}

func (self *Queue_t) RangeBack(f func(interface{}) bool) {
	self.mx.Lock()
	self.q.RangeBack(f)
	self.mx.Unlock()
}

func (self *Queue_t) Close() List {
	self.mx.Lock()
	self.q = NewClosed(self.q.Close())
	self.mx.Unlock()
	return nil
}
