package buffer

import (
	"container/list"
	"container/ring"
	"io"
)

type semaphore chan int

// acquire n resources
func (s semaphore) Acquire(n int) {
	e := n
	for i := 0; i < n; i++ {
		s <- e
	}
}

// release n resources
func (s semaphore) Release(n int) {
	for i := 0; i < n; i++ {
		<-s
	}
}

type Buffer struct {
	bufferSize   int
	writePos     int
	readers      *list.List
	WritePointer *ring.Ring
}

func NewBuffer(elements int) Buffer {
	buf := Buffer{
		bufferSize:   elements,
		readers:      list.New(),
		WritePointer: ring.New(elements),
	}

	return buf
}

func (b *Buffer) Write(p []byte) (n int, err error) {
	b.incrementAllReaders()
	b.WritePointer.Value = p
	b.WritePointer = b.WritePointer.Next()
	b.writePos++
	return len(p), nil
}

func (b *Buffer) incrementAllReaders() {
	for e := b.readers.Front(); e != nil; e = e.Next() {
		e.Value.(*lector).semaphore.Release(1)
	}
}

func (b *Buffer) NewReader() io.Reader {
	l := &lector{
		semaphore:   make(semaphore, b.bufferSize),
		ReadPointer: b.WritePointer,
	}
	b.readers.PushBack(l)
	return l
}

// ------
// Reader
// ------

type lector struct {
	semaphore   semaphore
	ReadPointer *ring.Ring
}

func (l *lector) Read(p []byte) (n int, err error) {
	l.semaphore.Acquire(1)
	copied := copy(p, l.ReadPointer.Value.([]byte))
	l.ReadPointer = l.ReadPointer.Next()

	return copied, nil
}
