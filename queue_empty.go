//
//
//

package queue

import (
	"sync"

	list "github.com/ondi/go-circular"
)

type Empty_t[Value_t any] struct {
	buf     *list.List_t[Value_t]
	reader  *sync.Cond
	writer  *sync.Cond
	readers int
	writers int
	limit   int
	state   int
}

// reader should Signal() to wake up Wait()
func (self *Empty_t[Value_t]) PushFront(value Value_t) int {
	self.writers++
	for self.state == 1 && (self.buf.Size() == 1 || self.readers == 0) {
		self.reader.Wait()
	}
	self.writers--
	if self.state == 1 && self.buf.PushFront(value) {
		self.writer.Broadcast()
		return 0
	}
	return self.state
}

// reader should Signal() to wake up Wait()
func (self *Empty_t[Value_t]) PushBack(value Value_t) int {
	self.writers++
	for self.state == 1 && (self.buf.Size() == 1 || self.readers == 0) {
		self.reader.Wait()
	}
	self.writers--
	if self.state == 1 && self.buf.PushBack(value) {
		self.writer.Broadcast()
		return 0
	}
	return self.state
}

func (self *Empty_t[Value_t]) PopFront() (Value_t, int) {
	self.readers++
	self.reader.Broadcast() // Signal() writer abount reader
	for self.state == 1 && self.buf.Size() == 0 {
		self.writer.Wait()
	}
	self.readers--
	value, ok := self.buf.PopFront()
	if ok {
		self.reader.Broadcast()
		return value, 0
	}
	return value, self.state
}

func (self *Empty_t[Value_t]) PopBack() (Value_t, int) {
	self.readers++
	self.reader.Broadcast() // Signal() writer abount reader
	for self.state == 1 && self.buf.Size() == 0 {
		self.writer.Wait()
	}
	self.readers--
	value, ok := self.buf.PopBack()
	if ok {
		self.reader.Broadcast()
		return value, 0
	}
	return value, self.state
}

func (self *Empty_t[Value_t]) PushFrontNoLock(value Value_t) int {
	self.writers++
	for self.state == 1 && self.buf.Size() == 1 && self.readers >= self.writers {
		self.reader.Wait() // Broadcast required
	}
	self.writers--
	if self.state == 1 && self.readers > 0 && self.buf.PushFront(value) {
		self.writer.Broadcast()
		return 0
	}
	return self.state
}

func (self *Empty_t[Value_t]) PushBackNoLock(value Value_t) int {
	self.writers++
	for self.state == 1 && self.buf.Size() == 1 && self.readers >= self.writers {
		self.reader.Wait() // Broadcast required
	}
	self.writers--
	if self.state == 1 && self.readers > 0 && self.buf.PushBack(value) {
		self.writer.Broadcast()
		return 0
	}
	return self.state
}

func (self *Empty_t[Value_t]) PopFrontNoLock() (Value_t, int) {
	self.readers++
	self.reader.Broadcast() // Signal() writer abount reader
	for self.state == 1 && self.buf.Size() == 0 && self.writers >= self.readers {
		self.writer.Wait() // Broadcast required
	}
	self.readers--
	value, ok := self.buf.PopFront()
	if ok {
		self.reader.Broadcast()
		return value, 0
	}
	return value, self.state
}

func (self *Empty_t[Value_t]) PopBackNoLock() (Value_t, int) {
	self.readers++
	self.reader.Broadcast() // Signal() writer abount reader
	for self.state == 1 && self.buf.Size() == 0 && self.writers >= self.readers {
		self.writer.Wait() // Broadcast required
	}
	self.readers--
	value, ok := self.buf.PopBack()
	if ok {
		self.reader.Broadcast()
		return value, 0
	}
	return value, self.state
}

func (self *Empty_t[Value_t]) Size() int {
	return self.buf.Size()
}

func (self *Empty_t[Value_t]) Readers() int {
	return self.readers
}

func (self *Empty_t[Value_t]) Writers() int {
	return self.writers
}

func (self *Empty_t[Value_t]) RangeFront(f func(Value_t) bool) {
	self.buf.RangeFront(f)
}

func (self *Empty_t[Value_t]) RangeBack(f func(Value_t) bool) {
	self.buf.RangeBack(f)
}

func (self *Empty_t[Value_t]) Close() {
	self.state = -1
	self.writer.Broadcast()
	self.reader.Broadcast()
}
