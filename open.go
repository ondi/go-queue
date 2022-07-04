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
	mx      sync.Locker
	readers int
	writers int
	limit   int
	closed  int
}

func NewOpen[Value_t any](mx sync.Locker, limit int) Queue[Value_t] {
	self := &Open_t[Value_t]{
		buf:    list.New[Value_t](limit + 1),
		reader: sync.NewCond(mx),
		writer: sync.NewCond(mx),
		mx:     mx,
		limit:  limit,
	}
	return self
}

func (self *Open_t[Value_t]) PushFront(value Value_t) int {
	self.writers++
	for self.buf.Size() > self.limit || self.buf.Size() == self.limit && self.readers == 0 {
		self.reader.Wait()
	}
	self.writers--
	self.buf.PushFront(value)
	self.writer.Signal()
	return self.closed
}

func (self *Open_t[Value_t]) PushBack(value Value_t) int {
	self.writers++
	for self.buf.Size() > self.limit || self.buf.Size() == self.limit && self.readers == 0 {
		self.reader.Wait()
	}
	self.writers--
	self.buf.PushBack(value)
	self.writer.Signal()
	return self.closed
}

func (self *Open_t[Value_t]) PushFrontNoWait(value Value_t) int {
	if self.buf.Size() > self.limit || self.buf.Size() == self.limit && self.readers == 0 {
		return 1
	}
	self.buf.PushFront(value)
	self.writer.Signal()
	return 0
}

func (self *Open_t[Value_t]) PushBackNoWait(value Value_t) int {
	if self.buf.Size() > self.limit || self.buf.Size() == self.limit && self.readers == 0 {
		return 1
	}
	self.buf.PushBack(value)
	self.writer.Signal()
	return 0
}

func (self *Open_t[Value_t]) PopFront() (Value_t, int) {
	self.readers++
	self.reader.Signal()
	for self.buf.Size() == 0 {
		self.writer.Wait()
	}
	self.readers--
	value, _ := self.buf.PopFront()
	return value, self.closed
}

func (self *Open_t[Value_t]) PopBack() (Value_t, int) {
	self.readers++
	self.reader.Signal()
	for self.buf.Size() == 0 {
		self.writer.Wait()
	}
	self.readers--
	value, _ := self.buf.PopBack()
	return value, self.closed
}

func (self *Open_t[Value_t]) PopFrontNoWait() (v Value_t, res int) {
	self.readers++
	self.reader.Signal()
	for self.buf.Size() == 0 && self.writers >= self.readers {
		self.writer.Wait()
	}
	self.readers--
	if value, ok := self.buf.PopFront(); ok {
		return value, 0
	}
	res = 1
	return
}

func (self *Open_t[Value_t]) PopBackNoWait() (v Value_t, res int) {
	self.readers++
	self.reader.Signal()
	for self.buf.Size() == 0 && self.writers >= self.readers {
		self.writer.Wait()
	}
	self.readers--
	if value, ok := self.buf.PopBack(); ok {
		return value, 0
	}
	res = 1
	return
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

func (self *Open_t[Value_t]) Close() (buf List[Value_t]) {
	self.closed = -1
	// readers may read data after close
	buf = self.buf
	// create buffer for pending readers and writers
	self.buf = list.New[Value_t](self.readers + self.writers + 1)
	// release pending readers (waiting for buf.Size() > 0)
	var temp Value_t
	for i := 0; i < self.readers; i++ {
		self.buf.PushBack(temp)
	}
	self.writer.Broadcast()
	// release pending writers (waiting for buf.Size() <= limit)
	self.limit = self.buf.Size() - 1
	self.reader.Broadcast()
	return
}
