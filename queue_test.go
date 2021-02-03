package queue

import (
	"sync/atomic"
	"testing"
	"time"

	"gotest.tools/assert"
)

func TestQueue1(t *testing.T) {
	var i interface{}
	var ok int

	q := New(2)
	ok = q.PushBack("lalala")
	assert.Assert(t, ok == 0)
	ok = q.PushFront("bububu")
	assert.Assert(t, ok == 0)
	ok = q.PushBackNoWait("kukuku")
	assert.Assert(t, ok == 1)
	ok = q.PushFrontNoWait("jujuju")
	assert.Assert(t, ok == 1)

	i, _ = q.PopBack()
	assert.Assert(t, i == "lalala")
	i, _ = q.PopBack()
	assert.Assert(t, i == "bububu")
}

func TestQueue2(t *testing.T) {
	var i interface{}
	var ok int

	q := New(4)
	ok = q.PushBack("lalala")
	assert.Assert(t, ok == 0)
	ok = q.PushFront("bububu")
	assert.Assert(t, ok == 0)
	q.Close()
	ok = q.PushBack("kukuku")
	assert.Assert(t, ok == -1)
	ok = q.PushFront("jujuju")
	assert.Assert(t, ok == -1)

	i, _ = q.PopBack()
	assert.Assert(t, i == "lalala")
	i, _ = q.PopBack()
	assert.Assert(t, i == "bububu")
}

func TestQueue3(t *testing.T) {
	var i interface{}
	var ok int

	q := New(2)
	ok = q.PushBack("lalala")
	assert.Assert(t, ok == 0)
	ok = q.PushFront("bububu")
	assert.Assert(t, ok == 0)

	ok = q.PushBackNoWait("lalala")
	assert.Assert(t, ok == 1)
	ok = q.PushFrontNoWait("bububu")
	assert.Assert(t, ok == 1)

	q.Close()
	ok = q.PushBack("kukuku")
	assert.Assert(t, ok == -1)
	ok = q.PushFront("jujuju")
	assert.Assert(t, ok == -1)

	i, _ = q.PopBack()
	assert.Assert(t, i == "lalala")
	i, _ = q.PopBack()
	assert.Assert(t, i == "bububu")
}

func Push(q Queue, current *int32, my int32) {
	for {
		if atomic.LoadInt32(current) == my {
			atomic.AddInt32(current, 1)
			q.PushBack(my)
			return
		}
		time.Sleep(time.Millisecond)
	}
}

func TestQueue4(t *testing.T) {
	q := New(0)
	var count int32
	go Push(q, &count, 4)
	go Push(q, &count, 3)
	go Push(q, &count, 2)
	go Push(q, &count, 1)
	go Push(q, &count, 0)

	for atomic.LoadInt32(&count) != 5 {
		time.Sleep(time.Millisecond)
	}

	temp, ok := q.PopFrontNoWait()
	assert.Assert(t, ok == 0)
	assert.Assert(t, temp.(int32) == 0)

	temp, ok = q.PopFrontNoWait()
	assert.Assert(t, ok == 0)
	assert.Assert(t, temp.(int32) == 1)

	temp, ok = q.PopFrontNoWait()
	assert.Assert(t, ok == 0)
	assert.Assert(t, temp.(int32) == 2)

	temp, ok = q.PopFrontNoWait()
	assert.Assert(t, ok == 0)
	assert.Assert(t, temp.(int32) == 3)

	temp, ok = q.PopFrontNoWait()
	assert.Assert(t, ok == 0)
	assert.Assert(t, temp.(int32) == 4)
}

func TestQueue5(t *testing.T) {
	q := New(0)

	go func() {
		ok := q.PushBack("lalala")
		assert.Assert(t, ok == 0)
	}()

	i, ok := q.PopBack()
	assert.Assert(t, ok == 0)
	assert.Assert(t, i == "lalala")
}

func TestQueue6(t *testing.T) {
	q := New(0)

	go func() {
		ok := q.PushFront("lalala")
		assert.Assert(t, ok == 0)
	}()

	i, ok := q.PopFront()
	assert.Assert(t, ok == 0)
	assert.Assert(t, i == "lalala")
}

func Benchmark_queue1(b *testing.B) {
	b.ReportAllocs()

	q := New(b.N)
	for i := 0; i < b.N; i++ {
		q.PushBack("lalala")
	}
}

func Benchmark_queue2(b *testing.B) {
	b.ReportAllocs()

	q := New(b.N)

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

	q := New(b.N)

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
