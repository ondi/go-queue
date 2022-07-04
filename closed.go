//
//
//

package queue

type Closed_t[Value_t any] struct {
	buf List[Value_t]
}

func NewClosed[Value_t any](buf List[Value_t]) Queue[Value_t] {
	return &Closed_t[Value_t]{buf: buf}
}

func (self *Closed_t[Value_t]) PushFront(Value_t) int {
	return -1
}

func (self *Closed_t[Value_t]) PushBack(Value_t) int {
	return -1
}

func (self *Closed_t[Value_t]) PushFrontNoWait(Value_t) int {
	return -1
}

func (self *Closed_t[Value_t]) PushBackNoWait(Value_t) int {
	return -1
}

func (self *Closed_t[Value_t]) PopFront() (v Value_t, res int) {
	if value, ok := self.buf.PopFront(); ok {
		return value, 0
	}
	res = -1
	return
}

func (self *Closed_t[Value_t]) PopBack() (v Value_t, res int) {
	if value, ok := self.buf.PopBack(); ok {
		return value, 0
	}
	res = -1
	return
}

func (self *Closed_t[Value_t]) PopFrontNoWait() (v Value_t, res int) {
	if value, ok := self.buf.PopFront(); ok {
		return value, 0
	}
	res = -1
	return
}

func (self *Closed_t[Value_t]) PopBackNoWait() (v Value_t, res int) {
	if value, ok := self.buf.PopBack(); ok {
		return value, 0
	}
	res = -1
	return
}

func (self *Closed_t[Value_t]) Size() int {
	return self.buf.Size()
}

func (self *Closed_t[Value_t]) Readers() int {
	return 0
}

func (self *Closed_t[Value_t]) Writers() int {
	return 0
}

func (self *Closed_t[Value_t]) RangeFront(f func(Value_t) bool) {
	self.buf.RangeFront(f)
}

func (self *Closed_t[Value_t]) RangeBack(f func(Value_t) bool) {
	self.buf.RangeBack(f)
}

func (self *Closed_t[Value_t]) Close() List[Value_t] {
	return self.buf
}
