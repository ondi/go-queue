//
//
//

package queue

import (
	"sync"

	// list "github.com/ondi/go-list"
	list "github.com/ondi/go-circular"
)

type Beffered_t[Value_t any] struct {
	buf     *list.List_t[Value_t]
	reader  *sync.Cond
	writer  *sync.Cond
	readers int
	writers int
	limit   int
	open    bool
}

func (self *Beffered_t[Value_t]) PushFront(value Value_t) bool {
	self.writers++
	for self.open && self.buf.Size() == self.limit {
		self.reader.Wait()
	}
	self.writers--
	if self.open && self.buf.PushFront(value) {
		self.writer.Broadcast()
		return true
	}
	return false
}

func (self *Beffered_t[Value_t]) PushBack(value Value_t) bool {
	self.writers++
	for self.open && self.buf.Size() == self.limit {
		self.reader.Wait()
	}
	self.writers--
	if self.open && self.buf.PushBack(value) {
		self.writer.Broadcast()
		return true
	}
	return false
}

func (self *Beffered_t[Value_t]) PopFront() (value Value_t, ok bool) {
	self.readers++
	for self.open && self.buf.Size() == 0 {
		self.writer.Wait()
	}
	self.readers--
	if value, ok = self.buf.PopFront(); ok {
		self.reader.Broadcast()
	}
	return
}

func (self *Beffered_t[Value_t]) PopBack() (value Value_t, ok bool) {
	self.readers++
	for self.open && self.buf.Size() == 0 {
		self.writer.Wait()
	}
	self.readers--
	if value, ok = self.buf.PopBack(); ok {
		self.reader.Broadcast()
	}
	return
}

func (self *Beffered_t[Value_t]) PushFrontNoLock(value Value_t) bool {
	self.writers++
	for self.open && self.buf.Size() == self.limit && self.readers >= self.writers {
		self.reader.Wait() // Broadcast required
	}
	self.writers--
	if self.open && self.buf.PushFront(value) {
		self.writer.Broadcast()
		return true
	}
	return false
}

func (self *Beffered_t[Value_t]) PushBackNoLock(value Value_t) bool {
	self.writers++
	for self.open && self.buf.Size() == self.limit && self.readers >= self.writers {
		self.reader.Wait() // Broadcast required
	}
	self.writers--
	if self.open && self.buf.PushBack(value) {
		self.writer.Broadcast()
		return true
	}
	return false
}

func (self *Beffered_t[Value_t]) PopFrontNoLock() (value Value_t, ok bool) {
	self.readers++
	for self.open && self.buf.Size() == 0 && self.writers >= self.readers {
		self.writer.Wait() // Broadcast required
	}
	self.readers--
	if value, ok = self.buf.PopFront(); ok {
		self.reader.Broadcast()
	}
	return
}

func (self *Beffered_t[Value_t]) PopBackNoLock() (value Value_t, ok bool) {
	self.readers++
	for self.open && self.buf.Size() == 0 && self.writers >= self.readers {
		self.writer.Wait() // Broadcast required
	}
	self.readers--
	if value, ok = self.buf.PopBack(); ok {
		self.reader.Broadcast()
	}
	return
}

func (self *Beffered_t[Value_t]) RangeFront(f func(Value_t) bool) {
	self.buf.RangeFront(f)
}

func (self *Beffered_t[Value_t]) RangeBack(f func(Value_t) bool) {
	self.buf.RangeBack(f)
}

func (self *Beffered_t[Value_t]) Readers() int {
	return self.readers
}

func (self *Beffered_t[Value_t]) Writers() int {
	return self.writers
}

func (self *Beffered_t[Value_t]) Limit() int {
	return self.limit
}

func (self *Beffered_t[Value_t]) Size() int {
	return self.buf.Size()
}

func (self *Beffered_t[Value_t]) Close() {
	self.open = false
	self.writer.Broadcast()
	self.reader.Broadcast()
}

func (self *Beffered_t[Value_t]) Closed() bool {
	return self.open == false
}
