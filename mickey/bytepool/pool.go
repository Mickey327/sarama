package bytepool

import (
	"io"
	"sync"

	"github.com/IBM/sarama/mickey/boolswitch"
)

const (
	defaultChunkCapacity    = 120000
	defaultChunkMaxCapacity = 200000
)

var (
	globalBytePool = newBytePool(
		defaultChunkCapacity,
		defaultChunkMaxCapacity,
	)
	UseDecompressBytePoolSwitch = boolswitch.NewDisabled()
)

func newBytePool(chunkCapacity, chunkMaxCapacity int) *bytePool {
	if chunkCapacity <= 0 {
		chunkCapacity = defaultChunkCapacity
	}
	if chunkMaxCapacity <= 0 {
		chunkMaxCapacity = defaultChunkMaxCapacity
	}

	return &bytePool{
		chunkMaxCapacity: chunkMaxCapacity,
		pool: &sync.Pool{
			New: func() interface{} { return make([]byte, 0, chunkCapacity) },
		},
	}
}

type bytePool struct {
	chunkMaxCapacity int
	pool             *sync.Pool
}

func (b *bytePool) Get() []byte {
	return b.pool.Get().([]byte)
}

func (b *bytePool) Put(sl []byte) {
	if len(sl) > b.chunkMaxCapacity {
		return
	}

	sl = sl[:0]
	b.pool.Put(sl)
}

var ReadAll = newPoolReader(globalBytePool).ReadAll

type poolReader struct {
	pool *bytePool
}

func newPoolReader(pool *bytePool) *poolReader {
	return &poolReader{
		pool: pool,
	}
}

func (rd *poolReader) ReadAll(r io.Reader) ([]byte, error) {
	chunk := rd.pool.Get()
	defer rd.pool.Put(chunk)

	for {
		if len(chunk) == cap(chunk) {
			chunk = append(chunk, 0)[:len(chunk)]
		}

		n, err := r.Read(chunk[len(chunk):cap(chunk)])
		chunk = chunk[:len(chunk)+n]
		if err != nil {
			if err == io.EOF {
				err = nil
			}

			cpy := make([]byte, len(chunk))
			copy(cpy, chunk)
			return cpy, err
		}
	}
}

type BytePoolOptions struct {
	ChunkCapacity    int
	ChunkMaxCapacity int
}

func DefaultBytePoolOptions() BytePoolOptions {
	return BytePoolOptions{
		ChunkCapacity:    defaultChunkCapacity,
		ChunkMaxCapacity: defaultChunkMaxCapacity,
	}
}

func ConfigureDecompressBytePool(options BytePoolOptions) {
	globalBytePool = newBytePool(
		options.ChunkCapacity,
		options.ChunkMaxCapacity,
	)
	UseDecompressBytePoolSwitch.Enable()
}
