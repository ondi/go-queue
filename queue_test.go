package queue

import (
	"testing"
	"time"

	"gotest.tools/assert"
)

func TestWritersWaitReaders1(t *testing.T) {
	q := NewSync[int32](0)

	go q.PushBack(4)
	go q.PushBack(3)
	go q.PushBack(2)
	go q.PushBack(1)
	go q.PushBack(0)

	for q.Writers() != 5 {
		time.Sleep(time.Millisecond)
	}

	_, ok := q.PopFront()
	assert.Assert(t, ok == 0, ok)

	_, ok = q.PopFront()
	assert.Assert(t, ok == 0, ok)

	_, ok = q.PopFront()
	assert.Assert(t, ok == 0, ok)

	_, ok = q.PopFront()
	assert.Assert(t, ok == 0, ok)

	_, ok = q.PopFront()
	assert.Assert(t, ok == 0, ok)
}

func TestWritersWaitReaders2(t *testing.T) {
	q := NewSync[int32](0)

	go q.PushFront(4)
	go q.PushFront(3)
	go q.PushFront(2)
	go q.PushFront(1)
	go q.PushFront(0)

	for q.Writers() != 5 {
		time.Sleep(time.Millisecond)
	}

	_, ok := q.PopBack()
	assert.Assert(t, ok == 0, ok)

	_, ok = q.PopBack()
	assert.Assert(t, ok == 0, ok)

	_, ok = q.PopBack()
	assert.Assert(t, ok == 0, ok)

	_, ok = q.PopBack()
	assert.Assert(t, ok == 0, ok)

	_, ok = q.PopBack()
	assert.Assert(t, ok == 0, ok)
}

func TestWritersWaitReaders3(t *testing.T) {
	q := NewSync[int32](0)

	go q.PushBack(4)
	go q.PushBack(3)
	go q.PushBack(2)
	go q.PushBack(1)
	go q.PushBack(0)

	for q.Writers() != 5 {
		time.Sleep(time.Millisecond)
	}

	_, ok := q.PopFrontNoLock()
	assert.Assert(t, ok == 0, ok)

	_, ok = q.PopFrontNoLock()
	assert.Assert(t, ok == 0, ok)

	_, ok = q.PopFrontNoLock()
	assert.Assert(t, ok == 0, ok)

	_, ok = q.PopFrontNoLock()
	assert.Assert(t, ok == 0, ok)

	_, ok = q.PopFrontNoLock()
	assert.Assert(t, ok == 0, ok)
}

func TestWritersWaitReaders4(t *testing.T) {
	q := NewSync[int32](0)

	go q.PushFront(4)
	go q.PushFront(3)
	go q.PushFront(2)
	go q.PushFront(1)
	go q.PushFront(0)

	for q.Writers() != 5 {
		time.Sleep(time.Millisecond)
	}

	_, ok := q.PopBackNoLock()
	assert.Assert(t, ok == 0, ok)

	_, ok = q.PopBackNoLock()
	assert.Assert(t, ok == 0, ok)

	_, ok = q.PopBackNoLock()
	assert.Assert(t, ok == 0, ok)

	_, ok = q.PopBackNoLock()
	assert.Assert(t, ok == 0, ok)

	_, ok = q.PopBackNoLock()
	assert.Assert(t, ok == 0, ok)
}

func PopBack(q Queue[int32], t *testing.T) {
	_, ok := q.PopBack()
	assert.Assert(t, ok == 0, ok)
}

func PopFront(q Queue[int32], t *testing.T) {
	_, ok := q.PopFront()
	assert.Assert(t, ok == 0, ok)
}

func TestReadersWaitWriters1(t *testing.T) {
	q := NewSync[int32](0)

	go PopBack(q, t)
	go PopBack(q, t)
	go PopBack(q, t)
	go PopBack(q, t)
	go PopBack(q, t)

	for q.Readers() != 5 {
		time.Sleep(time.Millisecond)
	}

	ok := q.PushFront(1)
	assert.Assert(t, ok == 0, ok)

	ok = q.PushFront(2)
	assert.Assert(t, ok == 0, ok)

	ok = q.PushFront(3)
	assert.Assert(t, ok == 0, ok)

	ok = q.PushFront(4)
	assert.Assert(t, ok == 0, ok)

	ok = q.PushFront(5)
	assert.Assert(t, ok == 0, ok)
}

func TestReadersWaitWriters2(t *testing.T) {
	q := NewSync[int32](0)

	go PopFront(q, t)
	go PopFront(q, t)
	go PopFront(q, t)
	go PopFront(q, t)
	go PopFront(q, t)

	for q.Readers() != 5 {
		time.Sleep(time.Millisecond)
	}

	ok := q.PushBack(1)
	assert.Assert(t, ok == 0, ok)

	ok = q.PushBack(2)
	assert.Assert(t, ok == 0, ok)

	ok = q.PushBack(3)
	assert.Assert(t, ok == 0, ok)

	ok = q.PushBack(4)
	assert.Assert(t, ok == 0, ok)

	ok = q.PushBack(5)
	assert.Assert(t, ok == 0, ok)
}

func TestReadersWaitWriters3(t *testing.T) {
	q := NewSync[int32](0)

	go PopBack(q, t)
	go PopBack(q, t)
	go PopBack(q, t)
	go PopBack(q, t)
	go PopBack(q, t)

	for q.Readers() != 5 {
		time.Sleep(time.Millisecond)
	}

	ok := q.PushFrontNoLock(1)
	assert.Assert(t, ok == 0, ok)

	ok = q.PushFrontNoLock(2)
	assert.Assert(t, ok == 0, ok)

	ok = q.PushFrontNoLock(3)
	assert.Assert(t, ok == 0, ok)

	ok = q.PushFrontNoLock(4)
	assert.Assert(t, ok == 0, ok)

	ok = q.PushFrontNoLock(5)
	assert.Assert(t, ok == 0, ok)
}

func TestReadersWaitWriters4(t *testing.T) {
	q := NewSync[int32](0)

	go PopFront(q, t)
	go PopFront(q, t)
	go PopFront(q, t)
	go PopFront(q, t)
	go PopFront(q, t)

	for q.Readers() != 5 {
		time.Sleep(time.Millisecond)
	}

	ok := q.PushBackNoLock(1)
	assert.Assert(t, ok == 0, ok)

	ok = q.PushBackNoLock(2)
	assert.Assert(t, ok == 0, ok)

	ok = q.PushBackNoLock(3)
	assert.Assert(t, ok == 0, ok)

	ok = q.PushBackNoLock(4)
	assert.Assert(t, ok == 0, ok)

	ok = q.PushBackNoLock(5)
	assert.Assert(t, ok == 0, ok)
}

func TestNoLock1(t *testing.T) {
	q := NewSync[int](5)

	ok := q.PushBackNoLock(1)
	assert.Assert(t, ok == 0, ok)
	ok = q.PushBackNoLock(2)
	assert.Assert(t, ok == 0, ok)
	ok = q.PushBackNoLock(3)
	assert.Assert(t, ok == 0, ok)
	ok = q.PushBackNoLock(4)
	assert.Assert(t, ok == 0, ok)
	ok = q.PushBackNoLock(5)
	assert.Assert(t, ok == 0, ok)

	temp, ok := q.PopFrontNoLock()
	assert.Assert(t, ok == 0 && temp == 1, ok)
	temp, ok = q.PopFrontNoLock()
	assert.Assert(t, ok == 0 && temp == 2, ok)
	temp, ok = q.PopFrontNoLock()
	assert.Assert(t, ok == 0 && temp == 3, ok)
	temp, ok = q.PopFrontNoLock()
	assert.Assert(t, ok == 0 && temp == 4, ok)
	temp, ok = q.PopFrontNoLock()
	assert.Assert(t, ok == 0 && temp == 5, ok)
}

func TestNoLock2(t *testing.T) {
	q := NewSync[int](0)

	ok := q.PushBackNoLock(1)
	assert.Assert(t, ok == 1, ok)

	temp, ok := q.PopFrontNoLock()
	assert.Assert(t, ok == 1 && temp == 0, ok)
}

func TestNoLock3(t *testing.T) {
	q := NewSync[int](0)

	ok := q.PushFrontNoLock(1)
	assert.Assert(t, ok == 1, ok)

	temp, ok := q.PopBackNoLock()
	assert.Assert(t, ok == 1 && temp == 0, ok)
}

func TestSize1(t *testing.T) {
	q := NewSync[string](2)

	ok := q.PushBack("lalala")
	assert.Assert(t, ok == 0, ok)
	ok = q.PushFront("bububu")
	assert.Assert(t, ok == 0, ok)
	ok = q.PushBackNoLock("kukuku")
	assert.Assert(t, ok == 1, ok)
	ok = q.PushFrontNoLock("jujuju")
	assert.Assert(t, ok == 1, ok)

	i, _ := q.PopBack()
	assert.Assert(t, i == "lalala", i)
	i, _ = q.PopBack()
	assert.Assert(t, i == "bububu", i)
}

func TestClose1(t *testing.T) {
	q := NewSync[string](4)

	ok := q.PushBack("lalala")
	assert.Assert(t, ok == 0, ok)
	ok = q.PushFront("bububu")
	assert.Assert(t, ok == 0, ok)
	q.Close()
	ok = q.PushBack("kukuku")
	assert.Assert(t, ok == -1, ok)
	ok = q.PushFront("jujuju")
	assert.Assert(t, ok == -1, ok)

	i, _ := q.PopBack()
	assert.Assert(t, i == "lalala", i)
	i, _ = q.PopBack()
	assert.Assert(t, i == "bububu", i)
}

func TestClose2(t *testing.T) {
	q := NewSync[string](2)

	ok := q.PushBack("lalala")
	assert.Assert(t, ok == 0, ok)
	ok = q.PushFront("bububu")
	assert.Assert(t, ok == 0, ok)

	ok = q.PushBackNoLock("lalala")
	assert.Assert(t, ok == 1, ok)
	ok = q.PushFrontNoLock("bububu")
	assert.Assert(t, ok == 1, ok)

	q.Close()
	ok = q.PushBack("kukuku")
	assert.Assert(t, ok == -1, ok)
	ok = q.PushFront("jujuju")
	assert.Assert(t, ok == -1, ok)

	i, _ := q.PopBack()
	assert.Assert(t, i == "lalala", i)
	i, _ = q.PopBack()
	assert.Assert(t, i == "bububu", i)
}

func Benchmark_queue1(b *testing.B) {
	b.ReportAllocs()

	q := NewSync[string](b.N)
	for i := 0; i < b.N; i++ {
		q.PushBack("lalala")
	}
}

func Benchmark_queue2(b *testing.B) {
	b.ReportAllocs()

	q := NewSync[string](b.N)

	b.RunParallel(func(pb *testing.PB) {
		var ok int
		for pb.Next() {
			if ok = q.PushBack("lalala"); ok != 0 {
				b.Fatal("WRITE ERROR")
			}
		}
	})

	b.RunParallel(func(pb *testing.PB) {
		var ok int
		for pb.Next() {
			if _, ok = q.PopFront(); ok != 0 {
				b.Fatal("READ ERROR")
			}
		}
	})
}

func Benchmark_queue3(b *testing.B) {
	b.ReportAllocs()

	q := NewSync[string](b.N)

	b.RunParallel(func(pb *testing.PB) {
		var ok int
		for pb.Next() {
			if ok = q.PushFront("lalala"); ok != 0 {
				b.Fatal("WRITE ERROR")
			}
		}
	})

	b.RunParallel(func(pb *testing.PB) {
		var ok int
		for pb.Next() {
			if _, ok = q.PopBack(); ok != 0 {
				b.Fatal("READ ERROR")
			}
		}
	})
}

func Benchmark_channel1(b *testing.B) {
	b.ReportAllocs()

	q := make(chan interface{}, b.N)
	for i := 0; i < b.N; i++ {
		q <- "lalala"
	}
}

func Benchmark_channel2(b *testing.B) {
	b.ReportAllocs()

	q := make(chan interface{}, b.N)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			q <- "lalala"
		}
	})

	b.RunParallel(func(pb *testing.PB) {
		var temp interface{}
		for pb.Next() {
			temp = <-q
		}
		_ = temp
	})
}
