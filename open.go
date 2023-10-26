//
//
//

package queue

import (
	"sync"

	// list "github.com/ondi/go-list"
	list "github.com/ondi/go-circular"
)

type Open_t[Value_t any] struct {
	buf     *list.List_t[Value_t]
	reader  *sync.Cond
	writer  *sync.Cond
	readers int
	writers int
	limit   int
	state   int
}

func NewOpen[Value_t any](mx sync.Locker, limit int) Queue[Value_t] {
	self := &Open_t[Value_t]{
		reader: sync.NewCond(mx),
		writer: sync.NewCond(mx),
		limit:  limit,
		state:  1,
	}
	if limit == 0 {
		self.buf = list.New[Value_t](1)
	} else {
		self.buf = list.New[Value_t](limit)
	}
	return self
}

func (self *Open_t[Value_t]) PushFront(value Value_t) int {
	self.writers++
	for self.state == 1 && (self.buf.Size() > self.limit || self.buf.Size() == self.limit && self.readers == 0) {
		self.reader.Wait()
	}
	self.writers--
	if self.state == 1 && self.buf.PushFront(value) {
		self.writer.Broadcast()
		return 0
	}
	return self.state
}

func (self *Open_t[Value_t]) PushBack(value Value_t) int {
	self.writers++
	for self.state == 1 && (self.buf.Size() > self.limit || self.buf.Size() == self.limit && self.readers == 0) {
		self.reader.Wait()
	}
	self.writers--
	if self.state == 1 && self.buf.PushBack(value) {
		self.writer.Broadcast()
		return 0
	}
	return self.state
}

func (self *Open_t[Value_t]) PopFront() (Value_t, int) {
	self.readers++
	self.reader.Signal()
	for self.state == 1 && self.buf.Size() == 0 {
		self.writer.Wait()
	}
	self.readers--
	value, ok := self.buf.PopFront()
	if ok {
		self.reader.Signal()
		return value, 0
	}
	return value, self.state
}

func (self *Open_t[Value_t]) PopBack() (Value_t, int) {
	self.readers++
	self.reader.Signal()
	for self.state == 1 && self.buf.Size() == 0 {
		self.writer.Wait()
	}
	self.readers--
	value, ok := self.buf.PopBack()
	if ok {
		self.reader.Signal()
		return value, 0
	}
	return value, self.state
}

func (self *Open_t[Value_t]) PushFrontNoWait(value Value_t) int {
	self.writers++
	for self.state == 1 && self.buf.Size() != 0 && self.readers >= self.writers {
		self.reader.Wait()
	}
	self.writers--
	if self.state == 1 && self.buf.PushFront(value) {
		self.writer.Broadcast()
		return 0
	}
	return self.state
}

func (self *Open_t[Value_t]) PushBackNoWait(value Value_t) int {
	self.writers++
	for self.state == 1 && self.buf.Size() != 0 && self.readers >= self.writers {
		self.reader.Wait()
	}
	self.writers--
	if self.state == 1 && self.buf.PushBack(value) {
		self.writer.Broadcast()
		return 0
	}
	return self.state
}

func (self *Open_t[Value_t]) PopFrontNoWait() (Value_t, int) {
	self.readers++
	self.reader.Signal()
	for self.state == 1 && self.buf.Size() == 0 && self.writers >= self.readers {
		self.writer.Wait() // Broadcast required
	}
	self.readers--
	value, ok := self.buf.PopFront()
	if ok {
		self.reader.Signal()
		return value, 0
	}
	return value, self.state
}

func (self *Open_t[Value_t]) PopBackNoWait() (Value_t, int) {
	self.readers++
	self.reader.Signal()
	for self.state == 1 && self.buf.Size() == 0 && self.writers >= self.readers {
		self.writer.Wait() // Broadcast required
	}
	self.readers--
	value, ok := self.buf.PopBack()
	if ok {
		self.reader.Signal()
		return value, 0
	}
	return value, self.state
}

func (self *Open_t[Value_t]) Size() int {
	return self.buf.Size()
}

func (self *Open_t[Value_t]) Readers() int {
	return self.readers
}

func (self *Open_t[Value_t]) Writers() int {
	return self.writers
}

func (self *Open_t[Value_t]) RangeFront(f func(Value_t) bool) {
	self.buf.RangeFront(f)
}

func (self *Open_t[Value_t]) RangeBack(f func(Value_t) bool) {
	self.buf.RangeBack(f)
}

func (self *Open_t[Value_t]) Close() {
	self.state = -1
	self.writer.Broadcast()
	self.reader.Broadcast()
}
