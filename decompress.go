package sarama

import (
	"bytes"
	"fmt"
	"io"
	"sync"

	"github.com/IBM/sarama/mickey/bytepool"
	snappy "github.com/eapache/go-xerial-snappy"
	"github.com/klauspost/compress/gzip"
	"github.com/pierrec/lz4/v4"
)

var (
	lz4ReaderPool = sync.Pool{
		New: func() interface{} {
			return lz4.NewReader(nil)
		},
	}

	gzipReaderPool sync.Pool
)

func decompress(cc CompressionCodec, data []byte) ([]byte, error) {
	switch cc {
	case CompressionNone:
		return data, nil
	case CompressionGZIP:
		var err error
		reader, ok := gzipReaderPool.Get().(*gzip.Reader)
		if !ok {
			reader, err = gzip.NewReader(bytes.NewReader(data))
		} else {
			err = reader.Reset(bytes.NewReader(data))
		}

		if err != nil {
			return nil, err
		}

		defer gzipReaderPool.Put(reader)

		if bytepool.UseDecompressBytePoolSwitch.Enabled() {
			return bytepool.ReadAll(reader)
		}

		return io.ReadAll(reader)
	case CompressionSnappy:
		return snappy.Decode(data)
	case CompressionLZ4:
		reader, ok := lz4ReaderPool.Get().(*lz4.Reader)
		if !ok {
			reader = lz4.NewReader(bytes.NewReader(data))
		} else {
			reader.Reset(bytes.NewReader(data))
		}
		defer lz4ReaderPool.Put(reader)

		if bytepool.UseDecompressBytePoolSwitch.Enabled() {
			return bytepool.ReadAll(reader)
		}

		return io.ReadAll(reader)
	case CompressionZSTD:
		return zstdDecompress(ZstdDecoderParams{}, nil, data)
	default:
		return nil, PacketDecodingError{fmt.Sprintf("invalid compression specified (%d)", cc)}
	}
}
