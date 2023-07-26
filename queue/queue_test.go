// Copyright (c) 2013-2017, Peter H. Froehlich. All rights reserved.
// Use of this source code is governed by a BSD-style license
// that can be found in the LICENSE file.

package queue

import (
	"bytes"
	"container/list"
	"fmt"
	"math/rand"
	"testing"
)

func ensureEmpty(t *testing.T, q *Queue[*int]) {
	if l := q.Len(); l != 0 {
		t.Errorf("q.Len() = %d, want %d", l, 0)
	}
	if e := q.Front(); e != nil {
		t.Errorf("q.Front() = %v, want %v", e, nil)
	}
	if e := q.Back(); e != nil {
		t.Errorf("q.Back() = %v, want %v", e, nil)
	}
	if e := q.PopFront(); e != nil {
		t.Errorf("q.PopFront() = %v, want %v", e, nil)
	}
	if e := q.PopBack(); e != nil {
		t.Errorf("q.PopBack() = %v, want %v", e, nil)
	}
}

func TestNew(t *testing.T) {
	q := New[*int]()
	ensureEmpty(t, q)
}

func ensureSingleton(t *testing.T, q *Queue[*int]) {
	if l := q.Len(); l != 1 {
		t.Errorf("q.Len() = %d, want %d", l, 1)
	}
	if e := q.Front(); *e != 42 {
		t.Errorf("q.Front() = %v, want %v", *e, 42)
	}
	if e := q.Back(); *e != 42 {
		t.Errorf("q.Back() = %v, want %v", *e, 42)
	}
}

func TestSingleton(t *testing.T) {
	i := 42
	q := New[*int]()
	ensureEmpty(t, q)
	q.PushFront(&i)
	ensureSingleton(t, q)
	q.PopFront()
	ensureEmpty(t, q)
	q.PushBack(&i)
	ensureSingleton(t, q)
	q.PopBack()
	ensureEmpty(t, q)
	q.PushFront(&i)
	ensureSingleton(t, q)
	q.PopBack()
	ensureEmpty(t, q)
	q.PushBack(&i)
	ensureSingleton(t, q)
	q.PopFront()
	ensureEmpty(t, q)
}

func TestDuos(t *testing.T) {
	i := 42
	j := 43
	q := New[*int]()
	ensureEmpty(t, q)
	q.PushFront(&i)
	ensureSingleton(t, q)
	q.PushBack(&j)
	if l := q.Len(); l != 2 {
		t.Errorf("q.Len() = %d, want %d", l, 2)
	}
	if e := q.Front(); *e != 42 {
		t.Errorf("q.Front() = %v, want %v", *e, 42)
	}
	if e := q.Back(); *e != 43 {
		t.Errorf("q.Back() = %v, want %v", *e, 43)
	}
}

func ensureLength(t *testing.T, q *Queue[*int], len int) {
	if l := q.Len(); l != len {
		t.Errorf("q.Len() = %d, want %d", l, len)
	}
}

func TestZeroValue(t *testing.T) {
	v := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	var q Queue[*int]
	q.PushFront(&v[0])
	ensureLength(t, &q, 1)
	q.PushFront(&v[1])
	ensureLength(t, &q, 2)
	q.PushFront(&v[2])
	ensureLength(t, &q, 3)
	q.PushFront(&v[3])
	ensureLength(t, &q, 4)
	q.PushFront(&v[4])
	ensureLength(t, &q, 5)

	q.PushBack(&v[5])
	ensureLength(t, &q, 6)
	q.PushBack(&v[6])
	ensureLength(t, &q, 7)
	q.PushBack(&v[7])
	ensureLength(t, &q, 8)
	q.PushBack(&v[8])
	ensureLength(t, &q, 9)

	// Output should be: [5 4 3 2 1 6 7 8 9]
	wantIdx := []int{4, 3, 2, 1, 0, 5, 6, 7, 8}
	var want bytes.Buffer
	want.WriteByte('[')
	for k, i := range wantIdx {
		want.WriteString(fmt.Sprintf("%v", &v[i]))
		if k < len(wantIdx)-1 {
			want.WriteByte(' ')
		}
	}
	want.WriteByte(']')
	if s := q.String(); s != want.String() {
		t.Errorf("q.String() = %s, want %s", s, want.String())
	}
}

func TestGrowShrink1(t *testing.T) {
	var q Queue[*int]
	for i := 0; i < size; i++ {
		i := i
		q.PushBack(&i)
		ensureLength(t, &q, i+1)
	}
	for i := 0; q.Len() > 0; i++ {
		x := q.PopFront()
		if *x != i {
			t.Errorf("q.PopFront() = %d, want %d", x, i)
		}
		ensureLength(t, &q, size-i-1)
	}
}
func TestGrowShrink2(t *testing.T) {
	var q Queue[*int]
	for i := 0; i < size; i++ {
		i := i
		q.PushFront(&i)
		ensureLength(t, &q, i+1)
	}
	for i := 0; q.Len() > 0; i++ {
		x := q.PopBack()
		if *x != i {
			t.Errorf("q.PopBack() = %d, want %d", *x, i)
		}
		ensureLength(t, &q, size-i-1)
	}
}

const size = 1024

func BenchmarkPushFrontQueue(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var q Queue[int]
		for n := 0; n < size; n++ {
			q.PushFront(n)
		}
	}
}
func BenchmarkPushFrontList(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var q list.List
		for n := 0; n < size; n++ {
			q.PushFront(n)
		}
	}
}

func BenchmarkPushBackQueue(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var q Queue[int]
		for n := 0; n < size; n++ {
			q.PushBack(n)
		}
	}
}
func BenchmarkPushBackList(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var q list.List
		for n := 0; n < size; n++ {
			q.PushBack(n)
		}
	}
}
func BenchmarkPushBackChannel(b *testing.B) {
	for i := 0; i < b.N; i++ {
		q := make(chan int, size)
		for n := 0; n < size; n++ {
			q <- n
		}
		close(q)
	}
}

var rands []float32

func makeRands() {
	if rands != nil {
		return
	}
	rando := rand.New(rand.NewSource(64738))
	for i := 0; i < 4*size; i++ {
		rands = append(rands, rando.Float32())
	}
}
func BenchmarkRandomQueue(b *testing.B) {
	makeRands()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var q Queue[int]
		for n := 0; n < 4*size; n += 4 {
			if rands[n] < 0.8 {
				q.PushBack(n)
			}
			if rands[n+1] < 0.8 {
				q.PushFront(n)
			}
			if rands[n+2] < 0.5 {
				q.PopFront()
			}
			if rands[n+3] < 0.5 {
				q.PopBack()
			}
		}
	}
}
func BenchmarkRandomList(b *testing.B) {
	makeRands()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var q list.List
		for n := 0; n < 4*size; n += 4 {
			if rands[n] < 0.8 {
				q.PushBack(n)
			}
			if rands[n+1] < 0.8 {
				q.PushFront(n)
			}
			if rands[n+2] < 0.5 {
				if e := q.Front(); e != nil {
					q.Remove(e)
				}
			}
			if rands[n+3] < 0.5 {
				if e := q.Back(); e != nil {
					q.Remove(e)
				}
			}
		}
	}
}

func BenchmarkGrowShrinkQueue(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var q Queue[int]
		for n := 0; n < size; n++ {
			q.PushBack(i)
		}
		for n := 0; n < size; n++ {
			q.PopFront()
		}
	}
}
func BenchmarkGrowShrinkList(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var q list.List
		for n := 0; n < size; n++ {
			q.PushBack(i)
		}
		for n := 0; n < size; n++ {
			if e := q.Front(); e != nil {
				q.Remove(e)
			}
		}
	}
}
