//
//
//

package queue

type Closed_t struct {
	buf List
}

func NewClosed(buf List) Queue {
	return &Closed_t{buf: buf}
}

func (self *Closed_t) PushFront(interface{}) int {
	return -1
}

func (self *Closed_t) PushBack(interface{}) int {
	return -1
}

func (self *Closed_t) PushFrontNoWait(interface{}) int {
	return -1
}

func (self *Closed_t) PushBackNoWait(interface{}) int {
	return -1
}

func (self *Closed_t) PopFront() (interface{}, int) {
	if value, ok := self.buf.PopFront(); ok {
		return value, 0
	}
	return nil, -1
}

func (self *Closed_t) PopBack() (interface{}, int) {
	if value, ok := self.buf.PopBack(); ok {
		return value, 0
	}
	return nil, -1
}

func (self *Closed_t) PopFrontNoWait() (interface{}, int) {
	if value, ok := self.buf.PopFront(); ok {
		return value, 0
	}
	return nil, -1
}

func (self *Closed_t) PopBackNoWait() (interface{}, int) {
	if value, ok := self.buf.PopBack(); ok {
		return value, 0
	}
	return nil, -1
}

func (self *Closed_t) Size() int {
	return self.buf.Size()
}

func (self *Closed_t) Readers() int {
	return 0
}

func (self *Closed_t) Writers() int {
	return 0
}

func (self *Closed_t) RangeFront(f func(interface{}) bool) {
	self.buf.RangeFront(f)
}

func (self *Closed_t) RangeBack(f func(interface{}) bool) {
	self.buf.RangeBack(f)
}

func (self *Closed_t) Close() List {
	return self.buf
}
