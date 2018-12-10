package buffer

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSingleWriteRead(t *testing.T) {
	b := NewBuffer(5)
	r := b.NewReader()

	go func() {
		bw, err := b.Write([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
		assert.Equal(t, bw, 10, "expect to write 10 bytes")
		assert.Nil(t, err)
	}()

	buf := make([]byte, 512)
	br, err := r.Read(buf)
	assert.Equal(t, br, 10, "expected to read 10 bytes")
	assert.Nil(t, err)
}

func TestWriteWithoutReader(t *testing.T) {
	b := NewBuffer(5)

	var i int
	for i = 0; i < 10; i++ {
		bw, err := b.Write([]byte{byte(i)})
		assert.Equal(t, bw, 1, "expect to write 1 byte")
		assert.Nil(t, err)
	}

	assert.Equal(t, i, 10, "expected to have written 10 times")
}

func TestWriteOverflowAndRead(t *testing.T) {
	b := NewBuffer(5)

	var i int
	for i = 0; i < 10; i++ {
		bw, err := b.Write([]byte{byte(i)})
		assert.Equal(t, bw, 1, "expected to write 1 byte")
		assert.Nil(t, err)
	}

	r := b.NewReader()
	for j := 0; j < 5; j++ {
		buf := make([]byte, 512)
		br, err := r.Read(buf)
		assert.Equal(t, br, 1, "expected to read 1 byte")
		assert.Nil(t, err)
	}

	acq := r.(*lector).semaphore.TryAcquire(1)
	assert.Equal(t, acq, false, "expected semaphore to be completed")
}
