package queue

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"
)

func Example_queue1() {
	var i interface{}
	var ok int

	q := New(2)
	ok = q.PushBack("lalala")
	fmt.Printf("%v\n", ok)
	ok = q.PushFront("bububu")
	fmt.Printf("%v\n", ok)
	ok = q.PushBackNoWait("kukuku")
	fmt.Printf("%v\n", ok)
	ok = q.PushFrontNoWait("jujuju")
	fmt.Printf("%v\n", ok)

	i, _ = q.PopBack()
	fmt.Printf("%v\n", i)
	i, _ = q.PopBack()
	fmt.Printf("%v\n", i)
	// Output:
	// 0
	// 0
	// 1
	// 1
	// lalala
	// bububu
}

func Example_queue2() {
	var i interface{}
	var ok int

	q := New(4)
	ok = q.PushBack("lalala")
	fmt.Printf("%v\n", ok)
	ok = q.PushFront("bububu")
	fmt.Printf("%v\n", ok)
	q.Close()
	ok = q.PushBack("kukuku")
	fmt.Printf("%v\n", ok)
	ok = q.PushFront("jujuju")
	fmt.Printf("%v\n", ok)

	i, _ = q.PopBack()
	fmt.Printf("%v\n", i)
	i, _ = q.PopBack()
	fmt.Printf("%v\n", i)
	// Output:
	// 0
	// 0
	// -1
	// -1
	// lalala
	// bububu
}

func Example_queue3() {
	var i interface{}
	var ok int

	q := New(2)
	ok = q.PushBack("lalala")
	fmt.Printf("%v\n", ok)
	ok = q.PushFront("bububu")
	fmt.Printf("%v\n", ok)

	ok = q.PushBackNoWait("lalala")
	fmt.Printf("%v\n", ok)
	ok = q.PushFrontNoWait("bububu")
	fmt.Printf("%v\n", ok)

	q.Close()
	ok = q.PushBack("kukuku")
	fmt.Printf("%v\n", ok)
	ok = q.PushFront("jujuju")
	fmt.Printf("%v\n", ok)

	i, _ = q.PopBack()
	fmt.Printf("%v\n", i)
	i, _ = q.PopBack()
	fmt.Printf("%v\n", i)
	// Output:
	// 0
	// 0
	// 1
	// 1
	// -1
	// -1
	// lalala
	// bububu
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

func Example_queue4() {
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
	fmt.Printf("%v %v\n", ok, temp)

	temp, ok = q.PopFrontNoWait()
	fmt.Printf("%v %v\n", ok, temp)

	temp, ok = q.PopFrontNoWait()
	fmt.Printf("%v %v\n", ok, temp)

	temp, ok = q.PopFrontNoWait()
	fmt.Printf("%v %v\n", ok, temp)

	temp, ok = q.PopFrontNoWait()
	fmt.Printf("%v %v\n", ok, temp)
	// Output:
	// 0 0
	// 0 1
	// 0 2
	// 0 3
	// 0 4
}

func Test_queue3(t *testing.T) {
	q := New(0)

	go func() {
		ok := q.PushBack("lalala")
		if ok != 0 {
			t.Fatalf("PushBack: %v", ok)
		}
	}()

	i, ok := q.PopBack()
	if ok != 0 || i.(string) != "lalala" {
		t.Fatalf("PopBack: %v %v", i, ok)
	}
}

func Test_queue4(t *testing.T) {
	q := New(0)

	go func() {
		ok := q.PushFront("lalala")
		if ok != 0 {
			t.Fatalf("PushBack: %v", ok)
		}
	}()

	i, ok := q.PopFront()
	if ok != 0 || i.(string) != "lalala" {
		t.Fatalf("PopBack: %v %v", i, ok)
	}
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
