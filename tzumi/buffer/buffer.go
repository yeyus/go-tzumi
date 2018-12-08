package buffer

import (
	"container/list"
	"container/ring"
	"context"
	"golang.org/x/sync/semaphore"
	"io"
)

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
		semaphore:   semaphore.NewWeighted(int64(b.bufferSize)),
		ReadPointer: b.WritePointer,
	}
	b.readers.PushBack(l)
	return l
}

// ------
// Reader
// ------

type lector struct {
	semaphore   *semaphore.Weighted
	ReadPointer *ring.Ring
}

func (l *lector) Read(p []byte) (n int, err error) {
Start:
	err = l.semaphore.Acquire(context.TODO(), 1)
	if err != nil {
		return 0, err
	}

	if l.ReadPointer.Value == nil {
		goto Start
	}
	copied := copy(p, l.ReadPointer.Value.([]byte))
	l.ReadPointer = l.ReadPointer.Next()

	return copied, nil
}
