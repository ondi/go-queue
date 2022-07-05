package queue

import (
	"sync/atomic"
	"testing"
	"time"

	"gotest.tools/assert"
)

func TestQueue1(t *testing.T) {
	var i string
	var ok int

	q := New[string](2)
	ok = q.PushBack("lalala")
	assert.Assert(t, ok == 0, ok)
	ok = q.PushFront("bububu")
	assert.Assert(t, ok == 0, ok)
	ok = q.PushBackNoWait("kukuku")
	assert.Assert(t, ok == 1, ok)
	ok = q.PushFrontNoWait("jujuju")
	assert.Assert(t, ok == 1, ok)

	i, _ = q.PopBack()
	assert.Assert(t, i == "lalala", i)
	i, _ = q.PopBack()
	assert.Assert(t, i == "bububu", i)
}

func TestQueue2(t *testing.T) {
	var i string
	var ok int

	q := New[string](4)
	ok = q.PushBack("lalala")
	assert.Assert(t, ok == 0, ok)
	ok = q.PushFront("bububu")
	assert.Assert(t, ok == 0, ok)
	q.Close()
	ok = q.PushBack("kukuku")
	assert.Assert(t, ok == -1, ok)
	ok = q.PushFront("jujuju")
	assert.Assert(t, ok == -1, ok)

	i, _ = q.PopBack()
	assert.Assert(t, i == "lalala", i)
	i, _ = q.PopBack()
	assert.Assert(t, i == "bububu", i)
}

func TestQueue3(t *testing.T) {
	var i string
	var ok int

	q := New[string](2)
	ok = q.PushBack("lalala")
	assert.Assert(t, ok == 0, ok)
	ok = q.PushFront("bububu")
	assert.Assert(t, ok == 0, ok)

	ok = q.PushBackNoWait("lalala")
	assert.Assert(t, ok == 1, ok)
	ok = q.PushFrontNoWait("bububu")
	assert.Assert(t, ok == 1, ok)

	q.Close()
	ok = q.PushBack("kukuku")
	assert.Assert(t, ok == -1, ok)
	ok = q.PushFront("jujuju")
	assert.Assert(t, ok == -1, ok)

	i, _ = q.PopBack()
	assert.Assert(t, i == "lalala", i)
	i, _ = q.PopBack()
	assert.Assert(t, i == "bububu", i)
}

func Push(q Queue[int32], current *int32, my int32) {
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
	q := New[int32](0)
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
	assert.Assert(t, ok == 0, ok)
	assert.Assert(t, temp == 0, temp)

	temp, ok = q.PopFrontNoWait()
	assert.Assert(t, ok == 0, ok)
	assert.Assert(t, temp == 1, temp)

	temp, ok = q.PopFrontNoWait()
	assert.Assert(t, ok == 0, ok)
	assert.Assert(t, temp == 2, temp)

	temp, ok = q.PopFrontNoWait()
	assert.Assert(t, ok == 0, ok)
	assert.Assert(t, temp == 3, temp)

	temp, ok = q.PopFrontNoWait()
	assert.Assert(t, ok == 0, ok)
	assert.Assert(t, temp == 4, temp)
}

func TestQueue5(t *testing.T) {
	q := New[string](0)

	go func() {
		ok := q.PushBack("lalala")
		assert.Assert(t, ok == 0, ok)
	}()

	i, ok := q.PopBack()
	assert.Assert(t, ok == 0, ok)
	assert.Assert(t, i == "lalala", i)
}

func TestQueue6(t *testing.T) {
	q := New[string](0)

	go func() {
		ok := q.PushFront("lalala")
		assert.Assert(t, ok == 0, ok)
	}()

	i, ok := q.PopFront()
	assert.Assert(t, ok == 0, ok)
	assert.Assert(t, i == "lalala", i)
}

func Benchmark_queue1(b *testing.B) {
	b.ReportAllocs()

	q := New[string](b.N)
	for i := 0; i < b.N; i++ {
		q.PushBack("lalala")
	}
}

func Benchmark_queue2(b *testing.B) {
	b.ReportAllocs()

	q := New[string](b.N)

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

	q := New[string](b.N)

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
