//
//
//

package queue

import (
	"runtime"
	"sync"

	// list "github.com/ondi/go-list"
	list "github.com/ondi/go-circular"
)

type Open_t struct {
	buf     *list.List_t
	reader  *sync.Cond
	writer  *sync.Cond
	mx      sync.Locker
	readers int
	writers int
	limit   int
	closed  int
}

func NewOpen(mx sync.Locker, limit int) (self *Open_t) {
	self = &Open_t{}
	self.buf = list.New(limit + 1)
	self.reader = sync.NewCond(mx)
	self.writer = sync.NewCond(mx)
	self.mx = mx
	self.limit = limit
	return
}

func (self *Open_t) PushFront(value interface{}) int {
	self.writers++
	for self.buf.Size() > self.limit || (self.buf.Size() == self.limit && self.readers == 0) {
		self.reader.Wait()
	}
	self.buf.PushFront(value)
	self.writer.Signal()
	self.writers--
	return self.closed
}

func (self *Open_t) PushBack(value interface{}) int {
	self.writers++
	for self.buf.Size() > self.limit || (self.buf.Size() == self.limit && self.readers == 0) {
		self.reader.Wait()
	}
	self.buf.PushBack(value)
	self.writer.Signal()
	self.writers--
	return self.closed
}

func (self *Open_t) PushFrontNoWait(value interface{}) int {
	if self.buf.Size() > self.limit || self.buf.Size() == self.limit && self.readers == 0 {
		return 1
	}
	self.buf.PushFront(value)
	self.writer.Signal()
	return 0
}

func (self *Open_t) PushBackNoWait(value interface{}) int {
	if self.buf.Size() > self.limit || self.buf.Size() == self.limit && self.readers == 0 {
		return 1
	}
	self.buf.PushBack(value)
	self.writer.Signal()
	return 0
}

func (self *Open_t) PopFront() (interface{}, int) {
	self.readers++
	self.reader.Signal()
	for self.buf.Size() == 0 {
		self.writer.Wait()
	}
	self.readers--
	value, _ := self.buf.PopFront()
	return value, self.closed
}

func (self *Open_t) PopBack() (interface{}, int) {
	self.readers++
	self.reader.Signal()
	for self.buf.Size() == 0 {
		self.writer.Wait()
	}
	self.readers--
	value, _ := self.buf.PopBack()
	return value, self.closed
}

func (self *Open_t) PopFrontNoWait() (interface{}, int) {
	self.readers++
	self.reader.Signal()
	for self.buf.Size() == 0 && self.writers > 0 {
		// self.writer.Wait()
		self.mx.Unlock()
		runtime.Gosched()
		self.mx.Lock()
	}
	self.readers--
	if value, ok := self.buf.PopFront(); ok {
		return value, 0
	}
	return nil, 1
}

func (self *Open_t) PopBackNoWait() (interface{}, int) {
	self.readers++
	self.reader.Signal()
	for self.buf.Size() == 0 && self.writers > 0 {
		// self.writer.Wait()
		self.mx.Unlock()
		runtime.Gosched()
		self.mx.Lock()
	}
	self.readers--
	if value, ok := self.buf.PopBack(); ok {
		return value, 0
	}
	return nil, 1
}

func (self *Open_t) Size() int {
	return self.buf.Size()
}

func (self *Open_t) Readers() int {
	return self.readers
}

func (self *Open_t) Writers() int {
	return self.writers
}

func (self *Open_t) RangeFront(f func(interface{}) bool) {
	self.buf.RangeFront(f)
}

func (self *Open_t) RangeBack(f func(interface{}) bool) {
	self.buf.RangeBack(f)
}

func (self *Open_t) Close() (buf List) {
	self.closed = -1
	// readers may read data after close
	buf = self.buf
	// create buffer for pending readers and writers
	self.buf = list.New(self.readers + self.writers + 1)
	// release pending readers (waiting for buf.Size() > 0)
	for i := 0; i < self.readers; i++ {
		self.buf.PushBack(nil)
	}
	self.writer.Broadcast()
	// release pending writers (waiting for buf.Size() <= limit)
	self.limit = self.buf.Size() - 1
	self.reader.Broadcast()
	return
}
